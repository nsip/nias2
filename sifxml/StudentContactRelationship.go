package sifxml


    type StudentContactRelationship struct {
        StudentContactRelationshipRefId IdRefType `xml:"StudentContactRelationshipRefId,attr"`
      StudentPersonalRefId RefIdType `xml:"StudentPersonalRefId"`
      StudentContactPersonalRefId RefIdType `xml:"StudentContactPersonalRefId"`
      Relationship RelationshipType `xml:"Relationship"`
      ParentRelationshipStatus string `xml:"ParentRelationshipStatus"`
      HouseholdList HouseholdList `xml:"HouseholdList"`
      ContactFlags ContactFlagsType `xml:"ContactFlags"`
      MainlySpeaksEnglishAtHome string `xml:"MainlySpeaksEnglishAtHome"`
      ContactSequence string `xml:"ContactSequence"`
      ContactSequenceSource string `xml:"ContactSequenceSource"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    