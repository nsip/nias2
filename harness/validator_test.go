package main

import (
	"os"
	//"io/ioutil"
	"testing"
	//"bytes"
        //"strconv"
	//"time"
        //"encoding/json"
	"log"
	//"strings"

	//"github.com/nats-io/nats"
	//lib "github.com/nsip/nias-go-naplan-registration/lib"
        //"github.com/wildducktheories/go-csv"
        "github.com/dghubble/sling"
	Nias2 "github.com/nsip/nias2/lib"

)

// var web_ec *nats.EncodedConn

/*
func csv2nats(csvstring string) ([]map[string]string, error) {
	// input csv, and output a slice of records to push to NATS. the same way that aggregator Post("/naplan/reg/:stateID") does
	reader := csv.WithIoReader(ioutil.NopCloser(bytes.NewReader([]byte(csvstring))))
	records, err := csv.ReadAll(reader)
	ret := make([]map[string]string, len(records))
	for i, r := range records {
		ret[i] = r.AsMap()
		ret[i]["OriginalLine"] = strconv.Itoa(i + 1)
		ret[i]["TxID"] = "dummyTxID"
	}
	return ret, err
}
*/

func TestSexMissingMandatory(t *testing.T) {
	test_harness(t, "../unit_test_files/1studentsMissingMandatorySex.csv", "Sex", "Sex must be one of the following")
}
/*
func TestSexInvalid(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1studentsInvalidSex.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "Sex", "Sex must be one of the following")
}

func TestYearLevelPrep(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1students1YearLevelPrep.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "BirthDate/TestLevel", "Year level supplied is P, does not match expected test level")
}

func TestYearLevelF(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1students2YearLevelF.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "BirthDate/YearLevel", "Student Year Level (yr F) does not match year level derived from BirthDate")
}

func TestFutureBirthdate(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1studentsFutureBirthDates.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "BirthDate/YearLevel", "Year Level calculated from BirthDate does not fall within expected NAPLAN year level ranges")
}

func TestMissingParent2LOTE(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1students2MissingParent2LOTE.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "Parent2LOTE", "Language Code is not one of ASCL 1267.0 codeset")
}

func TestACARAIDandStateBlank(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1studentsACARAIDandStateBlank.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "StateTerritory", "StateTerritory must be one of the following")
}

func TestBirthdateYearLevel(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1studentsBirthdateYearLevel.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "BirthDate/YearLevel/TestLevel", "does not match year level derived from BirthDate")
}

func TestACARAIDandStateMismatch(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1studentsACARAIDandStateMismatch.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "ASLSchoolId", "is a valid ID, but not for")
}

func TestMissingSurname(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1studentsMissingSurname.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "FamilyName", "FamilyName is required")
}

func TestEmptySurname(t *testing.T) {
	data, err := ioutil.ReadFile("../unit_test_files/1studentsEmptySurname.csv")
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	test_harness(t, string(data), "ASLSchoolId", "is a valid ID, but not for")
}
*/

type Issue struct {
    txid  string `json:"TxID"`
    msgid   string `json:"MsgID"`
}


func test_harness(t *testing.T, filename string, errfield string, errdescription string) {
/*
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Error %s", err)
	}
*/
        port := Nias2.NiasConfig.WebServerPort
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	issues := new([]Issue)
	resp, err := sling.New().Post(":"+port+"/naplan/reg/validate").Body(f).ReceiveSuccess(issues)
	log.Println(resp)
	log.Println()
	log.Println(issues)



/*
	records, err := csv2nats(csv)
	if err != nil {
		t.Fatalf("Error %s", err)
	}
	log.Printf("records received: %v", len(records))
	sub, err := natsconn.Nc.SubscribeSync("validation.errors") 
	if err != nil {
		t.Fatalf("Error %s", err)
	}
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
	os.Exit(m.Run())
}
