// webserver.go
package nias2

// handles all web interactions with users

import (
	gcsv "encoding/csv"
	"encoding/json"
	"encoding/xml"
	"github.com/labstack/echo"
	//"github.com/labstack/echo/engine/fasthttp"
	"bytes"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	ms "github.com/mitchellh/mapstructure"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats"
	"github.com/nats-io/nuid"
	"github.com/twinj/uuid"
	"github.com/wildducktheories/go-csv"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

var VALIDATION_ROUTE = NiasConfig.ValidationRoute
var req_ec *nats.EncodedConn
var req_conn stan.Conn

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

// truncate the record by removing items that have blank entries.
// this prevents the validation from throwing validation exceptions
// for fields that are not mandatory but included as empty in the
// dataset
func removeBlanks(m map[string]string) map[string]string {

	reducedmap := make(map[string]string)
	for key, val := range m {
		if val != "" {
			reducedmap[key] = val
		}
	}
	return reducedmap
}

// generic publish routine that handles different requirements
// of the 3 possible message infrastrucutres
func publish(msg *NiasMessage) {

	switch NiasConfig.MsgTransport {
	case "MEM":
		req_chan <- *msg
	case "NATS":
		req_ec.Publish(REQUEST_TOPIC, msg)
	case "STAN":
		req_conn.Publish(REQUEST_TOPIC, EncodeNiasMessage(msg))
	default:
		req_chan <- *msg
	}

}

// read csv file as stream and post records onto processing queue
func enqueueCSV(file multipart.File) (IngestResponse, error) {

	ir := IngestResponse{}

	reader := csv.WithIoReader(file)
	defer reader.Close()

	i := 0
	txid := nuid.Next()
	for record := range reader.C() {

		i = i + 1

		regr := RegistrationRecord{}
		r := removeBlanks(record.AsMap())
		decode_err := ms.Decode(r, &regr)
		if decode_err != nil {
			return ir, decode_err
		}

		msg := &NiasMessage{}
		msg.Body = regr
		msg.SeqNo = strconv.Itoa(i)
		msg.TxID = txid
		msg.MsgID = nuid.Next()
		msg.Target = VALIDATION_PREFIX
		msg.Route = VALIDATION_ROUTE

		publish(msg)

	}

	ir.Records = i
	ir.TxID = txid

	return ir, nil

}

// read xml file as stream and post records onto processing queue
func enqueueXMLforNAPLANValidation(file multipart.File) (IngestResponse, error) {

	ir := IngestResponse{}

	decoder := xml.NewDecoder(file)
	total := 0
	txid := nuid.Next()
	var inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			if inElement == "StudentPersonal" {

				total++

				var rr RegistrationRecord
				decode_err := decoder.DecodeElement(&rr, &se)
				if decode_err != nil {
					return ir, decode_err
				}

				msg := &NiasMessage{}
				msg.Body = rr
				msg.SeqNo = strconv.Itoa(total)
				msg.TxID = txid
				msg.MsgID = nuid.Next()
				msg.Target = VALIDATION_PREFIX
				msg.Route = VALIDATION_ROUTE

				// xml_ec.Publish(REQUEST_TOPIC, msg)
				// xml_conn.Publish(REQUEST_TOPIC, EncodeNiasMessage(msg))
				// req_chan <- msg
				publish(msg)

			}
		default:
		}

	}

	ir.Records = total
	ir.TxID = txid

	return ir, nil

}

// read xml file as stream and post records onto processing queue
func enqueueXML(file multipart.File, usecase string, route []string) (IngestResponse, error) {

	ir := IngestResponse{}
	v := XMLContainer{"none"}

	var b bytes.Buffer
	decoder := xml.NewDecoder(file)
	encoder := xml.NewEncoder(&b)
	child := false
	total := 0
	txid := nuid.Next()
	for {
		t, _ := decoder.Token()
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

				msg := &NiasMessage{}
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

				publish(msg)
			}
			child = true

		default:
		}
	}

	ir.Records = total
	ir.TxID = txid
	return ir, nil

}

// start the server
func (nws *NIASWebServer) Run() {

	log.Println("Connecting to message bus")
	switch NiasConfig.MsgTransport {
	case "NATS":
		req_ec = CreateNATSConnection()
	case "STAN":
		req_conn, _ = stan.Connect(NIAS_CLUSTER_ID, nuid.Next())
	}
	sms_route := make([]string, 1)
	sms_route[0] = "sif2graph"

	log.Println("Initialising uuid generator")
	config := uuid.StateSaverConfig{SaveReport: true, SaveSchedule: 30 * time.Minute}
	uuid.SetupFileSystemStateSaver(config)
	log.Println("UUID generator initialised.")

	log.Println("Loading xml templates")
	fp := path.Join("templates", "studentpersonals.tmpl")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatalf("Unable to parse xml conversion template, service aborting...")
	}
	sptmpl = tmpl
	log.Println("XML conversion template loaded ok.")

	// create the web service framework
	e := echo.New()

	// handler for data file ingest
	e.Post("/sifxml/ingest", func(c echo.Context) error {
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
			if ir, err = enqueueXML(src, STORE_AND_FORWARD_PREFIX, nil); err != nil {
				return err
			}
		} else {
			return c.String(http.StatusBadRequest, "File submitted is not .xml")
		}

		log.Println("ir: ", ir)
		return c.JSON(http.StatusAccepted, ir)
	})

	// handler for data file store as graph
	e.Post("/sifxml/store", func(c echo.Context) error {
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
			if ir, err = enqueueXML(src, SIF_MEMORY_STORE_PREFIX, sms_route); err != nil {
				return err
			}
		} else {
			return c.String(http.StatusBadRequest, "File submitted is not .xml")
		}

		log.Println("ir: ", ir)
		return c.JSON(http.StatusAccepted, ir)
	})

	// handler for validation of NAPLAN
	e.Post("/naplan/reg/validate", func(c echo.Context) error {

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
		if strings.Contains(file.Filename, ".csv") {
			if ir, err = enqueueCSV(src); err != nil {
				return err
			}
		} else if strings.Contains(file.Filename, ".xml") {
			if ir, err = enqueueXMLforNAPLANValidation(src); err != nil {
				return err
			}
		} else {

			return c.String(http.StatusBadRequest, "File submitted is not .csv or .xml")
		}

		log.Println("ir: ", ir)
		return c.JSON(http.StatusAccepted, ir)

	})

	// handler for csv-xml conversion
	e.Post("/naplan/reg/convert", func(c echo.Context) error {

		// get the file from the input form
		file, err := c.FormFile("conversionFile")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// check it's a csv file
		if !strings.Contains(file.Filename, ".csv") {
			return c.String(http.StatusBadRequest, "File must be of type .csv")
		}

		// create outbound file name
		fname := file.Filename
		rplcr := strings.NewReplacer(".csv", ".xml")
		xml_fname := rplcr.Replace(fname)

		// read the csv file
		reader := csv.WithIoReader(src)
		records, err := csv.ReadAll(reader)
		if err != nil {
			return err
		}

		// create valid sif guids
		sprsnls := make([]map[string]string, 0)
		for _, r := range records {
			r := r.AsMap()
			r = removeBlanks(r)
			r["SIFuuid"] = uuid.NewV4().String()
			sprsnls = append(sprsnls, r)
		}

		// set headers to 'force' file download where appropriate
		c.Response().Header().Set("Content-Disposition", "attachment; filename="+xml_fname)
		c.Response().Header().Set("Content-Type", "application/xml")

		// apply the template & write results to the client
		if err := sptmpl.Execute(c.Response().Writer(), sprsnls); err != nil {
			return err
		}

		return nil

	})

	// monitoring endpoint for validation progress
	e.Get("/naplan/reg/status/:txid", func(c echo.Context) error {

		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().WriteHeader(http.StatusOK)

		txid := c.Param("txid")
		reply := GetTrackingData(txid)

		sm, _ := json.Marshal(reply)
		suffix := string(sm) + "\n\n"
		if _, err := c.Response().Write([]byte("data: " + suffix)); err != nil {
			log.Println(err)
		}

		return nil

	})

	// get validation analysis results
	e.Get("/naplan/reg/results/:txid", func(c echo.Context) error {

		msgs, err := GetTxData(c.Param("txid"), VALIDATION_PREFIX, false)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, msgs)

	})

	// get filtered text
	e.Get("/sifxml/ingest/low/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"low::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.Get("/sifxml/ingest/medium/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"medium::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.Get("/sifxml/ingest/high/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"high::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.Get("/sifxml/ingest/extreme/:txid", func(c echo.Context) error {
		msgs, err := GetTxData(c.Param("txid"), STORE_AND_FORWARD_PREFIX+"extreme::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.Get("/sifxml/ingest/low", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"low::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.Get("/sifxml/ingest/medium", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"medium::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.Get("/sifxml/ingest/high", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"high::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})
	e.Get("/sifxml/ingest/extreme", func(c echo.Context) error {
		msgs, err := GetTxData("", STORE_AND_FORWARD_PREFIX+"extreme::", false)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, msgs)
	})

	// get the validation errors data for a given transaction as a downloadable csv file
	e.Get("/naplan/reg/report/:txID/:fname", func(c echo.Context) error {

		txID := c.Param("txID")

		// get filename from params
		fname := c.Param("fname")
		rplcr := strings.NewReplacer(".csv", "_error_report.csv", ".xml", "_error_report.csv")
		rfname := rplcr.Replace(fname)

		c.Response().Header().Set("Content-Disposition", "attachment; filename="+rfname)
		c.Response().Header().Set("Content-Type", "text/csv")

		w := gcsv.NewWriter(c.Response().Writer())

		// write the headers
		hdr := []string{"Original File Line No. where error occurred",
			"Validation Type",
			"Field that failed validation",
			"Error Description"}

		if err := w.Write(hdr); err != nil {
			log.Println("error writing headers to csv:", err)
		}

		data, err := GetTxData(txID, VALIDATION_PREFIX, true)
		if err != nil {
			log.Println("Error fetching report data: ", err)
			return err
		}

		for _, record := range data {
			// cast from interface type returned by gettxdata
			ve := record.(ValidationError)
			if err := w.Write(ve.ToSlice()); err != nil {
				log.Println("error writing record to csv:", err)
			}
		}

		w.Flush()

		if err := w.Error(); err != nil {
			log.Println("Error constructing csv report:", err)
			return err
		}

		return nil

	})

	// static resources
	e.Static("/", "public")

	// homepage
	e.File("/", "public/index.html")
	e.File("/nias", "public/index.html")

	e.Use(mw.Logger())
	e.Use(mw.Recover())
	log.Println("Starting web-ui services...")
	port := NiasConfig.WebServerPort
	log.Println("Service is listening on localhost:" + port)

	//e.Run(fasthttp.New(":" + port))
	e.Run(standard.New(":1325"))

}
