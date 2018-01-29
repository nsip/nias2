package sifxml


    type CalendarSummary struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      LocalId LocalIdType `xml:"LocalId"`
      Description string `xml:"Description"`
      DaysInSession string `xml:"DaysInSession"`
      StartDate string `xml:"StartDate"`
      EndDate string `xml:"EndDate"`
      FirstInstructionDate string `xml:"FirstInstructionDate"`
      LastInstructionDate string `xml:"LastInstructionDate"`
      GraduationDate GraduationDateType `xml:"GraduationDate"`
      InstructionalMinutes string `xml:"InstructionalMinutes"`
      MinutesPerDay string `xml:"MinutesPerDay"`
      YearLevels YearLevelsType `xml:"YearLevels"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    