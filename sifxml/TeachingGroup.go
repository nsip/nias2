package sifxml


    type TeachingGroup struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      LocalId LocalIdType `xml:"LocalId"`
      ShortName string `xml:"ShortName"`
      LongName string `xml:"LongName"`
      GroupType string `xml:"GroupType"`
      Set string `xml:"Set"`
      Block string `xml:"Block"`
      CurriculumLevel string `xml:"CurriculumLevel"`
      SchoolInfoRefId RefIdType `xml:"SchoolInfoRefId"`
      SchoolLocalId LocalIdType `xml:"SchoolLocalId"`
      SchoolCourseInfoRefId RefIdType `xml:"SchoolCourseInfoRefId"`
      SchoolCourseLocalId LocalIdType `xml:"SchoolCourseLocalId"`
      TimeTableSubjectRefId RefIdType `xml:"TimeTableSubjectRefId"`
      TimeTableSubjectLocalId LocalIdType `xml:"TimeTableSubjectLocalId"`
      Semester string `xml:"Semester"`
      StudentList StudentListType `xml:"StudentList"`
      TeacherList TeacherListType `xml:"TeacherList"`
      MinClassSize string `xml:"MinClassSize"`
      MaxClassSize string `xml:"MaxClassSize"`
      TeachingGroupPeriodList TeachingGroupPeriodListType `xml:"TeachingGroupPeriodList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    