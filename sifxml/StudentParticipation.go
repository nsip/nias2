package sifxml


    type StudentParticipation struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      StudentParticipationAsOfDate string `xml:"StudentParticipationAsOfDate"`
      ProgramType string `xml:"ProgramType"`
      ProgramFundingSources ProgramFundingSourcesType `xml:"ProgramFundingSources"`
      ManagingSchool StudentParticipation_ManagingSchool `xml:"ManagingSchool"`
      ReferralDate string `xml:"ReferralDate"`
      ReferralSource ReferralSourceType `xml:"ReferralSource"`
      ProgramStatus ProgramStatusType `xml:"ProgramStatus"`
      GiftedEligibilityCriteria string `xml:"GiftedEligibilityCriteria"`
      EvaluationParentalConsentDate string `xml:"EvaluationParentalConsentDate"`
      EvaluationDate string `xml:"EvaluationDate"`
      EvaluationExtensionDate string `xml:"EvaluationExtensionDate"`
      ExtensionComments string `xml:"ExtensionComments"`
      ReevaluationDate string `xml:"ReevaluationDate"`
      ProgramEligibilityDate string `xml:"ProgramEligibilityDate"`
      ProgramPlanDate string `xml:"ProgramPlanDate"`
      ProgramPlanEffectiveDate string `xml:"ProgramPlanEffectiveDate"`
      NOREPDate string `xml:"NOREPDate"`
      PlacementParentalConsentDate string `xml:"PlacementParentalConsentDate"`
      ProgramPlacementDate string `xml:"ProgramPlacementDate"`
      ExtendedSchoolYear string `xml:"ExtendedSchoolYear"`
      ExtendedDay string `xml:"ExtendedDay"`
      ProgramAvailability ProgramAvailabilityType `xml:"ProgramAvailability"`
      EntryPerson string `xml:"EntryPerson"`
      StudentSpecialEducationFTE string `xml:"StudentSpecialEducationFTE"`
      ParticipationContact string `xml:"ParticipationContact"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type StudentParticipation_ManagingSchool struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
