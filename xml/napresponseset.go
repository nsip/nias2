package xml

type NAPResponseSet struct {
	ResponseID            string      `xml:"RefId,attr"`
	ReportExclusionFlag   string      `xml:"ReportExclusionFlag,omitempty"`
	CalibrationSampleFlag string      `xml:"CalibrationSampleFlag,omitempty"`
	EquatingSampleFlag    string      `xml:"EquatingSampleFlag,omitempty"`
	PathTakenForDomain    string      `xml:"PathTakenForDomain,omitempty"`
	ParallelTest          string      `xml:"ParallelTest,omitempty"`
	StudentID             string      `xml:"StudentPersonalRefId,omitempty"`
	PSI                   string      `xml:"PlatformStudentIdentifier"`
	TestID                string      `xml:"NAPTestRefId,omitempty"`
	TestLocalID           string      `xml:"NAPTestLocalId"`
	DomainScore           DomainScore `xml:"DomainScore,omitempty"`
	TestletList           TestletList `xml:"TestletList"`
}

type TestletList struct {
	Testlet []NAPResponseSet_Testlet `xml:"Testlet"`
}

type DomainScore struct {
	RawScore                      string                   `xml:"RawScore"`
	ScaledScoreValue              string                   `xml:"ScaledScoreValue"`
	ScaledScoreLogitValue         string                   `xml:"ScaledScoreLogitValue"`
	ScaledScoreStandardError      string                   `xml:"ScaledScoreStandardError"`
	ScaledScoreLogitStandardError string                   `xml:"ScaledScoreLogitStandardError"`
	StudentDomainBand             string                   `xml:"StudentDomainBand"`
	StudentProficiency            string                   `xml:"StudentProficiency,omitempty"`
	PlausibleScaledValueList      PlausibleScaledValueList `xml:"PlausibleScaledValueList"`
}

type PlausibleScaledValueList struct {
	PlausibleScaledValue []string `xml:"PlausibleScaledValue"`
}

type NAPResponseSet_Testlet struct {
	NapTestletRefId   string           `xml:"NAPTestletRefId,omitempty"`
	NapTestletLocalId string           `xml:"NAPTestletLocalId"`
	TestletScore      string           `xml:"TestletSubScore,omitempty"`
	ItemResponseList  ItemResponseList `xml:"ItemResponseList"`
}

type ItemResponseList struct {
	ItemResponse []NAPResponseSet_ItemResponse `xml:"ItemResponse"`
}

type NAPResponseSet_Subscore struct {
	SubscoreType  string `xml:"SubscoreType,omitempty"`
	SubscoreValue string `xml:"SubscoreValue,omitempty"`
}

type NAPResponseSet_ItemResponse struct {
	ItemRefID           string       `xml:"NAPTestItemRefId,omitempty"`
	LocalID             string       `xml:"NAPTestItemLocalId"`
	Response            string       `xml:"Response,omitempty"`
	ResponseCorrectness string       `xml:"ResponseCorrectness"`
	Score               string       `xml:"Score,omitempty"`
	LapsedTimeItem      string       `xml:"LapsedTimeItem,omitempty"`
	SequenceNumber      string       `xml:"SequenceNumber"`
	ItemWeight          string       `xml:"ItemWeight"`
	SubscoreList        SubscoreList `xml:"SubscoreList,omitempty"`
	Item                NAPTestItem  `xml:"-"`
}

type SubscoreList struct {
	Subscore []NAPResponseSet_Subscore `xml:"Subscore,omitempty"`
}

func (t NAPResponseSet) GetHeaders() []string {
	return []string{"ResponseID", "ReportExclusionFlag", "CalibrationSampleFlag", "EquatingSampleFlag",
		"PathTakenForDomain", "ParallelTest", "StudentID", "PSI", "TestID", "TestLocalID",
		"RawScore", "ScaledScoreValue", "ScaledScoreLogitValue", "ScaledScoreStandardError", "ScaledScoreLogitStandardError",
		"StudentDomainBand", "StudentProficiency"}
}

func (t NAPResponseSet) GetSlice() []string {
	return []string{t.ResponseID, t.ReportExclusionFlag, t.CalibrationSampleFlag, t.EquatingSampleFlag,
		t.PathTakenForDomain, t.ParallelTest, t.StudentID, t.PSI, t.TestID, t.TestLocalID,
		t.DomainScore.RawScore, t.DomainScore.ScaledScoreValue, t.DomainScore.ScaledScoreLogitValue,
		t.DomainScore.ScaledScoreStandardError, t.DomainScore.ScaledScoreLogitStandardError,
		t.DomainScore.StudentDomainBand, t.DomainScore.StudentProficiency}
}
