package sifxml


    type NAPTest struct {
        RefId RefIdType `xml:"RefId,attr"`
      TestContent NAPTestContentType `xml:"TestContent"`
      SIF_Metadata SIF_Metadata `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElements `xml:"SIF_ExtendedElements"`
      
      }
    