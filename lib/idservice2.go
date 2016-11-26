// idservice2.go
package nias2

import (
	"fmt"
	"github.com/streamrail/concurrent-map"
	"github.com/tylertreat/BoomFilters"
	"log"
)

// core bloom filter with large window for duplicates
var ibf = boom.NewInverseBloomFilter(200000)

// threadsafe map to keep list of locations where data has been seen
var prevLocations = cmap.New()

// alternative idservice based on bloom filters
// much faster than lookup version, but lose traceability of which record
// is duplicated; if absolute speed is required use this version
// this is not the default - more info for users on which records
// are duplicated assumed to be better than raw performance for now
type IDService2 struct{}

// create a new id service instance
func NewIDService2() (*IDService2, error) {
	ids := IDService2{}
	return &ids, nil
}

// implement the nias Service interface
func (ids *IDService2) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)

	rr, ok := req.Body.(RegistrationRecord)
	if !ok {
		log.Println("IDService received a message that is not a RegistrationRecord, ignoring")
		return responses, nil
	}

	// create the lookup key
	id := fmt.Sprintf("%s%s%s%s%s%s", req.TxID,
		rr.LocalId, rr.ASLSchoolId, rr.FamilyName, rr.GivenName, rr.BirthDate)

	if !ibf.TestAndAdd([]byte(id)) {
		// we haven't seen this id before
	} else {
		desc := "This record has duplicates in the dataset"
		ve := ValidationError{
			Description:  desc,
			Field:        "Multiple (see description)",
			OriginalLine: req.SeqNo,
			Vtype:        "identity",
		}
		r := NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)
	}

	return responses, nil
}
