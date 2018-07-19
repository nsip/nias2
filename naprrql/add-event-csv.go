package naprrql

import (
	"encoding/csv"
	"encoding/json"
	exml "encoding/xml"
	"fmt"
	"github.com/nsip/nias2/xml"
	"github.com/syndtr/goleveldb/leveldb"
	csvutils "github.com/wildducktheories/go-csv"
	"log"
	"os"
)

func AddEventCSV(csvFileName string) error {
	PNPEvents, err := ReadEventCSV(csvFileName)
	batch := new(leveldb.Batch)
	updateEvents(PNPEvents, batch)
	batcherr := db.Write(batch, nil)
	if batcherr != nil {
		log.Fatalln("batch error: ", batcherr)
	}
	batch.Reset()
	return err
}

func updateEvents(PNPEvents map[string]xml.Adjustment, batch *leveldb.Batch) error {
	for k, v := range PNPEvents {
		studentEventIds := getIdentifiers("event_by_PSI_ACARAID_LocalTestID:" + k)
		if len(studentEventIds) < 1 {
			continue
		}
		eventObjs, err := getObjects(studentEventIds)
		if err != nil {
			return err
		}
		for _, eventObj := range eventObjs {
			event := eventObj.(xml.NAPEvent)
			// TODO: make this merge
			event.Adjustment = v
			b, _ := json.Marshal(event)
			event1 := xml.NAPEvent{}
			json.Unmarshal(b, &event1)
			out, _ := exml.MarshalIndent(event1, "", "  ")
			log.Printf("%#v\n", string(out))
			PutEvent(event, batch)
		}
	}
	return nil
}

func ReadEventCSV(csvFileName string) (map[string]xml.Adjustment, error) {
	ret := make(map[string]xml.Adjustment)
	f, err := os.Open(csvFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open csv file: " + csvFileName)
	}

	csvr := csv.NewReader(f)
	reader := csvutils.WithCsvReader(csvr, f)
	defer f.Close()
	defer reader.Close()
	for record := range reader.C() {
		r := record.AsMap()
		key := r["PSI"] + ":" + r["SchoolID"] + ":" + r["NAPTestLocalID"]
		var adjustment xml.Adjustment
		if err := json.Unmarshal([]byte(r["Adjustment"]), &adjustment); err != nil {
			return nil, fmt.Errorf("csv file processing error: %s", err)
		}
		ret[key] = adjustment
	}
	return ret, nil

}

func PutEvent(e xml.NAPEvent, batch *leveldb.Batch) {
	ge := GobEncoder{}
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

	// event_by_student_test:{student}:{test}:{id} = {id}
	key = []byte("event_by_student_test:" + e.SPRefID + ":" + e.TestID + ":" + e.EventID)
	batch.Put(key, []byte(e.EventID))

	// {test}:NAPEventStudentLink-type:{school}:{id} = {id}
	key = []byte(e.TestID + ":NAPEventStudentLink:" + e.SchoolRefID + ":" + e.EventID)
	batch.Put(key, []byte(e.EventID))

	// {test}:NAPEventStudentLink-type:{school}:{id} = {id}
	key = []byte("event_by_PSI_ACARAID_LocalTestID:" + e.PSI + ":" + e.SchoolID + ":" + e.NAPTestLocalID)
	batch.Put(key, []byte(e.EventID))

}
