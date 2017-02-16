// reportgenerator.go
package naprr

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/nsip/nias2/lib"
	"log"
)

type ReportGenerator struct {
	sc stan.Conn
	ge GobEncoder
}

func NewReportGenerator() *ReportGenerator {
	rg := ReportGenerator{
		sc: CreateSTANConnection(),
		ge: GobEncoder{},
	}
	return &rg
}

//
// routines to build the required reports
//

// generate codeframe objects (currently as per VCAA requirements)
// generated only once as represents strucure of test not school-level data
func (rg *ReportGenerator) GenerateCodeFrameData(nd *NAPLANData) {

	count := 0
	cfds := make([]CodeFrameDataSet, 0)

	for _, codeframe := range nd.Codeframes {
		for _, cf_testlet := range codeframe.TestletList.Testlet {
			tl := nd.Testlets[cf_testlet.NAPTestletRefId]
			// log.Printf("\t%s", tl.TestletContent.TestletName)
			for _, cf_item := range cf_testlet.TestItemList.TestItem {
				ti := nd.Items[cf_item.TestItemRefId]
				// log.Printf("\t\t%s", ti.TestItemContent.ItemName)
				cfd := CodeFrameDataSet{
					Test:    nd.Tests[codeframe.NAPTestRefId],
					Testlet: tl,
					Item:    ti,
				}
				cfds = append(cfds, cfd)
			}
		}
	}

	count = len(cfds)

	// publish the records
	for _, cfd := range cfds {
		payload, err := rg.ge.Encode(cfd)
		if err != nil {
			log.Println("unable to encode codeframe: ", err)
		}
		// log.Printf("\t%s - %s - %s", cfd.Test.TestContent.TestDomain,
		// 	cfd.Testlet.TestletContent.TestletName, cfd.Item.TestItemContent.ItemName)
		rg.sc.Publish("reports.cframe", payload)
	}

	// finish the transaction - completion msg
	txu := lib.TxStatusUpdate{TxComplete: true}
	gtxu, err := rg.ge.Encode(txu)
	if err != nil {
		log.Println("unable to encode txu codeframe report: ", err)
	}
	rg.sc.Publish("reports.cframe", gtxu)

	log.Printf("codeframe records %d: ", count)

}

// generate domain score objects
func (rg *ReportGenerator) GenerateDomainScoreData(nd *NAPLANData, sd *SchoolData) {

	count := 0

	for _, response := range sd.Responses {
		rds := ResponseDataSet{
			Test:     nd.Tests[response.TestID],
			Response: response,
		}
		// log.Printf("sc_score_summ:\n\n%v\n\n%v\n\n", rds.Test, rds.Response)

		payload, err := rg.ge.Encode(rds)
		if err != nil {
			log.Println("unable to encode domain scores: ", err)
		}
		rg.sc.Publish("reports."+sd.ACARAId+".dscores", payload)

		count++
	}

	// finish the transaction - completion msg
	txu := lib.TxStatusUpdate{TxComplete: true}
	gtxu, err := rg.ge.Encode(txu)
	if err != nil {
		log.Println("unable to encode txu domain scores report: ", err)
	}
	rg.sc.Publish("reports."+sd.ACARAId+".dscores", gtxu)

	log.Printf("domain score records %d: ", count)

}

// generate school summary objects
func (rg *ReportGenerator) GenerateSchoolScoreSummaryData(nd *NAPLANData, sd *SchoolData) {

	count := 0

	for _, scoresummary := range sd.ScoreSummaries {
		scsumm := ScoreSummaryDataSet{
			Test: nd.Tests[scoresummary.NAPTestRefId],
			Summ: scoresummary,
		}
		// log.Printf("sc_score_summ:\n\n%v\n\n%v\n\n", scsumm.Test, scsumm.Summ)

		payload, err := rg.ge.Encode(scsumm)
		if err != nil {
			log.Println("unable to encode sch. summ: ", err)
		}
		rg.sc.Publish("reports."+sd.ACARAId+".scsumm", payload)

		count++
	}

	// finish the transaction - completion msg
	txu := lib.TxStatusUpdate{TxComplete: true}
	gtxu, err := rg.ge.Encode(txu)
	if err != nil {
		log.Println("unable to encode txu score summary report: ", err)
	}
	rg.sc.Publish("reports."+sd.ACARAId+".scsumm", gtxu)

	log.Printf("score summary records %d: ", count)

}

// generate the participation summary objects
func (rg *ReportGenerator) GenerateParticipationData(nd *NAPLANData, sd *SchoolData) {

	count := 0
	studentEvents := make(map[string][]EventInfo)

	for _, event := range sd.Events {
		ei := EventInfo{Event: event, Test: nd.Tests[event.TestID]}
		infos, ok := studentEvents[event.SPRefID]
		if !ok {
			infos = make([]EventInfo, 0)
		}
		infos = append(infos, ei)
		studentEvents[event.SPRefID] = infos
	}

	for _, student := range sd.Students {
		pds := ParticipationDataSet{
			Student:    sd.Students[student.RefId],
			School:     sd.SchoolInfos[student.ASLSchoolId],
			EventInfos: studentEvents[student.RefId],
			Summary:    make(map[string]string),
		}
		for _, ei := range pds.EventInfos {
			pds.Summary[ei.Test.TestContent.TestDomain] = ei.Event.ParticipationCode
		}
		payload, err := rg.ge.Encode(pds)
		if err != nil {
			log.Println("unable to encode pds: ", err)
		}
		rg.sc.Publish("reports."+sd.ACARAId+".particip", payload)

		count++
	}

	// finish the transaction - completion msg
	txu := lib.TxStatusUpdate{TxComplete: true}
	gtxu, err := rg.ge.Encode(txu)
	if err != nil {
		log.Println("unable to encode txu particip. report: ", err)
	}
	rg.sc.Publish("reports."+sd.ACARAId+".particip", gtxu)

	log.Printf("particpation records %d: ", count)
}

//
//
