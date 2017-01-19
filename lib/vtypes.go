// common types used in validation processes
// shared types for passing validation errors
// between services

// All types to be encoded directly onto NATS as message bodies must be defined here, and registered
// with the gob encode in the init() routine

package nias2

import (
	"encoding/gob"
	"encoding/xml"
	"github.com/nsip/nias2/go_SifMessage"
)

func init() {
	// make gob encoder aware of local types
	gob.Register(ValidationError{})
	gob.Register(RegistrationRecord{})
	gob.Register(GraphStruct{})
	gob.Register(SifStudentPersonal{})
}

type RegistrationRecord struct {
	// XML Configuration
	// XMLName            xml.Name `xml:"StudentPersonal"`
	// Important fields
	ASLSchoolId      string `json:",omitempty" xml:"MostRecent>SchoolACARAId"`
	AddressLine1     string `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line1"`
	AddressLine2     string `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line2"`
	BirthDate        string `json:",omitempty" xml:"PersonInfo>Demographics>BirthDate"`
	ClassGroup       string `json:",omitempty" xml:"MostRecent>ClassCode"`
	CountryOfBirth   string `json:",omitempty" xml:"PersonInfo>Demographics>CountryOfBirth"`
	DiocesanId       string `json:",omitempty"`
	EducationSupport string `json:",omitempty" xml:"EducationSupport"`
	FFPOS            string `json:",omitempty" xml:"MostRecent>FFPOS"`
	//FTE                       float32 `json:",omitempty" xml:"MostRecent>FTE"`
	FTE                       string `json:",omitempty" xml:"MostRecent>FTE"`
	FamilyName                string `json:",omitempty" xml:"PersonInfo>Name>FamilyName"`
	GivenName                 string `json:",omitempty" xml:"PersonInfo>Name>GivenName"`
	HomeSchooledStudent       string `json:",omitempty" xml:"HomeSchooledStudent"`
	IndigenousStatus          string `json:",omitempty" xml:"PersonInfo>Demographics>IndigenousStatus"`
	JurisdictionId            string `json:",omitempty"`
	LBOTE                     string `json:",omitempty" xml:"PersonInfo>Demographics>LBOTE"`
	LocalCampusId             string `json:",omitempty" xml:"MostRecent>LocalCampusId"`
	LocalId                   string `json:",omitempty" xml:"LocalId"`
	Locality                  string `json:",omitempty" xml:"PersonInfo>AddressList>Address>City"`
	MainSchoolFlag            string `json:",omitempty" xml:"MostRecent>MembershipType"`
	MiddleName                string `json:",omitempty" xml:"PersonInfo>Name>MiddleName"`
	NationalId                string `json:",omitempty"`
	OfflineDelivery           string `json:",omitempty" xml:"OfflineDelivery"`
	OtherId                   string `json:",omitempty"`
	OtherSchoolId             string `json:",omitempty" xml:"MostRecent>OtherEnrollmentSchoolACARAId"`
	Parent1LOTE               string `json:",omitempty" xml:"MostRecent>Parent1Language"`
	Parent1NonSchoolEducation string `json:",omitempty" xml:"MostRecent>Parent1NonSchoolEducation"`
	Parent1Occupation         string `json:",omitempty" xml:"MostRecent>Parent1EmploymentType"`
	Parent1SchoolEducation    string `json:",omitempty" xml:"MostRecent>Parent1SchoolEducationLevel"`
	Parent2LOTE               string `json:",omitempty" xml:"MostRecent>Parent2Language"`
	Parent2NonSchoolEducation string `json:",omitempty" xml:"MostRecent>Parent2NonSchoolEducation"`
	Parent2Occupation         string `json:",omitempty" xml:"MostRecent>Parent2EmploymentType"`
	Parent2SchoolEducation    string `json:",omitempty" xml:"MostRecent>Parent2SchoolEducationLevel"`
	PlatformId                string `json:",omitempty"`
	Postcode                  string `json:",omitempty" xml:"PersonInfo>AddressList>Address>PostalCode`
	PreferredName             string `json:",omitempty" xml:"PersonInfo>Name>PreferredGivenName"`
	PreviousDiocesanId        string `json:",omitempty"`
	PreviousJurisdictionId    string `json:",omitempty"`
	PreviousLocalId           string `json:",omitempty"`
	PreviousNationalId        string `json:",omitempty"`
	PreviousOtherId           string `json:",omitempty"`
	PreviousPlatformId        string `json:",omitempty"`
	PreviousSectorId          string `json:",omitempty"`
	PreviousTAAId             string `json:",omitempty"`
	ReportingSchoolId         string `json:",omitempty" xml:"MostRecent>ReportingSchoolId"`
	SchoolLocalId             string `json:",omitempty" xml:"MostRecent>SchoolLocalId"`
	SectorId                  string `json:",omitempty"`
	Sensitive                 string `json:",omitempty" xml:"Sensitive"`
	Sex                       string `json:",omitempty" xml:"PersonInfo>Demographics>Sex"`
	StateTerritory            string `json:",omitempty" xml:"PersonInfo>AddressList>Address>StateProvince"`
	StudentLOTE               string `json:",omitempty" xml:"PersonInfo>Demographics>LanguageList>Language>Code"`
	TAAId                     string `json:",omitempty"`
	TestLevel                 string `json:",omitempty" xml:"MostRecent>TestLevel>Code"`
	VisaCode                  string `json:",omitempty" xml:"PersonInfo>Demographics>VisaSubClass"`
	YearLevel                 string `json:",omitempty" xml:"MostRecent>YearLevel>Code"`
}

// struct to handle reporting of validation errors found in
// naplan registration files
type ValidationError struct {
	Field        string `json:"errField"`     // the field that has an error
	Description  string `json:"description"`  // error description
	OriginalLine string `json:"originalLine"` // input file record line that has the error
	Vtype        string `json:"validationType"`
}

// helper method for writing out csv encoding of error reports
func (ve *ValidationError) ToSlice() []string {

	return []string{ve.OriginalLine, ve.Vtype, ve.Field, ve.Description}
}

// information extracted out of SIF for graph
type GraphStruct struct {
	Guid          string            // RefID of object
	EquivalentIds []string          // equivalent Ids
	OtherIds      map[string]string // map of OtherId type to OtherId
	Type          string            // object type
	Links         []string          // list of related ids
	Label         string            // human readable label
}

// from metaleap/go-xsd conversion
type SifStudentPersonal struct {
	XMLName xml.Name `xml:"StudentPersonal"`
	go_SifMessage.TStudentPersonalType
}
