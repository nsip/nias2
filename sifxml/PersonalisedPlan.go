package sifxml


    type PersonalisedPlan struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      PersonalisedPlanCategory string `xml:"PersonalisedPlanCategory"`
      PersonalisedPlanStartDate string `xml:"PersonalisedPlanStartDate"`
      PersonalisedPlanEndDate string `xml:"PersonalisedPlanEndDate"`
      PersonalisedPlanReviewDate string `xml:"PersonalisedPlanReviewDate"`
      PersonalisedPlanNotes string `xml:"PersonalisedPlanNotes"`
      DocumentList WellbeingDocumentListType `xml:"DocumentList"`
      AssociatedAttachment string `xml:"AssociatedAttachment"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    