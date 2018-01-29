package sifxml


    type GradingAssignment struct {
        RefId RefIdType `xml:"RefId,attr"`
      TeachingGroupRefId IdRefType `xml:"TeachingGroupRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      GradingCategory string `xml:"GradingCategory"`
      Description string `xml:"Description"`
      PointsPossible string `xml:"PointsPossible"`
      CreateDate string `xml:"CreateDate"`
      DueDate string `xml:"DueDate"`
      Weight string `xml:"Weight"`
      MaxAttemptsAllowed string `xml:"MaxAttemptsAllowed"`
      DetailedDescriptionURL string `xml:"DetailedDescriptionURL"`
      DetailedDescriptionBinary string `xml:"DetailedDescriptionBinary"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    