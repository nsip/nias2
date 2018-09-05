package xml

import (
	"encoding/xml"
)

type NAPEvent struct {
	XMLName                xml.Name           `xml:"NAPEventStudentLink"`
	EventID                string             `xml:"RefId,attr"`
	SPRefID                string             `xml:"StudentPersonalRefId,omitempty"`
	PSI                    string             `xml:"PlatformStudentIdentifier"`
	SchoolRefID            string             `xml:"SchoolInfoRefId"`
	SchoolID               string             `xml:"SchoolACARAId,omitempty"`
	TestID                 string             `xml:"NAPTestRefId,omitempty"`
	NAPTestLocalID         string             `xml:"NAPTestLocalId"`
	SchoolSector           string             `xml:"SchoolSector"`
	System                 string             `xml:"System,omitempty"`
	SchoolGeolocation      string             `xml:"SchoolGeolocation,omitempty"`
	ReportingSchoolName    string             `xml:"ReportingSchoolName,omitempty"`
	NAPJurisdiction        string             `xml:"NAPJurisdiction,omitempty"`
	ParticipationCode      string             `xml:"ParticipationCode"`
	ParticipationText      string             `xml:"ParticipationText"`
	Device                 string             `xml:"Device,omitempty"`
	Date                   string             `xml:"Date,omitempty"`
	StartTime              string             `xml:"StartTime,omitempty"`
	LapsedTimeTest         string             `xml:"LapsedTimeTest,omitempty"`
	ExemptionReason        string             `xml:"ExemptionReason,omitempty"`
	PersonalDetailsChanged string             `xml:"PersonalDetailsChanged,omitempty"`
	PossibleDuplicate      string             `xml:"PossibleDuplicate,omitempty"`
	DOBRange               string             `xml:"DOBRange,omitempty"`
	TestDisruptionList     TestDisruptionList `xml:"TestDisruptionList,omitempty"`
	Adjustment             Adjustment         `xml:"Adjustment,omitempty"`
}

type TestDisruptionList struct {
	TestDisruption []TestDisruption `xml:"TestDisruption,omitempty"`
}

type TestDisruption struct {
	Event string `xml:"Event,omitempty"`
}

type Adjustment struct {
	PNPCodelist PNPCodelist `xml:"PNPCodeList,omitempty"`
	BookletType string      `xml:"BookletType,omitempty"`
}

type PNPCodelist struct {
	PNPCode []string `xml:"PNPCode,omitempty"`
}

func (t NAPEvent) GetHeaders() []string {
	return []string{
		"EventID", "SPRefID", "PSI", "SchoolRefID", "SchoolID", "TestID", "NAPTestLocalID", "SchoolSector",
		"System", "SchoolGeolocation", "ReportingSchoolName", "NAPJurisdiction", "ParticipationCode",
		"ParticipationText", "Device", "Date", "StartTime", "LapsedTimeTest", "ExemptionReason",
		"PersonalDetailsChanged", "PossibleDuplicate", "DOBRange", "BookletType"}
}

func (t NAPEvent) GetSlice() []string {
	return []string{
		t.EventID, t.SPRefID, t.PSI, t.SchoolRefID, t.SchoolID, t.TestID, t.NAPTestLocalID, t.SchoolSector,
		t.System, t.SchoolGeolocation, t.ReportingSchoolName, t.NAPJurisdiction, t.ParticipationCode,
		t.ParticipationText, t.Device, t.Date, t.StartTime, t.LapsedTimeTest, t.ExemptionReason,
		t.PersonalDetailsChanged, t.PossibleDuplicate, t.DOBRange, t.Adjustment.BookletType}
}
