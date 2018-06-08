package xml

type NAPTestlet struct {
	TestletID      string         `xml:"RefId,attr"`
	NAPTestRefId   string         `xml:"NAPTestRefId,omitempty"`
	NAPTestLocalId string         `xml:"NAPTestLocalId"`
	TestletContent TestletContent `xml:"TestletContent"`
	TestItemList   TestItemList   `xml:"TestItemList"`
}

type TestletContent struct {
	LocalId             string `xml:"NAPTestletLocalId"`
	TestletName         string `xml:"TestletName"`
	Node                string `xml:"Node,omitempty"`
	LocationInStage     string `xml:"LocationInStage,omitempty"`
	TestletMaximumScore string `xml:"TestletMaximumScore"`
}

type TestItemList struct {
	TestItem []NAPTestlet_TestItem `xml:"TestItem"`
}

type NAPTestlet_TestItem struct {
	TestItemRefId   string `xml:"TestItemRefId"`
	TestItemLocalId string `xml:"TestItemLocalId"`
	SequenceNumber  string `xml:"SequenceNumber"`
}

func (t NAPTestlet) GetHeaders() []string {
	return []string{
		"TestletID", "NAPTestRefId", "LocalId", "NAPTestLocalId", "TestletName", "Node",
		"LocationInStage", "TestletMaximumScore"}
}

func (t NAPTestlet) GetSlice() []string {
	return []string{
		t.TestletID, t.NAPTestRefId, t.NAPTestLocalId, t.TestletContent.LocalId,
		t.TestletContent.TestletName, t.TestletContent.Node, t.TestletContent.LocationInStage,
		t.TestletContent.TestletMaximumScore}
}
