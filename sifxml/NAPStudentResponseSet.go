package sifxml


    type NAPStudentResponseSet struct {
        RefId RefIdType `xml:"RefId,attr"`
      ReportExclusionFlag string `xml:"ReportExclusionFlag"`
      CalibrationSampleFlag string `xml:"CalibrationSampleFlag"`
      EquatingSampleFlag string `xml:"EquatingSampleFlag"`
      PathTakenForDomain string `xml:"PathTakenForDomain"`
      ParallelTest string `xml:"ParallelTest"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      PlatformStudentIdentifier LocalId `xml:"PlatformStudentIdentifier"`
      NAPTestRefId IdRefType `xml:"NAPTestRefId"`
      NAPTestLocalId LocalId `xml:"NAPTestLocalId"`
      DomainScore DomainScoreType `xml:"DomainScore"`
      TestletList NAPStudentResponseTestletListType `xml:"TestletList"`
      SIF_Metadata SIF_Metadata `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElements `xml:"SIF_ExtendedElements"`
      
      }
    