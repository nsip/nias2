package sifxml


    type StudentAttendanceTimeList struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      Date string `xml:"Date"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      AttendanceTimes AttendanceTimesType `xml:"AttendanceTimes"`
      PeriodAttendances PeriodAttendancesType `xml:"PeriodAttendances"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    