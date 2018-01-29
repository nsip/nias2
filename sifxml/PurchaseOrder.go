package sifxml


    type PurchaseOrder struct {
        RefId RefIdType `xml:"RefId,attr"`
      FormNumber string `xml:"FormNumber"`
      VendorInfoRefId IdRefType `xml:"VendorInfoRefId"`
      ChargedLocationInfoRefId IdRefType `xml:"ChargedLocationInfoRefId"`
      EmployeePersonalRefId IdRefType `xml:"EmployeePersonalRefId"`
      PurchasingItems PurchasingItemsType `xml:"PurchasingItems"`
      CreationDate string `xml:"CreationDate"`
      TaxRate string `xml:"TaxRate"`
      TaxAmount MonetaryAmountType `xml:"TaxAmount"`
      TotalAmount MonetaryAmountType `xml:"TotalAmount"`
      UpdateDate string `xml:"UpdateDate"`
      FullyDelivered string `xml:"FullyDelivered"`
      OriginalPurchaseOrderRefId IdRefType `xml:"OriginalPurchaseOrderRefId"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    