// schooldata.go
package naprr

import (
	"github.com/nsip/nias2/xml"
	"log"
)

// brings together all school-specific data for analysis and reports
//

type SchoolData struct {
	SchoolInfos    map[string]xml.SchoolInfo
	Events         map[string]xml.NAPEvent
	ScoreSummaries map[string]xml.NAPTestScoreSummary
	Responses      map[string]xml.NAPResponseSet
	Students       map[string]xml.RegistrationRecord
	ACARAId        string
}

func NewSchoolData(acaraid string) *SchoolData {
	sd := SchoolData{
		SchoolInfos:    make(map[string]xml.SchoolInfo),
		Events:         make(map[string]xml.NAPEvent),
		ScoreSummaries: make(map[string]xml.NAPTestScoreSummary),
		Responses:      make(map[string]xml.NAPResponseSet),
		Students:       make(map[string]xml.RegistrationRecord),
		ACARAId:        acaraid,
	}
	return &sd
}

func (sd *SchoolData) PrintSummary() {
	log.Printf("\nSchool: %s", sd.ACARAId)
	log.Printf("students: %d", len(sd.Students))
	log.Printf("events: %d", len(sd.Events))
	log.Printf("schoolinfos: %d", len(sd.SchoolInfos))
	log.Printf("school summaries: %d", len(sd.ScoreSummaries))
	log.Printf("responses: %d", len(sd.Responses))
	log.Printf("\n\n")

}
