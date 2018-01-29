package sifxml


    type ResourceUsage struct {
        RefId RefIdType `xml:"RefId,attr"`
      SchoolInfoRefId IdRefType `xml:"SchoolInfoRefId"`
      ResourceUsageContentType ResourceUsage_ResourceUsageContentType `xml:"ResourceUsageContentType"`
      ResourceReportColumnList ResourceUsage_ResourceReportColumnList `xml:"ResourceReportColumnList"`
      ResourceReportLineList ResourceUsage_ResourceReportLineList `xml:"ResourceReportLineList"`
      SIF_Metadata SIF_MetadataType `xml:"SIF_Metadata"`
      SIF_ExtendedElements SIF_ExtendedElementsType `xml:"SIF_ExtendedElements"`
      
      }
    type ResourceUsage_ResourceUsageContentType struct {
       Code string `xml:"Code"`
       LocalDescription string `xml:"LocalDescription"`
}
type ResourceUsage_ResourceReportColumnList struct {
      ResourceReportColumn []ResourceUsage_ResourceReportColumn `xml:"ResourceReportColumn"`
}
type ResourceUsage_ResourceReportLineList struct {
      ResourceReportLine []ResourceUsage_ResourceReportLine `xml:"ResourceReportLine"`
}
type ResourceUsage_ResourceReportColumn struct {
       ColumnName string `xml:"ColumnName"`
       ColumnDescription string `xml:"ColumnDescription"`
       ColumnDelimiter string `xml:"ColumnDelimiter"`
}
type ResourceUsage_ResourceReportLine struct {
      SIF_RefId ResourceUsage_SIF_RefId `xml:"SIF_RefId"`
       StartDate string `xml:"StartDate"`
       EndDate string `xml:"EndDate"`
       CurrentCost MonetaryAmountType `xml:"CurrentCost"`
       ReportRow string `xml:"ReportRow"`
}
type ResourceUsage_SIF_RefId struct {
      SIF_RefObject string `xml:"SIF_RefObject,attr"`
      Value string `xml:",chardata"`
}
