package main

// streaming XML parser.

import (
	goxml "encoding/xml"
	"flag"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/naprr"
	"github.com/nsip/nias2/xml"
	"log"
)

// var inputFile = flag.String("infile", "../../test_data/napcombined.xml", "Input file path")
var inputFile = flag.String("infile", "../../../test_data/naprr/master_nap.xml", "Input file path")

// var inputFile = flag.String("infile", "../../../test_data/naprr/12school.nap.xml", "Input file path")

func main() {
	flag.Parse()

	// map to hold student-school links temporarily
	// so student responses can be assigned to correct schools
	ss_link := make(map[string]string)

	// simple list of schools
	schools := make([]string, 0)

	// create a connection to the streaming server
	log.Println("Connecting to STAN server...")
	sc := naprr.CreateSTANConnection()
	defer sc.Close()

	// sc2 := naprr.CreateSTANConnection()
	// defer sc2.Close()

	gobenc := naprr.GobEncoder{}

	// open the data file for streaming read
	log.Println("Opening results data file...")
	xmlFile, err := naprr.OpenResultsFile(*inputFile)
	if err != nil {
		log.Fatalln("unable to open results data file: ", err)
	}

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
				gt, err := gobenc.Encode(t)
				if err != nil {
					log.Println("Unable to gob-encode nap test: ", err)
				}
				sc.Publish("meta", gt)
				totalTests++

			case "NAPTestlet":
				var tl xml.NAPTestlet
				decoder.DecodeElement(&tl, &se)
				gtl, err := gobenc.Encode(tl)
				if err != nil {
					log.Println("Unable to gob-encode nap testlet: ", err)
				}
				sc.Publish("meta", gtl)
				totalTestlets++

			case "NAPTestItem":
				var ti xml.NAPTestItem
				decoder.DecodeElement(&ti, &se)
				gti, err := gobenc.Encode(ti)
				if err != nil {
					log.Println("Unable to gob-encode nap test item: ", err)
				}
				sc.Publish("meta", gti)
				totalTestItems++

			case "NAPTestScoreSummary":
				var tss xml.NAPTestScoreSummary
				decoder.DecodeElement(&tss, &se)
				gtss, err := gobenc.Encode(tss)
				if err != nil {
					log.Println("Unable to gob-encode nap test-score-summary: ", err)
				}
				sc.Publish(tss.SchoolACARAId, gtss)
				totalTestScoreSummarys++

			case "NAPEventStudentLink":
				var e xml.NAPEvent
				decoder.DecodeElement(&e, &se)
				ge, err := gobenc.Encode(e)
				if err != nil {
					log.Println("Unable to gob-encode nap event link: ", err)
				}
				sc.Publish(e.SchoolID, ge)
				totalEvents++

			case "NAPStudentResponseSet":
				var r xml.NAPResponseSet
				decoder.DecodeElement(&r, &se)
				gr, err := gobenc.Encode(r)
				if err != nil {
					log.Println("Unable to gob-encode student response set: ", err)
				}
				sc.Publish("responses", gr)
				totalResponses++

			case "NAPCodeFrame":
				var cf xml.NAPCodeFrame
				decoder.DecodeElement(&cf, &se)
				gcf, err := gobenc.Encode(cf)
				if err != nil {
					log.Println("Unable to gob-encode nap codeframe: ", err)
				}
				sc.Publish("meta", gcf)
				totalCodeFrames++

			case "SchoolInfo":
				var si xml.SchoolInfo
				decoder.DecodeElement(&si, &se)
				gsi, err := gobenc.Encode(si)
				if err != nil {
					log.Println("Unable to gob-encode schoolinfo: ", err)
				}
				sc.Publish(si.ACARAId, gsi)
				// store school in local list
				schools = append(schools, si.ACARAId)
				totalSchools++

			case "StudentPersonal":
				var sp xml.RegistrationRecord
				decoder.DecodeElement(&sp, &se)
				gsp, err := gobenc.Encode(sp)
				if err != nil {
					log.Println("Unable to gob-encode studentpersonal: ", err)
				}
				// store linkage locally
				ss_link[sp.RefId] = sp.ASLSchoolId
				sc.Publish(sp.ASLSchoolId, gsp)
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

	log.Println("Finalising meta & responses...")

	// post end of stream message to responses q
	eot := lib.TxStatusUpdate{TxComplete: true}
	geot, err := gobenc.Encode(eot)
	if err != nil {
		log.Println("Unable to gob-encode tx complete message: ", err)
	}
	sc.Publish("responses", geot)
	sc.Publish("meta", geot)

	log.Println("Assigning responses to schools...")

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := gobenc.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case xml.NAPResponseSet:
			rs := m_if.(xml.NAPResponseSet)
			studentId := rs.StudentID
			schoolId := ss_link[studentId]
			sc.Publish(schoolId, m.Data)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			log.Printf("unknown message type in response assign handler: %v", m_if)
		}
		// log.Printf("message decoded from stan is:\n\n %+v\n\n", msg)

	}

	sub, err := sc.Subscribe("responses", mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("stan subsciption error response assignment: ", err)
	}

	<-txComplete

	log.Println("response assignment complete")

	log.Println("Finalising school data read transactions...")
	for _, school := range schools {
		log.Println("finalising ", school)
		sc.Publish(school, geot)
		sd := naprr.SchoolDetails{ACARAId: school}
		gsd, err := gobenc.Encode(sd)
		if err != nil {
			log.Println("Error creating school details message: ", err)
		}
		sc.Publish("schools", gsd)
	}
	sc.Publish("schools", geot)

	log.Println("ingestion complete\n\n")

}
