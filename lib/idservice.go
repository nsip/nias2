// idservice.go
package nias2

import (
	"log"
)

// identity checking service for registration records
// lookups are done using speicialised hash keys to meet the specifications
// for duplicate student records
// Entries are created in the ledis hashtable on first sight
// if the same key is retrieved in a subsequent message then identity
// is assumed to be a duplicate.

// set up the data structures used as hash keys
// simple duplicate check is, have we seen this userid for this school before
type IDSimpleKey struct {
	LocalId     string
	ASLSchoolId string
}

// this checks the user against a set of likely colliding matches
type IDExtendedKey struct {
	LocalId     string
	ASLSchoolId string
	FamilyName  string
	GivenName   string
	BirthDate   string
}

// implementation of the id service
type IDService struct{}

// create a new id service instance
func NewIDService() (*IDService, error) {
	ids := IDService{}

	return &ids, nil
}

// implement the nias Service interface
func (ids *IDService) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)

	_, ok := req.Body.(RegistrationRecord)
	if !ok {
		log.Println("IDService received a message that is not a RegistrationRecord, ignoring")
		return responses, nil
	}

	// see if this id has been seen before
	if SimpleAndComplexIDKeySetnx(req) {
		// if not save this id and move on
		return responses, nil
	}

	// if the record exists, check simple & complex matches
	if ol, err := ComplexIDKeySeen(req); err != nil {
		desc := "Potential duplicate of record: " + ol + "\n" +
			"based on matching: student local id, school asl id, family & given names and birthdate"
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

	} else if ol, err := SimpleIDKeySeen(req); err != nil {
		desc := "LocalID (Student) and ASL ID (School) are potential duplicate of record: " + ol

		ve := ValidationError{
			Description:  desc,
			Field:        "LocalID/ASL ID",
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

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
