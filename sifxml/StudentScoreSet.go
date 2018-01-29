package sifxml


    type StudentScoreSet struct {
        RefId RefIdType `xml:"RefId,attr"`
      ScoreMetric string `xml:"ScoreMetric"`
      AssessmentAdministrationRefId IdRefType `xml:"AssessmentAdministrationRefId"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      AssessmentRegistrationRefId IdRefType `xml:"AssessmentRegistrationRefId"`
      Scores ScoresType `xml:"Scores"`
      StartDateTime string `xml:"StartDateTime"`
      FinishDateTime string `xml:"FinishDateTime"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    