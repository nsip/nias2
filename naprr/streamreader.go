// streamreader.go
//
//
// utility routines to read naplan and school data from a
// stan stream

package naprr

import (
	"encoding/gob"
	"fmt"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
	"log"
	"os"
)

type StreamReader struct {
	sc stan.Conn
	ge GobEncoder
}

// track how many students matched between Yr 3 writing and core XML ingest
type Yr3WritingMatchReport struct {
	Matches         []string
	Yr3w_mismatches []string // students in yr3w ingest not in xml
	Xml_mismatches  []string // students in xml not in yr3 ingest
}

const META_STREAM = "meta"
const META_YR3W_STREAM = "meta_yr3w"
const RESULTS_YR3W_STREAM = "studentAndResults"
const REPORTS_CODEFRAME = "reports.cframe"
const REPORTS_YR3W = "reports.yr3w"
const REPORTS_YR3W_STATUS = "reports.yr3w.status"

func NewStreamReader() *StreamReader {
	sr := StreamReader{
		sc: CreateSTANConnection(),
		ge: GobEncoder{},
	}
	gob.Register(Yr3WritingMatchReport{})
	return &sr
}

func (sr *StreamReader) GetResultsByStudent() []ResultsByStudent {

	cfds := make([]ResultsByStudent, 0)

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	// get names of schools that have been processed by ingest
	// and create reports
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader codeframe message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case ResultsByStudent:
			cfd := m_if.(ResultsByStudent)
			cfds = append(cfds, cfd)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			// log.Printf("unknown message type in participation data handler: %v", m_if)
		}
	}

	sub, err := sr.sc.Subscribe(REPORTS_YR3W, mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error get students & results data: ", err)
	}

	<-txComplete

	log.Printf("Retrieved %d results by students records\n", len(cfds))
	return cfds

}

func (sr *StreamReader) GetCodeFrameData(stream_name string) []CodeFrameDataSet {

	cfds := make([]CodeFrameDataSet, 0)

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	// get names of schools that have been processed by ingest
	// and create reports
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader codeframe message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case CodeFrameDataSet:
			cfd := m_if.(CodeFrameDataSet)
			cfds = append(cfds, cfd)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			// log.Printf("unknown message type in participation data handler: %v", m_if)
		}
	}

	sub, err := sr.sc.Subscribe(stream_name, mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error get codeframe data: ", err)
	}

	<-txComplete

	return cfds

}

//
// get domain scores
//
func (sr *StreamReader) GetDomainScoreData(acaraid string) []ResponseDataSet {

	rds := make([]ResponseDataSet, 0)
	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	// get names of schools that have been processed by ingest
	// and create reports
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader domain score message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case ResponseDataSet:
			rd := m_if.(ResponseDataSet)
			rds = append(rds, rd)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			// log.Printf("unknown message type in participation data handler: %v", m_if)
		}
	}

	sub, err := sr.sc.Subscribe("reports."+acaraid+".dscores", mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error get domain scores data: ", err)
	}

	<-txComplete

	return rds

}

//
// get school score summaries
//
func (sr *StreamReader) GetScoreSummaryData(acaraid string) []ScoreSummaryDataSet {

	ssds := make([]ScoreSummaryDataSet, 0)

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	// get names of schools that have been processed by ingest
	// and create reports
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader score summ message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case ScoreSummaryDataSet:
			ssd := m_if.(ScoreSummaryDataSet)
			ssds = append(ssds, ssd)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			// log.Printf("unknown message type in participation data handler: %v", m_if)
		}
	}

	sub, err := sr.sc.Subscribe("reports."+acaraid+".scsumm", mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error get score summary data: ", err)
	}

	<-txComplete

	return ssds

}

//
// retrieve participation data for the given school
//
func (sr *StreamReader) GetParticipationData(acaraid string) []ParticipationDataSet {

	pds := make([]ParticipationDataSet, 0)

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	// get names of schools that have been processed by ingest
	// and create reports
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader participation message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case ParticipationDataSet:
			pd := m_if.(ParticipationDataSet)
			pds = append(pds, pd)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			// log.Printf("unknown message type in participation data handler: %v", m_if)
		}
	}

	sub, err := sr.sc.Subscribe("reports."+acaraid+".particip", mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error get participation data: ", err)
	}

	<-txComplete

	return pds

}

// returns simple list of all schools that have been processed
// from the results data set
//
// details are a slice of slices so that downstream processes
// such as concurrent file writes can work with a moderate batch size
// e.g. 100 at a time to prevent issues with too many open filehandles
// stan connections etc.
func (sr *StreamReader) GetSchoolDetails() [][]SchoolDetails {

	sds := make([]SchoolDetails, 0)

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	// get names of schools that have been processed by ingest
	// and create reports
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case SchoolDetails:
			sd := m_if.(SchoolDetails)
			sds = append(sds, sd)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			// log.Printf("unknown message type in stream reader school details handler: %v", m_if)
		}
	}

	sub, err := sr.sc.Subscribe("schools", mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error get school details: ", err)
	}

	<-txComplete

	// now chunk up the list of schools into sub-slices
	var sds_chunks [][]SchoolDetails
	chunkSize := 99

	for i := 0; i < len(sds); i += chunkSize {
		end := i + chunkSize

		if end > len(sds) {
			end = len(sds)
		}

		sds_chunks = append(sds_chunks, sds[i:end])
	}

	return sds_chunks

}

// NAPLAN data is the same for all schools, so can be retrieved once
// codeframe_stream: "meta" is ingested from xml, "meta_yr3w" from Year 3 writing files
func (sr *StreamReader) GetNAPLANData(codeframe_stream string) *NAPLANData {

	nd := NewNAPLANData()

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	mcb := func(m *stan.Msg) {
		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader: schooldata message decoding error: ", err)
			txComplete <- true
		}

		switch mtype := m_if.(type) {
		case xml.NAPTest:
			t := m_if.(xml.NAPTest)
			nd.Tests[t.TestID] = t
		case xml.NAPTestlet:
			tl := m_if.(xml.NAPTestlet)
			nd.Testlets[tl.TestletID] = tl
		case xml.NAPTestItem:
			ti := m_if.(xml.NAPTestItem)
			nd.Items[ti.ItemID] = ti
		case xml.NAPCodeFrame:
			cf := m_if.(xml.NAPCodeFrame)
			nd.Codeframes[cf.NAPTestRefId] = cf
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = mtype
			// log.Printf("unknown message type in stream reader meta handler: %v", m_if)
		}

	}

	log.Println("Ingesting for " + codeframe_stream)
	sub, err := sr.sc.Subscribe(codeframe_stream, mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error meta channel: ", err)
	}

	<-txComplete

	return nd

}

// Get all the students and results in the stream
func (sr *StreamReader) GetStudentAndResultsData(student_ids map[string]string, NaprrConfig naprr_config) *StudentAndResultsData {

	srd := NewStudentAndResultsData()

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader: schooldata message decoding error: ", err)
			txComplete <- true
		}

		switch mtype := m_if.(type) {
		case xml.RegistrationRecord:
			t := m_if.(xml.RegistrationRecord)
			srd.Students[t.RefId] = &t
		case xml.NAPEvent:
			e := m_if.(xml.NAPEvent)
			srd.Events[e.SPRefID] = &e
		case xml.NAPResponseSet:
			rs := m_if.(xml.NAPResponseSet)
			srd.ResponseSets[rs.StudentID] = &rs
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = mtype
			// log.Printf("unknown message type in stream reader meta handler: %v", m_if)
		}

	}

	sub, err := sr.sc.Subscribe(RESULTS_YR3W_STREAM, mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error meta channel: ", err)
	}

	<-txComplete

	log.Printf("Retrieved %d students and results records\n", len(srd.Students))
	srd = sr.remapStudents(srd, student_ids, NaprrConfig)
	return srd

}

// given the mapping of keys to XML RefIDs, remap all references to made up RegistrationRecords from fixed format file, to XML records already ingested
func (sr *StreamReader) remapStudents(srd *StudentAndResultsData, student_ids map[string]string, NaprrConfig naprr_config) *StudentAndResultsData {
	yr3wmatches := Yr3WritingMatchReport{}
	matches := make(map[string]int)
	newStudents := make(map[string]*xml.RegistrationRecord)
	log.Printf("%v\n", NaprrConfig)
	log.Println(NaprrConfig.Yr3WStudentMatch)
	for _, sp := range srd.Students {
		student_key := StudentKeyLookup(*sp, NaprrConfig.Yr3WStudentMatch)
		if newRefId, ok := student_ids[student_key]; ok {
			refid := sp.RefId
			log.Printf("Mapped student %s from Yr 3 Writing from %s to matching student %s in XML ingest", student_key, refid, newRefId)
			//append(yr3wmatches.matches, student_key)
			matches[student_key] = 1
			// eliminating this record from the list: not passing on to newStudents
			if _, ok := srd.Events[refid]; ok {
				log.Printf("%v\n", srd.Events[refid])
				(srd.Events[refid]).SPRefID = newRefId
				log.Printf("%v\n", srd.Events[refid])
			}
			if _, ok := srd.ResponseSets[refid]; ok {
				(srd.ResponseSets[refid]).StudentID = newRefId
			}
		} else {
			log.Printf("No match in mapping student %s from Yr 3 Writing to matching student in XML ingest", student_key)
			yr3wmatches.Yr3w_mismatches = append(yr3wmatches.Yr3w_mismatches, student_key)
			newStudents[sp.RefId] = sp
		}
	}
	for k, _ := range student_ids {
		if _, ok := matches[k]; !ok {
			yr3wmatches.Xml_mismatches = append(yr3wmatches.Xml_mismatches, k)
		}
	}
	yr3wmatches.Matches = make([]string, 0, len(matches))
	for k := range matches {
		yr3wmatches.Matches = append(yr3wmatches.Matches, k)
	}
	log.Printf("%v\n", yr3wmatches)
	/*
		payload, err := sr.ge.Encode(yr3wmatches)
		if err != nil {
			log.Println("unable to encode yr 3 writing status report: ", err)
		}
	*/
	//sr.sc.Publish(REPORTS_YR3W_STATUS, payload)
	// create directory for the school
	fpath := "yr3w/"
	err := os.MkdirAll(fpath, os.ModePerm)
	check(err)

	// create the report data file in the output directory
	// delete any ecisting files and create empty new one
	fname := fpath + "codeframe_report.txt"
	err = os.RemoveAll(fname)
	f, err := os.Create(fname)
	check(err)
	defer f.Close()
	payload := fmt.Sprintf("Matches: %v\nYr3W only: %v\nXML only: %v\n",
		yr3wmatches.Matches, yr3wmatches.Yr3w_mismatches, yr3wmatches.Xml_mismatches)
	f.WriteString(payload)

	srd.Students = newStudents
	return srd
}

// for the school identified by the acaraid retrieves all of the
// raw results data objects
func (sr *StreamReader) GetSchoolData(acaraid string) *SchoolData {

	sd := NewSchoolData(acaraid)

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := sr.ge.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("message decoding error: ", err)
			txComplete <- true
		}

		switch mtype := m_if.(type) {
		case xml.SchoolInfo:
			si := m_if.(xml.SchoolInfo)
			sd.SchoolInfos[si.ACARAId] = si
		case xml.NAPEvent:
			e := m_if.(xml.NAPEvent)
			sd.Events[e.EventID] = e
		case xml.RegistrationRecord:
			sp := m_if.(xml.RegistrationRecord)
			sd.Students[sp.RefId] = sp
		case xml.NAPTestScoreSummary:
			tss := m_if.(xml.NAPTestScoreSummary)
			sd.ScoreSummaries[tss.SummaryID] = tss
		case xml.NAPResponseSet:
			rs := m_if.(xml.NAPResponseSet)
			sd.Responses[rs.ResponseID] = rs
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = mtype
			// log.Printf("unknown message type in stream reader schooldata handler: %v", m_if)
		}

	}

	sub, err := sr.sc.Subscribe(sd.ACARAId, mcb, stan.DeliverAllAvailable())
	if err != nil {
		log.Println("streamreader: stan subsciption error school data channel: ", err)
	}
	defer sub.Unsubscribe()

	<-txComplete

	return sd

}
