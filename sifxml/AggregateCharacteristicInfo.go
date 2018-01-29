package sifxml


    type AggregateCharacteristicInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      Description string `xml:"Description"`
      Definition string `xml:"Definition"`
      ElementName string `xml:"ElementName"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    