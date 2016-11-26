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

	Nias2 "github.com/nsip/nias2/lib"
	"github.com/siddontang/goredis"
	"menteslibres.net/gosexy/rest"
)

func TestPrivacy(t *testing.T) {
	test_harness_filecomp_privacy_xml(t, "../unit_test_files/1students.xml")
}
func TestSMS(t *testing.T) {
	test_harness_sms(t, "../unit_test_files/StudentPersonal.xml", "../unit_test_files/1StudentPersonal_Graph.json")
}
func TestSMS1(t *testing.T) {
	test_harness_sms(t, "../unit_test_files/Sif3AssessmentRegistration.xml", "../unit_test_files/1Sif3AssessmentRegistration_Graph.json")
}
func TestSif2Graph_StudentPersonal(t *testing.T) {
	sif2graph_harness(t, "../unit_test_files/StudentPersonal.xml", "../unit_test_files/1StudentPersonal_Graph.json")
}

func TestSif2Graph_Sif3AssessmentRegistration(t *testing.T) {
	sif2graph_harness(t, "../unit_test_files/Sif3AssessmentRegistration.xml", "../unit_test_files/1Sif3AssessmentRegistration_Graph.json")
}

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

/* compare two files */
func test_harness_filecomp_privacy_xml(t *testing.T, filename string) {
	var err error
	var sensitivities = [5]string{"none", "low", "medium", "high", "extreme"}

	bytebuf := []byte{}
	dat := []string{}

	txid, err := post_file(filename, "/sifxml/ingest")
	errcheck(t, err)
	log.Printf("txid: %s", txid)
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < len(sensitivities); i++ {
		err = customClient.Get(&bytebuf, "/sifxml/ingest/"+sensitivities[i]+"/"+txid, nil)
		errcheck(t, err)
		// we are getting back a JSON array
		err = json.Unmarshal(bytebuf, &dat)
		errcheck(t, err)
		log.Println(strings.Split(dat[0], "\n")[0])
		err = compare_files(strings.Join(dat, "\n"), filename+"."+sensitivities[i])
		errcheck(t, err)
	}

}

// compare the retrieved file in retvalue to the file in filename
func compare_files(retvalue string, filename string) error {
	var err error
	var re *regexp.Regexp
	dat1 := []byte{}

	dat1, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if re, err = regexp.Compile("(\\n|^)\\s+"); err != nil {
		return err
	}
	dat1 = re.ReplaceAll(dat1, []byte("\n"))
	retvalue1 := re.ReplaceAll([]byte(retvalue), []byte("\n"))
	if bytes.Compare(dat1, retvalue1) != 0 {
		return fmt.Errorf("output does not match file %s:\n=====\n%s\n====\n%s\n====\n", filename, string(dat1), string(retvalue1))
	}
	return nil
}

func sif2graph_harness(t *testing.T, filename string, json_filename string) {
	s2g, err := Nias2.NewSif2GraphService()
	errcheck(t, err)
	dat, err := ioutil.ReadFile(filename)
	errcheck(t, err)
	r := Nias2.NiasMessage{}
	r.TxID = nuid.Next()
	r.SeqNo = "1"
	dat1 := strings.Split(string(dat), "\n")
	r.Body = strings.Join(dat1[1:len(dat1)-2], "\n")
	ret, err := s2g.HandleMessage(&r)
	errcheck(t, err)
	if len(ret) < 1 {
		t.Fatalf("Error %s", err)
	}
	jsondat, err := ioutil.ReadFile(json_filename)
	errcheck(t, err)
	graphstruct := Nias2.GraphStruct{}
	err = json.Unmarshal(jsondat, &graphstruct)
	errcheck(t, err)
	log.Println(ret[0].Body.(Nias2.GraphStruct))
	if !reflect.DeepEqual(graphstruct, ret[0].Body.(Nias2.GraphStruct)) {
		t.Fatalf("Mapping of %s to SMS graph format did not match %s:\n%s", filename, json_filename, ret[0].Body.(Nias2.GraphStruct))
	}

}

// clear keys to be interrogated from redis, for preexisting runs
func clear_redis(t *testing.T, graphstruct Nias2.GraphStruct, ms *Nias2.MessageStore) {
	_, err := ms.C.Do("hdel", "labels", graphstruct.Guid)
	errcheck(t, err)
	_, err = ms.C.Do("srem", graphstruct.Type, graphstruct.Guid)
	errcheck(t, err)
	if localid, ok := graphstruct.OtherIds["LocalId"]; ok {
		_, err := ms.C.Do("hdel", "oid:"+localid, "LocalId")
		errcheck(t, err)
		_, err = ms.C.Do("srem", "other:ids", "oid:"+localid)
		errcheck(t, err)
	}
	if len(graphstruct.EquivalentIds) > 0 {
		_, err := ms.C.Do("srem", "equivalent:id"+graphstruct.EquivalentIds[0], graphstruct.Guid)
		errcheck(t, err)
	}
	if len(graphstruct.Links) > 0 {
		_, err := ms.C.Do("srem", graphstruct.Links[0], graphstruct.Guid)
		errcheck(t, err)
	}
}

func test_harness_sms(t *testing.T, filename string, json_filename string) {
	var err error

	//bytebuf := []byte{}
	//dat := []string{}

	log.Println("Starting")
	ms := Nias2.NewMessageStore()
	jsondat, err := ioutil.ReadFile(json_filename)
	errcheck(t, err)
	graphstruct := Nias2.GraphStruct{}
	err = json.Unmarshal(jsondat, &graphstruct)
	errcheck(t, err)
	Nias2.PrefixGraphStruct(&graphstruct, Nias2.SIF_MEMORY_STORE_PREFIX)
	clear_redis(t, graphstruct, ms)

	log.Println("Stage 1")
	_, err = post_file(filename, "/sifxml/store")
	errcheck(t, err)
	time.Sleep(2000 * time.Millisecond)
	log.Println("Stage 1a")

	label, err := goredis.String(ms.C.Do("hget", "labels", graphstruct.Guid))
	//errcheck(t, err)
	log.Println("Stage 1b")
	if label != graphstruct.Label {
		t.Fatalf("Label retrieved for %s is not expected %s, but %s", graphstruct.Guid, graphstruct.Label, label)
	}
	log.Printf("Label retrieved for %s is expected %s", graphstruct.Guid, label)

	log.Println("Stage 2")
	_, err = post_file(filename, "/sifxml/store")
	membership, err := goredis.Bool(ms.C.Do("sismember", graphstruct.Type, graphstruct.Guid))
	errcheck(t, err)
	if !membership {
		t.Fatalf("Guid %s is not member of set %s", graphstruct.Guid, graphstruct.Type)
	}
	log.Printf("Guid %s is member of set %s", graphstruct.Guid, graphstruct.Type)

	log.Println("Stage 3")
	if otherid, ok := graphstruct.OtherIds["LocalId"]; ok {
		localid, err := goredis.String(ms.C.Do("hget", "oid:"+otherid, "LocalId"))
		errcheck(t, err)
		if localid != graphstruct.Guid {
			t.Fatalf("oid:(%s){LocalId} does not map to %s", otherid, graphstruct.Guid)
		}
		log.Printf("oid:(%s){LocalId} does map to %s", otherid, graphstruct.Guid)

		membership, err = goredis.Bool(ms.C.Do("sismember", "other:ids", "oid:"+otherid))
		errcheck(t, err)
		if !membership {
			t.Fatalf("Local Id %s is not a member of set %s", otherid, "other:ids")
		}
		log.Printf("Local Id %s is a member of set %s", otherid, "other:ids")
	}

	log.Println("Stage 4")
	if len(graphstruct.Links) > 0 {
		membership, err = goredis.Bool(ms.C.Do("sismember", graphstruct.Links[0], graphstruct.Guid))
		errcheck(t, err)
		if !membership {
			t.Fatalf("Guid %s is not a member of set %s", graphstruct.Guid, graphstruct.Links[0])
		}
		log.Printf("Guid %s is not a member of set %s", graphstruct.Guid, graphstruct.Links[0])
	}
	log.Println("No error")
}
