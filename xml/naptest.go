package xml

type NAPTest struct {
	TestID      string      `xml:"RefId,attr"`
	TestContent TestContent `xml:"TestContent"`
}

type TestContent struct {
	LocalId           string            `xml:"NAPTestLocalId,omitempty"`
	TestName          string            `xml:"TestName,omitempty"`
	TestLevel         string            `xml:"TestLevel>Code"`
	TestDomain        string            `xml:"Domain"`
	TestYear          string            `xml:"TestYear"`
	StagesCount       string            `xml:"StagesCount,omitempty"`
	TestType          string            `xml:"TestType,omitempty"`
	DomainBands       DomainBands       `xml:"DomainBands,omitempty"`
	DomainProficiency DomainProficiency `xml:"DomainProficiency,omitempty"`
}

type DomainBands struct {
	Band1Lower  string `xml:"Band1Lower"`
	Band1Upper  string `xml:"Band1Upper"`
	Band2Lower  string `xml:"Band2Lower"`
	Band2Upper  string `xml:"Band2Upper"`
	Band3Lower  string `xml:"Band3Lower"`
	Band3Upper  string `xml:"Band3Upper"`
	Band4Lower  string `xml:"Band4Lower"`
	Band4Upper  string `xml:"Band4Upper"`
	Band5Lower  string `xml:"Band5Lower"`
	Band5Upper  string `xml:"Band5Upper"`
	Band6Lower  string `xml:"Band6Lower"`
	Band6Upper  string `xml:"Band6Upper"`
	Band7Lower  string `xml:"Band7Lower"`
	Band7Upper  string `xml:"Band7Upper"`
	Band8Lower  string `xml:"Band8Lower"`
	Band8Upper  string `xml:"Band8Upper"`
	Band9Lower  string `xml:"Band9Lower"`
	Band9Upper  string `xml:"Band9Upper"`
	Band10Lower string `xml:"Band10Lower"`
	Band10Upper string `xml:"Band10Upper"`
}

type DomainProficiency struct {
	Level1Lower string `xml:"Level1Lower"`
	Level1Upper string `xml:"Level1Upper"`
	Level2Lower string `xml:"Level2Lower"`
	Level2Upper string `xml:"Level2Upper"`
	Level3Lower string `xml:"Level3Lower"`
	Level3Upper string `xml:"Level3Upper"`
	Level4Lower string `xml:"Level4Lower"`
	Level4Upper string `xml:"Level4Upper"`
}

func (t NAPTest) GetHeaders() []string {
	return []string{
		"TestID", "LocalId", "TestName", "TestLevel", "TestDomain", "TestYear",
		"StagesCount", "DomainBand1Lower", "DomainBand1Upper",
		"DomainBand2Lower", "DomainBand2Upper",
		"DomainBand3Lower", "DomainBand3Upper",
		"DomainBand4Lower", "DomainBand4Upper",
		"DomainBand5Lower", "DomainBand5Upper",
		"DomainBand6Lower", "DomainBand6Upper",
		"DomainBand7Lower", "DomainBand7Upper",
		"DomainBand8Lower", "DomainBand8Upper",
		"DomainBand9Lower", "DomainBand9Upper",
		"DomainBand10Lower", "DomainBand10Upper",
		"ProficiencyLevel1Lower", "ProficiencyLevel1Upper",
		"ProficiencyLevel2Lower", "ProficiencyLevel2Upper",
		"ProficiencyLevel3Lower", "ProficiencyLevel3Upper",
		"ProficiencyLevel4Lower", "ProficiencyLevel4Upper"}
}

func (t NAPTest) GetSlice() []string {
	return []string{
		t.TestID, t.TestContent.LocalId, t.TestContent.TestName, t.TestContent.TestLevel,
		t.TestContent.TestDomain, t.TestContent.TestYear, t.TestContent.StagesCount,
		t.TestContent.DomainBands.Band1Lower, t.TestContent.DomainBands.Band1Upper,
		t.TestContent.DomainBands.Band2Lower, t.TestContent.DomainBands.Band2Upper,
		t.TestContent.DomainBands.Band3Lower, t.TestContent.DomainBands.Band3Upper,
		t.TestContent.DomainBands.Band4Lower, t.TestContent.DomainBands.Band4Upper,
		t.TestContent.DomainBands.Band5Lower, t.TestContent.DomainBands.Band5Upper,
		t.TestContent.DomainBands.Band6Lower, t.TestContent.DomainBands.Band6Upper,
		t.TestContent.DomainBands.Band7Lower, t.TestContent.DomainBands.Band7Upper,
		t.TestContent.DomainBands.Band8Lower, t.TestContent.DomainBands.Band8Upper,
		t.TestContent.DomainBands.Band9Lower, t.TestContent.DomainBands.Band9Upper,
		t.TestContent.DomainBands.Band10Lower, t.TestContent.DomainBands.Band10Upper,
		t.TestContent.DomainProficiency.Level1Lower, t.TestContent.DomainProficiency.Level1Upper,
		t.TestContent.DomainProficiency.Level2Lower, t.TestContent.DomainProficiency.Level2Upper,
		t.TestContent.DomainProficiency.Level3Lower, t.TestContent.DomainProficiency.Level3Upper,
		t.TestContent.DomainProficiency.Level4Lower, t.TestContent.DomainProficiency.Level4Upper}
}
