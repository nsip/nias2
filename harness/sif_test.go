package main

import (
	"encoding/json"
	//"log"
	"net/url"
	"os"
	"path"
	//"strings"
	//"errors"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
	"time"

	// Nias2 "github.com/nsip/nias2/lib"
	"menteslibres.net/gosexy/rest"
)

//var customClient *rest.Client

func TestPrivacy(t *testing.T) {
	test_harness_filecomp_privacy_xml(t, "../unit_test_files/1students.xml")
}

/* compare two files */
func test_harness_filecomp_privacy_xml(t *testing.T, filename string) {
	var f *os.File
	var err error
	var sensitivities = [4]string{"low", "medium", "high", "extreme"}

	bytebuf := []byte{}
	dat := []string{}

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
	if err = customClient.PostMultipart(&dst, "/sifxml/ingest", msg); err != nil {
		t.Fatalf("Error %s", err)
	}
	txid := dst["TxID"].(string)
	time.Sleep(200 * time.Millisecond)

	for i := 0; i < len(sensitivities); i++ {
		if err = customClient.Get(&bytebuf, "/sifxml/ingest/"+sensitivities[i]+"/"+txid, nil); err != nil {
			t.Fatalf("Error %s", err)
		}
		// we are getting back a JSON array
		if err = json.Unmarshal(bytebuf, &dat); err != nil {
			t.Fatalf("Error %s", err)
		}
		if err = compare_files(strings.Join(dat, "\n"), filename+"."+sensitivities[i]); err != nil {
			t.Fatalf("Error %s", err)
		}
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
