package naprrql

import (
	goxml "encoding/xml"
	"log"

	"github.com/nats-io/nuid"
	"github.com/nsip/nias2/xml"
	"github.com/syndtr/goleveldb/leveldb"
)

func IngestResultsFile(resultsFilePath string) {

	db := GetDB()
	ge := GobEncoder{}

	// open the data file for streaming read
	xmlFile, err := OpenResultsFile(resultsFilePath)
	if err != nil {
		log.Fatalln("unable to open results file")
	}

	log.Printf("Reading data file [%s]", resultsFilePath)

	batch := new(leveldb.Batch)

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
				gt, err := ge.Encode(t)
				if err != nil {
					log.Println("Unable to gob-encode nap test: ", err)
				}

				// {NAPTest} = object
				batch.Put([]byte(t.TestID), gt)

				// NAPTest-type:{id} = id
				key := []byte("NAPTest:" + t.TestID)
				batch.Put(key, []byte(t.TestID))

				totalTests++

			case "NAPTestlet":
				var tl xml.NAPTestlet
				decoder.DecodeElement(&tl, &se)
				gtl, err := ge.Encode(tl)
				if err != nil {
					log.Println("Unable to gob-encode nap testlet: ", err)
				}

				// {NAPTestlet} = object
				batch.Put([]byte(tl.TestletID), gtl)

				// NAPTestlet-type:{id} = {id}
				key := []byte("NAPTestlet:" + tl.TestletID)
				batch.Put(key, []byte(tl.TestletID))

				totalTestlets++

			case "NAPTestItem":
				var ti xml.NAPTestItem
				decoder.DecodeElement(&ti, &se)
				gti, err := ge.Encode(ti)
				if err != nil {
					log.Println("Unable to gob-encode nap test item: ", err)
				}

				// {NAPTestItem} = object
				batch.Put([]byte(ti.ItemID), gti)

				// NapTestItem-type:{id} = {id}
				key := []byte("NAPTestItem:" + ti.ItemID)
				batch.Put(key, []byte(ti.ItemID))

				totalTestItems++

			case "NAPTestScoreSummary":
				var tss xml.NAPTestScoreSummary
				decoder.DecodeElement(&tss, &se)

				gtss, err := ge.Encode(tss)
				if err != nil {
					log.Println("Unable to gob-encode nap test-score-summary: ", err)
				}

				// {NAPTestScoreSummary} = object
				batch.Put([]byte(tss.SummaryID), gtss)

				// NAPTestScoreSummary-type:{id} = {id}
				key := []byte("NAPTestScoreSummary:" + tss.SummaryID)
				batch.Put(key, []byte(tss.SummaryID))

				// {school}:NAPTestScoreSummary-type:{id} = {id}
				key = []byte(tss.SchoolInfoRefId + ":NAPTestScoreSummary:" + tss.SummaryID)
				batch.Put(key, []byte(tss.SummaryID))

				totalTestScoreSummarys++

			case "NAPEventStudentLink":
				var e xml.NAPEvent
				decoder.DecodeElement(&e, &se)
				gev, err := ge.Encode(e)
				if err != nil {
					log.Println("Unable to gob-encode nap event link: ", err)
				}

				// {NAPEventStudentLink} = object
				batch.Put([]byte(e.EventID), gev)

				// NAPEventStudentLink-type:{id} = {id}
				key := []byte("NAPEventStudentLink:" + e.EventID)
				batch.Put(key, []byte(e.EventID))

				// {school}:NAPEventStudentLink-type:{id} = {id}
				key = []byte(e.SchoolRefID + ":NAPEventStudentLink:" + e.EventID)
				batch.Put(key, []byte(e.EventID))

				// {student}:NAPEventStudentLink-type:{id} = {id}
				key = []byte(e.SPRefID + ":NAPEventStudentLink:" + e.EventID)
				batch.Put(key, []byte(e.EventID))

				totalEvents++

			case "NAPStudentResponseSet":
				var r xml.NAPResponseSet
				decoder.DecodeElement(&r, &se)

				gr, err := ge.Encode(r)
				if err != nil {
					log.Println("Unable to gob-encode student response set: ", err)
				}

				// {response-id} = object
				batch.Put([]byte(r.ResponseID), gr)

				// response-type:{id} = {id}
				key := []byte("NAPStudentResponseSet:" + r.ResponseID)
				batch.Put(key, []byte(r.ResponseID))

				// {test}:NAPStudentResponseSet-type:{student} = {id}
				key = []byte(r.TestID + ":NAPStudentResponseSet:" + r.StudentID)
				batch.Put(key, []byte(r.ResponseID))

				// responseset_by_student:{sprefid}:{id} = {id}
				key = []byte("responseset_by_student:" + r.StudentID + ":" + r.ResponseID)
				batch.Put(key, []byte(r.ResponseID))

				totalResponses++

			case "NAPCodeFrame":
				var cf xml.NAPCodeFrame
				decoder.DecodeElement(&cf, &se)
				gcf, err := ge.Encode(cf)
				if err != nil {
					log.Println("Unable to gob-encode nap codeframe: ", err)
				}

				// {NAPCodeFrame-id} = object
				batch.Put([]byte(cf.RefId), gcf)

				// NAPCodeFrame-type:{id} = {id}
				key := []byte("NAPCodeFrame:" + cf.RefId)
				batch.Put(key, []byte(cf.RefId))

				totalCodeFrames++

			case "SchoolInfo":
				var si xml.SchoolInfo
				decoder.DecodeElement(&si, &se)
				gsi, err := ge.Encode(si)
				if err != nil {
					log.Println("Unable to gob-encode schoolinfo: ", err)
				}

				// {SchoolInfo-id} = object
				batch.Put([]byte(si.RefId), gsi)

				// SchoolInfo-type:{id} = {id}
				key := []byte("SchoolInfo:" + si.RefId)
				batch.Put(key, []byte(si.RefId))

				// ASL lookup
				// {acara-id} = {refid}
				key = []byte(si.ACARAId + ":")
				batch.Put(key, []byte(si.RefId))

				// SchoolDetails lookup object
				// not a sif object so needs a guid
				sd_id := nuid.Next()
				key = []byte(sd_id)
				sd := SchoolDetails{
					SchoolName: si.SchoolName,
					ACARAId:    si.ACARAId,
				}
				gsd, err := ge.Encode(sd)
				if err != nil {
					log.Println("Unable to gob-encode schooldetails: ", err)
				}
				// {SchoolDetails-id} = object
				batch.Put(key, gsd)

				// SchoolDetails-type:{id} = {id}
				key = []byte("SchoolDetails:" + sd_id)
				batch.Put(key, []byte(sd_id))

				totalSchools++

			case "StudentPersonal":
				var sp xml.RegistrationRecord
				decoder.DecodeElement(&sp, &se)
				gsp, err := ge.Encode(sp)
				if err != nil {
					log.Println("Unable to gob-encode studentpersonal: ", err)
				}

				// {StudentPersonal-id} = object
				batch.Put([]byte(sp.RefId), gsp)

				// StudentPersonal-type:{id} = {id}
				key := []byte("StudentPersonal:" + sp.RefId)
				batch.Put(key, []byte(sp.RefId))

				// student_by_acaraid:{id} = {id}
				key = []byte("student_by_acaraid:" + sp.ASLSchoolId + ":" + sp.RefId)
				batch.Put(key, []byte(sp.RefId))

				totalStudents++

			}
		default:
		}

	}

	batcherr := db.Write(batch, nil)
	if batcherr != nil {
		log.Fatalln("batch error: ", batcherr)
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

	log.Printf("ingestion complete for [%s]", resultsFilePath)

}
