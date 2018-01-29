package sifxml


    type LearningStandardDocument struct {
        RefId RefIdType `xml:"RefId,attr"`
      Title string `xml:"Title"`
      Description string `xml:"Description"`
      RichDescription AbstractContentElementType `xml:"RichDescription"`
      Source string `xml:"Source"`
      Organizations OrganizationsType `xml:"Organizations"`
      Authors AuthorsType `xml:"Authors"`
      OrganizationContactPoint string `xml:"OrganizationContactPoint"`
      SubjectAreas ACStrandAreaListType `xml:"SubjectAreas"`
      DocumentStatus string `xml:"DocumentStatus"`
      DocumentDate string `xml:"DocumentDate"`
      LocalAdoptionDate string `xml:"LocalAdoptionDate"`
      LocalArchiveDate string `xml:"LocalArchiveDate"`
      EndOfLifeDate string `xml:"EndOfLifeDate"`
      Copyright LearningStandardDocument_Copyright `xml:"Copyright"`
      YearLevels YearLevelsType `xml:"YearLevels"`
      RepositoryDate string `xml:"RepositoryDate"`
      LearningStandardItemRefId IdRefType `xml:"LearningStandardItemRefId"`
      RelatedLearningStandards LearningStandardsDocumentType `xml:"RelatedLearningStandards"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type LearningStandardDocument_Copyright struct {
       Date string `xml:"Date"`
       Holder string `xml:"Holder"`
}
