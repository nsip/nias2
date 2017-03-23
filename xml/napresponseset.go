package xml

type NAPResponseSet struct {
	ResponseID            string `xml:"RefId,attr"`
	ReportExclusionFlag   string `xml:"ReportExclusionFlag"`
	CalibrationSampleFlag string `xml:"CalibrationSampleFlag"`
	EquatingSampleFlag    string `xml:"EquatingSampleFlag"`
	PathTakenForDomain    string `xml:"PathTakenForDomain"`
	ParallelTest          string `xml:"ParallelTest"`
	StudentID             string `xml:"StudentPersonalRefId"`
	PSI                   string `xml:"PlatformStudentIdentifier"`
	TestID                string `xml:"NAPTestRefId"`
	TestLocalID           string `xml:"NAPTestLocalId"`

	DomainScore struct {
		RawScore                      string `xml:"RawScore"`
		ScaledScoreValue              string `xml:"ScaledScoreValue"`
		ScaledScoreLogitValue         string `xml:"ScaledScoreLogitValue"`
		ScaledScoreStandardError      string `xml:"ScaledScoreStandardError"`
		ScaledScoreLogitStandardError string `xml:"ScaledScoreLogitStandardError"`
		StudentDomainBand             string `xml:"StudentDomainBand"`
		StudentProficiency            string `xml:"StudentProficiency"`
		PlausibleScaledValueList      struct {
			PlausibleScaledValue []string `xml:"PlausibleScaledValue"`
		} `xml:"PlausibleScaledValueList"`
	} `xml:"DomainScore"`

	TestletList struct {
		Testlet []NAPResponseSet_Testlet `xml:"Testlet"`
	} `xml:"TestletList"`
}

type NAPResponseSet_Testlet struct {
	NapTestletRefId   string `xml:"NAPTestletRefId"`
	NapTestletLocalId string `xml:"NAPTestletLocalId"`
	TestletScore      string `xml:"TestletSubScore"`
	ItemResponseList  struct {
		ItemResponse []NAPResponseSet_ItemResponse `xml:"ItemResponse"`
	} `xml:"ItemResponseList"`
}

type NAPResponseSet_Subscore struct {
	SubscoreType  string `xml:"SubscoreType"`
	SubscoreValue string `xml:"SubscoreValue"`
}

type NAPResponseSet_ItemResponse struct {
	ItemRefID           string `xml:"NAPTestItemRefId"`
	LocalID             string `xml:"NAPTestItemLocalId"`
	Response            string `xml:"Response"`
	ResponseCorrectness string `xml:"ResponseCorrectness"`
	Score               string `xml:"Score"`
	LapsedTimeItem      string `xml:"LapsedTimeItem"`
	SequenceNumber      string `xml:"SequenceNumber"`
	ItemWeight          string `xml:"ItemWeight"`
	SubscoreList        struct {
		Subscore []NAPResponseSet_Subscore `xml:"Subscore"`
	} `xml:"SubscoreList"`
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
