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
	Type  string `xml:"Type,attr,omitempty"`
	Value string `xml:",chardata"`
}

type OtherIdList struct {
	OtherId []XMLAttributeStruct `xml:"OtherId"`
}

// StudentPersonal for results reporting
/* Contents of OtherIdList are duplicated into separate fields which are not in the XML; e.g. DiocesanId.
On ingest of CSV, Unflatten() is used to populate OtherIdList.
On export to JSON or CSV, Flatten() is used to populate the duplicate fields. */
type RegistrationRecord struct {
	// XML Configuration
	XMLName xml.Name `xml:"StudentPersonal"`
	// Important fields
	RefId                        string      `json:",omitempty" xml:"RefId,attr"`
	LocalId                      string      `json:",omitempty" xml:"LocalId"`
	StateProvinceId              string      `json:",omitempty" xml:"StateProvinceId,omitempty"`
	OtherIdList                  OtherIdList `xml:OtherIdList,omitempty`
	NameType                     string      `json:",omitempty" xml:"-"`
	FamilyName                   string      `json:",omitempty" xml:"PersonInfo>Name>FamilyName,omitempty"`
	GivenName                    string      `json:",omitempty" xml:"PersonInfo>Name>GivenName,omitempty"`
	MiddleName                   string      `json:",omitempty" xml:"PersonInfo>Name>MiddleName,omitempty"`
	PreferredName                string      `json:",omitempty" xml:"PersonInfo>Name>PreferredGivenName,omitempty"`
	IndigenousStatus             string      `json:",omitempty" xml:"PersonInfo>Demographics>IndigenousStatus,omitempty"`
	Sex                          string      `json:",omitempty" xml:"PersonInfo>Demographics>Sex,omitempty"`
	BirthDate                    string      `json:",omitempty" xml:"PersonInfo>Demographics>BirthDate,omitempty"`
	CountryOfBirth               string      `json:",omitempty" xml:"PersonInfo>Demographics>CountryOfBirth,omitempty"`
	StudentLOTE                  string      `json:",omitempty" xml:"PersonInfo>Demographics>LanguageList>Language>Code,omitempty"`
	VisaCode                     string      `json:",omitempty" xml:"PersonInfo>Demographics>VisaSubClass,omitempty"`
	LBOTE                        string      `json:",omitempty" xml:"PersonInfo>Demographics>LBOTE,omitempty"`
	AddressRole                  string      `json:",omitempty" xml:"-"`
	AddressType                  string      `json:",omitempty" xml:"-"`
	AddressLine1                 string      `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line1,omitempty"`
	AddressLine2                 string      `json:",omitempty" xml:"PersonInfo>AddressList>Address>Street>Line2,omitempty"`
	Locality                     string      `json:",omitempty" xml:"PersonInfo>AddressList>Address>City,omitempty"`
	StateTerritory               string      `json:",omitempty" xml:"PersonInfo>AddressList>Address>StateProvince,omitempty"`
	Postcode                     string      `json:",omitempty" xml:"PersonInfo>AddressList>Address>PostalCode,omitempty"`
	SchoolLocalId                string      `json:",omitempty" xml:"MostRecent>SchoolLocalId,omitempty"`
	YearLevel                    string      `json:",omitempty" xml:"MostRecent>YearLevel>Code,omitempty"`
	FTE                          string      `json:",omitempty" xml:"MostRecent>FTE,omitempty"`
	Parent1LOTE                  string      `json:",omitempty" xml:"MostRecent>Parent1Language,omitempty"`
	Parent2LOTE                  string      `json:",omitempty" xml:"MostRecent>Parent2Language,omitempty"`
	Parent1Occupation            string      `json:",omitempty" xml:"MostRecent>Parent1EmploymentType,omitempty"`
	Parent2Occupation            string      `json:",omitempty" xml:"MostRecent>Parent2EmploymentType,omitempty"`
	Parent1SchoolEducation       string      `json:",omitempty" xml:"MostRecent>Parent1SchoolEducationLevel,omitempty"`
	Parent2SchoolEducation       string      `json:",omitempty" xml:"MostRecent>Parent2SchoolEducationLevel,omitempty"`
	Parent1NonSchoolEducation    string      `json:",omitempty" xml:"MostRecent>Parent1NonSchoolEducation,omitempty"`
	Parent2NonSchoolEducation    string      `json:",omitempty" xml:"MostRecent>Parent2NonSchoolEducation,omitempty"`
	LocalCampusId                string      `json:",omitempty" xml:"MostRecent>LocalCampusId,omitempty"`
	ASLSchoolId                  string      `json:",omitempty" xml:"MostRecent>SchoolACARAId,omitempty"`
	TestLevel                    string      `json:",omitempty" xml:"MostRecent>TestLevel>Code,omitempty"`
	Homegroup                    string      `json:",omitempty" xml:"MostRecent>Homegroup,omitempty"`
	ClassGroup                   string      `json:",omitempty" xml:"MostRecent>ClassCode,omitempty"`
	MainSchoolFlag               string      `json:",omitempty" xml:"MostRecent>MembershipType,omitempty"`
	FFPOS                        string      `json:",omitempty" xml:"MostRecent>FFPOS,omitempty"`
	ReportingSchoolId            string      `json:",omitempty" xml:"MostRecent>ReportingSchoolId,omitempty"`
	OtherEnrollmentSchoolACARAId string      `json:",omitempty" xml:"MostRecent>OtherEnrollmentSchoolACARAId,omitempty"`
	EducationSupport             string      `json:",omitempty" xml:"EducationSupport,omitempty"`
	HomeSchooledStudent          string      `json:",omitempty" xml:"HomeSchooledStudent,omitempty"`
	Sensitive                    string      `json:",omitempty" xml:"Sensitive,omitempty"`
	OfflineDelivery              string      `json:",omitempty" xml:"OfflineDelivery,omitempty"`
	DiocesanId                   string      `json:",omitempty" xml:"-"`
	JurisdictionId               string      `json:",omitempty" xml:"-"`
	NationalId                   string      `json:",omitempty" xml:"-"`
	OtherId                      string      `json:",omitempty" xml:"-"`
	PlatformId                   string      `json:",omitempty" xml:"-"`
	PreviousDiocesanId           string      `json:",omitempty" xml:"-"`
	PreviousLocalId              string      `json:",omitempty" xml:"-"`
	PreviousNationalId           string      `json:",omitempty" xml:"-"`
	PreviousOtherId              string      `json:",omitempty" xml:"-"`
	PreviousPlatformId           string      `json:",omitempty" xml:"-"`
	PreviousSectorId             string      `json:",omitempty" xml:"-"`
	PreviousStateProvinceId      string      `json:",omitempty" xml:"-"`
	PreviousJurisdictionId       string      `json:",omitempty" xml:"-"`
	PreviousTAAId                string      `json:",omitempty" xml:"-"`
	SectorId                     string      `json:",omitempty" xml:"-"`
	TAAId                        string      `json:",omitempty" xml:"-"`
	OtherIdLocality              string      `json:",omitempty" xml:"-"`
	DOBRange                     string      `json:",omitempty" xml:"-"`
	PersonalDetailsChanged       string      `json:",omitempty" xml:"-"`
	PossibleDuplicate            string      `json:",omitempty" xml:"-"`
	PsiOtherIdMismatch           string      `json:",omitempty" xml:"-"`
	OtherSchoolId                string      `json:",omitempty" xml:"-"`
}

// More complete SIF/XML representation
type SIFRegistrationRecord struct {
	XMLName                      xml.Name    `xml:"StudentPersonal"`
	RefId                        string      `json:",omitempty" xml:"RefId,attr"`
	LocalId                      string      `json:",omitempty" xml:"LocalId"`
	StateProvinceId              string      `json:",omitempty" xml:"StateProvinceId,omitempty"`
	OtherIdList                  OtherIdList `xml:OtherIdList,omitempty`
	Name                         SIFName     `json:",omitempty" xml:"PersonInfo>Name,omitempty"`
	IndigenousStatus             string      `json:",omitempty" xml:"PersonInfo>Demographics>IndigenousStatus,omitempty"`
	Sex                          string      `json:",omitempty" xml:"PersonInfo>Demographics>Sex,omitempty"`
	BirthDate                    string      `json:",omitempty" xml:"PersonInfo>Demographics>BirthDate,omitempty"`
	CountryOfBirth               string      `json:",omitempty" xml:"PersonInfo>Demographics>CountryOfBirth,omitempty"`
	StudentLOTE                  string      `json:",omitempty" xml:"PersonInfo>Demographics>LanguageList>Language>Code,omitempty"`
	VisaCode                     string      `json:",omitempty" xml:"PersonInfo>Demographics>VisaSubClass,omitempty"`
	LBOTE                        string      `json:",omitempty" xml:"PersonInfo>Demographics>LBOTE,omitempty"`
	Address                      Address     `json:",omitempty" xml:"PersonInfo>AddressList>Address"`
	SchoolLocalId                string      `json:",omitempty" xml:"MostRecent>SchoolLocalId,omitempty"`
	YearLevel                    string      `json:",omitempty" xml:"MostRecent>YearLevel>Code,omitempty"`
	FTE                          string      `json:",omitempty" xml:"MostRecent>FTE,omitempty"`
	Parent1LOTE                  string      `json:",omitempty" xml:"MostRecent>Parent1Language,omitempty"`
	Parent2LOTE                  string      `json:",omitempty" xml:"MostRecent>Parent2Language,omitempty"`
	Parent1Occupation            string      `json:",omitempty" xml:"MostRecent>Parent1EmploymentType,omitempty"`
	Parent2Occupation            string      `json:",omitempty" xml:"MostRecent>Parent2EmploymentType,omitempty"`
	Parent1SchoolEducation       string      `json:",omitempty" xml:"MostRecent>Parent1SchoolEducationLevel,omitempty"`
	Parent2SchoolEducation       string      `json:",omitempty" xml:"MostRecent>Parent2SchoolEducationLevel,omitempty"`
	Parent1NonSchoolEducation    string      `json:",omitempty" xml:"MostRecent>Parent1NonSchoolEducation,omitempty"`
	Parent2NonSchoolEducation    string      `json:",omitempty" xml:"MostRecent>Parent2NonSchoolEducation,omitempty"`
	LocalCampusId                string      `json:",omitempty" xml:"MostRecent>LocalCampusId,omitempty"`
	ASLSchoolId                  string      `json:",omitempty" xml:"MostRecent>SchoolACARAId,omitempty"`
	TestLevel                    string      `json:",omitempty" xml:"MostRecent>TestLevel>Code,omitempty"`
	Homegroup                    string      `json:",omitempty" xml:"MostRecent>Homegroup,omitempty"`
	ClassGroup                   string      `json:",omitempty" xml:"MostRecent>ClassCode,omitempty"`
	MainSchoolFlag               string      `json:",omitempty" xml:"MostRecent>MembershipType,omitempty"`
	FFPOS                        string      `json:",omitempty" xml:"MostRecent>FFPOS,omitempty"`
	ReportingSchoolId            string      `json:",omitempty" xml:"MostRecent>ReportingSchoolId,omitempty"`
	OtherEnrollmentSchoolACARAId string      `json:",omitempty" xml:"MostRecent>OtherEnrollmentSchoolACARAId,omitempty"`
	EducationSupport             string      `json:",omitempty" xml:"EducationSupport,omitempty"`
	HomeSchooledStudent          string      `json:",omitempty" xml:"HomeSchooledStudent,omitempty"`
	Sensitive                    string      `json:",omitempty" xml:"Sensitive,omitempty"`
	OfflineDelivery              string      `json:",omitempty" xml:"OfflineDelivery,omitempty"`
}

func (r *SIFRegistrationRecord) From_SIF() RegistrationRecord {
	ret := RegistrationRecord{
		RefId:                        r.RefId,
		LocalId:                      r.LocalId,
		StateProvinceId:              r.StateProvinceId,
		OtherIdList:                  r.OtherIdList,
		IndigenousStatus:             r.IndigenousStatus,
		Sex:                          r.Sex,
		BirthDate:                    r.BirthDate,
		CountryOfBirth:               r.CountryOfBirth,
		StudentLOTE:                  r.StudentLOTE,
		VisaCode:                     r.VisaCode,
		LBOTE:                        r.LBOTE,
		SchoolLocalId:                r.SchoolLocalId,
		YearLevel:                    r.YearLevel,
		FTE:                          r.FTE,
		Parent1LOTE:                  r.Parent1LOTE,
		Parent2LOTE:                  r.Parent2LOTE,
		Parent1Occupation:            r.Parent1Occupation,
		Parent2Occupation:            r.Parent2Occupation,
		Parent1SchoolEducation:       r.Parent1SchoolEducation,
		Parent2SchoolEducation:       r.Parent2SchoolEducation,
		Parent1NonSchoolEducation:    r.Parent1NonSchoolEducation,
		Parent2NonSchoolEducation:    r.Parent2NonSchoolEducation,
		LocalCampusId:                r.LocalCampusId,
		ASLSchoolId:                  r.ASLSchoolId,
		TestLevel:                    r.TestLevel,
		Homegroup:                    r.Homegroup,
		ClassGroup:                   r.ClassGroup,
		MainSchoolFlag:               r.MainSchoolFlag,
		FFPOS:                        r.FFPOS,
		ReportingSchoolId:            r.ReportingSchoolId,
		OtherEnrollmentSchoolACARAId: r.OtherEnrollmentSchoolACARAId,
		EducationSupport:             r.EducationSupport,
		HomeSchooledStudent:          r.HomeSchooledStudent,
		Sensitive:                    r.Sensitive,
		OfflineDelivery:              r.OfflineDelivery,
		NameType:                     r.Name.NameType,
		FamilyName:                   r.Name.FamilyName,
		GivenName:                    r.Name.GivenName,
		MiddleName:                   r.Name.MiddleName,
		PreferredName:                r.Name.PreferredName,
		AddressRole:                  r.Address.Role,
		AddressType:                  r.Address.Type,
		AddressLine1:                 r.Address.Street.Line1,
		AddressLine2:                 r.Address.Street.Line2,
		Locality:                     r.Address.City,
		StateTerritory:               r.Address.StateProvince,
		Postcode:                     r.Address.PostalCode,
	}
	return ret

}

func (r *RegistrationRecord) To_SIF() SIFRegistrationRecord {
	ret := SIFRegistrationRecord{
		RefId:                        r.RefId,
		LocalId:                      r.LocalId,
		StateProvinceId:              r.StateProvinceId,
		OtherIdList:                  r.OtherIdList,
		IndigenousStatus:             r.IndigenousStatus,
		Sex:                          r.Sex,
		BirthDate:                    r.BirthDate,
		CountryOfBirth:               r.CountryOfBirth,
		StudentLOTE:                  r.StudentLOTE,
		VisaCode:                     r.VisaCode,
		LBOTE:                        r.LBOTE,
		SchoolLocalId:                r.SchoolLocalId,
		YearLevel:                    r.YearLevel,
		FTE:                          r.FTE,
		Parent1LOTE:                  r.Parent1LOTE,
		Parent2LOTE:                  r.Parent2LOTE,
		Parent1Occupation:            r.Parent1Occupation,
		Parent2Occupation:            r.Parent2Occupation,
		Parent1SchoolEducation:       r.Parent1SchoolEducation,
		Parent2SchoolEducation:       r.Parent2SchoolEducation,
		Parent1NonSchoolEducation:    r.Parent1NonSchoolEducation,
		Parent2NonSchoolEducation:    r.Parent2NonSchoolEducation,
		LocalCampusId:                r.LocalCampusId,
		ASLSchoolId:                  r.ASLSchoolId,
		TestLevel:                    r.TestLevel,
		Homegroup:                    r.Homegroup,
		ClassGroup:                   r.ClassGroup,
		MainSchoolFlag:               r.MainSchoolFlag,
		FFPOS:                        r.FFPOS,
		ReportingSchoolId:            r.ReportingSchoolId,
		OtherEnrollmentSchoolACARAId: r.OtherEnrollmentSchoolACARAId,
		EducationSupport:             r.EducationSupport,
		HomeSchooledStudent:          r.HomeSchooledStudent,
		Sensitive:                    r.Sensitive,
		OfflineDelivery:              r.OfflineDelivery,
	}
	ret.Name = SIFName{NameType: r.NameType, FamilyName: r.FamilyName,
		GivenName: r.GivenName, MiddleName: r.MiddleName, PreferredName: r.PreferredName}
	ret.Address = Address{
		Role:          r.AddressRole,
		Type:          r.AddressType,
		Street:        Street{Line1: r.AddressLine1, Line2: r.AddressLine2},
		City:          r.Locality,
		StateProvince: r.StateTerritory,
		PostalCode:    r.Postcode,
	}
	return ret
}

type SIFName struct {
	NameType      string `json:",omitempty" xml:"Type,attr"`
	FamilyName    string `json:",omitempty" xml:"FamilyName,omitempty"`
	GivenName     string `json:",omitempty" xml:"GivenName,omitempty"`
	MiddleName    string `json:",omitempty" xml:"MiddleName,omitempty"`
	PreferredName string `json:",omitempty" xml:"PreferredGivenName,omitempty"`
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
		if id.Type == "Locality" {
			r.OtherIdLocality = id.Value
		}
		if id.Type == "DOBRange" {
			r.DOBRange = id.Value
		}
		if id.Type == "PersonalDetailsChanged" {
			r.PersonalDetailsChanged = id.Value
		}
		if id.Type == "PossibleDuplicate" {
			r.PossibleDuplicate = id.Value
		}
		if id.Type == "PsiOtherIdMismatch" {
			r.PsiOtherIdMismatch = id.Value
		}
		if id.Type == "OtherSchoolId" {
			r.OtherSchoolId = id.Value
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
	if r.OtherIdLocality != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"Locality", r.OtherIdLocality})
	}
	if r.PersonalDetailsChanged != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PersonalDetailsChanged", r.PersonalDetailsChanged})
	}
	if r.DOBRange != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"DOBRange", r.DOBRange})
	}
	if r.PossibleDuplicate != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PossibleDuplicate", r.PossibleDuplicate})
	}
	if r.PsiOtherIdMismatch != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"PsiOtherIdMismatch", r.PsiOtherIdMismatch})
	}
	if r.OtherSchoolId != "" {
		r.OtherIdList.OtherId = append(r.OtherIdList.OtherId, XMLAttributeStruct{"OtherSchoolId", r.OtherSchoolId})
	}

	return *r
}

// return key based on concatenation of named fields
// avoiding reflection for performance reasons
func (r RegistrationRecord) FieldsKey(keys []string) string {
	ret := ""
	for _, k := range keys {
		switch k {
		case "RefId":
			ret += r.RefId
		case "LocalId":
			ret += r.LocalId
		case "StateProvinceId":
			ret += r.StateProvinceId
		case "FamilyName":
			ret += r.FamilyName
		case "GivenName":
			ret += r.GivenName
		case "MiddleName":
			ret += r.MiddleName
		case "PreferredName":
			ret += r.PreferredName
		case "IndigenousStatus":
			ret += r.IndigenousStatus
		case "Sex":
			ret += r.Sex
		case "BirthDate":
			ret += r.BirthDate
		case "CountryOfBirth":
			ret += r.CountryOfBirth
		case "StudentLOTE":
			ret += r.StudentLOTE
		case "VisaCode":
			ret += r.VisaCode
		case "LBOTE":
			ret += r.LBOTE
		case "AddressLine1":
			ret += r.AddressLine1
		case "AddressLine2":
			ret += r.AddressLine2
		case "Locality":
			ret += r.Locality
		case "StateTerritory":
			ret += r.StateTerritory
		case "Postcode":
			ret += r.Postcode
		case "SchoolLocalId":
			ret += r.SchoolLocalId
		case "YearLevel":
			ret += r.YearLevel
		case "FTE":
			ret += r.FTE
		case "Parent1LOTE":
			ret += r.Parent1LOTE
		case "Parent2LOTE":
			ret += r.Parent2LOTE
		case "Parent1Occupation":
			ret += r.Parent1Occupation
		case "Parent2Occupation":
			ret += r.Parent2Occupation
		case "Parent1SchoolEducation":
			ret += r.Parent1SchoolEducation
		case "Parent2SchoolEducation":
			ret += r.Parent2SchoolEducation
		case "Parent1NonSchoolEducation":
			ret += r.Parent1NonSchoolEducation
		case "Parent2NonSchoolEducation":
			ret += r.Parent2NonSchoolEducation
		case "LocalCampusId":
			ret += r.LocalCampusId
		case "ASLSchoolId":
			ret += r.ASLSchoolId
		case "TestLevel":
			ret += r.TestLevel
		case "Homegroup":
			ret += r.Homegroup
		case "ClassGroup":
			ret += r.ClassGroup
		case "MainSchoolFlag":
			ret += r.MainSchoolFlag
		case "FFPOS":
			ret += r.FFPOS
		case "ReportingSchoolId":
			ret += r.ReportingSchoolId
		case "OtherSchoolId":
			ret += r.OtherSchoolId
		case "EducationSupport":
			ret += r.EducationSupport
		case "HomeSchooledStudent":
			ret += r.HomeSchooledStudent
		case "Sensitive":
			ret += r.Sensitive
		case "OfflineDelivery":
			ret += r.OfflineDelivery
		case "DiocesanId":
			ret += r.DiocesanId
		case "JurisdictionId":
			ret += r.JurisdictionId
		case "NationalId":
			ret += r.NationalId
		case "OtherId":
			ret += r.OtherId
		case "PlatformId":
			ret += r.PlatformId
		case "PreviousDiocesanId":
			ret += r.PreviousDiocesanId
		case "PreviousNationalId":
			ret += r.PreviousNationalId
		case "PreviousOtherId":
			ret += r.PreviousOtherId
		case "PreviousPlatformId":
			ret += r.PreviousPlatformId
		case "PreviousSectorId":
			ret += r.PreviousSectorId
		case "PreviousLocalId":
			ret += r.PreviousLocalId
		case "PreviousStateProvinceId":
			ret += r.PreviousStateProvinceId
		case "SectorId":
			ret += r.SectorId
		case "TAAId":
			ret += r.TAAId
		}
		ret += ":::"
	}
	return ret
}

// convenience method to return otherid by type
func (r RegistrationRecord) GetOtherId(idtype string) string {

	for _, id := range r.OtherIdList.OtherId {
		if strings.EqualFold(id.Type, idtype) {
			return id.Value
		}
	}

	return ""
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
		"YearLevel",
		"OtherIdLocality",
		"DOBRange",
		"PersonalDetailsChanged",
		"PossibleDuplicate",
		"PsiOtherIdMismatch",
	}
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
		r.GetOtherId("OtherSchoolId"),
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
		r.YearLevel,
		r.GetOtherId("Locality"),
		r.GetOtherId("DOBRange"),
		r.GetOtherId("PersonalDetailsChanged"),
		r.GetOtherId("PossibleDuplicate"),
		r.GetOtherId("PsiOtherIdMismatch"),
	}
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
