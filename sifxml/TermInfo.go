package sifxml


    type TermInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      StartDate string `xml:"StartDate"`
      EndDate string `xml:"EndDate"`
      Description string `xml:"Description"`
      RelativeDuration string `xml:"RelativeDuration"`
      TermCode string `xml:"TermCode"`
      Track string `xml:"Track"`
      TermSpan string `xml:"TermSpan"`
      MarkingTerm string `xml:"MarkingTerm"`
      SchedulingTerm string `xml:"SchedulingTerm"`
      AttendanceTerm string `xml:"AttendanceTerm"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    