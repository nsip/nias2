package sifxml


    type StudentSectionEnrollment struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SectionInfoRefId IdRefType `xml:"SectionInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      EntryDate string `xml:"EntryDate"`
      ExitDate string `xml:"ExitDate"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    