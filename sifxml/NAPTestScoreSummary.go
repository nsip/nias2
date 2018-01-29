package sifxml


    type NAPTestScoreSummary struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SchoolACARAId LocalIdType `xml:"SchoolACARAId"`
      NAPTestRefId IdRefType `xml:"NAPTestRefId"`
      NAPTestLocalId LocalId `xml:"NAPTestLocalId"`
      DomainNationalAverage string `xml:"DomainNationalAverage"`
      DomainSchoolAverage string `xml:"DomainSchoolAverage"`
      DomainJurisdictionAverage string `xml:"DomainJurisdictionAverage"`
      DomainTopNational60Percent string `xml:"DomainTopNational60Percent"`
      DomainBottomNational60Percent string `xml:"DomainBottomNational60Percent"`
      SIF_Metadata SIF_Metadata `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElements `xml:"SIF_ExtendedElements"`
      
      }
    