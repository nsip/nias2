package xml

// codeframe object: we only record the GUIDS
type NAPCodeFrame struct {
	RefId        string `xml:"RefId,attr"`
	NAPTestRefId string `xml:"NAPTestRefId"`
	TestletList  struct {
		Testlet []struct {
			NAPTestletRefId string `xml:"NAPTestletRefId"`
			TestItemList    struct {
				TestItem []struct {
					TestItemRefId string `xml:"TestItemRefId"`
				} `xml:"TestItem"`
			} `xml:"TestItemList"`
		} `xml:"Testlet"`
	} `xml:"TestletList"`
}
