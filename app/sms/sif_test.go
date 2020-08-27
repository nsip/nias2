package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"
	//"strings"
	//"errors"
	"bytes"
	"fmt"
	"github.com/nats-io/nuid"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gosexy/rest"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/sms"
	"github.com/nsip/nias2/xml"
	"github.com/siddontang/goredis"
)

var customClient *rest.Client

func errcheck(t *testing.T, err error, num int) {
	if err != nil {
		t.Fatalf("Error %d: %s", num, err)
	}
}

func TestPrivacy(t *testing.T) {
	test_harness_filecomp_privacy_xml(t, "../../unit_test_files/1students.xml")
}
func TestSMS(t *testing.T) {
	test_harness_sms(t, "../../unit_test_files/StudentPersonal.xml", "../../unit_test_files/1StudentPersonal_Graph.json")
}
func TestSMS1(t *testing.T) {
	test_harness_sms(t, "../../unit_test_files/1napstudentresponseset.xml", "../../unit_test_files/1napstudentresponseset_Graph.json")
}
func TestSif2Graph_StudentPersonal(t *testing.T) {
	sif2graph_harness(t, "../../unit_test_files/StudentPersonal.xml", "../../unit_test_files/1StudentPersonal_Graph.json")
}

func TestSif2Graph_1TeachingGroup(t *testing.T) {
	sif2graph_harness(t, "../../unit_test_files/1napstudentresponseset.xml", "../../unit_test_files/1napstudentresponseset_Graph.json")
}

func post_file(filename string, endpoint string) (string, error) {
	var f *os.File
	var err error
	if f, err = os.Open(filename); err != nil {
		log.Println("POST FILE 1")
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
		log.Println("POST FILE 2")
		return "", err
	}
	dst := map[string]interface{}{}
	log.Printf("%v\n", msg)
	if err = customClient.PostMultipart(&dst, endpoint, msg); err != nil {
		log.Println("POST FILE 3")
		return "", err
	}
	txid := dst["TxID"].(string)
	return txid, nil
}

/* compare two files */
func test_harness_filecomp_privacy_xml(t *testing.T, filename string) {
	var err error
	var sensitivities = [5]string{"none", "low", "medium", "high", "extreme"}

	bytebuf := []byte{}
	dat := []string{}

	txid, err := post_file(filename, "/sifxml/ingest")
	errcheck(t, err, 1)
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < len(sensitivities); i++ {
		err = customClient.Get(&bytebuf, "/sifxml/ingest/"+sensitivities[i]+"/"+txid, nil)
		errcheck(t, err, 2)
		// we are getting back a JSON array
		err = json.Unmarshal(bytebuf, &dat)
		log.Println(dat)
		errcheck(t, err, 3)
		err = compare_files(strings.Join(dat, "\n"), filename+"."+sensitivities[i])
		errcheck(t, err, 4)
	}
	log.Printf("%v\n", err)

}

// compare the retrieved file in retvalue to the file in filename
func compare_files(retvalue string, filename string) error {
	var err error
	var re, re2 *regexp.Regexp
	dat1 := []byte{}

	dat1, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if re, err = regexp.Compile("(\\n|^)\\s+"); err != nil {
		return err
	}
	if re2, err = regexp.Compile("[\\n|\\s]+$"); err != nil {
		return err
	}
	dat1 = re.ReplaceAll(dat1, []byte("\n"))
	retvalue1 := re.ReplaceAll([]byte(retvalue), []byte("\n"))
	dat1 = re2.ReplaceAll(dat1, []byte(""))
	retvalue1 = re2.ReplaceAll(retvalue1, []byte(""))
	log.Printf("COMP\n%v####%v####\n", string(dat1), string(retvalue1))
	if bytes.Compare(dat1, retvalue1) != 0 {
		return fmt.Errorf("output does not match file %s:\n=====\n%s\n====\n%s\n====\n", filename, string(dat1), string(retvalue1))
	}
	return nil
}

func sif2graph_harness(t *testing.T, filename string, json_filename string) {
	s2g, err := sms.NewSif2GraphService()
	errcheck(t, err, 5)
	dat, err := ioutil.ReadFile(filename)
	errcheck(t, err, 6)
	r := lib.NiasMessage{}
	r.TxID = nuid.Next()
	r.SeqNo = "1"
	dat1 := strings.Split(string(dat), "\n")
	r.Body = strings.Join(dat1[1:len(dat1)-2], "\n")
	ret, err := s2g.HandleMessage(&r)
	errcheck(t, err, 7)
	if len(ret) < 1 {
		t.Fatalf("Error %s", err)
	}
	jsondat, err := ioutil.ReadFile(json_filename)
	errcheck(t, err, 8)
	graphstruct := xml.GraphStruct{}
	err = json.Unmarshal(jsondat, &graphstruct)
	errcheck(t, err, 9)
	/*
		graphstruct1 := xml.GraphStruct{}
		err = json.Unmarshal(ret[0].Body.([]byte), &graphstruct1)
	*/
	graphstruct1 := ret[0].Body.(xml.GraphStruct)
	if !reflect.DeepEqual(graphstruct, graphstruct1) {
		t.Fatalf("Mapping of %s to SMS graph format did not match %s:\n%s", filename, json_filename, ret[0].Body.(xml.GraphStruct))
	}

}

// clear keys to be interrogated from redis, for preexisting runs
func clear_redis(t *testing.T, graphstruct xml.GraphStruct, ms *sms.Store) {
	_, err := ms.Ledis.Do("hdel", "labels", graphstruct.Guid)
	errcheck(t, err, 10)
	_, err = ms.Ledis.Do("srem", graphstruct.Type, graphstruct.Guid)
	errcheck(t, err, 11)
	if localid, ok := graphstruct.OtherIds["LocalId"]; ok {
		_, err := ms.Ledis.Do("hdel", "oid:"+localid, "LocalId")
		errcheck(t, err, 12)
		_, err = ms.Ledis.Do("srem", "other:ids", "oid:"+localid)
		errcheck(t, err, 13)
	}
	if len(graphstruct.EquivalentIds) > 0 {
		_, err := ms.Ledis.Do("srem", "equivalent:id"+graphstruct.EquivalentIds[0], graphstruct.Guid)
		errcheck(t, err, 14)
	}
	if len(graphstruct.Links) > 0 {
		_, err := ms.Ledis.Do("srem", graphstruct.Links[0], graphstruct.Guid)
		errcheck(t, err, 15)
	}
}

func test_harness_sms(t *testing.T, filename string, json_filename string) {
	var err error

	//bytebuf := []byte{}
	//dat := []string{}

	ms := sms.NewStore()
	jsondat, err := ioutil.ReadFile(json_filename)
	errcheck(t, err, 16)
	graphstruct := xml.GraphStruct{}
	err = json.Unmarshal(jsondat, &graphstruct)
	errcheck(t, err, 17)
	sms.PrefixGraphStruct(&graphstruct, sms.SIF_MEMORY_STORE_PREFIX)
	clear_redis(t, graphstruct, ms)

	_, err = post_file(filename, "/sifxml/store")
	errcheck(t, err, 18)
	time.Sleep(300 * time.Millisecond)

	label, err := goredis.String(ms.Ledis.Do("hget", "labels", graphstruct.Guid))
	//errcheck(t, err, 19)
	if label != graphstruct.Label {
		t.Fatalf("Label retrieved for %s is not expected %s, but %s", graphstruct.Guid, graphstruct.Label, label)
	}
	log.Printf("Label retrieved for %s is expected %s", graphstruct.Guid, label)

	_, err = post_file(filename, "/sifxml/store")
	membership, err := goredis.Bool(ms.Ledis.Do("sismember", graphstruct.Type, graphstruct.Guid))
	errcheck(t, err, 20)
	if !membership {
		t.Fatalf("Guid %s is not member of set %s", graphstruct.Guid, graphstruct.Type)
	}
	log.Printf("Guid %s is member of set %s", graphstruct.Guid, graphstruct.Type)

	if otherid, ok := graphstruct.OtherIds["LocalId"]; ok {
		localid, err := goredis.String(ms.Ledis.Do("hget", "oid:"+otherid, "LocalId"))
		errcheck(t, err, 21)
		if localid != graphstruct.Guid {
			t.Fatalf("oid:(%s){LocalId} does not map to %s", otherid, graphstruct.Guid)
		}
		log.Printf("oid:(%s){LocalId} does map to %s", otherid, graphstruct.Guid)

		membership, err = goredis.Bool(ms.Ledis.Do("sismember", "other:ids", "oid:"+otherid))
		errcheck(t, err, 22)
		if !membership {
			t.Fatalf("Local Id %s is not a member of set %s", otherid, "other:ids")
		}
		log.Printf("Local Id %s is a member of set %s", otherid, "other:ids")
	}

	log.Println("Stage 4")
	if len(graphstruct.Links) > 0 {
		membership, err = goredis.Bool(ms.Ledis.Do("sismember", graphstruct.Links[0], graphstruct.Guid))
		errcheck(t, err, 23)
		if !membership {
			t.Fatalf("Guid %s is not a member of set %s", graphstruct.Guid, graphstruct.Links[0])
		}
		log.Printf("Guid %s is not a member of set %s", graphstruct.Guid, graphstruct.Links[0])
	}
	log.Println("No error")
}

func TestMain(m *testing.M) {
	config := lib.LoadDefaultConfig()
	customClient, _ = rest.New("http://localhost:" + config.WebServerPort + "/")
	os.Exit(m.Run())
}
