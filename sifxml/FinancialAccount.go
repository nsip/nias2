package sifxml


    type FinancialAccount struct {
        RefId RefIdType `xml:"RefId,attr"`
      ParentAccountRefId IdRefType `xml:"ParentAccountRefId"`
      ChargedLocationInfoRefId IdRefType `xml:"ChargedLocationInfoRefId"`
      AccountNumber string `xml:"AccountNumber"`
      Name string `xml:"Name"`
      Description string `xml:"Description"`
      ClassType string `xml:"ClassType"`
      CreationDate string `xml:"CreationDate"`
      CreationTime string `xml:"CreationTime"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    