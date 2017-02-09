// streamreader.go
//
//
// utility routines to read naplan and school data from a
// stan stream

package naprr

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
	"log"
)

type StreamReader struct {
	sc stan.Conn
	ge GobEncoder
}

func NewStreamReader() *StreamReader {
	sr := StreamReader{
		sc: CreateSTANConnection(),
		ge: GobEncoder{},
	}
	return &sr
}

func (sr *StreamReader) GetCodeFrameData() []CodeframeDataSet {

	cfds := make([]CodeframeDataSet, 0)

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	// get names of schools that have been processed by ingest
	// and create reports
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := gobenc.Decode(m.Data, &m_if)
		if err != nil {
			log.Println("streamreader codeframe message decoding error: ", err)
			txComplete <- true
		}

		switch t := m_if.(type) {
		case CodeframeDataSet:
			cfd := m_if.(CodeframeDataSet)
			cfds = append(cfds, cfd)
		case lib.TxStatusUpdate:
			txComplete <- true
		default:
			_ = t
			// log.Printf("unknown message type in participation data handler: %v", m_if)
		}
	}

	sub, err := sc.Subscribe("reports.cframe", mcb, stan.DeliverAllAvailable())
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
		err := gobenc.Decode(m.Data, &m_if)
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

	sub, err := sc.Subscribe("reports."+acaraid+".dscores", mcb, stan.DeliverAllAvailable())
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
		err := gobenc.Decode(m.Data, &m_if)
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

	sub, err := sc.Subscribe("reports."+acaraid+".scsumm", mcb, stan.DeliverAllAvailable())
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
		err := gobenc.Decode(m.Data, &m_if)
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

	sub, err := sc.Subscribe("reports."+acaraid+".particip", mcb, stan.DeliverAllAvailable())
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
		err := gobenc.Decode(m.Data, &m_if)
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

	sub, err := sc.Subscribe("schools", mcb, stan.DeliverAllAvailable())
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
func (sr *StreamReader) GetNAPLANData() *NAPLANData {

	nd := NewNAPLANData()

	// signal channel to notify asynch stan stream read is complete
	txComplete := make(chan bool)

	// main message handling callback for the stan stream
	mcb := func(m *stan.Msg) {

		// as we don't know message type ([]byte slice on wire) decode as interface
		// then assert type dynamically
		var m_if interface{}
		err := gobenc.Decode(m.Data, &m_if)
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

	sub, err := sc.Subscribe("meta", mcb, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Println("streamreader: stan subsciption error meta channel: ", err)
	}

	<-txComplete

	return nd

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
		err := gobenc.Decode(m.Data, &m_if)
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

	sub, err := sc.Subscribe(sd.ACARAId, mcb, stan.DeliverAllAvailable())
	if err != nil {
		log.Println("streamreader: stan subsciption error school data channel: ", err)
	}
	defer sub.Unsubscribe()

	<-txComplete

	return sd

}
