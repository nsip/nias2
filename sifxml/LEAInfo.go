package sifxml


    type LEAInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      LocalId LocalIdType `xml:"LocalId"`
      StateProvinceId StateProvinceIdType `xml:"StateProvinceId"`
      CommonwealthId string `xml:"CommonwealthId"`
      LEAName string `xml:"LEAName"`
      LEAURL string `xml:"LEAURL"`
      EducationAgencyType AgencyType `xml:"EducationAgencyType"`
      LEAContactList LEAContactListType `xml:"LEAContactList"`
      PhoneNumberList PhoneNumberListType `xml:"PhoneNumberList"`
      AddressList AddressListType `xml:"AddressList"`
      OperationalStatus OperationalStatusType `xml:"OperationalStatus"`
      JurisdictionLowerHouse string `xml:"JurisdictionLowerHouse"`
      SLA string `xml:"SLA"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    