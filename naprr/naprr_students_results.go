// naprr_naplandata.go
package naprr

import (
	"github.com/nsip/nias2/xml"
	// "log"
)

// structs to hold student metadata around a test
// assumes one test (for now)

// each map indexed on Student RefID, so that RefID can be changed quickly
type StudentAndResultsData struct {
	Students     map[string]xml.RegistrationRecord
	Events       map[string]xml.NAPEvent
	ResponseSets map[string]xml.NAPResponseSet
}

func NewStudentAndResultsData() *StudentAndResultsData {
	nd := StudentAndResultsData{
		Students:     make(map[string]xml.RegistrationRecord),
		Events:       make(map[string]xml.NAPEvent),
		ResponseSets: make(map[string]xml.NAPResponseSet),
	}
	return &nd
}
