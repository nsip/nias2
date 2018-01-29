package sifxml


    type StudentAttendanceSummary struct {
        StudentAttendanceSummaryRefId IdRefType `xml:"StudentAttendanceSummaryRefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      StartDate string `xml:"StartDate"`
      EndDate string `xml:"EndDate"`
      StartDay string `xml:"StartDay"`
      EndDay string `xml:"EndDay"`
      FTE string `xml:"FTE"`
      DaysAttended string `xml:"DaysAttended"`
      ExcusedAbsences string `xml:"ExcusedAbsences"`
      UnexcusedAbsences string `xml:"UnexcusedAbsences"`
      DaysTardy string `xml:"DaysTardy"`
      DaysInMembership string `xml:"DaysInMembership"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    