package xml

type NAPEvent struct {
	EventID                string `xml:"RefId,attr"`
	SPRefID                string `xml:"StudentPersonalRefId"`
	PSI                    string `xml:"PlatformStudentIdentifier"`
	SchoolRefID            string `xml:"SchoolInfoRefId"`
	SchoolID               string `xml:"SchoolACARAId"`
	TestID                 string `xml:"NAPTestRefId"`
	NAPTestLocalID         string `xml:"NAPTestLocalId"`
	SchoolSector           string `xml:"SchoolSector"`
	System                 string `xml:"System"`
	SchoolGeolocation      string `xml:"SchoolGeolocation"`
	ReportingSchoolName    string `xml:"ReportingSchoolName"`
	JurisdictionID         string `xml:"JurisdictionID"`
	ParticipationCode      string `xml:"ParticipationCode"`
	ParticipationText      string `xml:"ParticipationText"`
	Device                 string `xml:"Device"`
	Date                   string `xml:"Date"`
	StartTime              string `xml:"StartTime"`
	LapsedTimeTest         string `xml:"LapsedTimeTest"`
	ExemptionReason        string `xml:"ExemptionReason"`
	PersonalDetailsChanged string `xml:"PersonalDetailsChanged"`
	PossibleDuplicate      string `xml:"PossibleDuplicate"`
	DOBRange               string `xml:"DOBRange"`
	TestDisruptionList     struct {
		TestDisruption []struct {
			Event string `xml:"Event"`
		} `xml:"TestDisruption"`
	} `xml:"TestDisruptionList"`
	Adjustment struct {
		PNPCodelist struct {
			PNPCode []string `xml:"PNPCode"`
		} `xml:"PNPCodeList"`
		BookletType string `xml:"BookletType"`
	} `xml:"Adjustment"`
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
