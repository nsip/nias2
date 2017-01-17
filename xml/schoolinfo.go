package xml

type SchoolInfo struct {
	RefId           string `xml:"RefId,attr"`
	LocalId         string `xml:"LocalId"`
	StateProvinceId string `xml:"StateProvinceId"`
	ACARAId         string `xml:"ACARAId"`
	SchoolName      string `xml:"SchoolName"`
	LEAInfoRefId    string `xml:"LEAInfoRefId"`
	OtherLEA        string `xml:"OtherLEA"`
	SchoolDistrict  string `xml:"SchoolDistrict"`
	SchoolType      string `xml:"SchoolType"`
	StudentCount    string // non xml field added for reporting

	SchoolFocusList struct {
		SchoolFocus []string `xml:"SchoolFocus"`
	} `xml:"SchoolFocusList"`
	SchoolURL string `xml:"SchoolURL"`

	PrincipalInfo struct {
		ContactName struct {
			Type       string `xml:"Type,attr"`
			Title      string `xml:"Title"`
			FamilyName string `xml:"FamilyName"`
			GivenName  string `xml:"GivenName"`
			MiddleName string `xml:"MiddleName"`
			Suffix     string `xml:"Suffix"`
			FullName   string `xml:"FullName"`
		} `xml:"ContactName"`
		ContactTitle string `xml:"ContactTitle"`
	} `xml:"PrincipalInfo"`

	SchoolContactList struct {
		SchoolContact []struct {
			PublishInDirectory string `xml:"PublishInDirectory"`
			ContactInfo        struct {
				Name struct {
					Type       string `xml:"Type,attr"`
					Title      string `xml:"Title"`
					FamilyName string `xml:"FamilyName"`
					GivenName  string `xml:"GivenName"`
					MiddleName string `xml:"MiddleName"`
					Suffix     string `xml:"Suffix"`
					FullName   string `xml:"FullName"`
				} `xml:"Name"`

				PositionTitle string `xml:"PositionTitle"`
				Role          string `xml:"Role"`

				Address struct {
					Type   string `xml:"Type,attr"`
					Role   string `xml:"Role,attr"`
					Street struct {
						Line1 string `xml:"Line1"`
					} `xml:"Street"`
					City          string `xml:"City"`
					StateProvince string `xml:"StateProvince"`
					Country       string `xml:"Country"`
					PostalCode    string `xml:"PostalCode"`
					GridLocation  struct {
						Latitude  string `xml:"Latitude"`
						Longitude string `xml:"Longitude"`
					} `xml:"GridLocation"`
				} `xml:"Address"`

				EmailList struct {
					Email []struct {
						Type    string `xml:"Type,attr"`
						Address string `xml:"Email"`
					}
				} `xml:"EmailList"`

				PhoneNumberList struct {
					PhoneNumber []struct {
						Type         string `xml:"Type,attr"`
						Number       string `xml:"Number"`
						Extension    string `xml:"Extension"`
						ListedStatus string `xml:"ListedStatus"`
					} `xml:"PhoneNumber"`
				} `xml:"PhoneNumberList"`
			} `xml:"ContactInfo"`
		} `xml:"SchoolContact"`
	} `xml:"SchoolContactList"`

	PhoneNumberList struct {
		PhoneNumber []struct {
			Type   string `xml:"Type,attr"`
			Number string `xml:"Number"`
		} `xml:"PhoneNumber"`
	} `xml:"PhoneNumberList"`

	SessionType string `xml:"SessionType"`

	YearLevels struct {
		YearLevel []struct {
			Code string `xml:"Code"`
		} `xml:"YearLevel"`
	} `xml:"YearLevels"`

	ARIA              string `xml:"ARIA"`
	OperationalStatus string `xml:"OperationalStatus"`
	FederalElectorate string `xml:"FederalElectorate"`

	Campus struct {
		SchoolCampusId string `xml:"SchoolCampusId"`
		CampusType     string `xml:"CampusType"`
		AdminStatus    string `xml:"AdminStatus"`
	} `xml:"Campus"`

	SchoolSector             string `xml:"SchoolSector"`
	IndependentSchool        string `xml:"IndependentSchool"`
	NonGovSystemicStatus     string `xml:"NonGovSystemicStatus"`
	System                   string `xml:"System"`
	ReligiousAffiliation     string `xml:"ReligiousAffiliation"`
	SchoolGeographicLocation string `xml:"SchoolGeographicLocation"`
	LocalGovernmentArea      string `xml:"LocalGovernmentArea"`
	JurisdictionLowerHouse   string `xml:"JurisdictionLowerHouse"`
	SLA                      string `xml:"SLA"`
	SchoolCoEdStatus         string `xml:"SchoolCoEdStatus"`

	SchoolGroupList struct {
		SchoolGroup []string `xml:"SchoolGroup"`
	} `xml:"SchoolGroupList"`
}
