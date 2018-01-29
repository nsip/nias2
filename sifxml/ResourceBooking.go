package sifxml


    type ResourceBooking struct {
        RefId RefIdType `xml:"RefId,attr"`
      ResourceRefId ResourceBooking_ResourceRefId `xml:"ResourceRefId"`
      ResourceLocalId LocalIdType `xml:"ResourceLocalId"`
      StartDateTime string `xml:"StartDateTime"`
      FinishDateTime string `xml:"FinishDateTime"`
      FromPeriod LocalIdType `xml:"FromPeriod"`
      ToPeriod LocalIdType `xml:"ToPeriod"`
      Booker IdRefType `xml:"Booker"`
      Reason string `xml:"Reason"`
      ScheduledActivityRefId IdRefType `xml:"ScheduledActivityRefId"`
      KeepOld string `xml:"KeepOld"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type ResourceBooking_ResourceRefId struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
