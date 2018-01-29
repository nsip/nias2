package sifxml


    type StudentPersonal struct {
        RefId RefIdType `xml:"RefId,attr"`
      AlertMessages AlertMessagesType `xml:"AlertMessages"`
      MedicalAlertMessages MedicalAlertMessagesType `xml:"MedicalAlertMessages"`
      LocalId LocalId `xml:"LocalId"`
      StateProvinceId StateProvinceId `xml:"StateProvinceId"`
      ElectronicIdList ElectronicIdListType `xml:"ElectronicIdList"`
      OtherIdList OtherIdListType `xml:"OtherIdList"`
      PersonInfo PersonInfoType `xml:"PersonInfo"`
      ProjectedGraduationYear ProjectedGraduationYearType `xml:"ProjectedGraduationYear"`
      OnTimeGraduationYear OnTimeGraduationYearType `xml:"OnTimeGraduationYear"`
      GraduationDate GraduationDateType `xml:"GraduationDate"`
      MostRecent StudentMostRecentContainerType `xml:"MostRecent"`
      AcceptableUsePolicy string `xml:"AcceptableUsePolicy"`
      GiftedTalented string `xml:"GiftedTalented"`
      EconomicDisadvantage string `xml:"EconomicDisadvantage"`
      ESL string `xml:"ESL"`
      ESLDateAssessed string `xml:"ESLDateAssessed"`
      YoungCarersRole string `xml:"YoungCarersRole"`
      Disability string `xml:"Disability"`
      IntegrationAide string `xml:"IntegrationAide"`
      EducationSupport string `xml:"EducationSupport"`
      HomeSchooledStudent string `xml:"HomeSchooledStudent"`
      Sensitive string `xml:"Sensitive"`
      OfflineDelivery string `xml:"OfflineDelivery"`
      PrePrimaryEducation string `xml:"PrePrimaryEducation"`
      FirstAUSchoolEnrollment string `xml:"FirstAUSchoolEnrollment"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    