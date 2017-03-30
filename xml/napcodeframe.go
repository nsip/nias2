package xml

// codeframe object: we only record the GUIDS
type NAPCodeFrame struct {
	RefId        string `xml:"RefId,attr"`
	NAPTestRefId string `xml:"NAPTestRefId"`
	TestletList  struct {
		Testlet []NAPCodeFrame_Testlet `xml:"Testlet"`
	} `xml:"TestletList"`
}

type NAPCodeFrame_Testlet struct {
	NAPTestletRefId string `xml:"NAPTestletRefId"`
	TestItemList    struct {
		TestItem []NAPCodeFrame_TestItem `xml:"TestItem"`
	} `xml:"TestItemList"`
}

type NAPCodeFrame_TestItem struct {
	TestItemRefId string `xml:"TestItemRefId"`
}
