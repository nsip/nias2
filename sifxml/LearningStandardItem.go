package sifxml


    type LearningStandardItem struct {
        RefId RefIdType `xml:"RefId,attr"`
      Resources LResourcesType `xml:"Resources"`
      StandardSettingBody LearningStandardItem_StandardSettingBody `xml:"StandardSettingBody"`
      StandardHierarchyLevel LearningStandardItem_StandardHierarchyLevel `xml:"StandardHierarchyLevel"`
      PredecessorItems LearningStandardsType `xml:"PredecessorItems"`
      StatementCodes StatementCodesType `xml:"StatementCodes"`
      Statements StatementsType `xml:"Statements"`
      YearLevels YearLevelsType `xml:"YearLevels"`
      ACStrandSubjectArea ACStrandSubjectAreaType `xml:"ACStrandSubjectArea"`
      StandardIdentifier LearningStandardItem_StandardIdentifier `xml:"StandardIdentifier"`
      LearningStandardDocumentRefId IdRefType `xml:"LearningStandardDocumentRefId"`
      RelatedLearningStandardItems LearningStandardItem_RelatedLearningStandardItems `xml:"RelatedLearningStandardItems"`
      Level4 string `xml:"Level4"`
      Level5 string `xml:"Level5"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type LearningStandardItem_StandardSettingBody struct {
       Country CountryType `xml:"Country"`
       StateProvince StateProvinceType `xml:"StateProvince"`
       SettingBodyName string `xml:"SettingBodyName"`
}
type LearningStandardItem_StandardHierarchyLevel struct {
       Number string `xml:"Number"`
       Description string `xml:"Description"`
}
type LearningStandardItem_StandardIdentifier struct {
       YearCreated string `xml:"YearCreated"`
       ACStrandSubjectArea ACStrandSubjectAreaType `xml:"ACStrandSubjectArea"`
       StandardNumber string `xml:"StandardNumber"`
       YearLevels YearLevelsType `xml:"YearLevels"`
       Benchmark string `xml:"Benchmark"`
       YearLevel YearLevelType `xml:"YearLevel"`
       IndicatorNumber string `xml:"IndicatorNumber"`
      AlternateIdentificationCodes LearningStandardItem_AlternateIdentificationCodes `xml:"AlternateIdentificationCodes"`
       Organization string `xml:"Organization"`
}
type LearningStandardItem_RelatedLearningStandardItems struct {
      LearningStandardItemRefId []LearningStandardItem_LearningStandardItemRefId `xml:"LearningStandardItemRefId"`
}
type LearningStandardItem_AlternateIdentificationCodes struct {
       AlternateIdentificationCode []string `xml:"AlternateIdentificationCode"`
}
type LearningStandardItem_LearningStandardItemRefId struct {
      RelationshipType string `xml:"RelationshipType,attr"`
      Value IdRefType `xml:",chardata"`
}
