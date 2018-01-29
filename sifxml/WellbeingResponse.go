package sifxml


    type WellbeingResponse struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      Date string `xml:"Date"`
      WellbeingResponseStartDate string `xml:"WellbeingResponseStartDate"`
      WellbeingResponseEndDate string `xml:"WellbeingResponseEndDate"`
      WellbeingResponseCategory string `xml:"WellbeingResponseCategory"`
      WellbeingResponseNotes string `xml:"WellbeingResponseNotes"`
      PersonInvolvementList PersonInvolvementListType `xml:"PersonInvolvementList"`
      Suspension SuspensionContainerType `xml:"Suspension"`
      Detention DetentionContainerType `xml:"Detention"`
      PlanRequired PlanRequiredContainerType `xml:"PlanRequired"`
      Award AwardContainerType `xml:"Award"`
      OtherResponse OtherWellbeingResponseContainerType `xml:"OtherResponse"`
      DocumentList WellbeingDocumentListType `xml:"DocumentList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    