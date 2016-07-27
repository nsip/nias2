package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	Nias2 "github.com/nsip/nias2/lib"
	"menteslibres.net/gosexy/rest"
)

var customClient *rest.Client

func TestSexMissingMandatory(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatorySex.csv", "Sex", "Sex is required")
}

func TestSexInvalid(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidSex.csv", "Sex", "Sex must be one of the following")
}

func TestYearLevelPrep(t *testing.T) {
	test_harness(t, "../unit_test_files/1students1YearLevelPrep.csv", "BirthDate/TestLevel", "Year level supplied is P, does not match expected test level")
}

func TestYearLevelF(t *testing.T) {
	test_harness(t, "../unit_test_files/1students2YearLevelF.csv", "BirthDate/YearLevel", "Student Year Level (yr F) does not match year level derived from BirthDate")
}

func TestFutureBirthdate(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsFutureBirthDates.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestMissingParent2LOTE(t *testing.T) {
	test_harness(t, "../unit_test_files/1students2MissingParent2LOTE.csv", "Parent2LOTE", "Parent2LOTE is required")
}

func TestACARAIDandStateBlank(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsACARAIDandStateBlank.csv", "ASLSchoolId", "ASLSchoolId is required")
}

func TestBirthdateYearLevel(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsBirthdateYearLevel.csv", "BirthDate/YearLevel/TestLevel", "does not match year level derived from BirthDate")
}

func TestACARAIDandStateMismatch(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsACARAIDandStateMismatch.csv", "ASLSchoolId", "is a valid ID, but not for")
}

func TestMissingSurname(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingSurname.csv", "FamilyName", "FamilyName is required")
}

func TestEmptySurname(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsEmptySurname.csv", "FamilyName", "FamilyName is required")
}

func TestInvalidVisaClass(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidVisaSubClass.csv", "VisaCode", "is not one of known values from")
}

func TestMalformedPSI(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentMalformedPlatformStudentID.csv", "PlatformId", "PlatformId is not in correct format")
}

func TestCommaAddressField(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsCommaAddressField.csv", "", "")
}

func TestDupGivenLastNameDOBDiffACARAId(t *testing.T) {
	test_harness(t, "../unit_test_files/2studentsDupGivenLastNameDOBDiffACARAId.csv", "", "")
}

func TestExceedCharLengthsSurname(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsExceedCharLengthsSurname.csv", "FamilyName", "String length must be less than or equal to 40")
}

func TestExceedCharLengthsAddress(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsExceedCharLengthsAddress.csv", "AddressLine1", "String length must be less than or equal to 40")
}

func TestExceedCharLengthsGivenName(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsExceedCharLengthsGivenName.csv", "GivenName", "String length must be less than or equal to 40")
}

func TestExceedLengthHomeGrp(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsExceedLengthHomeGrp.csv", "Homegroup", "String length must be less than or equal to 10")
}

func TestInvalidAcaraId(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidAcaraId.csv", "ASLSchoolId", "not found in ASL list of valid IDs")
}

func TestInvalidCountryCodes(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidCountryCodes.csv", "CountryOfBirth", "Country Code is not one of SACC 1269.0 codeset")
}

func TestInvalidDateFormat(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidDateFormat.csv", "BirthDate", "Date provided does not parse correctly for yyyy-mm-dd")
}

func TestInvalidLanguageCodes(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidLanguageCodes.csv", "StudentLOTE", "Language Code is not one of ASCL 1267.0 codeset")
}

func TestInvalidValuesLBOTE(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesLBOTE.csv", "LBOTE", "LBOTE must be one of the following")
}

func TestInvalidValuesOfflineDelivery(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesOfflineDelivery.csv", "OfflineDelivery", "OfflineDelivery must be one of the following")
}

func TestInvalidValuesMainSchoolFlag(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesMainSchoolFlag.csv", "MainSchoolFlag", "MainSchoolFlag must be one of the following")
}

func TestInvalidValuesParent1LOTE(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidParent1LOTE.csv", "Parent1LOTE", "Language Code is not one of ASCL 1267.0 codeset")
}

func TestInvalidInvalidValuesFFPOS(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesFFPOS.csv", "FFPOS", "FFPOS must be one of the following")
}

func TestInvalidValuesParent2Occupation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesParent2Occupation.csv", "Parent2Occupation", "Parent2Occupation must be one of the following")
}

func TestInvalidValuesParent2NonSchoolEducation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesParent2NonSchoolEducation.csv", "Parent2NonSchoolEducation", "Parent2NonSchoolEducation must be one of the following")
}

func TestInvalidValuesParent2SchoolEducation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesParent2SchoolEducation.csv", "Parent2SchoolEducation", "Parent2SchoolEducation must be one of the following")
}

func TestInvalidValuesHomeSchooledStudent(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesHomeSchooledStudent.csv", "HomeSchooledStudent", "HomeSchooledStudent must be one of the following")
}

func TestInvalidValuesYearLevel(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsInvalidValuesYearLevel.csv", "YearLevel", "YearLevel must be one of the following")
}

func TestMissingMandatoryParent1SchlEduc(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent1SchlEduc.csv", "Parent1SchoolEducation", "Parent1SchoolEducation is required")
}

func TestMissingMandatoryASLSchoolID(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryASLSchoolID.csv", "ASLSchoolId", "ASLSchoolId is required")
}

func TestMissingMandatoryDOB(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryDOB.csv", "BirthDate", "BirthDate is required")
}

func TestMissingMandatoryFFPOS(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryFFPOS.csv", "FFPOS", "FFPOS is required")
}

func TestMissingMandatoryGivenName(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryGivenName.csv", "GivenName", "GivenName is required")
}

func TestMissingMandatoryIndigenousStatus(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryIndigenousStatus.csv", "IndigenousStatus", "IndigenousStatus is required")
}

func TestMissingMandatoryLocalId(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryLocalId.csv", "LocalId", "LocalId is required")
}

func TestMissingMandatoryParent1LOTE(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent1LOTE.csv", "Parent1LOTE", "Parent1LOTE is required")
}

func TestMissingMandatoryParent1NonSchoolEducation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent1NonSchlEduc.csv", "Parent1NonSchoolEducation", "Parent1NonSchoolEducation is required")
}

func TestMissingMandatoryParent1Occupation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent1Occupation.csv", "Parent1Occupation", "Parent1Occupation is required")
}

func TestMissingMandatoryParent2LOTE(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent2LOTE.csv", "Parent2LOTE", "Parent2LOTE is required")
}

func TestMissingMandatoryParent2NonSchoolEducation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent2NonSchoolEduc.csv", "Parent2NonSchoolEducation", "Parent2NonSchoolEducation is required")
}

func TestMissingMandatoryParent2Occupation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent2Occupation.csv", "Parent2Occupation", "Parent2Occupation is required")
}

func TestMissingMandatoryParent2SchoolEducation(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryParent2SchlEduc.csv", "Parent2SchoolEducation", "Parent2SchoolEducation is required")
}

func TestMissingMandatoryCountryOfBirth(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryStudentCountryBirth.csv", "CountryOfBirth", "CountryOfBirth is required")
}

func TestMissingMandatoryStudentLOTE(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryStudentLOTE.csv", "StudentLOTE", "StudentLOTE is required")
}

func TestMissingMandatoryTestLevel(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryTestLevel.csv", "TestLevel", "TestLevel is required")
}

func TestMissingMandatoryYearLevel(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatoryYearLevel.csv", "YearLevel", "YearLevel is required")
}

func TestOptionalMissing(t *testing.T) {
	test_harness(t, "../unit_test_files/4studentsOptionalMissing.csv", "", "")
}

func TestOutsideAgeRange3(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsOutsideAgeRange3.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRange5(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsOutsideAgeRange5.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRange7(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsOutsideAgeRange7.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRange9(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsOutsideAgeRange9.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRangeUG(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsOutsideAgeRangeUG.csv", "BirthDate/TestLevel/YearLevel", "Year level supplied is UG, will result in SRM warning flag for test level")
}

func TestSameStudentIdTwoDifferentSchoolId(t *testing.T) {
	test_harness(t, "../unit_test_files/2studentsSameStudentIdTwoDifferentSchoolId.csv", "", "")
}

func TestUngradedValuesUGJunSec(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsUngradedValuesUGJunSec.csv", "BirthDate/TestLevel/YearLevel", "Year level supplied is UG, will result in SRM warning flag for test level")
}

func TestUngradedValuesUGSnrSec(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsUngradedValuesUGSnrSec.csv", "BirthDate/TestLevel/YearLevel", "Year level supplied is UG, will result in SRM warning flag for test level")
}

func TestUngradedValuesUGPri(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsUngradedValuesUGPri.csv", "BirthDate/TestLevel/YearLevel", "Year level supplied is UG, will result in SRM warning flag for test level")
}

func TestUngradedValuesUGSec(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsUngradedValuesUGSec.csv", "BirthDate/TestLevel/YearLevel", "Year level supplied is UG, will result in SRM warning flag for test level")
}

func TestUnusualCountryCodes(t *testing.T) {
	test_harness(t, "../unit_test_files/3studentsUnusualCountryCodes.csv", "", "")
}

func TestUnusualLanguageCodes(t *testing.T) {
	test_harness(t, "../unit_test_files/3studentsUnusualLanguageCodes.csv", "", "")
}

func TestWrongChecksumPlatformId(t *testing.T) {
	test_harness(t, "../unit_test_files/1WrongChecksumPlatformId.csv", "PlatformId", "Platform ID has incorrect checksum")
}

func TestWrongChecksumPreviousPlatformId(t *testing.T) {
	test_harness(t, "../unit_test_files/1WrongChecksumPrevPlatformId.csv", "PreviousPlatformId", "Previous Platform ID has incorrect checksum")
}

func TestYearLevelTestLevelMismatch(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsYearLevelTestLevelMismatch.csv", "BirthDate/TestLevel", "does not match year level derived from BirthDate")
}

func TestDupGivenLastNameDOBCARAId(t *testing.T) {
	for i := 0; i < 20; i++ {
		test_harness(t, "../unit_test_files/2studentsDupGivenLastNameDOBSchool.csv", "Multiple (see description)", "otential duplicate of record")
	}
}

func TestDuplicateStudentOneSchool(t *testing.T) {
	for i := 0; i < 20; i++ {
		test_harness(t, "../unit_test_files/2studentsDuplicateStudentOneSchool.csv", "LocalID/ASL ID", "otential duplicate of record")
	}
}

/* if errfield is nil, we expect test to pass */
func test_harness(t *testing.T, filename string, errfield string, errdescription string) {
	var f *os.File
	var err error
	bytebuf := []byte{}
	dat := []map[string]string{}

	if f, err = os.Open(filename); err != nil {
		t.Fatalf("Error %s", err)
	}
	defer f.Close()
	files := rest.FileMap{
		"validationFile": []rest.File{{
			Name:   path.Base(f.Name()),
			Reader: f},
		},
	}
	requestVariables := url.Values{"name": {path.Base(f.Name())}}
	msg, err := rest.NewMultipartMessage(requestVariables, files)
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	dst := map[string]interface{}{}
	if err = customClient.PostMultipart(&dst, "/naplan/reg/validate", msg); err != nil {
		t.Fatalf("Error %s", err)
	}
	txid := dst["TxID"].(string)
	time.Sleep(100 * time.Millisecond)
	if err = customClient.Get(&bytebuf, "/naplan/reg/results/"+txid, nil); err != nil {
		t.Fatalf("Error %s", err)
	}
	// we are getting back a JSON array
	if err = json.Unmarshal(bytebuf, &dat); err != nil {
		t.Fatalf("Error %s", err)
	}
	log.Println(dat)
	if errfield == "" {
		if len(dat) > 0 {
			t.Fatalf("Expected no error, got error in %s: %s", dat[0]["errfield"], dat[0]["description"])
		}
	} else {
		if len(dat) < 1 {
			t.Fatalf("Expected error field %s, got no error", errfield)
		} else {
			if dat[0]["errField"] != errfield {
				t.Fatalf("Expected error field %s, got field %s", errfield, dat[0]["errField"])
			}
			if !strings.Contains(dat[0]["description"], errdescription) {
				t.Fatalf("Expected error description %s, got description %s", errdescription, dat[0]["description"])
			}
		}
	}
}

func TestMain(m *testing.M) {
	customClient, _ = rest.New("http://localhost:" + Nias2.NiasConfig.WebServerPort + "/")
	os.Exit(m.Run())
}
