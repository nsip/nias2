package sifxml


    type TimeTable struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      LocalId LocalIdType `xml:"LocalId"`
      Title string `xml:"Title"`
      DaysPerCycle string `xml:"DaysPerCycle"`
      PeriodsPerDay string `xml:"PeriodsPerDay"`
      TeachingPeriodsPerDay string `xml:"TeachingPeriodsPerDay"`
      SchoolLocalId LocalIdType `xml:"SchoolLocalId"`
      SchoolName string `xml:"SchoolName"`
      TimeTableCreationDate string `xml:"TimeTableCreationDate"`
      StartDate string `xml:"StartDate"`
      EndDate string `xml:"EndDate"`
      TimeTableDayList TimeTableDayListType `xml:"TimeTableDayList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    