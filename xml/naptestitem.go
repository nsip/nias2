package xml

type NAPTestItem struct {
	ItemID          string          `xml:"RefId,attr"`
	TestItemContent TestItemContent `xml:"TestItemContent"`
}

type TestItemContent struct {
	NAPTestItemLocalId        string                 `xml:"NAPTestItemLocalId"`
	ItemName                  string                 `xml:"ItemName"`
	ItemType                  string                 `xml:"ItemType"`
	Subdomain                 string                 `xml:"Subdomain"`
	WritingGenre              string                 `xml:"WritingGenre,omitempty"`
	ItemDescriptor            string                 `xml:"ItemDescriptor"`
	ReleasedStatus            string                 `xml:"ReleasedStatus"`
	MarkingType               string                 `xml:"MarkingType"`
	MultipleChoiceOptionCount string                 `xml:"MultipleChoiceOptionCount,omitempty"`
	CorrectAnswer             string                 `xml:"CorrectAnswer,omitempty"`
	MaximumScore              string                 `xml:"MaximumScore"`
	ItemDifficulty            string                 `xml:"ItemDifficulty"`
	ItemDifficultyLogit5      string                 `xml:"ItemDifficultyLogit5"`
	ItemDifficultyLogit62     string                 `xml:"ItemDifficultyLogit62"`
	ItemDifficultyLogit5SE    string                 `xml:"ItemDifficultyLogit5SE"`
	ItemDifficultyLogit62SE   string                 `xml:"ItemDifficultyLogit62SE"`
	ItemProficiencyBand       string                 `xml:"ItemProficiencyBand"`
	ItemProficiencyLevel      string                 `xml:"ItemProficiencyLevel,omitempty"`
	ExemplarURL               string                 `xml:"ExemplarURL,omitempty"`
	ItemSubstitutedForList    ItemSubstitutedForList `xml:"ItemSubstitutedForList,omitempty"`
	ContentDescriptionList    ContentDescriptionList `xml:"ContentDescriptionList,omitempty"`
	StimulusList              StimulusList           `xml:"StimulusList,omitempty"`
	NAPWritingRubricList      NAPWritingRubricList   `xml:"NAPWritingRubricList,omitempty"`
}

type ItemSubstitutedForList struct {
	SubstituteItem []SubstituteItem `xml:"SubstituteItem,omitempty"`
}

type SubstituteItem struct {
	SubstituteItemRefId string      `xml:"SubstituteItemRefId"`
	LocalId             string      `xml:"SubstituteItemLocalId,omitempty"`
	PNPCodeList         PNPCodeList `xml:"PNPCodeList,omitempty"`
}

type PNPCodeList struct {
	PNPCode []string `xml:"PNPCode,omitempty"`
}

type ContentDescriptionList struct {
	ContentDescription []string `xml:"ContentDescription,omitempty"`
}

type StimulusList struct {
	Stimulus []Stimulus `xml:"Stimulus,omitempty"`
}

type Stimulus struct {
	LocalId        string `xml:"StimulusLocalId"`
	TextGenre      string `xml:"TextGenre,omitempty"`
	TextType       string `xml:"TextType,omitempty"`
	WordCount      string `xml:"WordCount,omitempty"`
	TextDescriptor string `xml:"TextDescriptor"`
	Content        string `xml:"Content"`
}

type NAPWritingRubricList struct {
	NAPWritingRubric []NAPWritingRubric `xml:"NAPWritingRubric,omitempty"`
}

type NAPWritingRubric struct {
	RubricType string    `xml:"RubricType"`
	ScoreList  ScoreList `xml:"ScoreList,omitempty"`
	Descriptor string    `xml:"Descriptor,omitempty"`
}

type ScoreList struct {
	Score []Score `xml:"Score,omitempty"`
}

type Score struct {
	MaxScoreValue        string               `xml:"MaxScoreValue"`
	ScoreDescriptionList ScoreDescriptionList `xml:"ScoreDescriptionList"`
}

type ScoreDescriptionList struct {
	ScoreDescription []ScoreDescription `xml:"ScoreDescription,omitempty"`
}

type ScoreDescription struct {
	ScoreValue string `xml:"ScoreValue"`
	Descriptor string `xml:"Descriptor"`
}

func (t NAPTestItem) GetHeaders() []string {
	return []string{"ItemID", "NAPTestItemLocalId", "ItemName", "ItemType", "Subdomain", "WritingGenre",
		"ItemDescriptor", "ReleasedStatus", "MarkingType", "MultipleChoiceOptionCount", "CorrectAnswer",
		"MaximumScore", "ItemDifficulty", "ItemDifficultyLogit5", "ItemDifficultyLogit62",
		"ItemDifficultyLogit5SE", "ItemDifficultyLogit62SE", "ItemProficiencyBand", "ItemProficiencyLevel", "ExemplarURL"}
}

func (t NAPTestItem) GetSlice() []string {
	return []string{t.ItemID, t.TestItemContent.NAPTestItemLocalId, t.TestItemContent.ItemName, t.TestItemContent.ItemType,
		t.TestItemContent.Subdomain, t.TestItemContent.WritingGenre, t.TestItemContent.ItemDescriptor,
		t.TestItemContent.ReleasedStatus, t.TestItemContent.MarkingType, t.TestItemContent.MultipleChoiceOptionCount,
		t.TestItemContent.CorrectAnswer, t.TestItemContent.MaximumScore, t.TestItemContent.ItemDifficulty,
		t.TestItemContent.ItemDifficultyLogit5, t.TestItemContent.ItemDifficultyLogit62,
		t.TestItemContent.ItemDifficultyLogit5SE, t.TestItemContent.ItemDifficultyLogit62SE,
		t.TestItemContent.ItemProficiencyBand, t.TestItemContent.ItemProficiencyLevel, t.TestItemContent.ExemplarURL}
}
