package xml

type SchoolInfo struct {
	RefId                    string                  `xml:"RefId,attr"`
	LocalId                  string                  `xml:"LocalId,omitempty"`
	StateProvinceId          string                  `xml:"StateProvinceId,omitempty"`
	ACARAId                  string                  `xml:"ACARAId,omitempty"`
	SchoolName               string                  `xml:"SchoolName"`
	LEAInfoRefId             string                  `xml:"LEAInfoRefId,omitempty"`
	OtherLEA                 string                  `xml:"OtherLEA,omitempty"`
	SchoolDistrict           string                  `xml:"SchoolDistrict,omitempty"`
	SchoolType               string                  `xml:"SchoolType,omitempty"`
	StudentCount             string                  // non xml field added for reporting
	SchoolFocusList          SchoolFocusList         `xml:"SchoolFocusList,omitempty"`
	SchoolURL                string                  `xml:"SchoolURL,omitempty"`
	PrincipalInfo            PrincipalInfo           `xml:"PrincipalInfo,omitempty"`
	SchoolContactList        SchoolContactList       `xml:"SchoolContactList,omitempty"`
	AddressList              AddressList             `xml:"AddressList,omitempty"`
	PhoneNumberList          PhoneNumberList         `xml:"PhoneNumberList,omitempty"`
	SessionType              string                  `xml:"SessionType,omitempty"`
	YearLevels               YearLevelList           `xml:"YearLevels,omitempty"`
	ARIA                     string                  `xml:"ARIA,omitempty"`
	OperationalStatus        string                  `xml:"OperationalStatus,omitempty"`
	FederalElectorate        string                  `xml:"FederalElectorate,omitempty"`
	Campus                   Campus                  `xml:"Campus,omitempty"`
	SchoolSector             string                  `xml:"SchoolSector"`
	IndependentSchool        string                  `xml:"IndependentSchool,omitempty"`
	NonGovSystemicStatus     string                  `xml:"NonGovSystemicStatus,omitempty"`
	System                   string                  `xml:"System,omitempty"`
	ReligiousAffiliation     string                  `xml:"ReligiousAffiliation,omitempty"`
	SchoolGeographicLocation string                  `xml:"SchoolGeographicLocation,omitempty"`
	LocalGovernmentArea      string                  `xml:"LocalGovernmentArea,omitempty"`
	JurisdictionLowerHouse   string                  `xml:"JurisdictionLowerHouse,omitempty"`
	SLA                      string                  `xml:"SLA,omitempty"`
	SchoolCoEdStatus         string                  `xml:"SchoolCoEdStatus,omitempty"`
	BoardingSchoolStatus     string                  `xml:"BoardingSchoolStatus,omitempty"`
	YearLevelEnrollmentList  YearLevelEnrollmentList `xml:"YearLevelEnrollmentList,omitempty"`
	TotalEnrollments         TotalEnrollments        `xml:"TotalEnrollments,omitempty"`
	Entity_Open              string                  `xml:"Entity_Open,omitempty"`
	Entity_Close             string                  `xml:"Entity_Close,omitempty"`
	SchoolGroupList          SchoolGroupList         `xml:"SchoolGroupList,omitempty"`
	SchoolTimeZone           string                  `xml:"SchoolTimeZone,omitempty"`
}

type YearLevelEnrollmentList struct {
	YearLevelEnrollment []YearLevelEnrollment `xml:"YearLevelEnrollment"`
}

type YearLevelEnrollment struct {
	Year       string `xml:"YearLevel>Code"`
	Enrollment string `xml:"Enrollment"`
}

type TotalEnrollments struct {
	Girls         string `xml:"Girls,omitempty"`
	Boys          string `xml:"Boys,omitempty"`
	TotalStudents string `xml:"TotalStudents,omitempty"`
}

type AddressList struct {
	Address []string `xml:"Address"`
}

type SchoolGroupList struct {
	SchoolGroup []string `xml:"SchoolGroup"`
}

type Campus struct {
	SchoolCampusId string `xml:"SchoolCampusId"`
	CampusType     string `xml:"CampusType"`
	AdminStatus    string `xml:"AdminStatus"`
}

type YearLevelList struct {
	YearLevel []YearLevel `xml:"YearLevel"`
}

type YearLevel struct {
	Code string `xml:"Code"`
}

type PhoneNumberList struct {
	PhoneNumber []PhoneNumber `xml:"PhoneNumber"`
}

type PhoneNumber struct {
	Type   string `xml:"Type,attr"`
	Number string `xml:"Number"`
}

type SchoolFocusList struct {
	SchoolFocus []string `xml:"SchoolFocus"`
}

type NameType struct {
	Type       string `xml:"Type,attr"`
	Title      string `xml:"Title"`
	FamilyName string `xml:"FamilyName"`
	GivenName  string `xml:"GivenName"`
	MiddleName string `xml:"MiddleName"`
	Suffix     string `xml:"Suffix"`
	FullName   string `xml:"FullName"`
}

type EmailList struct {
	Email Email
}

type Email struct {
	Type    string `xml:"Type,attr"`
	Address string `xml:"Email"`
}

type FullPhoneNumberList struct {
	PhoneNumber []FullPhoneNumber `xml:"PhoneNumber"`
}

type FullPhoneNumber struct {
	Type         string `xml:"Type,attr"`
	Number       string `xml:"Number"`
	Extension    string `xml:"Extension"`
	ListedStatus string `xml:"ListedStatus"`
}

type PrincipalInfo struct {
	ContactName  NameType `xml:"ContactName"`
	ContactTitle string   `xml:"ContactTitle"`
}

type SchoolContactList struct {
	SchoolContact []SchoolContact `xml:"SchoolContact"`
}

type SchoolContact struct {
	PublishInDirectory string      `xml:"PublishInDirectory"`
	ContactInfo        ContactInfo `xml:"ContactInfo"`
}

type ContactInfo struct {
	Name            NameType            `xml:"Name"`
	PositionTitle   string              `xml:"PositionTitle"`
	Role            string              `xml:"Role"`
	Address         Address             `xml:"Address"`
	EmailList       EmailList           `xml:"EmailList"`
	PhoneNumberList FullPhoneNumberList `xml:"PhoneNumberList"`
}

type Address struct {
	Type          string       `xml:"Type,attr"`
	Role          string       `xml:"Role,attr"`
	Street        Street       `xml:"Street"`
	City          string       `xml:"City"`
	StateProvince string       `xml:"StateProvince"`
	Country       string       `xml:"Country"`
	PostalCode    string       `xml:"PostalCode"`
	GridLocation  GridLocation `xml:"GridLocation"`
}

type Street struct {
	Line1 string `xml:"Line1"`
}

type GridLocation struct {
	Latitude  string `xml:"Latitude"`
	Longitude string `xml:"Longitude"`
}
