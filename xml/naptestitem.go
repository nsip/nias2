package xml

type NAPTestItem struct {
	ItemID          string `xml:"RefId,attr"`
	TestItemContent struct {
		NAPTestItemLocalId        string `xml:"NAPTestItemLocalId"`
		ItemName                  string `xml:"ItemName"`
		ItemType                  string `xml:"ItemType"`
		Subdomain                 string `xml:"Subdomain"`
		WritingGenre              string `xml:"WritingGenre"`
		ItemDescriptor            string `xml:"ItemDescriptor"`
		ReleasedStatus            string `xml:"ReleasedStatus"`
		MarkingType               string `xml:"MarkingType"`
		MultipleChoiceOptionCount string `xml:"MultipleChoiceOptionCount"`
		CorrectAnswer             string `xml:"CorrectAnswer"`
		MaximumScore              string `xml:"MaximumScore"`
		ItemDifficulty            string `xml:"ItemDifficulty"`
		ItemDifficultyLogit5      string `xml:"ItemDifficultyLogit5"`
		ItemDifficultyLogit62     string `xml:"ItemDifficultyLogit62"`
		ItemDifficultyLogit5SE    string `xml:"ItemDifficultyLogit5SE"`
		ItemDifficultyLogit62SE   string `xml:"ItemDifficultyLogit62SE"`
		ItemProficiencyBand       string `xml:"ItemProficiencyBand"`
		ItemProficiencyLevel      string `xml:"ItemProficiencyLevel"`
		ExemplarURL               string `xml:"ExemplarURL"`

		ItemSubstitutedForList struct {
			SubstituteItem []struct {
				SubstituteItemRefId string `xml:"SubstituteItemRefId"`
				LocalId             string `xml:"SubstituteItemLocalId"`
				PNPCodeList         struct {
					PNPCode []string `xml:"PNPCode"`
				} `xml:"PNPCodeList"`
			} `xml:"SubstituteItem"`
		} `xml:"ItemSubstitutedForList"`

		ContentDescriptionList struct {
			ContentDescription []string `xml:"ContentDescription"`
		} `xml:"ContentDescriptionList"`

		StimulusList struct {
			Stimulus []struct {
				LocalId        string `xml:"StimulusLocalId"`
				TextGenre      string `xml:"TextGenre"`
				TextType       string `xml:"TextType"`
				WordCount      string `xml:"WordCount"`
				TextDescriptor string `xml:"TextDescriptor"`
				Content        string `xml:"Content"`
			} `xml:"Stimulus"`
		} `xml:"StimulusList"`

		NAPWritingRubricList struct {
			NAPWritingRubric []struct {
				RubricType string `xml:"RubricType"`
				Descriptor string `xml:"Descriptor"`
				ScoreList  struct {
					Score []struct {
						MaxScoreValue        string `xml:"MaxScoreValue"`
						ScoreDescriptionList struct {
							ScoreDescription []struct {
								ScoreValue string `xml:"ScoreValue"`
								Descriptor string `xml:"Descriptor"`
							} `xml:"ScoreDescription"`
						} `xml:"ScoreDescriptionList"`
					} `xml:"Score"`
				} `xml:"ScoreList"`
			} `xml:"NAPWritingRubric"`
		} `xml:"NAPWritingRubricList"`
	} `xml:"TestItemContent"`
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
