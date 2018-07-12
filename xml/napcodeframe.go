package xml

// codeframe object: we only record the GUIDS
type NAPCodeFrame struct {
	RefId        string                   `xml:"RefId,attr"`
	NAPTestRefId string                   `xml:"NAPTestRefId"`
	TestContent  TestContent              `xml:"TestContent"` // added for XML roundtrip
	TestletList  NAPCodeFrame_TestletList `xml:"TestletList"`
}

type NAPCodeFrame_TestItemList struct {
	TestItem []NAPCodeFrame_TestItem `xml:"TestItem"`
}
type NAPCodeFrame_TestletList struct {
	Testlet []NAPCodeFrame_Testlet `xml:"Testlet"`
}

type NAPCodeFrame_Testlet struct {
	NAPTestletRefId string                    `xml:"NAPTestletRefId"`
	TestletContent  TestletContent            `xml:"TestletContent"` // added for XML roundtrip
	TestItemList    NAPCodeFrame_TestItemList `xml:"TestItemList"`
}

type NAPCodeFrame_TestItem struct {
	TestItemRefId   string          `xml:"TestItemRefId"`
	SequenceNumber  string          `xml:"SequenceNumber"`
	TestItemContent TestItemContent `xml:"TestItemContent"` // added for XML roundtrip
}
