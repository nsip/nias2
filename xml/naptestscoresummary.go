package xml

type NAPTestScoreSummary struct {
	SummaryID                     string `xml:"RefId,attr"`
	SchoolInfoRefId               string `xml:"SchoolInfoRefId,omitempty"`
	SchoolACARAId                 string `xml:"SchoolACARAId"`
	NAPTestRefId                  string `xml:"NAPTestRefId,omitempty"`
	NAPTestLocalId                string `xml:"NAPTestLocalId"`
	DomainNationalAverage         string `xml:"DomainNationalAverage,omitempty"`
	DomainSchoolAverage           string `xml:"DomainSchoolAverage,omitempty"`
	DomainJurisdictionAverage     string `xml:"DomainJurisdictionAverage,omitempty"`
	DomainTopNational60Percent    string `xml:"DomainTopNational60Percent,omitempty"`
	DomainBottomNational60Percent string `xml:"DomainBottomNational60Percent,omitempty"`
}

func (t NAPTestScoreSummary) GetHeaders() []string {
	return []string{
		"SummaryID", "SchoolInfoRefId", "SchoolACARAId", "NAPTestRefId", "NAPTestLocalId",
		"DomainNationalAverage", "DomainSchoolAverage", "DomainJurisdictionAverage",
		"DomainTopNational60Percent", "DomainBottomNational60Percent"}
}

func (t NAPTestScoreSummary) GetSlice() []string {
	return []string{
		t.SummaryID, t.SchoolInfoRefId, t.SchoolACARAId, t.NAPTestRefId, t.NAPTestLocalId,
		t.DomainNationalAverage, t.DomainSchoolAverage, t.DomainJurisdictionAverage,
		t.DomainTopNational60Percent, t.DomainBottomNational60Percent}
}
