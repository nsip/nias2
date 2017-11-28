package naprrql

import (
	"encoding/gob"
	"log"
	"strconv"
	"strings"

	"github.com/nsip/nias2/xml"
)

func init() {
	gob.Register(ParticipationDataSet{})
	gob.Register(EventInfo{})
	gob.Register(SchoolDetails{})
	gob.Register(ScoreSummaryDataSet{})
	gob.Register(ResponseDataSet{})
	gob.Register(EventResponseDataSet{})
	gob.Register(CodeFrameDataSet{})
	gob.Register(ResultsByStudent{})
	gob.Register(GuidCheckDataSet{})
	gob.Register(CodeframeCheckDataSet{})
}

// convenience types for aggregating response information sets
// used in reporting and support types for sorting results.
//

// Codeframe validity check

type CodeframeCheckDataSet struct {
	ObjectID   string
	LocalID    string
	ObjectType string
	Test       xml.NAPTest
	Testlet    xml.NAPTestlet
	TestItem   xml.NAPTestItem
}

// struct for sorting support
type CodeframeCheckComparator []GuidCheckDataSet

// sort interface implementation for itemresponsedatasets
func (resps CodeframeCheckComparator) Len() int {
	return len(resps)
}

// sort interface implementation for responsedatasets
func (resps CodeframeCheckComparator) Swap(i, j int) {
	resps[i], resps[j] = resps[j], resps[i]
}

// sort interface implementation for responsedatasets
func (resps CodeframeCheckComparator) Less(i, j int) bool {
	return resps[i].Guid < resps[j].Guid
}

// GUID validity check
type GuidCheckDataSet struct {
	ObjectName    string
	ObjectType    string
	Guid          string
	ShouldPointTo string
	PointsTo      string
}

// struct for sorting support
type GuidCheckComparator []GuidCheckDataSet

// sort interface implementation for itemresponsedatasets
func (resps GuidCheckComparator) Len() int {
	return len(resps)
}

// sort interface implementation for responsedatasets
func (resps GuidCheckComparator) Swap(i, j int) {
	resps[i], resps[j] = resps[j], resps[i]
}

// sort interface implementation for responsedatasets
func (resps GuidCheckComparator) Less(i, j int) bool {
	return resps[i].Guid < resps[j].Guid
}

// aggregating type used for reporting item responses against items
type ItemResponseDataSet struct {
	Test              xml.NAPTest
	Testlet           xml.NAPTestlet
	TestItem          xml.NAPTestItem
	Student           xml.RegistrationRecord
	Response          xml.NAPResponseSet
	ParticipationCode string
	SchoolDetails     SchoolDetails
}

// struct for sorting support
type ItemResponseComparator []ItemResponseDataSet

// sort interface implementation for itemresponsedatasets
func (resps ItemResponseComparator) Len() int {
	return len(resps)
}

// sort interface implementation for responsedatasets
func (resps ItemResponseComparator) Swap(i, j int) {
	resps[i], resps[j] = resps[j], resps[i]
}

// sort interface implementation for responsedatasets
func (resps ItemResponseComparator) Less(i, j int) bool {
	return resps[i].TestItem.TestItemContent.ItemName < resps[j].TestItem.TestItemContent.ItemName
}

// aggregating type used for reporting domain scores and events
type EventResponseDataSet struct {
	Event         xml.NAPEvent
	Test          xml.NAPTest
	Student       xml.RegistrationRecord
	Response      xml.NAPResponseSet
	SchoolDetails SchoolDetails
}

// struct for sorting support
type EventResponseComparator []EventResponseDataSet

// sort interface implementation for responsedatasets
func (resps EventResponseComparator) Len() int {
	return len(resps)
}

// sort interface implementation for responsedatasets
func (resps EventResponseComparator) Swap(i, j int) {
	resps[i], resps[j] = resps[j], resps[i]
}

// sort interface implementation for responsedatasets
func (resps EventResponseComparator) Less(i, j int) bool {
	return resps[i].Test.TestContent.TestName < resps[j].Test.TestContent.TestName
}

// aggregating type used for reporting domain scores
type ResponseDataSet struct {
	Test     xml.NAPTest
	Student  xml.RegistrationRecord
	Response xml.NAPResponseSet
}

// struct for sorting support
type ResponseComparator []ResponseDataSet

// sort interface implementation for responsedatasets
func (resps ResponseComparator) Len() int {
	return len(resps)
}

// sort interface implementation for responsedatasets
func (resps ResponseComparator) Swap(i, j int) {
	resps[i], resps[j] = resps[j], resps[i]
}

// sort interface implementation for responsedatasets
func (resps ResponseComparator) Less(i, j int) bool {
	return resps[i].Test.TestContent.TestName < resps[j].Test.TestContent.TestName
}

// aggregating type used for producing school summary reports
type ScoreSummaryDataSet struct {
	Test xml.NAPTest
	Summ xml.NAPTestScoreSummary
}

// struct for sorting support
type ScoreSummaryComparator []ScoreSummaryDataSet

// sort interface implementation for summarydatasets
func (summs ScoreSummaryComparator) Len() int {
	return len(summs)
}

// sort interface implementation for summarydatasets
func (summs ScoreSummaryComparator) Swap(i, j int) {
	summs[i], summs[j] = summs[j], summs[i]
}

// sort interface implementation for summarydatasets
func (summs ScoreSummaryComparator) Less(i, j int) bool {
	return summs[i].Test.TestContent.TestName < summs[j].Test.TestContent.TestName
}

// helper type for summary
type ParticipationSummary struct {
	Domain            string
	ParticipationCode string
}

// reporting object for student participation
type ParticipationDataSet struct {
	Student    xml.RegistrationRecord
	School     xml.SchoolInfo
	EventInfos []EventInfo
	Summary    []ParticipationSummary
}

// helper type for test/event info
type EventInfo struct {
	Event xml.NAPEvent
	Test  xml.NAPTest
}

// make school id and name a type for transmission
type SchoolDetails struct {
	ACARAId    string
	SchoolName string
}

// summary object from codeframe
type CodeFrameDataSet struct {
	Test           xml.NAPTest
	Testlet        xml.NAPTestlet
	Item           xml.NAPTestItem
	SequenceNumber string
}

// aggregate all objects referencing students for a single test
type ResultsByStudent struct {
	Student     *xml.RegistrationRecord
	Event       *xml.NAPEvent
	ResponseSet *xml.NAPResponseSet
}

//
// helper method that walks the structure to find location
//
func (cfd CodeFrameDataSet) GetItemLocationInTestlet(itemrefid string) string {

	// check if the sequence no. is known
	for _, item := range cfd.Testlet.TestItemList.TestItem {
		if item.TestItemRefId == itemrefid {
			return item.SequenceNumber
		}
	}

	// if not see if the item is an alternative
	// get the alternative list of items
	// see if they have a sequence number in the testlet
	for _, altItem := range cfd.Item.TestItemContent.ItemSubstitutedForList.SubstituteItem {
		for _, item := range cfd.Testlet.TestItemList.TestItem {
			if item.TestItemRefId == altItem.SubstituteItemRefId {
				return item.SequenceNumber
			}
		}
	}

	return "unknown"
}

//
// helpers for deeply nested writing rubric details
//
func (cfd CodeFrameDataSet) GetWritingRubricDescriptor(rubrictype string) string {

	for _, rubric := range cfd.Item.TestItemContent.NAPWritingRubricList.NAPWritingRubric {
		if strings.EqualFold(rubric.RubricType, rubrictype) {
			return rubric.Descriptor
		}
	}

	return "unknown"
}

//
// helpers for deeply nested writing rubric details
//
func (cfd CodeFrameDataSet) GetWritingRubricMax(rubrictype string) string {

	for _, rubric := range cfd.Item.TestItemContent.NAPWritingRubricList.NAPWritingRubric {
		if strings.EqualFold(rubric.RubricType, rubrictype) {
			var max_score int
			for _, score := range rubric.ScoreList.Score {
				score_int, err := strconv.Atoi(score.MaxScoreValue)
				if err != nil {
					log.Println("Score for ", rubrictype, " cannot be converted to int: ", err)
				}
				max_score += score_int
			}
			return strconv.Itoa(max_score)
		}
	}

	return "unknown"

}

//
//
//
//
//
