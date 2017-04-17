package xml

type NAPTestItem struct {
	ItemID          string          `xml:"RefId,attr"`
	TestItemContent TestItemContent `xml:"TestItemContent"`
}

type TestItemContent struct {
	NAPTestItemLocalId        string `xml:"NAPTestItemLocalId"`
	ItemName                  string `xml:"ItemName,omitempty"`
	ItemType                  string `xml:"ItemType,omitempty"`
	Subdomain                 string `xml:"Subdomain,omitempty"`
	WritingGenre              string `xml:"WritingGenre,omitempty"`
	ItemDescriptor            string `xml:"ItemDescriptor,omitempty"`
	ReleasedStatus            string `xml:"ReleasedStatus,omitempty"`
	MarkingType               string `xml:"MarkingType,omitempty"`
	MultipleChoiceOptionCount string `xml:"MultipleChoiceOptionCount,omitempty"`
	CorrectAnswer             string `xml:"CorrectAnswer,omitempty"`
	MaximumScore              string `xml:"MaximumScore,omitempty"`
	ItemDifficulty            string `xml:"ItemDifficulty,omitempty"`
	ItemDifficultyLogit5      string `xml:"ItemDifficultyLogit5,omitempty"`
	ItemDifficultyLogit62     string `xml:"ItemDifficultyLogit62,omitempty"`
	ItemDifficultyLogit5SE    string `xml:"ItemDifficultyLogit5SE,omitempty"`
	ItemDifficultyLogit62SE   string `xml:"ItemDifficultyLogit62SE,omitempty"`
	ItemProficiencyBand       string `xml:"ItemProficiencyBand,omitempty"`
	ItemProficiencyLevel      string `xml:"ItemProficiencyLevel,omitempty"`
	ExemplarURL               string `xml:"ExemplarURL,omitempty"`

	ItemSubstitutedForList struct {
		SubstituteItem []struct {
			SubstituteItemRefId string `xml:"SubstituteItemRefId,omitempty"`
			LocalId             string `xml:"SubstituteItemLocalId,omitempty"`
			PNPCodeList         struct {
				PNPCode []string `xml:"PNPCode,omitempty"`
			} `xml:"PNPCodeList,omitempty"`
		} `xml:"SubstituteItem,omitempty"`
	} `xml:"ItemSubstitutedForList,omitempty"`

	ContentDescriptionList struct {
		ContentDescription []string `xml:"ContentDescription,omitempty"`
	} `xml:"ContentDescriptionList,omitempty"`

	StimulusList struct {
		Stimulus []struct {
			LocalId        string `xml:"StimulusLocalId"`
			TextGenre      string `xml:"TextGenre"`
			TextType       string `xml:"TextType"`
			WordCount      string `xml:"WordCount"`
			TextDescriptor string `xml:"TextDescriptor"`
			Content        string `xml:"Content"`
		} `xml:"Stimulus,omitempty"`
	} `xml:"StimulusList,omitempty"`

	NAPWritingRubricList struct {
		NAPWritingRubric []struct {
			RubricType string `xml:"RubricType,omitempty"`
			Descriptor string `xml:"Descriptor,omitempty"`
			ScoreList  struct {
				Score []struct {
					MaxScoreValue        string `xml:"MaxScoreValue,omitempty"`
					ScoreDescriptionList struct {
						ScoreDescription []struct {
							ScoreValue string `xml:"ScoreValue,omitempty"`
							Descriptor string `xml:"Descriptor,omitempty"`
						} `xml:"ScoreDescription,omitempty"`
					} `xml:"ScoreDescriptionList,omitempty"`
				} `xml:"Score,omitempty"`
			} `xml:"ScoreList,omitempty"`
		} `xml:"NAPWritingRubric,omitempty"`
	} `xml:"NAPWritingRubricList,omitempty"`
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
