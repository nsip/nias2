package main

import (
	"os"
	"path"
	"testing"
	"time"
	//"bytes"
	//"strconv"
	//"time"
	"encoding/json"
	"log"
	"net/url"
	"strings"

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
	test_harness(t, "../unit_test_files/1studentsACARAIDandStateMismatch.csv", "ASLSchoolID", "is a valid ID, but not for")
}

func TestMissingSurname(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingSurname.csv", "FamilyName", "FamilyName is required")
}

func TestEmptySurname(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsEmptySurname.csv", "FamilyName", "FamilyName is required")
}

func test_harness(t *testing.T, filename string, errfield string, errdescription string) {
	var f *os.File
	var err error

	if f, err = os.Open(filename); err != nil {
		t.Fatalf("Error %s", err)
	}
	defer f.Close()
	files := rest.FileMap{
		"validationFile": []rest.File{
			{
				Name:   path.Base(f.Name()),
				Reader: f,
			},
		},
	}
	requestVariables := url.Values{
		"name": {path.Base(f.Name())},
	}
	msg, err := rest.NewMultipartMessage(requestVariables, files)
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	dst := map[string]interface{}{}
	if err = customClient.PostMultipart(&dst, "/naplan/reg/validate", msg); err != nil {
		t.Fatalf("Error %s", err)
	}
	txid := dst["TxID"].(string)
	bytebuf := []byte{}
	//dat := []map[string]interface{}{}
	dat := []map[string]string{}
	time.Sleep(2 * time.Second)
	if err = customClient.Get(&bytebuf, "/naplan/reg/results/"+txid, nil); err != nil {
		t.Fatalf("Error %s", err)
	}
	// we are getting back a JSON array
	if err = json.Unmarshal(bytebuf, &dat); err != nil {
		t.Fatalf("Error %s", err)
	}
	log.Println(dat)
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

	/*
		for _, r := range records {
			log.Println(r)
			natsconn.Ec.Publish("validation.naplan", r)
		}
		msg, err := sub.NextMsg(5 * time.Second)
		if err != nil {
			t.Fatalf("Error %s", err)
		}
		dat := make(map[string]string)
		if err := json.Unmarshal(msg.Data, &dat); err != nil {
			t.Fatalf("Error unmarshalling json message: %s", err)
		}
		log.Println(dat)
		if dat["errField"] != errfield {
			t.Fatalf("Expected error field %s, got field %s", errfield, dat["errField"])
		}
		if !strings.Contains(dat["description"], errdescription) {
			t.Fatalf("Expected error description %s, got description %s", errdescription, dat["description"])
		}
	*/
}

func TestMain(m *testing.M) {
	customClient, _ = rest.New("http://localhost:" + Nias2.NiasConfig.WebServerPort + "/")
	os.Exit(m.Run())
}
