package xml

import (
	"encoding/gob"
	"encoding/xml"
	"strings"
)

func init() {
	// make gob encoder aware of local types
	gob.Register(RegistrationRecord{})
	gob.Register(GraphStruct{})
	// gob.Register(SifStudentPersonal{})
}

type XMLAttributeStruct struct {
	Type  string `xml:"Type,attr"`
	Value string `xml:",chardata"`
}

// StudentPersonal for results reporting
/* Contents of OtherIdList are duplicated into separate fields which are not in the XML; e.g. DiocesanId.
On ingest of CSV, Unflatten() is used to populate OtherIdList.
On export to JSON or CSV, Flatten() is used to populate the duplicate fields. */
type RegistrationRecord struct {
	// XML Configuration
	XMLName xml.Name `xml:"StudentPersonal"`
	// Important fields
	RefId           string `json:",omitempty" xml:"RefId,attr"`
	LocalId         string `json:",omitempty" xml:"LocalId,omitempty"`
	StateProvinceId string `json:",omitempty" xml:"StateProvinceId,omitempty"`
	OtherIdList     struct {
		OtherId []XMLAttributeStruct `xml:"OtherId"`
	} `xml:OtherIdList,omitempty`
	FamilyName                string `json:",omitempty" xml:"PersonInfo>Name>FamilyName,omitempty"`
	GivenName                 string `json:",omitempty" xml:"PersonInfo>Name>GivenName,omitempty"`
	MiddleName                string `json:",omitempty" xml:"PersonInfo>Name>MiddleName,omitempty"`
	PreferredName             string `json:",omitempty" xml:"PersonInfo>Name>PreferredGivenName,omitempty"`
	IndigenousStatus          string `json:",omitempty" xml:"PersonInfo>Demographics>IndigenousStatus,omitempty"`
	Sex                       string `json:",omitempty" xml:"PersonInfo>Demographics>Sex,omitempty"`
	BirthDate                 string `json:",omitempty" xml:"PersonInfo>Demographics>BirthDate,omitempty"`
	CountryOfBirth            string `json:",omitempty" xml:"PersonInfo>Demographics>CountryOfBirth,omitempty"`
	StudentLOTE               string `json:",omitempty" xml:"PersonInfo>Demographics>LanguageList>Language>Code,omitempty"`
	VisaCode                  string `json:",omitempty" xml:"PersonInfo>Demographics>VisaSubClass,omitempty"`
	LBOTE                     string `json:",omitempty" xml:"PersonInfo>Demographics>LBOTE,omitempty"`
	AddressLine1              string `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line1,omitempty"`
	AddressLine2              string `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line2,omitempty"`
	Locality                  string `json:",omitempty" xml:"PersonInfo>AddressList>Address>City,omitempty"`
	StateTerritory            string `json:",omitempty" xml:"PersonInfo>AddressList>Address>StateProvince,omitempty"`
	Postcode                  string `json:",omitempty" xml:"PersonInfo>AddressList>Address>PostalCode,omitempty"`
	SchoolLocalId             string `json:",omitempty" xml:"MostRecent>SchoolLocalId,omitempty"`
	YearLevel                 string `json:",omitempty" xml:"MostRecent>YearLevel>Code,omitempty"`
	FTE                       string `json:",omitempty" xml:"MostRecent>FTE,omitempty"`
	Parent1LOTE               string `json:",omitempty" xml:"MostRecent>Parent1Language,omitempty"`
	Parent2LOTE               string `json:",omitempty" xml:"MostRecent>Parent2Language,omitempty"`
	Parent1Occupation         string `json:",omitempty" xml:"MostRecent>Parent1EmploymentType,omitempty"`
	Parent2Occupation         string `json:",omitempty" xml:"MostRecent>Parent2EmploymentType,omitempty"`
	Parent1SchoolEducation    string `json:",omitempty" xml:"MostRecent>Parent1SchoolEducationLevel,omitempty"`
	Parent2SchoolEducation    string `json:",omitempty" xml:"MostRecent>Parent2SchoolEducationLevel,omitempty"`
	Parent1NonSchoolEducation string `json:",omitempty" xml:"MostRecent>Parent1NonSchoolEducation,omitempty"`
	Parent2NonSchoolEducation string `json:",omitempty" xml:"MostRecent>Parent2NonSchoolEducation,omitempty"`
	LocalCampusId             string `json:",omitempty" xml:"MostRecent>LocalCampusId,omitempty"`
	ASLSchoolId               string `json:",omitempty" xml:"MostRecent>SchoolACARAId,omitempty"`
	TestLevel                 string `json:",omitempty" xml:"MostRecent>TestLevel>Code,omitempty"`
	Homegroup                 string `json:",omitempty" xml:"MostRecent>Homegroup,omitempty"`
	ClassGroup                string `json:",omitempty" xml:"MostRecent>ClassCode,omitempty"`
	MainSchoolFlag            string `json:",omitempty" xml:"MostRecent>MembershipType,omitempty"`
	FFPOS                     string `json:",omitempty" xml:"MostRecent>FFPOS,omitempty"`
	ReportingSchoolId         string `json:",omitempty" xml:"MostRecent>ReportingSchoolId,omitempty"`
	OtherSchoolId             string `json:",omitempty" xml:"MostRecent>OtherEnrollmentSchoolACARAId,omitempty"`
	EducationSupport          string `json:",omitempty" xml:"EducationSupport,omitempty"`
	HomeSchooledStudent       string `json:",omitempty" xml:"HomeSchooledStudent,omitempty"`
	Sensitive                 string `json:",omitempty" xml:"Sensitive,omitempty"`
	OfflineDelivery           string `json:",omitempty" xml:"OfflineDelivery,omitempty"`
	DiocesanId                string `json:",omitempty" xml:"-"`
	JurisdictionId            string `json:",omitempty" xml:"-"`
	NationalId                string `json:",omitempty" xml:"-"`
	OtherId                   string `json:",omitempty" xml:"-"`
	PlatformId                string `json:",omitempty" xml:"-"`
	PreviousDiocesanId        string `json:",omitempty" xml:"-"`
	//PreviousJurisdictionId    string `json:",omitempty"`
	PreviousLocalId         string `json:",omitempty" xml:"-"`
	PreviousNationalId      string `json:",omitempty" xml:"-"`
	PreviousOtherId         string `json:",omitempty" xml:"-"`
	PreviousPlatformId      string `json:",omitempty" xml:"-"`
	PreviousSectorId        string `json:",omitempty" xml:"-"`
	PreviousStateProvinceId string `json:",omitempty" xml:"-"`
	PreviousTAAId           string `json:",omitempty" xml:"-"`
	SectorId                string `json:",omitempty" xml:"-"`
	TAAId                   string `json:",omitempty" xml:"-"`
}

// Flatten out Other IDs from XML into JSON/CSV flat structure
func (r *RegistrationRecord) Flatten() RegistrationRecord {
	for _, id := range r.OtherIdList.OtherId {
		if id.Type == "DiocesanStudentId" {
			r.DiocesanId = id.Value
		}
		if id.Type == "NationalStudentId" {
			r.NationalId = id.Value
		}
		if id.Type == "OtherId" {
			r.OtherId = id.Value
		}
		if id.Type == "NAPPlatformStudentId" {
			r.PlatformId = id.Value
		}
		if id.Type == "PreviousDiocesanStudentId" {
			r.PreviousDiocesanId = id.Value
		}
		if id.Type == "PreviousLocalSchoolStudentId" {
			r.PreviousLocalId = id.Value
		}
		if id.Type == "PreviousNationalStudentId" {
			r.PreviousNationalId = id.Value
		}
		if id.Type == "PreviousOtherId" {
			r.PreviousOtherId = id.Value
		}
		if id.Type == "PreviousNAPPlatformStudentId" {
			r.PreviousPlatformId = id.Value
		}
		if id.Type == "PreviousSectorStudentId" {
			r.PreviousSectorId = id.Value
		}
		if id.Type == "PreviousStateProvinceId" {
			r.PreviousStateProvinceId = id.Value
		}
		if id.Type == "PreviousTAAStudentId" {
			r.PreviousTAAId = id.Value
		}
		if id.Type == "SectorStudentId" {
			r.SectorId = id.Value
		}
		if id.Type == "TAAStudentId" {
			r.TAAId = id.Value
		}
	}
	return *r
}

// Unflatten out Other IDs from JSON/CSV flat structure into XML structure
func (r *RegistrationRecord) Unflatten() RegistrationRecord {
	r.OtherIdList.OtherId = make([]XMLAttributeStruct, 0)
	if r.DiocesanId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"DiocesanStudentId", r.DiocesanId})
	}
	if r.NationalId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"NationalStudentId", r.NationalId})
	}
	if r.OtherId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"OtherStudentId", r.OtherId})
	}
	if r.PlatformId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"NAPPlatformStudentId", r.PlatformId})
	}
	if r.PreviousDiocesanId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousDiocesanStudentId", r.PreviousDiocesanId})
	}
	if r.PreviousLocalId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousLocalSchoolStudentId", r.PreviousLocalId})
	}
	if r.PreviousNationalId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousNationalStudentId", r.PreviousNationalId})
	}
	if r.PreviousOtherId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousOtherId", r.PreviousOtherId})
	}
	if r.PreviousPlatformId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousNAPPlatformStudentId", r.PreviousPlatformId})
	}
	if r.PreviousSectorId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousSectorStudentId", r.PreviousSectorId})
	}
	if r.PreviousStateProvinceId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousStateProvinceId", r.PreviousStateProvinceId})
	}
	if r.PreviousTAAId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PreviousTAAStudentId", r.PreviousTAAId})
	}
	if r.SectorId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"SectorStudentId", r.SectorId})
	}
	if r.TAAId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"TAAStudentId", r.TAAId})
	}
	return *r
}

// convenience method to return otherid by type
func (r RegistrationRecord) GetOtherId(idtype string) string {

	for _, id := range r.OtherIdList.OtherId {
		if strings.EqualFold(id.Type, idtype) {
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
		"ClassGroup",
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
		r.ClassGroup,
		r.CountryOfBirth,
		//r.DiocesanId,
		r.GetOtherId("DiocesanId"),
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

// information extracted out of SIF for graph
type GraphStruct struct {
	Guid          string            // RefID of object
	EquivalentIds []string          // equivalent Ids
	OtherIds      map[string]string // map of OtherId type to OtherId
	Type          string            // object type
	Links         []string          // list of related ids
	Label         string            // human readable label
}
