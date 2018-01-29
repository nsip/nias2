package sifxml


    type Debtor struct {
        RefId RefIdType `xml:"RefId,attr"`
      BilledEntity Debtor_BilledEntity `xml:"BilledEntity"`
      AddressList string `xml:"AddressList"`
      BillingName string `xml:"BillingName"`
      BillingNote string `xml:"BillingNote"`
      Discount string `xml:"Discount"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type Debtor_BilledEntity struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
