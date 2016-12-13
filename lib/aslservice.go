// aslservice.go
package nias2

// validation service to check acara school ids in registration records
//
//
// uses a threadsafe set-filter seeded with the ASL data to avoid calls to db
// which makes significant perfromance differece for large datasets

import (
	// "fmt"
	"github.com/wildducktheories/go-csv"
	"gopkg.in/fatih/set.v0"
	"log"
	"os"
)

type ASLService struct{}

// create set-filter instances, set thread-safe
var ibf_aslid = set.New()
var ibf_aslstateid = set.New()

// create a new id service instance
func NewASLService() (*ASLService, error) {
	asls := ASLService{}
	return &asls, nil
}

// implement the service interface
func (asls *ASLService) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)

	rr, ok := req.Body.(RegistrationRecord)
	if !ok {
		log.Println("ASLService received a message that is not a RegistrationRecord, ignoring")
		return responses, nil
	}

	if !ASLKeyExists(rr.ASLSchoolId) {
		desc := "ASL ID " + rr.ASLSchoolId + " not found in ASL list of valid IDs"
		ve := ValidationError{
			Description:  desc,
			Field:        "ASLSchoolId",
			OriginalLine: req.SeqNo,
			Vtype:        "asl",
		}
		r := NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

	} else {

		if rr.StateTerritory != "" {
			if !ASLStateKeyExists(rr.ASLSchoolId, rr.StateTerritory) {
				desc := "ASL ID " + rr.ASLSchoolId + " is a valid ID, but not for " + rr.StateTerritory
				ve := ValidationError{
					Description:  desc,
					Field:        "ASLSchoolId",
					OriginalLine: req.SeqNo,
					Vtype:        "asl",
				}
				r := NiasMessage{}
				r.TxID = req.TxID
				r.SeqNo = req.SeqNo
				r.Target = VALIDATION_PREFIX
				r.Body = ve
				responses = append(responses, r)

			}

		}

	}

	return responses, nil
}

// see if this ASLid is in the known list
func ASLKeyExists(aslid string) bool {

	return ibf_aslid.Has(aslid)

}

// see if this ASLid is identified with the correct state/territory
func ASLStateKeyExists(aslid, stateid string) bool {

	id := aslid + stateid
	return ibf_aslstateid.Has(id)

}

// load ASL data into the set filters
func LoadASLLookupData() {
	f, err := os.Open("./schoolslist/asl_schools_20161213.csv")
	reader := csv.WithIoReader(f)
	records, err := csv.ReadAll(reader)
	log.Printf("ASL records read: %v", len(records))
	if err != nil {
		log.Fatalf("Unable to open ASL schools file, service aborting...")
	}

	for _, r := range records {

		r := r.AsMap()

		id := r["ACARA ID"]
		id_state := id + r["State"]

		// log.Println("ASLID:", id, "ASL_STATE:", id_state)

		ibf_aslid.Add(id)
		ibf_aslstateid.Add(id_state)

	}
	log.Println("...all ASL records loaded into filter for validation")

}
