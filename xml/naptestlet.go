package xml

type NAPTestlet struct {
	TestletID      string `xml:"RefId,attr"`
	NAPTestRefId   string `xml:"NAPTestRefId"`
	TestletContent struct {
		LocalId             string `xml:"NAPTestletLocalId"`
		NAPTestLocalId      string `xml:"NAPTestLocalId"`
		TestletName         string `xml:"TestletName"`
		Node                string `xml:"Node"`
		LocationInStage     string `xml:"LocationInStage"`
		TestletMaximumScore string `xml:"TestletMaximumScore"`
	} `xml:"TestletContent"`
	TestItemList struct {
		TestItem []struct {
			TestItemRefId   string `xml:"TestItemRefId"`
			TestItemLocalId string `xml:"TestItemLocalId"`
			SequenceNumber  string `xml:"SequenceNumber"`
		} `xml:"TestItem"`
	} `xml:"TestItemList"`
}

func (t NAPTestlet) GetHeaders() []string {
	return []string{
		"TestletID", "NAPTestRefId", "LocalId", "NAPTestLocalId", "TestletName", "Node",
		"LocationInStage", "TestletMaximumScore"}
}

func (t NAPTestlet) GetSlice() []string {
	return []string{
		t.TestletID, t.NAPTestRefId, t.TestletContent.LocalId, t.TestletContent.NAPTestLocalId,
		t.TestletContent.TestletName, t.TestletContent.Node, t.TestletContent.LocationInStage,
		t.TestletContent.TestletMaximumScore}
}
