package sifxml


    type TimeTableSubject struct {
        RefId RefIdType `xml:"RefId,attr"`
      SubjectLocalId LocalIdType `xml:"SubjectLocalId"`
      AcademicYear YearLevelType `xml:"AcademicYear"`
      AcademicYearRange YearRangeType `xml:"AcademicYearRange"`
      CourseLocalId LocalIdType `xml:"CourseLocalId"`
      SchoolCourseInfoRefId RefIdType `xml:"SchoolCourseInfoRefId"`
      Faculty string `xml:"Faculty"`
      SubjectShortName string `xml:"SubjectShortName"`
      SubjectLongName string `xml:"SubjectLongName"`
      SubjectType string `xml:"SubjectType"`
      ProposedMaxClassSize string `xml:"ProposedMaxClassSize"`
      ProposedMinClassSize string `xml:"ProposedMinClassSize"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolLocalId LocalIdType `xml:"SchoolLocalId"`
      Semester string `xml:"Semester"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      OtherCodeList OtherCodeListType `xml:"OtherCodeList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    