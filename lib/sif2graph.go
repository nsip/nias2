package nias2

import (
	"github.com/beevik/etree"
	//"log"
	"errors"
	"fmt"
	"strings"
)

// information extracted out of SIF for graph
type GraphStruct struct {
	Guid          string            // RefID of object
	EquivalentIds []string          // equivalent Ids
	OtherIds      map[string]string // map of OtherId type to OtherId
	Type          string            // object type
	Links         []string          // list of related ids
	Label         string            // human readable label
}

// implementation of the psi service
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

func parseSIF(s2g *Sif2GraphService, xml string) (*GraphStruct, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return nil, err
	}
	root := doc.Root()
	if root == nil {
		return nil, errors.New("XML has no root")
	}
	label, err := extract_label(doc, root)
	if err != nil {
		return nil, err
	}
	ret := GraphStruct{
		Guid:          root.SelectAttrValue("RefId", ""),
		EquivalentIds: make([]string, 0),
		OtherIds:      make(map[string]string),
		Type:          root.Tag,
		Links:         make([]string, 0),
		Label:         label,
	}
	extractLinks(ret.Links, root)
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

	return &ret, nil
}

// recursively extract links from xml: elements suffixed with RefId, attributes suffixed with RefId,
// and SIF_RefObject attributes
func extractLinks(links []string, root *etree.Element) {
	for _, attr := range root.Attr {
		if strings.HasSuffix(attr.Key, "RefId") {
			links = append(links, attr.Value)
		}
		if attr.Key == "SIF_RefObject" {
			links = append(links, root.Text())
		}
	}
	if strings.HasSuffix(root.Tag, "RefId") {
		links = append(links, root.Text())
	}
	children := root.ChildElements()
	for _, elem := range children {
		extractLinks(links, elem)
	}
}

/*
func main() {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<StudentPersonal RefId="7C834EA9EDA12090347F83297E1C290C">
  <AlertMessages>
      <AlertMessage Type="Legal">Mother is legal guardian</AlertMessage>
        </AlertMessages>
  <MedicalAlertMessages>
      <MedicalAlertMessage Severity="Severe">Student has Peanut Allergy</MedicalAlertMessage>
          <MedicalAlertMessage Severity="Moderate">Student has Diabetes</MedicalAlertMessage>
    </MedicalAlertMessages>
      <LocalId>S1234567</LocalId>
        <StateProvinceId>ABC1234</StateProvinceId>
  <ElectronicIdList>
      <ElectronicId Type="03">ZZZZZZ21</ElectronicId>
          <ElectronicId Type="03">ZZZZZZ22</ElectronicId>
    </ElectronicIdList>
      <OtherIdList>
          <OtherId Type="PreviousNAPPlatformStudentId">888rdgf</OtherId>
      <OtherId Type="DiocesanStudentId">1234</OtherId>
        </OtherIdList>
  <PersonInfo>
      <Name Type="LGL">
            <FamilyName>Smith</FamilyName>
          <GivenName>Fred</GivenName>
        <FullName>Fred Smith</FullName>
    </Name>
        <OtherNames>
      <Name Type="AKA">
              <FamilyName>Anderson</FamilyName>
              <GivenName>Samuel</GivenName>
              <FullName>Samuel Anderson</FullName>
            </Name>
          <Name Type="PRF">
          <FamilyName>Rowinski</FamilyName>
          <GivenName>Sam</GivenName>
          <FullName>Sam Rowinski </FullName>
        </Name>
    </OtherNames>
        <Demographics>
      <IndigenousStatus>3</IndigenousStatus>
            <Sex>1</Sex>
          <BirthDate>1990-09-26</BirthDate>
        <BirthDateVerification>1004</BirthDateVerification>
      <PlaceOfBirth>Clayton</PlaceOfBirth>
            <StateOfBirth>VIC</StateOfBirth>
          <CountryOfBirth>1101</CountryOfBirth>
        <CountriesOfCitizenship>
        <CountryOfCitizenship>8104</CountryOfCitizenship>
        <CountryOfCitizenship>1101</CountryOfCitizenship>
      </CountriesOfCitizenship>
            <CountriesOfResidency>
            <CountryOfResidency>8104</CountryOfResidency>
            <CountryOfResidency>1101</CountryOfResidency>
          </CountriesOfResidency>
        <CountryArrivalDate>1990-09-26</CountryArrivalDate>
      <AustralianCitizenshipStatus>1</AustralianCitizenshipStatus>
            <EnglishProficiency>
            <Code>1</Code>
          </EnglishProficiency>
        <LanguageList>
        <Language>
          <Code>0001</Code>
            <LanguageType>1</LanguageType>
            </Language>
          </LanguageList>
        <DwellingArrangement>
        <Code>1671</Code>
      </DwellingArrangement>
            <Religion>
            <Code>2013</Code>
          </Religion>
        <ReligiousEventList>
        <ReligiousEvent>
          <Type>Baptism</Type>
            <Date>2000-09-01</Date>
            </ReligiousEvent>
            <ReligiousEvent>
              <Type>Christmas</Type>
                <Date>2009-12-24</Date>
        </ReligiousEvent>
      </ReligiousEventList>
            <ReligiousRegion>The Religion Region</ReligiousRegion>
          <PermanentResident>P</PermanentResident>
        <VisaSubClass>101</VisaSubClass>
      <VisaStatisticalCode>05</VisaStatisticalCode>
          </Demographics>
      <AddressList>
            <Address Type="0123" Role="2382">
            <Street>
              <Line1>Unit1/10</Line1>
                <Line2>Barkley Street</Line2>
        </Street>
        <City>Yarra Glenn</City>
        <StateProvince>VIC</StateProvince>
        <Country>1101</Country>
        <PostalCode>9999</PostalCode>
      </Address>
            <Address Type="0123A" Role="013A">
            <Street>
              <Line1>34 Term Address Street</Line1>
              </Street>
              <City>Home Town</City>
              <StateProvince>WA</StateProvince>
              <Country>1101</Country>
              <PostalCode>9999</PostalCode>
            </Address>
        </AddressList>
    <PhoneNumberList>
          <PhoneNumber Type="0096">
          <Number>03 9637-2289</Number>
          <Extension>72289</Extension>
          <ListedStatus>Y</ListedStatus>
        </PhoneNumber>
      <PhoneNumber Type="0888">
              <Number>0437-765-234</Number>
              <ListedStatus>N</ListedStatus>
            </PhoneNumber>
        </PhoneNumberList>
    <EmailList>
          <Email Type="01">fsmith@yahoo.com</Email>
        <Email Type="02">freddy@gmail.com</Email>
    </EmailList>
      </PersonInfo>
        <ProjectedGraduationYear>2014</ProjectedGraduationYear>
  <OnTimeGraduationYear>2012</OnTimeGraduationYear>
    <MostRecent>
        <SchoolLocalId>S1234567</SchoolLocalId>
    <HomeroomLocalId>hr12345</HomeroomLocalId>
        <YearLevel>
      <Code>P</Code>
          </YearLevel>
      <FTE>0.5</FTE>
          <Parent1Language>1201</Parent1Language>
      <Parent2Language>1201</Parent2Language>
          <LocalCampusId>D</LocalCampusId>
      <SchoolACARAId>VIC687</SchoolACARAId>
          <Homegroup>7A</Homegroup>
      <ClassCode>English 7D</ClassCode>
          <MembershipType>02</MembershipType>
      <FFPOS>2</FFPOS>
          <ReportingSchoolId>VIC670</ReportingSchoolId>
      <OtherEnrollmentSchoolACARAId>VIC6273</OtherEnrollmentSchoolACARAId>
        </MostRecent>
  <AcceptableUsePolicy>Y</AcceptableUsePolicy>
    <EconomicDisadvantage>N</EconomicDisadvantage>
      <ESL>U</ESL>
        <YoungCarersRole>N</YoungCarersRole>
  <Disability>N</Disability>
    <IntegrationAide>N</IntegrationAide>
      <EducationSupport>N</EducationSupport>
        <HomeSchooledStudent>N</HomeSchooledStudent>
  <Sensitive>N</Sensitive>
  </StudentPersonal>`
	s2g, _ := NewSif2GraphService()
	fmt.Println(parseSIF(s2g, xml))
}
*/

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
	r.Target = SIF_MEMORY_STORE_PREFIX + "::"
	r.Body = *out
	responses = append(responses, r)
	return responses, nil
}
