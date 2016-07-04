// aslservice.go
package nias2

// validation service to check acara school ids in registration records

import (
	"github.com/wildducktheories/go-csv"
	"log"
	"os"
)

type ASLService struct{}

// lookup structure
type ASLSearchTerms struct {
	AcaraID string
	State   string
}

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

	if !ASLKeyExists(req) {
		desc := "ASL ID " + rr.ASLSchoolId + " not found in ASL list of valid IDs"
		ve := ValidationError{
			Description:  desc,
			Field:        "ASLSchoolID",
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
			stateid, err := GetASLValue(req)
			if err != nil {
				return responses, err
			}
			if stateid != rr.StateTerritory {
				desc := "ASL ID " + rr.ASLSchoolId + " is a valid ID, but not for " + rr.StateTerritory
				ve := ValidationError{
					Description:  desc,
					Field:        "ASLSchoolID",
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

// utility function to add ASL data to lookup database,
// only needs to be invoked once at application startup
func LoadASLLookupData() {

	f, err := os.Open("./schoolslist/asl_schools_20160321.csv")
	reader := csv.WithIoReader(f)
	records, err := csv.ReadAll(reader)
	log.Printf("ASL records read: %v", len(records))
	if err != nil {
		log.Fatalf("Unable to open ASL schools file, service aborting...")
	}

	for _, r := range records {

		r := r.AsMap()

		s := ASLSearchTerms{
			AcaraID: r["ACARA ID"],
			State:   r["State"],
		}
		SetASLValue(s)
	}
	log.Println("...all ASL records loaded for validation")

}
