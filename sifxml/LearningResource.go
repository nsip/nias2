package sifxml


    type LearningResource struct {
        RefId RefIdType `xml:"RefId,attr"`
      Name string `xml:"Name"`
      Author string `xml:"Author"`
      Contacts ContactsType `xml:"Contacts"`
      Location LearningResource_Location `xml:"Location"`
      Status string `xml:"Status"`
      Description string `xml:"Description"`
      YearLevels YearLevelsType `xml:"YearLevels"`
      SubjectAreas ACStrandAreaListType `xml:"SubjectAreas"`
      MediaTypes MediaTypesType `xml:"MediaTypes"`
      UseAgreement string `xml:"UseAgreement"`
      AgreementDate string `xml:"AgreementDate"`
      Approvals ApprovalsType `xml:"Approvals"`
      Evaluations EvaluationsType `xml:"Evaluations"`
      Components ComponentsType `xml:"Components"`
      LearningStandards LearningStandardsType `xml:"LearningStandards"`
      LearningResourcePackageRefId IdRefType `xml:"LearningResourcePackageRefId"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type LearningResource_Location struct {
      ReferenceType string `xml:"ReferenceType,attr"`
      Value string `xml:",chardata"`
}
