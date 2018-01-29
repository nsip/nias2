package sifxml


    type SectionInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolCourseInfoRefId IdRefType `xml:"SchoolCourseInfoRefId"`
      LocalId LocalIdType `xml:"LocalId"`
      Description string `xml:"Description"`
      SchoolYear SchoolYear `xml:"SchoolYear"`
      TermInfoRefId IdRefType `xml:"TermInfoRefId"`
      MediumOfInstruction MediumOfInstructionType `xml:"MediumOfInstruction"`
      LanguageOfInstruction LanguageOfInstructionType `xml:"LanguageOfInstruction"`
      LocationOfInstruction LocationOfInstructionType `xml:"LocationOfInstruction"`
      SummerSchool string `xml:"SummerSchool"`
      SchoolCourseInfoOverride SchoolCourseInfoOverrideType `xml:"SchoolCourseInfoOverride"`
      CourseSectionCode string `xml:"CourseSectionCode"`
      SectionCode string `xml:"SectionCode"`
      CountForAttendance string `xml:"CountForAttendance"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    