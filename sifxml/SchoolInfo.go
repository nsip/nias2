package sifxml


    type SchoolInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      LocalId LocalIdType `xml:"LocalId"`
      StateProvinceId StateProvinceIdType `xml:"StateProvinceId"`
      CommonwealthId string `xml:"CommonwealthId"`
      ACARAId string `xml:"ACARAId"`
      OtherIdList OtherIdListType `xml:"OtherIdList"`
      SchoolName string `xml:"SchoolName"`
      LEAInfoRefId RefIdType `xml:"LEAInfoRefId"`
      OtherLEA SchoolInfo_OtherLEA `xml:"OtherLEA"`
      SchoolDistrict string `xml:"SchoolDistrict"`
      SchoolDistrictLocalId LocalIdType `xml:"SchoolDistrictLocalId"`
      SchoolType string `xml:"SchoolType"`
      SchoolFocusList SchoolFocusListType `xml:"SchoolFocusList"`
      SchoolURL SchoolURLType `xml:"SchoolURL"`
      SchoolEmailList EmailListType `xml:"SchoolEmailList"`
      PrincipalInfo PrincipalInfoType `xml:"PrincipalInfo"`
      SchoolContactList SchoolContactListType `xml:"SchoolContactList"`
      AddressList AddressListType `xml:"AddressList"`
      PhoneNumberList PhoneNumberListType `xml:"PhoneNumberList"`
      SessionType string `xml:"SessionType"`
      YearLevels YearLevelsType `xml:"YearLevels"`
      ARIA string `xml:"ARIA"`
      OperationalStatus OperationalStatusType `xml:"OperationalStatus"`
      FederalElectorate string `xml:"FederalElectorate"`
      Campus CampusContainerType `xml:"Campus"`
      SchoolSector string `xml:"SchoolSector"`
      IndependentSchool string `xml:"IndependentSchool"`
      NonGovSystemicStatus string `xml:"NonGovSystemicStatus"`
      System string `xml:"System"`
      ReligiousAffiliation string `xml:"ReligiousAffiliation"`
      SchoolGeographicLocation string `xml:"SchoolGeographicLocation"`
      LocalGovernmentArea string `xml:"LocalGovernmentArea"`
      JurisdictionLowerHouse string `xml:"JurisdictionLowerHouse"`
      SLA string `xml:"SLA"`
      SchoolCoEdStatus string `xml:"SchoolCoEdStatus"`
      BoardingSchoolStatus string `xml:"BoardingSchoolStatus"`
      YearLevelEnrollmentList YearLevelEnrollmentListType `xml:"YearLevelEnrollmentList"`
      TotalEnrollments TotalEnrollmentsType `xml:"TotalEnrollments"`
      Entity_Open string `xml:"Entity_Open"`
      Entity_Close string `xml:"Entity_Close"`
      SchoolGroupList SchoolGroupListType `xml:"SchoolGroupList"`
      SchoolTimeZone string `xml:"SchoolTimeZone"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type SchoolInfo_OtherLEA struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value RefIdType `xml:",chardata"`
}
