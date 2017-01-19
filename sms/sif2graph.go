package sms

import (
	//"encoding/json"
	"errors"
	"github.com/beevik/etree"
	"strings"
)

// implementation of the sif2graph service
type Sif2GraphService struct {
	Paths map[string]etree.Path // map of paths to search
}

// create a new sif2graph service instance
func NewSif2GraphService() (*Sif2GraphService, error) {
	s2g := Sif2GraphService{}
	s2g.Paths = make(map[string]etree.Path)
	s2g.Paths["LocalId"] = etree.MustCompilePath("//LocalId")
	s2g.Paths["StateProvinceId"] = etree.MustCompilePath("//StateProvinceId")
	s2g.Paths["ACARAId"] = etree.MustCompilePath("//ACARAId")
	s2g.Paths["ElectronicId"] = etree.MustCompilePath("//ElectronicIdList/ElectronicId")
	s2g.Paths["OtherId"] = etree.MustCompilePath("//OtherIdList/OtherId")
	return &s2g, nil
}

// extract human readable label based on SIF object.
func extract_label(doc *etree.Document, root *etree.Element) (string, error) {
	var ret string
	switch root.Tag {
	case "Activity":
		ret = doc.FindElement("//Title").Text()
	case "AggregateStatisticInfo":
		ret = doc.FindElement("//StatisticName").Text()
	case "Assessment", "Sif3Assessment":
		ret = doc.FindElement("//Name").Text()
	case "AssessmentAdministration", "Sif3AssessmentAdministration":
		ret = doc.FindElement("//AdministrationName").Text()
	case "Sif3AssessmentAsset":
		ret = doc.FindElement("//AssetName").Text()
	case "AssessmentForm", "Sif3AssessmentForm":
		ret = doc.FindElement("//FormName").Text()
	case "AssessmentItem", "Sif3AssessmentItem":
		ret = doc.FindElement("//ItemLabel").Text()
	case "Sif3AssessmentRubric":
		ret = doc.FindElement("//RubricName").Text()
	case "Sif3AssessmentScoreTable":
		ret = doc.FindElement("//ScoreTableName").Text()
	case "Sif3AssessmentSection":
		ret = doc.FindElement("//SectionName").Text()
	case "Sif3AssessmentSession":
		ret = doc.FindElement("//SessionName").Text()
	case "AssessmentSubTest":
		ret = doc.FindElement("//Name").Text()
	case "Sif3AssessmentSubTest":
		ret = doc.FindElement("//SubTestName").Text()
	case "CalendarDate":
		ret = root.SelectAttrValue("@Date", "")
	case "CalendarSummary":
		ret = doc.FindElement("//LocalId").Text()
	case "ChargedLocationInfo":
		ret = doc.FindElement("//Name").Text()
	case "Debtor":
		bname := doc.FindElement("//BillingName")
		if bname != nil {
			ret = bname.Text()
		}
	case "EquipmentInfo":
		ret = doc.FindElement("//Name").Text()
	case "FinancialAccount":
		ret = doc.FindElement("//AccountNumber").Text()
	case "GradingAssignment":
		ret = doc.FindElement("//Description").Text()
	case "Invoice":
		fnumber := doc.FindElement("//FormNumber")
		if fnumber != nil {
			ret = fnumber.Text()
		}
	case "LEAInfo":
		ret = doc.FindElement("//LEAName").Text()
	case "LearningResource":
		ret = doc.FindElement("//Name").Text()
	case "LearningStandardDocument":
		ret = doc.FindElement("//Title").Text()
	case "LearningStandardItem":
		ret = doc.FindElement("//StatementCodes/StatementCode[1]").Text()
	case "PaymentReceipt":
		ret = doc.FindElement("//ReceivedTransactionId").Text()
	case "PurchaseOrder":
		ret = doc.FindElement("//FormNumber").Text()
	case "ReportAuthorityInfo":
		ret = doc.FindElement("//AuthorityName").Text()
	case "ResourceBooking":
		ret = doc.FindElement("//ResourceLocalId").Text()
	case "RoomInfo":
		ret = doc.FindElement("//RoomNumber").Text()
	case "ScheduledActivity":
		ret = doc.FindElement("//ActivityName").Text()
	case "SchoolCourse":
		ret = doc.FindElement("//CourseCode").Text()
	case "SchoolInfo":
		ret = doc.FindElement("//SchoolName").Text()
	case "SectionInfo":
		ret = doc.FindElement("//LocalId").Text()
	case "StaffPersonal", "StudentContactPersonal", "StudentPersonal":
		fname := doc.FindElement("//PersonInfo/Name/FullName")
		if fname == nil {
			ret1 := doc.FindElement("//PersonInfo/Name/GivenName").Text()
			ret2 := doc.FindElement("//PersonInfo/Name/FamilyName").Text()
			ret = ret1 + " " + ret2
		} else {
			ret = fname.Text()
		}
	case "StudentActivityInfo":
		ret = doc.FindElement("//Title").Text()
	case "TeachingGroup":
		ret = doc.FindElement("//ShortName").Text()
	case "TermInfo":
		ret = doc.FindElement("//TermCode").Text()
	case "TimeTable":
		ret = doc.FindElement("//Title").Text()
	case "TimeTableCell":
		ret = doc.FindElement("//DayId").Text() + ":" + doc.FindElement("//PeriodId").Text()
	case "TimeTableSubject":
		ret = doc.FindElement("//SubjectLocalId").Text()
	case "VendorInfo":
		fname := doc.FindElement("//Name/FullName")
		if fname == nil {
			ret1 := doc.FindElement("//Name/GivenName").Text()
			ret2 := doc.FindElement("//Name/FamilyName").Text()
			ret = ret1 + " " + ret2
		} else {
			ret = fname.Text()
		}

	default:
		ret = root.SelectAttrValue("RefId", "")
	}
	return ret, nil
}

func parseSIF(s2g *Sif2GraphService, xml string) (GraphStruct, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return GraphStruct{}, err
	}
	root := doc.Root()
	if root == nil {
		return GraphStruct{}, errors.New("XML has no root")
	}
	label, err := extract_label(doc, root)
	if err != nil {
		return GraphStruct{}, err
	}
	guid := root.SelectAttrValue("RefId", "")
	ret := GraphStruct{
		Guid:          guid,
		EquivalentIds: make([]string, 0),
		OtherIds:      make(map[string]string),
		Type:          root.Tag,
		Links:         make([]string, 0),
		Label:         label,
	}
	ret.Links = extractLinks(ret.Links, root, guid)
	ids := doc.FindElementsPath(s2g.Paths["LocalId"])
	for _, elem := range ids {
		ret.OtherIds["LocalId"] = elem.Text()
	}
	ids = doc.FindElementsPath(s2g.Paths["StateProvinceId"])
	for _, elem := range ids {
		ret.OtherIds["StateProvinceId"] = elem.Text()
	}
	ids = doc.FindElementsPath(s2g.Paths["ACARAId"])
	for _, elem := range ids {
		ret.OtherIds["ACARAId"] = elem.Text()
	}
	ids = doc.FindElementsPath(s2g.Paths["ElectronicId"])
	for _, elem := range ids {
		ret.OtherIds["ElectronicId:"+elem.SelectAttrValue("Type", "")] = elem.Text()
	}
	ids = doc.FindElementsPath(s2g.Paths["OtherId"])
	for _, elem := range ids {
		ret.OtherIds[elem.SelectAttrValue("Type", "")] = elem.Text()
	}

	return ret, nil
}

// recursively extract links from xml: elements suffixed with RefId, attributes suffixed with RefId,
// and SIF_RefObject attributes. Do not append a RefId if it equals guid
func extractLinks(links []string, root *etree.Element, guid string) []string {
	for _, attr := range root.Attr {
		if strings.HasSuffix(attr.Key, "RefId") {
			links = conditional_append(links, attr.Value, guid)
		}
		if attr.Key == "SIF_RefObject" {
			links = conditional_append(links, root.Text(), guid)
		}
	}
	if strings.HasSuffix(root.Tag, "RefId") {
		links = conditional_append(links, root.Text(), guid)
	}
	children := root.ChildElements()
	for _, elem := range children {
		links = extractLinks(links, elem, guid)
	}
	return links
}

// append
func conditional_append(slice []string, element string, skip string) []string {
	if element != skip {
		slice = append(slice, element)
	}
	return slice
}

// implement the nias Service interface
func (s2g *Sif2GraphService) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)
	out, err := parseSIF(s2g, req.Body.(string))
	if err != nil {
		return nil, err
	}
	r := NiasMessage{}
	r.TxID = req.TxID
	r.SeqNo = req.SeqNo
	r.Target = SIF_MEMORY_STORE_PREFIX
	r.Body = out
	// having to JSON marshal this: the struct is blocking NATS2 from proccessing it
	//r.Body, _ = json.Marshal(out)
	responses = append(responses, r)
	return responses, nil
}
