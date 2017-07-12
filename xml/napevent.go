package xml

type NAPEvent struct {
	EventID                string             `xml:"RefId,attr"`
	SPRefID                string             `xml:"StudentPersonalRefId"`
	PSI                    string             `xml:"PlatformStudentIdentifier,omitempty"`
	SchoolRefID            string             `xml:"SchoolInfoRefId,omitempty"`
	SchoolID               string             `xml:"SchoolACARAId,omitempty"`
	TestID                 string             `xml:"NAPTestRefId"`
	NAPTestLocalID         string             `xml:"NAPTestLocalId,omitempty"`
	SchoolSector           string             `xml:"SchoolSector,omitempty"`
	System                 string             `xml:"System,omitempty"`
	SchoolGeolocation      string             `xml:"SchoolGeolocation,omitempty"`
	ReportingSchoolName    string             `xml:"ReportingSchoolName,omitempty"`
	JurisdictionID         string             `xml:"JurisdictionID,omitempty"`
	ParticipationCode      string             `xml:"ParticipationCode"`
	ParticipationText      string             `xml:"ParticipationText,omitempty"`
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
		"System", "SchoolGeolocation", "ReportingSchoolName", "JurisdictionID", "ParticipationCode",
		"ParticipationText", "Device", "Date", "StartTime", "LapsedTimeTest", "ExemptionReason",
		"PersonalDetailsChanged", "PossibleDuplicate", "DOBRange", "BookletType"}
}

func (t NAPEvent) GetSlice() []string {
	return []string{
		t.EventID, t.SPRefID, t.PSI, t.SchoolRefID, t.SchoolID, t.TestID, t.NAPTestLocalID, t.SchoolSector,
		t.System, t.SchoolGeolocation, t.ReportingSchoolName, t.JurisdictionID, t.ParticipationCode,
		t.ParticipationText, t.Device, t.Date, t.StartTime, t.LapsedTimeTest, t.ExemptionReason,
		t.PersonalDetailsChanged, t.PossibleDuplicate, t.DOBRange, t.Adjustment.BookletType}
}
