package sifxml


    type StaffAssignment struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      StaffPersonalRefId IdRefType `xml:"StaffPersonalRefId"`
      Description string `xml:"Description"`
      PrimaryAssignment string `xml:"PrimaryAssignment"`
      JobStartDate string `xml:"JobStartDate"`
      JobEndDate string `xml:"JobEndDate"`
      JobFTE string `xml:"JobFTE"`
      JobFunction string `xml:"JobFunction"`
      EmploymentStatus string `xml:"EmploymentStatus"`
      StaffSubjectList StaffSubjectListType `xml:"StaffSubjectList"`
      StaffActivity StaffActivityExtensionType `xml:"StaffActivity"`
      YearLevels YearLevelsType `xml:"YearLevels"`
      CasualReliefTeacher string `xml:"CasualReliefTeacher"`
      Homegroup string `xml:"Homegroup"`
      House string `xml:"House"`
      CalendarSummaryList CalendarSummaryListType `xml:"CalendarSummaryList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    