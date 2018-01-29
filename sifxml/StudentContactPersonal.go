package sifxml


    type StudentContactPersonal struct {
        RefId RefIdType `xml:"RefId,attr"`
      LocalId LocalIdType `xml:"LocalId"`
      OtherIdList OtherIdListType `xml:"OtherIdList"`
      PersonInfo PersonInfoType `xml:"PersonInfo"`
      EmploymentType string `xml:"EmploymentType"`
      SchoolEducationalLevel EducationalLevelType `xml:"SchoolEducationalLevel"`
      NonSchoolEducation string `xml:"NonSchoolEducation"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    