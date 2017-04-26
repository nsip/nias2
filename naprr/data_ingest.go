package naprr

import (
	goxml "encoding/xml"
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
	"sync"

	"github.com/nats-io/go-nats-streaming"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
)

type naprr_config struct {
	// fields on which to match student in Yr3W ingest with student in XML ingest
	Yr3WStudentMatch []string
	Loaded           bool
}

func LoadConfig() naprr_config {
	ncfg := naprr_config{Loaded: false}
	if _, err := toml.DecodeFile("naprr.toml", &ncfg); err != nil {
		log.Fatalln("Unable to read naprr config, aborting.", err)
	}
	ncfg.Loaded = true
	return ncfg
}

type DataIngest struct {
	sc stan.Conn
	ge GobEncoder
	sr *StreamReader
}

func NewDataIngest() *DataIngest {
	di := DataIngest{sc: CreateSTANConnection(), ge: GobEncoder{}, sr: NewStreamReader()}
	return &di
}

func (di *DataIngest) Close() {
	di.sc.Close()
}

func (di *DataIngest) Run() {

	xmlFiles := parseXMLFileDirectory()

	var wg sync.WaitGroup

	for _, xmlFile := range xmlFiles {
		wg.Add(1)
		go di.ingestResultsFile(xmlFile, &wg)
	}

	wg.Wait()

	di.finaliseTransactions()

	log.Println("All data files read, ingest complete.")
}

func (di *DataIngest) RunSynchronous(FilePath string) {
	di.ingestResultsFile(FilePath, nil)
	di.finaliseTransactions()
}

// Given a student record and a list of fields, return a colon-delimited key giving those field values
func studentKeyLookup(r xml.RegistrationRecord, fields []string) string {
	// we could use Go Reflect, but for performance, let's not
	key := ""
	for _, field := range fields {
		switch field {
		case "GivenName":
			key = key + "::" + r.GivenName
		case "FamilyName":
			key = key + "::" + r.FamilyName
		case "MiddleName":
			key = key + "::" + r.MiddleName
		case "PreferredName":
			key = key + "::" + r.PreferredName
		case "LocalId":
			key = key + "::" + r.LocalId
		case "StateProvinceId":
			key = key + "::" + r.StateProvinceId
		case "DiocesanId":
			key = key + "::" + r.DiocesanId
		case "NationalId":
			key = key + "::" + r.NationalId
		case "PlatformId":
			key = key + "::" + r.PlatformId
		}
	}
	return key
}

func parseXMLFileDirectory() []string {

	files, _ := filepath.Glob("./in/*.zip")

	if len(files) == 0 {
		log.Fatalln("No results data zip files found in input folder.")
	}

	return files

}

func (di *DataIngest) ingestResultsFile(resultsFilePath string, wg *sync.WaitGroup) {

	// create a connection to the streaming server
	log.Println("Connecting to STAN server...")

	config := LoadConfig()
	// map to hold student-school links temporarily
	// so student responses can be assigned to correct schools
	ss_link := make(map[string]string)
	// map to hold student identities temporarily, so that Yr3 Writing links to students can be made later
	student_ids := make(map[string]string)

	// simple list of schools
	// schools := make([]SchoolDetails, 0)

	// open the data file for streaming read
	log.Printf("Opening results data file %s...", resultsFilePath)
	xmlFile, err := OpenResultsFile(resultsFilePath)
	if err != nil {
		log.Fatalln("unable to open results data file: ", err)
	}
	xmlFile, err = OpenResultsFile(resultsFilePath)

	log.Println("Reading data file...")

	decoder := goxml.NewDecoder(xmlFile)
	totalTests := 0
	totalTestlets := 0
	totalTestItems := 0
	totalTestScoreSummarys := 0
	totalEvents := 0
	totalResponses := 0
	totalCodeFrames := 0
	totalSchools := 0
	totalStudents := 0
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case goxml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...handle by type
			switch inElement {
			case "NAPTest":
				var t xml.NAPTest
				decoder.DecodeElement(&t, &se)
				gt, err := di.ge.Encode(t)
				if err != nil {
					log.Println("Unable to gob-encode nap test: ", err)
				}
				di.sc.Publish(META_STREAM, gt)
				totalTests++

			case "NAPTestlet":
				var tl xml.NAPTestlet
				decoder.DecodeElement(&tl, &se)
				gtl, err := di.ge.Encode(tl)
				if err != nil {
					log.Println("Unable to gob-encode nap testlet: ", err)
				}
				di.sc.Publish(META_STREAM, gtl)
				totalTestlets++

			case "NAPTestItem":
				var ti xml.NAPTestItem
				decoder.DecodeElement(&ti, &se)
				gti, err := di.ge.Encode(ti)
				if err != nil {
					log.Println("Unable to gob-encode nap test item: ", err)
				}
				di.sc.Publish(META_STREAM, gti)
				totalTestItems++

			case "NAPTestScoreSummary":
				var tss xml.NAPTestScoreSummary
				decoder.DecodeElement(&tss, &se)
				gtss, err := di.ge.Encode(tss)
				if err != nil {
					log.Println("Unable to gob-encode nap test-score-summary: ", err)
				}
				di.sc.Publish(tss.SchoolACARAId, gtss)
				totalTestScoreSummarys++

			case "NAPEventStudentLink":
				var e xml.NAPEvent
				decoder.DecodeElement(&e, &se)
				ge, err := di.ge.Encode(e)
				if err != nil {
					log.Println("Unable to gob-encode nap event link: ", err)
				}
				di.sc.Publish(e.SchoolID, ge)
				totalEvents++

			case "NAPStudentResponseSet":
				var r xml.NAPResponseSet
				decoder.DecodeElement(&r, &se)
				gr, err := di.ge.Encode(r)
				if err != nil {
					log.Println("Unable to gob-encode student response set: ", err)
				}
				di.sc.Publish("responses", gr)
				totalResponses++

			case "NAPCodeFrame":
				var cf xml.NAPCodeFrame
				decoder.DecodeElement(&cf, &se)
				gcf, err := di.ge.Encode(cf)
				if err != nil {
					log.Println("Unable to gob-encode nap codeframe: ", err)
				}
				di.sc.Publish(META_STREAM, gcf)
				totalCodeFrames++

			case "SchoolInfo":
				var si xml.SchoolInfo
				decoder.DecodeElement(&si, &se)
				gsi, err := di.ge.Encode(si)
				if err != nil {
					log.Println("Unable to gob-encode schoolinfo: ", err)
				}
				di.sc.Publish(si.ACARAId, gsi)
				// store school in local list
				sd := SchoolDetails{ACARAId: si.ACARAId, SchoolName: si.SchoolName}
				gsd, err := di.ge.Encode(sd)
				if err != nil {
					log.Println("Unable to gob-encode school-details: ", err)
				}
				di.sc.Publish("schools", gsd)

				totalSchools++

			case "StudentPersonal":
				var sp xml.RegistrationRecord
				decoder.DecodeElement(&sp, &se)
				gsp, err := di.ge.Encode(sp)
				if err != nil {
					log.Println("Unable to gob-encode studentpersonal: ", err)
				}
				// store linkage locally
				ss_link[sp.RefId] = sp.ASLSchoolId
				student_key := studentKeyLookup(sp, config.Yr3WStudentMatch)
				student_ids[student_key] = sp.RefId
				di.sc.Publish(sp.ASLSchoolId, gsp)
				totalStudents++

			}
		default:
		}

	}

	log.Println("Data file read complete...")
	log.Printf("Total tests: %d \n", totalTests)
	log.Printf("Total codeframes: %d \n", totalCodeFrames)
	log.Printf("Total testlets: %d \n", totalTestlets)
	log.Printf("Total test items: %d \n", totalTestItems)
	log.Printf("Total test score summaries: %d \n", totalTestScoreSummarys)
	log.Printf("Total events: %d \n", totalEvents)
	log.Printf("Total responses: %d \n", totalResponses)
	log.Printf("Total schools: %d \n", totalSchools)
	log.Printf("Total students: %d \n", totalStudents)

	log.Println("Finalising test metadata & responses...")

	// post end of stream message to responses queue
	eot := lib.TxStatusUpdate{TxComplete: true}
	geot, err := di.ge.Encode(eot)
	if err != nil {
		log.Println("Unable to gob-encode tx complete message: ", err)
	}
	di.sc.Publish("responses", geot)
	di.sc.Publish(META_STREAM, geot)

	di.assignResponsesToSchools(ss_link)

	log.Println("response assignment complete")

	log.Printf("ingestion complete for %s", resultsFilePath)

	if wg != nil {
		wg.Done()
	}

}

func (di *DataIngest) assignResponsesToSchools(ss_link map[string]string) {

	log.Println("Assigning responses to schools...")

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := di.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case xml.NAPResponseSet:
			rs := m_if.(xml.NAPResponseSet)
			studentId := rs.StudentID
			schoolId := ss_link[studentId]
			di.sc.Publish(schoolId, m.Data)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			log.Printf("unknown message type in response assign handler: %v", m_if)
		}
		// log.Printf("message decoded from stan is:\n\n %+v\n\n", msg)

	}

	sub, err := di.sc.Subscribe("responses", mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("stan subsciption error - response assignment: ", err)
	}

	<-txComplete

	log.Println("All reponses assigned to schools.")

}

func (di *DataIngest) finaliseTransactions() {

	log.Println("Finalising data read transactions...")

	// end of tx marker message
	eot := lib.TxStatusUpdate{TxComplete: true}
	geot, err := di.ge.Encode(eot)
	if err != nil {
		log.Println("Unable to gob-encode tx complete message: ", err)
	}

	// finalise the known list of schools
	di.sc.Publish("schools", geot)

	// then use the list to finalise each school data stream
	schools := di.sr.GetSchoolDetails()

	for _, subslice := range schools {
		for _, school := range subslice {
			log.Println("finalising ", school)
			di.sc.Publish(school.ACARAId, geot)

		}
	}

	log.Println("All transactions finalised.")

}
