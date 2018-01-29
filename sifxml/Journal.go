package sifxml


    type Journal struct {
        RefId RefIdType `xml:"RefId,attr"`
      DebitFinancialAccountRefId IdRefType `xml:"DebitFinancialAccountRefId"`
      CreditFinancialAccountRefId IdRefType `xml:"CreditFinancialAccountRefId"`
      OriginatingTransactionRefId Journal_OriginatingTransactionRefId `xml:"OriginatingTransactionRefId"`
      Amount MonetaryAmountType `xml:"Amount"`
      GSTCodeOriginal string `xml:"GSTCodeOriginal"`
      GSTCodeReplacement string `xml:"GSTCodeReplacement"`
      Note string `xml:"Note"`
      CreatedDate string `xml:"CreatedDate"`
      ApprovedDate string `xml:"ApprovedDate"`
      CreatedBy string `xml:"CreatedBy"`
      ApprovedBy string `xml:"ApprovedBy"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type Journal_OriginatingTransactionRefId struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value IdRefType `xml:",chardata"`
}
