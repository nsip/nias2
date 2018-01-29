package sifxml


    type StudentActivityInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      Title string `xml:"Title"`
      Description string `xml:"Description"`
      StudentActivityType StudentActivityType `xml:"StudentActivityType"`
      StudentActivityLevel string `xml:"StudentActivityLevel"`
      YearLevels YearLevelsType `xml:"YearLevels"`
      CurricularStatus string `xml:"CurricularStatus"`
      Location LocationType `xml:"Location"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    