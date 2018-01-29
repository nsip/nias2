package sifxml


    type NAPCodeFrame struct {
        RefId RefIdType `xml:"RefId,attr"`
      NAPTestRefId IdRefType `xml:"NAPTestRefId"`
      TestContent NAPTestContentType `xml:"TestContent"`
      TestletList NAPCodeFrameTestletListType `xml:"TestletList"`
      SIF_Metadata SIF_Metadata `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElements `xml:"SIF_ExtendedElements"`
      
      }
    