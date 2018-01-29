package sifxml


    type StaffPersonal struct {
        RefId RefIdType `xml:"RefId,attr"`
      LocalId LocalIdType `xml:"LocalId"`
      StateProvinceId StateProvinceIdType `xml:"StateProvinceId"`
      ElectronicIdList ElectronicIdListType `xml:"ElectronicIdList"`
      OtherIdList OtherIdListType `xml:"OtherIdList"`
      PersonInfo PersonInfoType `xml:"PersonInfo"`
      Title string `xml:"Title"`
      EmploymentStatus string `xml:"EmploymentStatus"`
      MostRecent StaffMostRecentContainerType `xml:"MostRecent"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    