package sifxml


    type RoomInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      LocalId LocalIdType `xml:"LocalId"`
      RoomNumber string `xml:"RoomNumber"`
      StaffList StaffListType `xml:"StaffList"`
      Description string `xml:"Description"`
      Building string `xml:"Building"`
      HomeroomNumber string `xml:"HomeroomNumber"`
      Size string `xml:"Size"`
      Capacity string `xml:"Capacity"`
      PhoneNumber PhoneNumberType `xml:"PhoneNumber"`
      RoomType string `xml:"RoomType"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    