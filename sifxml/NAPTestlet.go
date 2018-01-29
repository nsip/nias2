package sifxml


    type NAPTestlet struct {
        RefId RefIdType `xml:"RefId,attr"`
      NAPTestRefId IdRefType `xml:"NAPTestRefId"`
      NAPTestLocalId LocalId `xml:"NAPTestLocalId"`
      TestletContent NAPTestletContentType `xml:"TestletContent"`
      TestItemList NAPTestItemListType `xml:"TestItemList"`
      SIF_Metadata SIF_Metadata `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElements `xml:"SIF_ExtendedElements"`
      
      }
    