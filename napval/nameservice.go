// namevalidservice.go
package napval

import (
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
	"log"
	"regexp"
)

// name validator service

// implementation of the id service
type NameValidService struct {
}

var invalidCharRegexp = regexp.MustCompile("[^" + config.LegalNameChars + "]")

// create a new id service instance
func NewNameValidService() (*NameValidService, error) {
	num := NameValidService{}
	return &num, nil
}

// implement the nias Service interface
func (num *NameValidService) HandleMessage(req *lib.NiasMessage) ([]lib.NiasMessage, error) {

	responses := make([]lib.NiasMessage, 0)

	rr, ok := req.Body.(xml.RegistrationRecord)
	if !ok {
		log.Println("NameValidService received a message that is not a RegistrationRecord, ignoring")
		return responses, nil
	}

	if len(invalidCharRegexp.FindString(rr.FamilyName)) > 0 {
		desc := "Family Name contains suspect character"
		responses = add_name_error(responses, desc, "FamilyName", req)
	}
	if len(invalidCharRegexp.FindString(rr.GivenName)) > 0 {
		desc := "Given Name contains suspect character"
		responses = add_name_error(responses, desc, "GivenName", req)
	}
	if len(invalidCharRegexp.FindString(rr.MiddleName)) > 0 {
		desc := "Middle Name contains suspect character"
		responses = add_name_error(responses, desc, "MiddleName", req)
	}
	if len(invalidCharRegexp.FindString(rr.PreferredName)) > 0 {
		desc := "Preferred Name contains suspect character"
		responses = add_name_error(responses, desc, "PreferredName", req)
	}

	return responses, nil

}

func add_name_error(responses []lib.NiasMessage, desc string, field string, req *lib.NiasMessage) []lib.NiasMessage {
	ve := ValidationError{
		Description:  desc,
		Field:        field,
		OriginalLine: req.SeqNo,
		Vtype:        "name",
		Severity:     "error",
	}
	r := lib.NiasMessage{}
	r.TxID = req.TxID
	r.SeqNo = req.SeqNo
	// r.Target = VALIDATION_PREFIX
	r.Body = ve
	responses = append(responses, r)
	return responses
}
