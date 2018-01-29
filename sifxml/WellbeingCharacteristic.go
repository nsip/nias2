package sifxml


    type WellbeingCharacteristic struct {
        RefId RefIdType `xml:"RefId,attr"`
      StudentPersonalRefId IdRefType `xml:"StudentPersonalRefId"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      WellbeingCharacteristicClassification string `xml:"WellbeingCharacteristicClassification"`
      WellbeingCharacteristicStartDate string `xml:"WellbeingCharacteristicStartDate"`
      WellbeingCharacteristicEndDate string `xml:"WellbeingCharacteristicEndDate"`
      WellbeingCharacteristicReviewDate string `xml:"WellbeingCharacteristicReviewDate"`
      WellbeingCharacteristicNotes string `xml:"WellbeingCharacteristicNotes"`
      WellbeingCharacteristicCategory string `xml:"WellbeingCharacteristicCategory"`
      WellbeingCharacteristicSubCategory string `xml:"WellbeingCharacteristicSubCategory"`
      LocalCharacteristicCode LocalId `xml:"LocalCharacteristicCode"`
      CharacteristicSeverity string `xml:"CharacteristicSeverity"`
      DailyManagement string `xml:"DailyManagement"`
      EmergencyManagement string `xml:"EmergencyManagement"`
      EmergencyResponsePlan string `xml:"EmergencyResponsePlan"`
      Trigger string `xml:"Trigger"`
      ConfidentialFlag string `xml:"ConfidentialFlag"`
      Alert string `xml:"Alert"`
      MedicationList MedicationListType `xml:"MedicationList"`
      DocumentList WellbeingDocumentListType `xml:"DocumentList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    