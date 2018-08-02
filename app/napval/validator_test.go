package main

import (
	"bytes"
	"encoding/json"
	//"log"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/nsip/nias2/napval"
	"menteslibres.net/gosexy/rest"
)

var customClient *rest.Client

func post_file(filename string, endpoint string) (string, error) {
	var f *os.File
	var err error
	if f, err = os.Open(filename); err != nil {
		return "", err
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
		return "", err
	}
	dst := map[string]interface{}{}
	if err = customClient.PostMultipart(&dst, endpoint, msg); err != nil {
		return "", err
	}
	txid := dst["TxID"].(string)
	return txid, nil
}

func TestRepeatPSIwithinSchool(t *testing.T) {
	test_harness(t, "../../unit_test_files/5StudentsDuplicatePlatformID.csv", "PSI/ASL ID", "Platform Student ID (Student) and ASL ID (School) are potential duplicate")
}

func TestAddress(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsAddress.csv", "AddressLine1", "Additional property AddressLine1 is not allowed")
}

func TestState(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsState.csv", "StateTerritory", "Additional property StateTerritory is not allowed")
}

func TestKurdistan(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsKurdistan.csv", "", "")
}

func TestSexMissingMandatory(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatorySex.csv", "Sex", "Sex is required")
}

func TestSexInvalid(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidSex.csv", "Sex", "Sex must be one of the following")
}

func TestYearLevelPrep(t *testing.T) {
	test_harness(t, "../../unit_test_files/1students1YearLevelPrep.csv", "BirthDate/TestLevel", "Year Level P does not match Test level")
}

func TestYearLevelF(t *testing.T) {
	test_harness(t, "../../unit_test_files/1students2YearLevelF.csv", "BirthDate/TestLevel", "Year Level F does not match Test level")
}

func TestYearLevelP(t *testing.T) {
	test_harness(t, "../../unit_test_files/1students1YearLevelP.csv", "BirthDate/YearLevel", "")
}

func TestFutureBirthdate(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsFutureBirthDates.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestMissingParent2LOTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1students2MissingParent2LOTE.csv", "Parent2LOTE", "Must be present if other Parent2 fields are present")
}

func TestACARAIDandStateBlank(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsACARAIDandStateBlank.csv", "ASLSchoolId", "ASLSchoolId is required")
}

func TestBirthdateYearLevel(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsBirthdateYearLevel.csv", "BirthDate/YearLevel/TestLevel", "does not match year level derived from BirthDate")
}

/*
func TestACARAIDandStateMismatch(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsACARAIDandStateMismatch.csv", "ASLSchoolId", "is a valid ID, but not for")
}
*/

func TestMissingSurname(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingSurname.csv", "FamilyName", "FamilyName is required")
}

func TestEmptySurname(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsEmptySurname.csv", "FamilyName", "FamilyName is required")
}

func TestInvalidVisaClass(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidVisaSubClass.csv", "VisaCode", "is not one of known values from")
}

func TestMalformedPSI(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentMalformedPlatformStudentID.csv", "PlatformId", "PlatformId is not in correct format")
}

func TestCommaAddressField(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsCommaAddressField.csv", "", "")
}

func TestDupGivenLastNameDOBDiffACARAId(t *testing.T) {
	test_harness(t, "../../unit_test_files/2studentsDupGivenLastNameDOBDiffACARAId.csv", "Multiple (see description)", "Potential duplicate")
}

func TestExceedCharLengthsSurname(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsExceedCharLengthsSurname.csv", "FamilyName", "String length must be less than or equal to 40")
}

/*
func TestExceedCharLengthsAddress(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsExceedCharLengthsAddress.csv", "AddressLine1", "String length must be less than or equal to 40")
}
*/

func TestExceedCharLengthsGivenName(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsExceedCharLengthsGivenName.csv", "GivenName", "String length must be less than or equal to 40")
}

func TestInvalidAcaraId(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidAcaraId.csv", "ASLSchoolId", "not found in ASL list of valid IDs")
}

func TestInvalidCountryCodes(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidCountryCodes.csv", "CountryOfBirth", "Country Code is not one of SACC 1269.0 codeset")
}

func TestInvalidDateFormat(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidDateFormat.csv", "BirthDate", "Date provided does not parse correctly for yyyy-mm-dd")
}

func TestInvalidLanguageCodes(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidLanguageCodes.csv", "StudentLOTE", "Language Code is not one of ASCL 1267.0 codeset")
}

func TestInvalidValuesLBOTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesLBOTE.csv", "LBOTE", "LBOTE must be one of the following")
}

func TestInvalidValuesOfflineDelivery(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesOfflineDelivery.csv", "OfflineDelivery", "OfflineDelivery must be one of the following")
}

func TestInvalidValuesMainSchoolFlag(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesMainSchoolFlag.csv", "MainSchoolFlag", "MainSchoolFlag must be one of the following")
}

func TestInvalidValuesParent1LOTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidParent1LOTE.csv", "Parent1LOTE", "Language Code is not one of ASCL 1267.0 codeset")
}

func TestInvalidInvalidValuesFFPOS(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesFFPOS.csv", "FFPOS", "FFPOS must be one of the following")
}

func TestInvalidValuesParent2Occupation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesParent2Occupation.csv", "Parent2Occupation", "Parent2Occupation must be one of the following")
}

func TestInvalidValuesParent2NonSchoolEducation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesParent2NonSchoolEducation.csv", "Parent2NonSchoolEducation", "Parent2NonSchoolEducation must be one of the following")
}

func TestInvalidValuesParent2SchoolEducation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesParent2SchoolEducation.csv", "Parent2SchoolEducation", "Parent2SchoolEducation must be one of the following")
}

func TestInvalidValuesHomeSchooledStudent(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesHomeSchooledStudent.csv", "HomeSchooledStudent", "HomeSchooledStudent must be one of the following")
}

func TestInvalidValuesYearLevel(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsInvalidValuesYearLevel.csv", "YearLevel", "YearLevel must be one of the following")
}

func TestMissingMandatoryParent1SchlEduc(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent1SchlEduc.csv", "Parent1SchoolEducation", "Parent1SchoolEducation is required")
}

func TestMissingMandatoryASLSchoolID(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryASLSchoolID.csv", "ASLSchoolId", "ASLSchoolId is required")
}

func TestMissingMandatoryDOB(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryDOB.csv", "BirthDate", "BirthDate is required")
}

func TestMissingMandatoryFFPOS(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryFFPOS.csv", "FFPOS", "FFPOS is required")
}

func TestMissingMandatoryGivenName(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryGivenName.csv", "GivenName", "GivenName is required")
}

func TestMissingMandatoryIndigenousStatus(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryIndigenousStatus.csv", "IndigenousStatus", "IndigenousStatus is required")
}

func TestMissingMandatoryLocalId(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryLocalId.csv", "LocalId", "LocalId is required")
}

func TestMissingMandatoryParent1LOTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent1LOTE.csv", "Parent1LOTE", "Parent1LOTE is required")
}

func TestMissingMandatoryParent1NonSchoolEducation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent1NonSchlEduc.csv", "Parent1NonSchoolEducation", "Parent1NonSchoolEducation is required")
}

func TestMissingMandatoryParent1Occupation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent1Occupation.csv", "Parent1Occupation", "Parent1Occupation is required")
}

func TestMissingMandatoryParent2LOTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent2LOTE.csv", "Parent2LOTE", "Must be present if other Parent2 fields are present")
}

func TestMissingMandatoryParent2NonSchoolEducation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent2NonSchoolEduc.csv", "Parent2NonSchoolEducation", "Must be present if other Parent2 fields are present")
}

func TestMissingMandatoryParent2Occupation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent2Occupation.csv", "Parent2Occupation", "Must be present if other Parent2 fields are present")
}

func TestMissingMandatoryParent2SchoolEducation(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryParent2SchlEduc.csv", "Parent2SchoolEducation", "Must be present if other Parent2 fields are present")
}

func TestMissingMandatoryCountryOfBirth(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryStudentCountryBirth.csv", "CountryOfBirth", "CountryOfBirth is required")
}

func TestMissingMandatoryStudentLOTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryStudentLOTE.csv", "StudentLOTE", "StudentLOTE is required")
}

func TestMissingMandatoryTestLevel(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryTestLevel.csv", "TestLevel", "TestLevel is required")
}

func TestMissingMandatoryYearLevel(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsMissingMandatoryYearLevel.csv", "YearLevel", "YearLevel is required")
}

func TestOptionalMissing(t *testing.T) {
	test_harness(t, "../../unit_test_files/4studentsOptionalMissing.csv", "", "")
}

func TestOutsideAgeRange3(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsOutsideAgeRange3.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRange5(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsOutsideAgeRange5.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRange7(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsOutsideAgeRange7.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRange9(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsOutsideAgeRange9.csv", "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestOutsideAgeRangeUG(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsOutsideAgeRangeUG.csv", "BirthDate/TestLevel", "Year level supplied is UG, will result in SRM warning flag for test level")
}

func TestSameStudentIdTwoDifferentSchoolId(t *testing.T) {
	test_harness(t, "../../unit_test_files/2studentsSameStudentIdTwoDifferentSchoolId.csv", "", "")
}

func TestUngradedValuesUGJunSec(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsUngradedValuesUGJunSec.csv", "BirthDate/YearLevel", "Student Year Level (yr UGJunSec) does not match year level derived from BirthDate")
}

func TestUngradedValuesUGSnrSec(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsUngradedValuesUGSnrSec.csv", "BirthDate/YearLevel", "Student Year Level (yr UGSnrSec) does not match year level derived from BirthDate")
}

func TestUngradedValuesUGPri(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsUngradedValuesUGPri.csv", "BirthDate/YearLevel", "Student Year Level (yr UGPri) does not match year level derived from BirthDate")
}

func TestUngradedValuesUGSec(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsUngradedValuesUGSec.csv", "BirthDate/YearLevel", "Student Year Level (yr UGSec) does not match year level derived from BirthDate")
}

func TestUnusualCountryCodes(t *testing.T) {
	test_harness(t, "../../unit_test_files/3studentsUnusualCountryCodes.csv", "", "")
}

func TestUnusualLanguageCodes(t *testing.T) {
	test_harness(t, "../../unit_test_files/3studentsUnusualLanguageCodes.csv", "", "")
}

func TestWrongChecksumPlatformId(t *testing.T) {
	test_harness(t, "../../unit_test_files/1WrongChecksumPlatformId.csv", "PlatformId", "Platform ID has incorrect checksum")
}

func TestWrongChecksumPreviousPlatformId(t *testing.T) {
	test_harness(t, "../../unit_test_files/1WrongChecksumPrevPlatformId.csv", "PreviousPlatformId", "Previous Platform ID has incorrect checksum")
}

func TestYearLevelTestLevelMismatch(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsYearLevelTestLevelMismatch.csv", "BirthDate/TestLevel", "does not match Test level ")
}

/*
func TestExtraneousNotPermittedField(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsExtraneousNotPermittedField.csv", "BirthDate/TestLevel", "does not match year level derived from BirthDate")
}
*/

func TestExtraneousPermittedField(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsExtraneousPermittedField.csv", "PossibleDuplicate", "Additional property PossibleDuplicate is not allowed")
}

func TestMaximumFTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1MaximumFTE.csv", "FTE", "FTE is greater than 1")
}

func TestMinimumFTE(t *testing.T) {
	test_harness(t, "../../unit_test_files/1MinimumFTE.csv", "FTE", "FTE is less than 0")
}

func TestDupGivenLastNameDOBCARAId(t *testing.T) {
	for i := 0; i < 20; i++ {
		test_harness(t, "../../unit_test_files/2studentsDupGivenLastNameDOBSchool.csv", "Multiple (see description)", "otential duplicate of record")
	}
}

func TestDuplicateStudentOneSchool(t *testing.T) {
	for i := 0; i < 20; i++ {
		test_harness(t, "../../unit_test_files/2studentsDuplicateStudentOneSchool.csv", "LocalID/ASL ID", "otential duplicate of record")
	}
}

func TestSuspectCharactersFamilyName(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsSuspectCharacterFamilyName.csv", "FamilyName", "Family Name contains suspect character")
}

func TestSuspectCharactersGivenName(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsSuspectCharacterGivenName.csv", "GivenName", "Given Name contains suspect character")
}

func TestSuspectCharactersMiddleName(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsSuspectCharacterMiddleName.csv", "MiddleName", "Middle Name contains suspect character")
}

func TestSuspectCharactersPreferredName(t *testing.T) {
	test_harness(t, "../../unit_test_files/1studentsSuspectCharacterPreferredName.csv", "PreferredName", "Preferred Name contains suspect character")
}

func errcheck(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error %s", err)
	}
}

/* if errfield is nil, we expect test to pass */
func test_harness(t *testing.T, filename string, errfield string, errdescription string) {
	var err error
	bytebuf := []byte{}
	dat := map[string]string{}

	txid, err := post_file(filename, "naplan/reg/validate")
	errcheck(t, err)
	time.Sleep(100 * time.Millisecond)

	err = customClient.Get(&bytebuf, "naplan/reg/results/"+txid, nil)
	errcheck(t, err)
	//log.Println("+++" + string(bytebuf))

	// we are getting back line delimited JSON
	lines := bytes.Split(bytebuf, []byte{'\n'})
	// grab first error

	//log.Println(">>>" + string(lines[0]) + "<<<")
	if len(lines[0]) == 0 {
		err = nil
	} else {
		//err = json.Unmarshal(bytebuf, &dat)
		err = json.Unmarshal(lines[0], &dat)
	}
	errcheck(t, err)
	// we are getting back a JSON array
	/*
		for i := 0; i < len(lines); i++ {
			log.Println("\t" + string(lines[i]))
			log.Println(dat)
		}
	*/
	if errfield == "" {
		if len(dat) > 0 {
			t.Fatalf("Expected no error, got error in %s: %s", dat["errfield"], dat["description"])
		}
	} else {
		if len(dat) < 1 {
			t.Fatalf("Expected error field %s, got no error", errfield)
		} else {
			if dat["errField"] != errfield {
				t.Fatalf("Expected error field %s, got field %s", errfield, dat["errField"])
			}
			if !strings.Contains(dat["description"], errdescription) {
				t.Fatalf("Expected error description %s, got description %s", errdescription, dat["description"])
			}
		}
	}
}

func TestMain(m *testing.M) {
	config := napval.LoadNAPLANConfig()
	customClient, _ = rest.New("http://localhost:" + config.WebServerPort + "/")
	os.Exit(m.Run())
}
