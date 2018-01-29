package sifxml


    type ReportAuthorityInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      AuthorityName string `xml:"AuthorityName"`
      AuthorityId string `xml:"AuthorityId"`
      AuthorityDepartment string `xml:"AuthorityDepartment"`
      AuthorityLevel string `xml:"AuthorityLevel"`
      ContactInfo ContactInfoType `xml:"ContactInfo"`
      Address AddressType `xml:"Address"`
      PhoneNumber PhoneNumberType `xml:"PhoneNumber"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    