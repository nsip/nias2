// webserver.go
package sms

// handles all web interactions with users

import (
	"bufio"
	"io/ioutil"
	//gcsv "encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo"
	//"github.com/labstack/echo/engine/fasthttp"
	"bytes"
	//"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	//ms "github.com/mitchellh/mapstructure"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nuid"
	"github.com/nsip/nias2/lib"
	"github.com/twinj/uuid"
	//"github.com/wildducktheories/go-csv"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	//"time"
)

var defaultconfig = lib.LoadDefaultConfig()
var VALIDATION_ROUTE = defaultconfig.ValidationRoute
var SSF_ROUTE = defaultconfig.SSFRoute
var SMS_ROUTE = defaultconfig.SMSRoute
var req_ec *nats.EncodedConn
var req_conn stan.Conn
var stan_conn stan.Conn
var tt *lib.TransactionTracker // = lib.NewTransactionTracker(defaultconfig.TxReportInterval)

var UI_LIMIT int

// rendering template for csv-xml conversion
var sptmpl *template.Template

type NIASWebServer struct{}

// standard response to successful file upload
type IngestResponse struct {
	TxID    string
	Records int
}

// Struct to contain single XML record
type XMLContainer struct {
	Value string `xml:",innerxml"`
}

// generic publish routine that handles different requirements
// of the 3 possible message infrastrucutres
func publish(msg *lib.NiasMessage) {

	req_ec.Publish(lib.REQUEST_TOPIC, msg)
	/*
		switch lib.DefaultConfig.MsgTransport {
		case "MEM":
			req_chan <- *msg
		case "NATS":
			req_ec.Publish(REQUEST_TOPIC, msg)
		case "STAN":
			req_conn.Publish(REQUEST_TOPIC, EncodeNiasMessage(msg))
		default:
			req_chan <- *msg
		}
	*/

}

// check header of csv file for plausibility
func checkHeaderCSVforNAPLANValidation(s string) error {
	// we have grabbed the first 1000 bytes of the csv file.
	// Grab its first line, and split by comma
	var lines []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	jsondat, _ := ioutil.ReadFile("schemas/core.json")
	var jsonschema map[string]interface{}
	jsonschema = nil
	json.Unmarshal(jsondat, &jsonschema)
	errfields := ""
	//log.Printf("%v\n", jsonschema)
	for _, field := range strings.Split(lines[0], ",") {
		_, prs := jsonschema["properties"].(map[string]interface{})[field]
		if !prs {
			errfields = errfields + " " + field
		}
	}
	if len(errfields) > 0 {
		return fmt.Errorf("%s is not an expected registration field", errfields)
	}
	return nil
}

func validateSIFXML(fileheader *multipart.FileHeader) error {
	file, _ := fileheader.Open()
	defer file.Close()
	log.Println("Validating data file...")
	filename := uuid.NewV4().String()
	tmp_xmlfile := "/tmp/" + filename + ".tmp.xml"
	log.Println(tmp_xmlfile)
	outfile, err := os.Create(tmp_xmlfile)
	if err != nil {
		return err
	}
	defer outfile.Close()
	_, err = io.Copy(outfile, file)
	if err != nil {
		return err
	}

	subProcess := exec.Command("xmllint", "--noout", "--schema", "SIF_Message.xsd", tmp_xmlfile)
	if err != nil {
		log.Printf("1: %v\n", err)
		return err
	}
	out, err := subProcess.CombinedOutput()
	if err != nil {
		log.Printf("2: %s\n", string(out))
		log.Printf("2: %v\n", err)
		return err
	}
	log.Println(string(out))
	return nil
}

// read xml file as stream and post records onto processing queue
func enqueueXML(file multipart.File, usecase string, route []string) (IngestResponse, error) {

	ir := IngestResponse{}
	v := XMLContainer{"none"}
	log.Printf("enqueueXML %v", file)
	var b bytes.Buffer
	decoder := xml.NewDecoder(file)
	encoder := xml.NewEncoder(&b)
	child := false
	total := 0
	txid := nuid.Next()
	for {
		t, _ := decoder.Token()
		log.Printf("Token %v\n", t)
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if child {
				total++

				decode_err := decoder.DecodeElement(&v, &se)
				if decode_err != nil {
					return ir, decode_err
				}
				end := se.End()
				b.Reset()
				encoder.EncodeToken(t)
				encoder.Flush()

				msg := &lib.NiasMessage{}
				msg.Body = b.String() + v.Value
				b.Reset()
				encoder.EncodeToken(end)
				encoder.Flush()
				msg.Body = msg.Body.(string) + b.String()
				msg.SeqNo = strconv.Itoa(total)
				msg.TxID = txid
				msg.MsgID = nuid.Next()
				msg.Target = usecase
				msg.Route = route
				log.Println(msg.Body)
				publish(msg)
			}
			child = true

		default:
		}
	}

	ir.Records = total
	ir.TxID = txid
	// update the tx tracker
	tt.SetTxSize(txid, total)

	return ir, nil

}

// start the server
func (nws *NIASWebServer) Run(nats_cfg lib.NATSConfig) {

	log.Println("SMS: Connecting to message bus")
	req_ec = lib.CreateNATSConnection(nats_cfg)
	log.Println("SMS: Initialising uuid generator")
	uuid.Init()
	log.Println("SMS: UUID generator initialised.")

	tt = lib.NewTransactionTracker(defaultconfig.TxReportInterval, nats_cfg)
	log.Println("SMS: Loading xml templates")
	fp := path.Join("templates", "studentpersonals.tmpl")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatalf("Unable to parse xml conversion template, service aborting...")
	}
	sptmpl = tmpl
	log.Println("SMS: XML conversion template loaded ok.")

	//setup stan connection
	stan_conn, _ = stan.Connect(lib.NAP_VAL_CID, nuid.Next())

	// create the web service framework
	e := echo.New()

	// handler for data file ingest
	e.POST("/sifxml/ingest", func(c echo.Context) error {
		// get the file from the input form
		file, err := c.FormFile("validationFile")
		if err != nil {
			log.Println("Ingest #1")
			return err
		}
		src, err := file.Open()
		if err != nil {
			log.Println("Ingest #2")
			return err
		}
		defer src.Close()

		// read onto qs with appropriate handler
		var ir IngestResponse
		if strings.Contains(file.Filename, ".xml") {
			if err = validateSIFXML(file); err != nil {
				log.Println("Ingest #3")
				return c.String(http.StatusBadRequest, err.Error())
			}
			if ir, err = enqueueXML(src, STORE_AND_FORWARD_PREFIX, SSF_ROUTE); err != nil {
				log.Println("Ingest #4")
				return err
			}
		} else {
			log.Println("Ingest #5")
			return c.String(http.StatusBadRequest, "File submitted is not .xml")
		}

		log.Println("ir: ", ir)
		return c.JSON(http.StatusAccepted, ir)
	})

	// handler for data file store as graph
	e.POST("/sifxml/store", func(c echo.Context) error {
		// get the file from the input form
		file, err := c.FormFile("validationFile")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// read onto qs with appropriate handler
		var ir IngestResponse
		if strings.Contains(file.Filename, ".xml") {
			if err = validateSIFXML(file); err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			if ir, err = enqueueXML(src, SIF_MEMORY_STORE_PREFIX, SMS_ROUTE); err != nil {
				return err
			}
		} else {
			return c.String(http.StatusBadRequest, "File submitted is not .xml")
		}

		log.Println("ir: ", ir)
		return c.JSON(http.StatusAccepted, ir)
	})

	// get filtered text
	e.GET("/sifxml/ingest/none/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"none::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/low/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"low::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/medium/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"medium::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/high/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"high::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/extreme/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"extreme::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/none", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"none::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/low", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"low::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/medium", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"medium::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/high", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"high::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.GET("/sifxml/ingest/extreme", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"extreme::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})

	// static resources
	e.Static("/", "public")

	// homepage
	e.File("/", "public/index.html")
	e.File("/nias", "public/index.html")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	log.Println("SMS: Starting web-ui services...")
	port := defaultconfig.WebServerPort
	log.Println("SMS: Service is listening on localhost:" + port)

	// set upper bound for no. messages sent to web clients
	UI_LIMIT = defaultconfig.UIMessageLimit

	//e.Run(fasthttp.New(":" + port))
	//e.Run(standard.New(":" + port))
	e.Logger.Fatal(e.Start(":" + port))

}
