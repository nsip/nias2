package sifxml


    type StudentActivityParticipation struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      StudentActivityInfoRefId IdRefType `xml:"StudentActivityInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      ParticipationComment string `xml:"ParticipationComment"`
      StartDate string `xml:"StartDate"`
      EndDate string `xml:"EndDate"`
      Role string `xml:"Role"`
      RecognitionList RecognitionListType `xml:"RecognitionList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    