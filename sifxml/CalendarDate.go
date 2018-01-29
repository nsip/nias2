package sifxml


    type CalendarDate struct {
        CalendarDateRefId IdRefType `xml:"CalendarDateRefId,attr"`
      Date string `xml:"Date"`
      CalendarSummaryRefId IdRefType `xml:"CalendarSummaryRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      CalendarDateType CalendarDateInfoType `xml:"CalendarDateType"`
      CalendarDateNumber string `xml:"CalendarDateNumber"`
      StudentAttendance AttendanceInfoType `xml:"StudentAttendance"`
      TeacherAttendance AttendanceInfoType `xml:"TeacherAttendance"`
      AdministratorAttendance AttendanceInfoType `xml:"AdministratorAttendance"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    