package sifxml


    type WellbeingAlert struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      Date string `xml:"Date"`
      WellbeingAlertStartDate string `xml:"WellbeingAlertStartDate"`
      WellbeingAlertEndDate string `xml:"WellbeingAlertEndDate"`
      WellbeingAlertCategory string `xml:"WellbeingAlertCategory"`
      WellbeingAlertDescription string `xml:"WellbeingAlertDescription"`
      EnrolmentRestricted string `xml:"EnrolmentRestricted"`
      AlertAudience string `xml:"AlertAudience"`
      AlertSeverity string `xml:"AlertSeverity"`
      AlertKeyContact string `xml:"AlertKeyContact"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    