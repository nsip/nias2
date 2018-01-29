package sifxml


    type WellbeingAppeal struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      WellbeingResponseRefId IdRefType `xml:"WellbeingResponseRefId"`
      LocalAppealId LocalIdType `xml:"LocalAppealId"`
      AppealStatusCode string `xml:"AppealStatusCode"`
      Date string `xml:"Date"`
      AppealNotes string `xml:"AppealNotes"`
      AppealOutcome string `xml:"AppealOutcome"`
      DocumentList WellbeingDocumentListType `xml:"DocumentList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    