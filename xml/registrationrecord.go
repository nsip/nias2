package xml

import (
	"encoding/gob"
)

func init() {
	// make gob encoder aware of local types
	gob.Register(RegistrationRecord{})
}

// StudentPersonal for results reporting
type RegistrationRecord struct {
	// XML Configuration
	// XMLName            xml.Name `xml:"StudentPersonal"`
	// Important fields
	RefId               string `json:",omitempty" xml:"RefId,attr"`
	ASLSchoolId         string `json:",omitempty" xml:"MostRecent>SchoolACARAId"`
	AddressLine1        string `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line1"`
	AddressLine2        string `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line2"`
	BirthDate           string `json:",omitempty" xml:"PersonInfo>Demographics>BirthDate"`
	ClassCode           string `json:",omitempty" xml:"MostRecent>ClassCode"`
	CountryOfBirth      string `json:",omitempty" xml:"PersonInfo>Demographics>CountryOfBirth"`
	DiocesanId          string `json:",omitempty"`
	EducationSupport    string `json:",omitempty" xml:"EducationSupport"`
	FFPOS               string `json:",omitempty" xml:"MostRecent>FFPOS"`
	FTE                 string `json:",omitempty" xml:"MostRecent>FTE"`
	FamilyName          string `json:",omitempty" xml:"PersonInfo>Name>FamilyName"`
	GivenName           string `json:",omitempty" xml:"PersonInfo>Name>GivenName"`
	HomeSchooledStudent string `json:",omitempty" xml:"HomeSchooledStudent"`
	Homegroup           string `json:",omitempty" xml:"MostRecent>Homegroup"`
	IndigenousStatus    string `json:",omitempty" xml:"PersonInfo>Demographics>IndigenousStatus"`
	JurisdictionId      string `json:",omitempty"`
	LBOTE               string `json:",omitempty" xml:"PersonInfo>Demographics>LBOTE"`
	LocalCampusId       string `json:",omitempty" xml:"MostRecent>LocalCampusId"`
	LocalId             string `json:",omitempty" xml:"LocalId"`
	Locality            string `json:",omitempty" xml:"PersonInfo>AddressList>Address>City"`
	MainSchoolFlag      string `json:",omitempty" xml:"MostRecent>MembershipType"`
	MiddleName          string `json:",omitempty" xml:"PersonInfo>Name>MiddleName"`
	NationalId          string `json:",omitempty"`
	OfflineDelivery     string `json:",omitempty" xml:"OfflineDelivery"`
	OtherIdList         struct {
		OtherId []struct {
			Type  string `xml:"Type,attr"`
			Value string `xml:",chardata"`
		} `xml:"OtherId"`
	} `xml:OtherIdList`
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
	PreviousStateProvinceId   string `json:",omitempty"`
	PreviousTAAId             string `json:",omitempty"`
	ReportingSchoolId         string `json:",omitempty" xml:"MostRecent>ReportingSchoolId"`
	SchoolLocalId             string `json:",omitempty" xml:"MostRecent>SchoolLocalId"`
	SectorId                  string `json:",omitempty"`
	Sensitive                 string `json:",omitempty" xml:"Sensitive"`
	Sex                       string `json:",omitempty" xml:"PersonInfo>Demographics>Sex"`
	StateProvinceId           string `json:",omitempty"`
	StateTerritory            string `json:",omitempty" xml:"PersonInfo>AddressList>Address>StateProvince"`
	StudentLOTE               string `json:",omitempty" xml:"PersonInfo>Demographics>LanguageList>Language>Code"`
	TAAId                     string `json:",omitempty"`
	TestLevel                 string `json:",omitempty" xml:"MostRecent>TestLevel>Code"`
	VisaCode                  string `json:",omitempty" xml:"PersonInfo>Demographics>VisaSubClass"`
	YearLevel                 string `json:",omitempty" xml:"MostRecent>YearLevel>Code"`
}

// convenience method to return otherid by type
func (r RegistrationRecord) GetOtherId(idtype string) string {

	for _, id := range r.OtherIdList.OtherId {
		if id.Type == idtype {
			return id.Value
		}
	}

	return idtype
}

// convenience method for writing to csv
func (r RegistrationRecord) GetHeaders() []string {
	return []string{"ASLSchoolId",
		"AddressLine1",
		"AddressLine2",
		"BirthDate",
		"ClassCode",
		"CountryOfBirth",
		"DiocesanId",
		"EducationSupport",
		"FFPOS",
		"FTE",
		"FamilyName",
		"GivenName",
		"HomeSchooledStudent",
		"Homegroup",
		"IndigenousStatus",
		"JurisdictionId",
		"LBOTE",
		"LocalCampusId",
		"LocalId",
		"Locality",
		"MainSchoolFlag",
		"MiddleName",
		"NationalId",
		"OfflineDelivery",
		"OtherId",
		"OtherSchoolId",
		"Parent1LOTE",
		"Parent1NonSchoolEducation",
		"Parent1Occupation",
		"Parent1SchoolEducation",
		"Parent2LOTE",
		"Parent2NonSchoolEducation",
		"Parent2Occupation",
		"Parent2SchoolEducation",
		"PlatformId",
		"Postcode",
		"PreferredName",
		"PreviousDiocesanId",
		"PreviousJurisdictionId",
		"PreviousLocalId",
		"PreviousNationalId",
		"PreviousOtherId",
		"PreviousPlatformId",
		"PreviousSectorId",
		"PreviousStateProvinceId",
		"PreviousTAAId",
		"ReportingSchoolId",
		"SchoolLocalId",
		"SectorId",
		"Sensitive",
		"Sex",
		"StateProvinceId",
		"StateTerritory",
		"StudentLOTE",
		"TAAId",
		"TestLevel",
		"VisaCode",
		"YearLevel"}
}

// convenience method for writing to csv
func (r RegistrationRecord) GetSlice() []string {
	return []string{r.ASLSchoolId,
		r.AddressLine1,
		r.AddressLine2,
		r.BirthDate,
		r.ClassCode,
		r.CountryOfBirth,
		r.DiocesanId,
		r.EducationSupport,
		r.FFPOS,
		r.FTE,
		r.FamilyName,
		r.GivenName,
		r.HomeSchooledStudent,
		r.Homegroup,
		r.IndigenousStatus,
		r.JurisdictionId,
		r.LBOTE,
		r.LocalCampusId,
		r.LocalId,
		r.Locality,
		r.MainSchoolFlag,
		r.MiddleName,
		r.NationalId,
		r.OfflineDelivery,
		r.GetOtherId("OtherStudentId"),
		r.OtherSchoolId,
		r.Parent1LOTE,
		r.Parent1NonSchoolEducation,
		r.Parent1Occupation,
		r.Parent1SchoolEducation,
		r.Parent2LOTE,
		r.Parent2NonSchoolEducation,
		r.Parent2Occupation,
		r.Parent2SchoolEducation,
		r.PlatformId,
		r.Postcode,
		r.PreferredName,
		r.PreviousDiocesanId,
		r.PreviousJurisdictionId,
		r.PreviousLocalId,
		r.PreviousNationalId,
		r.PreviousOtherId,
		r.PreviousPlatformId,
		r.PreviousSectorId,
		r.PreviousStateProvinceId,
		r.PreviousTAAId,
		r.ReportingSchoolId,
		r.SchoolLocalId,
		r.SectorId,
		r.Sensitive,
		r.Sex,
		r.StateProvinceId,
		r.StateTerritory,
		r.StudentLOTE,
		r.TAAId,
		r.TestLevel,
		r.VisaCode,
		r.YearLevel}
}
