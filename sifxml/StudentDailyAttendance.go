package sifxml


    type StudentDailyAttendance struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      Date string `xml:"Date"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      DayValue string `xml:"DayValue"`
      AttendanceCode AttendanceCodeType `xml:"AttendanceCode"`
      AttendanceStatus string `xml:"AttendanceStatus"`
      TimeIn string `xml:"TimeIn"`
      TimeOut string `xml:"TimeOut"`
      AbsenceValue string `xml:"AbsenceValue"`
      AttendanceNote string `xml:"AttendanceNote"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    