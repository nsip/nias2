package sifxml


    type NAPTestItem struct {
        RefId RefIdType `xml:"RefId,attr"`
      TestItemContent NAPTestItemContentType `xml:"TestItemContent"`
      SIF_Metadata SIF_Metadata `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElements `xml:"SIF_ExtendedElements"`
      
      }
    