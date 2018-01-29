package sifxml


    type AggregateStatisticInfo struct {
        RefId RefIdType `xml:"RefId,attr"`
      StatisticName string `xml:"StatisticName"`
      CalculationRule AggregateStatisticInfo_CalculationRule `xml:"CalculationRule"`
      ApprovalDate string `xml:"ApprovalDate"`
      ExpirationDate string `xml:"ExpirationDate"`
      ExclusionRules ExclusionRulesType `xml:"ExclusionRules"`
      Source string `xml:"Source"`
      EffectiveDate string `xml:"EffectiveDate"`
      DiscontinueDate string `xml:"DiscontinueDate"`
      Location LocationType `xml:"Location"`
      Measure string `xml:"Measure"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type AggregateStatisticInfo_CalculationRule struct {
      Type string `xml:"Type,attr"`
      Value string `xml:",chardata"`
}
