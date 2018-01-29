package sifxml


    type PaymentReceipt struct {
        RefId RefIdType `xml:"RefId,attr"`
      TransactionType string `xml:"TransactionType"`
      InvoiceRefId IdRefType `xml:"InvoiceRefId"`
      VendorInfoRefId IdRefType `xml:"VendorInfoRefId"`
      DebtorRefId IdRefType `xml:"DebtorRefId"`
      ChargedLocationInfoRefId IdRefType `xml:"ChargedLocationInfoRefId"`
      TransactionDate string `xml:"TransactionDate"`
      TransactionAmount DebitOrCreditAmountType `xml:"TransactionAmount"`
      ReceivedTransactionId string `xml:"ReceivedTransactionId"`
      FinancialAccountRefIdList FinancialAccountRefIdListType `xml:"FinancialAccountRefIdList"`
      TransactionDescription string `xml:"TransactionDescription"`
      TaxRate string `xml:"TaxRate"`
      TaxAmount MonetaryAmountType `xml:"TaxAmount"`
      TransactionMethod string `xml:"TransactionMethod"`
      ChequeNumber string `xml:"ChequeNumber"`
      TransactionNote string `xml:"TransactionNote"`
      AccountingPeriod LocalIdType `xml:"AccountingPeriod"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    