package sifxml


    type VendorInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      Name string `xml:"Name"`
      ContactInfo ContactInfoType `xml:"ContactInfo"`
      CustomerId string `xml:"CustomerId"`
      ABN string `xml:"ABN"`
      RegisteredForGST string `xml:"RegisteredForGST"`
      PaymentTerms string `xml:"PaymentTerms"`
      BPay string `xml:"BPay"`
      BSB string `xml:"BSB"`
      AccountNumber string `xml:"AccountNumber"`
      AccountName string `xml:"AccountName"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    