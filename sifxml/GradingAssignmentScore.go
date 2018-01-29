package sifxml


    type GradingAssignmentScore struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      StudentPersonalLocalId LocalIdType `xml:"StudentPersonalLocalId"`
      TeachingGroupRefId IdRefType `xml:"TeachingGroupRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      GradingAssignmentRefId IdRefType `xml:"GradingAssignmentRefId"`
      ScorePoints string `xml:"ScorePoints"`
      ScorePercent string `xml:"ScorePercent"`
      ScoreLetter string `xml:"ScoreLetter"`
      ScoreDescription string `xml:"ScoreDescription"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    