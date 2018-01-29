package sifxml


    type Activity struct {
        RefId RefIdType `xml:"RefId,attr"`
      Title string `xml:"Title"`
      Preamble string `xml:"Preamble"`
      TechnicalRequirements TechnicalRequirementsType `xml:"TechnicalRequirements"`
      SoftwareRequirementList SoftwareRequirementListType `xml:"SoftwareRequirementList"`
      EssentialMaterials EssentialMaterialsType `xml:"EssentialMaterials"`
      LearningObjectives LearningObjectivesType `xml:"LearningObjectives"`
      LearningStandards LearningStandardsType `xml:"LearningStandards"`
      SubjectArea SubjectAreaType `xml:"SubjectArea"`
      Prerequisites PrerequisitesType `xml:"Prerequisites"`
      Students StudentsType `xml:"Students"`
      SourceObjects SourceObjectsType `xml:"SourceObjects"`
      Points string `xml:"Points"`
      ActivityTime ActivityTimeType `xml:"ActivityTime"`
      AssessmentRefId IdRefType `xml:"AssessmentRefId"`
      MaxAttemptsAllowed string `xml:"MaxAttemptsAllowed"`
      ActivityWeight string `xml:"ActivityWeight"`
      Evaluation Activity_Evaluation `xml:"Evaluation"`
      LearningResources LearningResourcesType `xml:"LearningResources"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type Activity_Evaluation struct {
      EvaluationType string `xml:"EvaluationType,attr"`
       Description string `xml:"Description"`
}
