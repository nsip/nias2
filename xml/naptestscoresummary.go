package xml

type NAPTestScoreSummary struct {
	SummaryID                     string `xml:"RefId,attr"`
	SchoolInfoRefId               string `xml:"SchoolInfoRefId"`
	SchoolACARAId                 string `xml:"SchoolACARAId"`
	NAPTestRefId                  string `xml:"NAPTestRefId"`
	NAPTestLocalId                string `xml:"NAPTestLocalId"`
	DomainNationalAverage         string `xml:"DomainNationalAverage"`
	DomainSchoolAverage           string `xml:"DomainSchoolAverage"`
	DomainJurisdictionAverage     string `xml:"DomainJurisdictionAverage"`
	DomainTopNational60Percent    string `xml:"DomainTopNational60Percent"`
	DomainBottomNational60Percent string `xml:"DomainBottomNational60Percent"`
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
