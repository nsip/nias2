package sifxml


    type StudentSchoolEnrollment struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      LocalId LocalId `xml:"LocalId"`
      MembershipType string `xml:"MembershipType"`
      TimeFrame string `xml:"TimeFrame"`
      SchoolYear SchoolYearType `xml:"SchoolYear"`
      EntryDate string `xml:"EntryDate"`
      EntryType StudentEntryContainerType `xml:"EntryType"`
      YearLevel YearLevelType `xml:"YearLevel"`
      Homeroom StudentSchoolEnrollment_Homeroom `xml:"Homeroom"`
      Advisor StudentSchoolEnrollment_Advisor `xml:"Advisor"`
      Counselor StudentSchoolEnrollment_Counselor `xml:"Counselor"`
      Homegroup string `xml:"Homegroup"`
      ACARASchoolId LocalIdType `xml:"ACARASchoolId"`
      ClassCode string `xml:"ClassCode"`
      TestLevel YearLevelType `xml:"TestLevel"`
      ReportingSchool string `xml:"ReportingSchool"`
      House string `xml:"House"`
      IndividualLearningPlan string `xml:"IndividualLearningPlan"`
      Calendar StudentSchoolEnrollment_Calendar `xml:"Calendar"`
      ExitDate string `xml:"ExitDate"`
      ExitStatus StudentExitStatusContainerType `xml:"ExitStatus"`
      ExitType StudentExitContainerType `xml:"ExitType"`
      FTE string `xml:"FTE"`
      FTPTStatus string `xml:"FTPTStatus"`
      FFPOS string `xml:"FFPOS"`
      CatchmentStatus CatchmentStatusContainerType `xml:"CatchmentStatus"`
      RecordClosureReason string `xml:"RecordClosureReason"`
      PromotionInfo PromotionInfoType `xml:"PromotionInfo"`
      PreviousSchool LocalIdType `xml:"PreviousSchool"`
      DestinationSchool LocalIdType `xml:"DestinationSchool"`
      StudentSubjectChoiceList StudentSubjectChoiceListType `xml:"StudentSubjectChoiceList"`
      StartedAtSchoolDate string `xml:"StartedAtSchoolDate"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type StudentSchoolEnrollment_Homeroom struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
type StudentSchoolEnrollment_Advisor struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
type StudentSchoolEnrollment_Counselor struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
type StudentSchoolEnrollment_Calendar struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
