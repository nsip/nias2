package sifxml


    type SchoolCourseInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolLocalId LocalIdType `xml:"SchoolLocalId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      TermInfoRefId IdRefType `xml:"TermInfoRefId"`
      CourseCode string `xml:"CourseCode"`
      StateCourseCode string `xml:"StateCourseCode"`
      DistrictCourseCode string `xml:"DistrictCourseCode"`
      SubjectAreaList SubjectAreaListType `xml:"SubjectAreaList"`
      CourseTitle string `xml:"CourseTitle"`
      Description string `xml:"Description"`
      InstructionalLevel string `xml:"InstructionalLevel"`
      CourseCredits string `xml:"CourseCredits"`
      CoreAcademicCourse string `xml:"CoreAcademicCourse"`
      GraduationRequirement string `xml:"GraduationRequirement"`
      Department string `xml:"Department"`
      CourseContent string `xml:"CourseContent"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    