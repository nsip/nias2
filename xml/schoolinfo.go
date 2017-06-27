package xml

type SchoolInfo struct {
	RefId                    string            `xml:"RefId,attr"`
	LocalId                  string            `xml:"LocalId"`
	StateProvinceId          string            `xml:"StateProvinceId"`
	ACARAId                  string            `xml:"ACARAId"`
	SchoolName               string            `xml:"SchoolName"`
	LEAInfoRefId             string            `xml:"LEAInfoRefId"`
	OtherLEA                 string            `xml:"OtherLEA"`
	SchoolDistrict           string            `xml:"SchoolDistrict"`
	SchoolType               string            `xml:"SchoolType"`
	StudentCount             string            // non xml field added for reporting
	SchoolFocusList          SchoolFocusList   `xml:"SchoolFocusList"`
	SchoolURL                string            `xml:"SchoolURL"`
	PrincipalInfo            PrincipalInfo     `xml:"PrincipalInfo"`
	SchoolContactList        SchoolContactList `xml:"SchoolContactList"`
	PhoneNumberList          PhoneNumberList   `xml:"PhoneNumberList"`
	SessionType              string            `xml:"SessionType"`
	YearLevels               YearLevelList     `xml:"YearLevels"`
	ARIA                     string            `xml:"ARIA"`
	OperationalStatus        string            `xml:"OperationalStatus"`
	FederalElectorate        string            `xml:"FederalElectorate"`
	Campus                   Campus            `xml:"Campus"`
	SchoolSector             string            `xml:"SchoolSector"`
	IndependentSchool        string            `xml:"IndependentSchool"`
	NonGovSystemicStatus     string            `xml:"NonGovSystemicStatus"`
	System                   string            `xml:"System"`
	ReligiousAffiliation     string            `xml:"ReligiousAffiliation"`
	SchoolGeographicLocation string            `xml:"SchoolGeographicLocation"`
	LocalGovernmentArea      string            `xml:"LocalGovernmentArea"`
	JurisdictionLowerHouse   string            `xml:"JurisdictionLowerHouse"`
	SLA                      string            `xml:"SLA"`
	SchoolCoEdStatus         string            `xml:"SchoolCoEdStatus"`
	SchoolGroupList          SchoolGroupList   `xml:"SchoolGroupList"`
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

type YearLevel []struct {
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
