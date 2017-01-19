package nias2

import (
	"encoding/xml"
	"github.com/beevik/etree"
	"github.com/nsip/nias2/go_SifMessage"
	"log"
)

// implementation of the psi service
type SifValidationService struct {
}

// create a new psi service instance
func NewSifValidationService() (*SifValidationService, error) {
	sif := SifValidationService{}
	return &sif, nil
}

// Filter out elements or attributes. Return array of filtered XML strings,
// with cumulative filtering
func parseToRecord(x string) (interface{}, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(x); err != nil {
		return nil, err
	}
	var err error
	var ret interface{}
	rootelem := doc.Root().Tag
	switch rootelem {
	case "StudentPersonal":
		ret = go_SifMessage.TStudentPersonalType{}
		err = xml.Unmarshal([]byte(x), &ret)
		log.Println(ret)
	default:
		log.Println("####  " + rootelem)
		ret = x
		err = nil
	}

	return ret, err
}

// implement the nias Service interface
func (pri *SifValidationService) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)
	out, err := parseToRecord(req.Body.(string))
	if err != nil {
		return nil, err
	}
	r := NiasMessage{}
	r.TxID = req.TxID
	r.SeqNo = req.SeqNo
	r.Target = STORE_AND_FORWARD_PREFIX + "toRec" + "::"
	r.Body = out
	responses = append(responses, r)
	return responses, nil
}
