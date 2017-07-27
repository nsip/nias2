package naprrql

import (
	goxml "encoding/xml"
	"log"

	"github.com/nsip/nias2/naprr"
	"github.com/nsip/nias2/xml"
	"github.com/syndtr/goleveldb/leveldb"
)

func IngestResultsFile(resultsFilePath string) {

	db := GetDB()
	ge := naprr.GobEncoder{}

	// open the data file for streaming read
	xmlFile, err := naprr.OpenResultsFile(resultsFilePath)
	if err != nil {
		log.Fatalln("unable to open results file")
	}

	log.Println("Reading data file...")

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
				// entries = badger.EntriesSet(entries, []byte(t.TestID), gt)
				batch.Put([]byte(t.TestID), gt)

				// NAPTest-type:{id} = id
				key := []byte("NAPTest:" + t.TestID)
				// entries = badger.EntriesSet(entries, key, []byte(t.TestID))
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
				// entries = badger.EntriesSet(entries, []byte(tl.TestletID), gtl)
				batch.Put([]byte(tl.TestletID), gtl)

				// NAPTestlet-type:{id} = {id}
				key := []byte("NAPTestlet:" + tl.TestletID)
				// entries = badger.EntriesSet(entries, key, []byte(tl.TestletID))
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
				// entries = badger.EntriesSet(entries, []byte(ti.ItemID), gti)
				batch.Put([]byte(ti.ItemID), gti)

				// NapTestItem-type:{id} = {id}
				key := []byte("NAPTestItem:" + ti.ItemID)
				// entries = badger.EntriesSet(entries, key, []byte(ti.ItemID))
				batch.Put(key, []byte(ti.ItemID))

				totalTestItems++

			case "NAPTestScoreSummary":
				var tss xml.NAPTestScoreSummary
				decoder.DecodeElement(&tss, &se)

				// jtss, err := jsonEncode(tss)
				// if err != nil {
				// 	log.Println("Unable to json-encode nap test-score-summary: ", err)
				// }

				gtss, err := ge.Encode(tss)
				if err != nil {
					log.Println("Unable to gob-encode nap test-score-summary: ", err)
				}

				// var tss2 interface{}
				// err = ge.Decode(gtss, &tss2)
				// if err != nil {
				// 	log.Println("summary did not decode: ", err)
				// 	// log.Printf("\n\ndecoded summary:\n\n%+v\n", tss2)
				// }
				// tss3, ok := tss2.(xml.NAPTestScoreSummary)
				// if !ok {
				// 	log.Println("summary type assertion failed")
				// 	log.Printf("\n\ndecoded summary:\n\n%+v\n", tss3)
				// }

				// {NAPTestScoreSummary} = object
				// entries = badger.EntriesSet(entries, []byte(tss.SummaryID), jtss)
				batch.Put([]byte(tss.SummaryID), gtss)

				// NAPTestScoreSummary-type:{id} = {id}
				key := []byte("NAPTestScoreSummary:" + tss.SummaryID)
				// entries = badger.EntriesSet(entries, key, []byte(tss.SummaryID))
				batch.Put(key, []byte(tss.SummaryID))

				// {school}:NAPTestScoreSummary-type:{id} = {id}
				key = []byte(tss.SchoolInfoRefId + ":NAPTestScoreSummary:" + tss.SummaryID)
				// entries = badger.EntriesSet(entries, key, []byte(tss.SummaryID))
				batch.Put(key, []byte(tss.SummaryID))

				totalTestScoreSummarys++

			case "NAPEventStudentLink":
				var e xml.NAPEvent
				decoder.DecodeElement(&e, &se)
				gev, err := ge.Encode(e)
				if err != nil {
					log.Println("Unable to gob-encode nap event link: ", err)
				}

				// var e2 xml.NAPEvent
				// err = ge.Decode(gev, &e2)
				// if err != nil {
				// 	log.Println("event did not decode: ", err)
				// 	// log.Printf("\n\ndecoded event:\n\n%+v\n", e2)
				// }

				// {NAPEventStudentLink} = object
				// entries = badger.EntriesSet(entries, []byte(e.EventID), gev)
				batch.Put([]byte(e.EventID), gev)

				// NAPEventStudentLink-type:{id} = {id}
				key := []byte("NAPEventStudentLink:" + e.EventID)
				// entries = badger.EntriesSet(entries, key, []byte(e.EventID))
				batch.Put(key, []byte(e.EventID))

				// {school}:NAPEventStudentLink-type:{id} = {id}
				key = []byte(e.SchoolRefID + ":NAPEventStudentLink:" + e.EventID)
				// entries = badger.EntriesSet(entries, key, []byte(e.EventID))
				batch.Put(key, []byte(e.EventID))

				totalEvents++

			case "NAPStudentResponseSet":
				var r xml.NAPResponseSet
				decoder.DecodeElement(&r, &se)

				// log.Printf("response: \n\n%v\n\n", r)

				// jr, err := jsonEncode(r)
				// if err != nil {
				// 	log.Println("Unable to json-encode student response set: ", err)
				// }

				// b := new(bytes.Buffer)
				// enc := gob.NewEncoder(b)
				// if err := enc.Encode(&r); err != nil {
				// 	log.Println("unable to gob-encode response set: ", err)
				// }

				gr, err := ge.Encode(r)
				if err != nil {
					log.Println("Unable to gob-encode student response set: ", err)
				}

				// var r2 xml.NAPResponseSet
				// err = ge.Decode(gr, &r2)
				// if err != nil {
				// 	log.Println("response did not decode: ", err)
				// 	// log.Printf("\n\ndecoded response:\n\n%+v\n", r2)
				// }

				// {response-id} = object
				// entries = badger.EntriesSet(entries, []byte(r.ResponseID), b.Bytes())
				batch.Put([]byte(r.ResponseID), gr)

				// response-type:{id} = {id}
				key := []byte("NAPStudentResponseSet:" + r.ResponseID)
				// entries = badger.EntriesSet(entries, key, []byte(r.ResponseID))
				batch.Put(key, []byte(r.ResponseID))

				// {test}:NAPStudentResponseSet-type:{student} = {id}
				key = []byte(r.TestID + ":NAPStudentResponseSet:" + r.StudentID)
				// entries = badger.EntriesSet(entries, key, []byte(r.ResponseID))
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
				// entries = badger.EntriesSet(entries, []byte(cf.RefId), gcf)
				batch.Put([]byte(cf.RefId), gcf)

				// NAPCodeFrame-type:{id} = {id}
				key := []byte("NAPCodeFrame:" + cf.RefId)
				// entries = badger.EntriesSet(entries, key, []byte(cf.RefId))
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
				// entries = badger.EntriesSet(entries, []byte(si.RefId), gsi)
				batch.Put([]byte(si.RefId), gsi)

				// SchoolInfo-type:{id} = {id}
				key := []byte("SchoolInfo:" + si.RefId)
				// entries = badger.EntriesSet(entries, key, []byte(si.RefId))
				batch.Put(key, []byte(si.RefId))

				// ASL lookup
				key = []byte(si.ACARAId)
				// entries = badger.EntriesSet(entries, key, []byte(si.RefId))
				batch.Put(key, []byte(si.RefId))

				totalSchools++

			case "StudentPersonal":
				var sp xml.RegistrationRecord
				decoder.DecodeElement(&sp, &se)
				gsp, err := ge.Encode(sp)
				if err != nil {
					log.Println("Unable to gob-encode studentpersonal: ", err)
				}

				// {StudentPersonal-id} = object
				// entries = badger.EntriesSet(entries, []byte(sp.RefId), gsp)
				batch.Put([]byte(sp.RefId), gsp)

				// StudentPersonal-type:{id} = {id}
				key := []byte("StudentPersonal:" + sp.RefId)
				// entries = badger.EntriesSet(entries, key, []byte(sp.RefId))
				batch.Put(key, []byte(sp.RefId))

				// {ASL-school-id}:StudentPersonal-type:{id} = {id}
				key = []byte(sp.ASLSchoolId + ":StudentPersonal:" + sp.RefId)
				// entries = badger.EntriesSet(entries, key, []byte(sp.RefId))
				batch.Put(key, []byte(sp.RefId))

				totalStudents++

			}
		default:
		}

	}

	// fmt.Println("entries:")
	// for _, entry := range entries {
	// 	fmt.Printf("key: %s\n val: %s\n\n", entry.Key, entry.Value)
	// }

	batcherr := db.Write(batch, nil)
	if batcherr != nil {
		log.Fatalln("batch error: ", batcherr)
	}
	log.Println("Closing db...")
	err = db.Close()
	if err != nil {
		log.Println("Error closing database: ", err)
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

	log.Printf("ingestion complete for %s", resultsFilePath)

}
