// webserver.go
package napval

// handles all web interactions with users

import (
	gcsv "encoding/csv"
	"encoding/json"
	"encoding/xml"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	ms "github.com/mitchellh/mapstructure"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nuid"
	"github.com/nsip/nias2/lib"
	nxml "github.com/nsip/nias2/xml"
	"github.com/twinj/uuid"
	"github.com/wildducktheories/go-csv"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

var VALIDATION_ROUTE = DefaultValidationConfig.ValidationRoute
var req_ec *nats.EncodedConn
var req_conn stan.Conn
var tt = lib.NewTransactionTracker(DefaultValidationConfig.TxReportInterval)
var stan_conn stan.Conn

var UI_LIMIT int

// rendering template for csv-xml conversion
var sptmpl *template.Template

type ValidationWebServer struct{}

// generic publish routine that handles different requirements
// of the 3 possible message infrastrucutres
func publish(msg *lib.NiasMessage) {

	req_ec.Publish(lib.REQUEST_TOPIC, msg)

}

//
// read csv file as stream and post records onto processing queue
//
func enqueueCSVforNAPLANValidation(file multipart.File) (lib.IngestResponse, error) {

	ir := lib.IngestResponse{}

	reader := csv.WithIoReader(file)
	defer reader.Close()

	i := 0
	txid := nuid.Next()
	for record := range reader.C() {

		i = i + 1

		regr := nxml.RegistrationRecord{}
		r := lib.RemoveBlanks(record.AsMap())
		decode_err := ms.Decode(r, &regr)
		if decode_err != nil {
			return ir, decode_err
		}

		msg := &lib.NiasMessage{}
		msg.Body = regr
		msg.SeqNo = strconv.Itoa(i)
		msg.TxID = txid
		msg.MsgID = nuid.Next()
		// msg.Target = VALIDATION_PREFIX
		msg.Route = VALIDATION_ROUTE

		publish(msg)

	}

	// create tx record to return to client
	ir.Records = i
	ir.TxID = txid

	// update the tx tracker
	tt.SetTxSize(txid, i)

	return ir, nil

}

//
// read xml file as stream and post records onto processing queue
//
func enqueueXMLforNAPLANValidation(file multipart.File) (lib.IngestResponse, error) {

	ir := lib.IngestResponse{}

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

				var rr nxml.RegistrationRecord
				decode_err := decoder.DecodeElement(&rr, &se)
				if decode_err != nil {
					return ir, decode_err
				}

				msg := &lib.NiasMessage{}
				msg.Body = rr
				msg.SeqNo = strconv.Itoa(total)
				msg.TxID = txid
				msg.MsgID = nuid.Next()
				// msg.Target = VALIDATION_PREFIX
				msg.Route = VALIDATION_ROUTE

				publish(msg)

			}
		default:
		}

	}

	// create tx record to return to client
	ir.Records = total
	ir.TxID = txid

	// update the tx tracker
	tt.SetTxSize(txid, total)

	return ir, nil

}

//
// start the server
//
func (vws *ValidationWebServer) Run() {

	log.Println("Connecting to message bus")
	req_ec = lib.CreateNATSConnection()

	log.Println("Initialising uuid generator")
	config := uuid.StateSaverConfig{SaveReport: true, SaveSchedule: 30 * time.Minute}
	uuid.SetupFileSystemStateSaver(config)
	log.Println("UUID generator initialised.")

	log.Println("Loading xml conversion templates")
	fp := path.Join("templates", "studentpersonals.tmpl")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		log.Fatalf("Unable to parse xml conversion template, service aborting...")
	}
	sptmpl = tmpl
	log.Println("XML conversion template loaded ok.")

	//setup stan connection
	stan_conn, _ = stan.Connect(lib.NAP_VAL_CID, nuid.Next())

	// create the web service framework
	e := echo.New()

	//
	// main handler for validation of NAPLAN registration data
	//
	e.POST("/naplan/reg/validate", func(c echo.Context) error {

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
		var ir lib.IngestResponse
		if strings.Contains(file.Filename, ".csv") {
			if ir, err = enqueueCSVforNAPLANValidation(src); err != nil {
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

	//
	// Once the data file is accepted, validation results are posted
	// to a stream. The stream id is returned to the client in the
	// IngestResponse when the file is posted.
	//
	// Attaching to the stream upgrades to a websocket,
	// two types of messages are pushed to the listening client:
	// status update messages about the progress of the transaction
	// and validation error/analysis messages
	//
	//
	e.GET("/naplan/reg/stream/:txid", func(c echo.Context) error {

		txid := c.Param("txid")

		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			var msgs_sent = 0
			var ui_limit_reached bool

			// this is the main callback that receives
			// messages from the datastore
			mcb := func(m *stan.Msg) {

				msg := lib.DecodeNiasMessage(m.Data)

				// convenience type to send results/progress data only
				// to web clients
				type VMessage struct {
					Type    string
					Payload interface{}
				}

				var vmsg VMessage
				var txsu lib.TxStatusUpdate

				switch t := msg.Body.(type) {
				case ValidationError:
					vmsg = VMessage{Type: "result", Payload: msg.Body.(ValidationError)}
					if !ui_limit_reached {
						err = websocket.JSON.Send(ws, vmsg)
						if err != nil {
							log.Fatal(err)
						}
					}
					msgs_sent++
					if msgs_sent >= UI_LIMIT {
						ui_limit_reached = true
					}

				case lib.TxStatusUpdate:
					txsu = msg.Body.(lib.TxStatusUpdate)
					if msgs_sent >= UI_LIMIT {
						// send ui limit flag as part of progress update
						txsu.UIComplete = true
					}
					vmsg = VMessage{Type: "progress", Payload: txsu}
					err = websocket.JSON.Send(ws, vmsg)
					if err != nil {
						log.Fatal(err)
					}
				default:
					_ = t
					vmsg = VMessage{Type: "unknown", Payload: ""}
					log.Printf("unknown message type in handler: %v", vmsg)
				}

				// log.Printf("message decoded from stan is:\n\n %+v\n\n", msg)

			}

			sub, err := stan_conn.Subscribe(txid, mcb, stan.DeliverAllAvailable())
			defer sub.Unsubscribe()
			if err != nil {
				log.Println("Error stan subscription stream-read: ", err)
			}

			for {
				// read loop
				cmsg := ""
				err = websocket.Message.Receive(ws, &cmsg)
				if err != nil {
					// eof means socket closed by client
					if err == io.EOF {
						break
					}
				} else {
					log.Println("Message sent to websocket handler caused error ", err)
					log.Printf("message was:\n%s\n", cmsg)
				}

			}

		}).ServeHTTP(c.Response(), c.Request())

		return nil

	})

	//
	// handler for csv-xml conversion
	//
	e.POST("/naplan/reg/convert", func(c echo.Context) error {

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
			r1 := lib.RemoveBlanks(r)
			r1["SIFuuid"] = uuid.NewV4().String()
			sprsnls = append(sprsnls, r1)
		}

		// set headers to 'force' file download where appropriate
		c.Response().Header().Set("Content-Disposition", "attachment; filename="+xml_fname)
		c.Response().Header().Set("Content-Type", "application/xml")

		// apply the template & write results to the client
		if err := sptmpl.Execute(c.Response().Writer, sprsnls); err != nil {
			return err
		}

		return nil

	})

	//
	// get validation analysis results - non-websocket, just json stream
	//
	e.GET("/naplan/reg/results/:txid", func(c echo.Context) error {

		txID := c.Param("txid")

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)

		// signal channel to notify asynch stan stream read is complete
		txComplete := make(chan bool)

		// main message handling callback for the stan stream
		mcb := func(m *stan.Msg) {

			msg := lib.DecodeNiasMessage(m.Data)

			// convenience type to send results/progress data only
			// to web clients
			type VMessage struct {
				Type    string
				Payload interface{}
			}

			switch t := msg.Body.(type) {
			case ValidationError:
				ve := msg.Body.(ValidationError)
				if err := json.NewEncoder(c.Response()).Encode(ve); err != nil {
					log.Println("error encoding json validationerror: ", err)
				}
				c.Response().Flush()
			case lib.TxStatusUpdate:
				txsu := msg.Body.(lib.TxStatusUpdate)
				if txsu.TxComplete {
					log.Println("Finished...")
					txComplete <- true
				}
			default:
				_ = t
				vmsg := VMessage{Type: "unknown", Payload: ""}
				log.Printf("unknown message type in handler: %v", vmsg)
			}

			// log.Printf("message decoded from stan is:\n\n %+v\n\n", msg)

		}

		sub, err := stan_conn.Subscribe(txID, mcb, stan.DeliverAllAvailable())
		defer sub.Unsubscribe()
		if err != nil {
			log.Println("stan subsciption error results-download: ", err)
			return err
		}

		<-txComplete

		return nil

	})

	//
	// get the validation errors data for a given transaction as a downloadable csv file
	//
	e.GET("/naplan/reg/report/:txid/:fname", func(c echo.Context) error {

		txID := c.Param("txid")

		// get filename from params
		fname := c.Param("fname")
		rplcr := strings.NewReplacer(".csv", "_error_report.csv", ".xml", "_error_report.csv")
		rfname := rplcr.Replace(fname)

		c.Response().Header().Set("Content-Disposition", "attachment; filename="+rfname)
		c.Response().Header().Set("Content-Type", "text/csv")

		w := gcsv.NewWriter(c.Response().Writer)

		// write the headers
		hdr := []string{"Original File Line No. where error occurred",
			"Validation Type",
			"Field that failed validation",
			"Error Description"}

		if err := w.Write(hdr); err != nil {
			log.Println("error writing headers to csv:", err)
		}

		// signal channel to notify asynch stan stream read is complete
		txComplete := make(chan bool)

		// main message handling callback for the stan stream
		mcb := func(m *stan.Msg) {

			msg := lib.DecodeNiasMessage(m.Data)

			// convenience type to send results/progress data only
			// to web clients
			type VMessage struct {
				Type    string
				Payload interface{}
			}

			switch t := msg.Body.(type) {
			case ValidationError:
				ve := msg.Body.(ValidationError)
				if err := w.Write(ve.ToSlice()); err != nil {
					log.Println("error writing record to csv:", err)
				}
			case lib.TxStatusUpdate:
				txsu := msg.Body.(lib.TxStatusUpdate)
				if txsu.TxComplete {
					log.Println("Finished...")
					txComplete <- true
				}
			default:
				_ = t
				vmsg := VMessage{Type: "unknown", Payload: ""}
				log.Printf("unknown message type in handler: %v", vmsg)
			}

			// log.Printf("message decoded from stan is:\n\n %+v\n\n", msg)

		}

		sub, err := stan_conn.Subscribe(txID, mcb, stan.DeliverAllAvailable())
		defer sub.Unsubscribe()
		if err != nil {
			log.Println("stan subsciption error csv-download: ", err)
			return err
		}

		<-txComplete

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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Println("Starting web-ui services...")
	port := DefaultValidationConfig.WebServerPort
	log.Println("Service is listening on localhost:" + port)

	// set upper bound for no. messages sent to web clients
	UI_LIMIT = DefaultValidationConfig.UIMessageLimit

	//e.Run(fasthttp.New(":" + port))
	// e.Run(standard.New(":" + port))
	e.Logger.Fatal(e.Start(":" + port))

}
