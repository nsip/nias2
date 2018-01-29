package sifxml


    type ChargedLocationInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      LocationType string `xml:"LocationType"`
      SiteCategory string `xml:"SiteCategory"`
      Name string `xml:"Name"`
      Description string `xml:"Description"`
      LocalId LocalIdType `xml:"LocalId"`
      StateProvinceId StateProvinceIdType `xml:"StateProvinceId"`
      ParentChargedLocationInfoRefId IdRefType `xml:"ParentChargedLocationInfoRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      AddressList AddressListType `xml:"AddressList"`
      PhoneNumberList PhoneNumberListType `xml:"PhoneNumberList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    