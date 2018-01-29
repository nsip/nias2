package sifxml


    type StudentPeriodAttendance struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      Date string `xml:"Date"`
      SessionInfoRefId IdRefType `xml:"SessionInfoRefId"`
      ScheduledActivityRefId IdRefType `xml:"ScheduledActivityRefId"`
      TimetablePeriod string `xml:"TimetablePeriod"`
      TimeIn string `xml:"TimeIn"`
      TimeOut string `xml:"TimeOut"`
      AttendanceCode AttendanceCodeType `xml:"AttendanceCode"`
      AttendanceStatus string `xml:"AttendanceStatus"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      AuditInfo AuditInfoType `xml:"AuditInfo"`
      AttendanceComment string `xml:"AttendanceComment"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    