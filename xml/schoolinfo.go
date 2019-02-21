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
	//Address []string `xml:"Address,omitempty"` // avoid shortcut for XML roundtripping!
	Address []Address `xml:"Address,omitempty"`
}

type SchoolGroupList struct {
	SchoolGroup []string `xml:"SchoolGroup,omitempty"`
}

type Campus struct {
	SchoolCampusId string `xml:"SchoolCampusId,omitempty"`
	CampusType     string `xml:"CampusType,omitempty"`
	AdminStatus    string `xml:"AdminStatus,omitempty"`
}

type YearLevelList struct {
	YearLevel []YearLevel `xml:"YearLevel,omitempty"`
}

type YearLevel struct {
	Code string `xml:"Code"`
}

type PhoneNumberList struct {
	PhoneNumber []PhoneNumber `xml:"PhoneNumber"`
}

type PhoneNumber struct {
	Type   string `xml:"Type,attr,omitempty"`
	Number string `xml:"Number,omitempty"`
}

type SchoolFocusList struct {
	SchoolFocus []string `xml:"SchoolFocus,omitempty"`
}

type NameType struct {
	Type       string `xml:"Type,attr"`
	Title      string `xml:"Title,omitempty"`
	FamilyName string `xml:"FamilyName,omitempty"`
	GivenName  string `xml:"GivenName,omitempty"`
	MiddleName string `xml:"MiddleName,omitempty"`
	Suffix     string `xml:"Suffix,omitempty"`
	FullName   string `xml:"FullName,omitempty"`
}

type EmailList struct {
	Email Email
}

type Email struct {
	Type    string `xml:"Type,attr,omitempty"`
	Address string `xml:"Email,omitempty"`
}

type FullPhoneNumberList struct {
	PhoneNumber []FullPhoneNumber `xml:"PhoneNumber,omitempty"`
}

type FullPhoneNumber struct {
	Type         string `xml:"Type,attr,omitempty"`
	Number       string `xml:"Number,omitempty"`
	Extension    string `xml:"Extension,omitempty"`
	ListedStatus string `xml:"ListedStatus,omitempty"`
}

type PrincipalInfo struct {
	ContactName  NameType `xml:"ContactName,omitempty"`
	ContactTitle string   `xml:"ContactTitle,omitempty"`
}

type SchoolContactList struct {
	SchoolContact []SchoolContact `xml:"SchoolContact,omitempty"`
}

type SchoolContact struct {
	PublishInDirectory string      `xml:"PublishInDirectory,omitempty"`
	ContactInfo        ContactInfo `xml:"ContactInfo,omitempty"`
}

type ContactInfo struct {
	Name            NameType            `xml:"Name,omitempty"`
	PositionTitle   string              `xml:"PositionTitle,omitempty"`
	Role            string              `xml:"Role,omitempty"`
	Address         Address             `xml:"Address,omitempty"`
	EmailList       EmailList           `xml:"EmailList,omitempty"`
	PhoneNumberList FullPhoneNumberList `xml:"PhoneNumberList,omitempty"`
}

type Address struct {
	Type          string       `xml:"Type,attr,omitempty"`
	Role          string       `xml:"Role,attr,omitempty"`
	Street        Street       `xml:"Street,omitempty"`
	City          string       `xml:"City,omitempty"`
	StateProvince string       `xml:"StateProvince,omitempty"`
	Country       string       `xml:"Country,omitempty"`
	PostalCode    string       `xml:"PostalCode,omitempty"`
	GridLocation  GridLocation `xml:"GridLocation,omitempty"`
}

type Street struct {
	Line1 string `xml:"Line1"`
	Line2 string `xml:"Line2"`
}

type GridLocation struct {
	Latitude  string `xml:"Latitude"`
	Longitude string `xml:"Longitude"`
}
