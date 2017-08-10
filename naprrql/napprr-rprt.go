package naprrql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/matt-farmer/json2csv"
	"github.com/tidwall/gjson"
)

var queryURL = "http://localhost:1329/graphql"

// creates reports based on gql queries
// by school & aggregate reports
//
func GenerateReports() {

	// test connection to server & retrieve schools list
	schools, err := getSchoolsList()
	if err != nil {
		log.Fatalln("Cannot connect to naprrql server: ", err)
	}

	err = runSchoolReports(schools)
	if err != nil {
		log.Printf("Error creating school reports: ", err)
	}

	err = runSystemReports(schools)
	if err != nil {
		log.Printf("Error creating system reports: ", err)
	}

}

//
// produces system (pan-school) reports
//
func runSystemReports(schools []string) error {

	systemTemplates := getTemplates("./system_templates/")

	variables := make(map[string]interface{})
	variables["acaraIDs"] = schools

	for filename, query := range systemTemplates {

		reportFilename := filepath.Base(filename)
		csvFilename := strings.Replace(reportFilename, ".gql", ".csv", 1)

		json, err := gqlQuery(query, variables)
		if err != nil {
			return err
		}

		results, err := json2csv.JSON2CSV(json.Value())
		if err != nil {
			return err
		}
		if len(results) == 0 {
			return nil
		}

		outFileDir := "./out"
		outFileName := outFileDir + "/" + csvFilename
		err = os.MkdirAll(outFileDir, os.ModePerm)
		outFile, err := os.Create(outFileName)
		defer outFile.Close()
		if err != nil {
			return err
		}

		headerStyle := json2csv.DotNotationStyle
		err = printCSV(outFile, results, headerStyle, false)
		if err != nil {
			return err
		}
		fmt.Println("Reporting file written: ", outFileName)

	}

	return nil

}

//
// produces per-school reports
//
func runSchoolReports(schools []string) error {

	schoolTemplates := getTemplates("./school_templates/")

	for _, school := range schools {

		variables := make(map[string]interface{})

		for filename, query := range schoolTemplates {

			reportFilename := filepath.Base(filename)
			csvFilename := strings.Replace(reportFilename, ".gql", ".csv", 1)

			variables["acaraIDs"] = school
			json, err := gqlQuery(query, variables)
			if err != nil {
				return err
			}

			results, err := json2csv.JSON2CSV(json.Value())
			if err != nil {
				return err
			}
			if len(results) == 0 {
				return nil
			}

			outFileDir := "./out/" + school
			outFileName := outFileDir + "/" + csvFilename
			err = os.MkdirAll(outFileDir, os.ModePerm)
			outFile, err := os.Create(outFileName)
			defer outFile.Close()
			if err != nil {
				return err
			}

			headerStyle := json2csv.DotNotationStyle
			err = printCSV(outFile, results, headerStyle, false)
			if err != nil {
				return err
			}
			fmt.Println("Reporting file written: ", outFileName)
		}
	}

	return nil

}

//
// make a gql query and return the json payload
//
func gqlQuery(query string, variables map[string]interface{}) (*gjson.Result, error) {

	// assemble the request & client to invoke it
	gqlReq := GQLRequest{Query: query, Variables: variables}
	jsonStr, err := json.Marshal(gqlReq)
	req, err := http.NewRequest("POST", queryURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	// connect to webserver, exp. backoff to ensure
	// server has time to come up before querying
	var resp *http.Response
	operation := func() error {
		resp, err = client.Do(req)
		return err // or an error
	}
	err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		// if we got here connection was made but there's a
		// different error, so pass up the stack
		return &gjson.Result{}, err
	}
	defer resp.Body.Close()

	// process the returned json, content taken from the
	// element in the returned payload identified in the query
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &gjson.Result{}, err
	}
	resultRoot := getDataPath(query)
	response := gjson.GetBytes(body, resultRoot)

	return &response, nil

}

//
// connect to the server & retrieve list of known acara-ids for schools
//
func getSchoolsList() ([]string, error) {

	schoolsList := make([]string, 0)

	// create gql query parameters
	query := `query schoolDetails{school_details{ACARAId}}`
	variables := make(map[string]interface{})

	json, err := gqlQuery(query, variables)
	if err != nil {
		return schoolsList, err
	}
	acaraids := json.Get("#.ACARAId")
	acaraids.ForEach(func(key, value gjson.Result) bool {
		schoolsList = append(schoolsList, value.String())
		return true // keep iterating
	})

	return schoolsList, nil

}

//
// read json response
//
func readJSON(r io.Reader) (interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

//
// parses a gql query to find the root node at which the
// data array will be returned
//
func getDataPath(query string) string {

	defRoot := "data"

	tokens := strings.SplitAfterN(query, "{", 3)

	// not an array query, flat members e.g. object counts
	if len(tokens) < 3 {
		return defRoot
	}

	// remove variables clause from first line of query
	// to find root element
	line1_tokens := strings.SplitN(tokens[1], "{", 2)
	query_tokens := strings.SplitN(line1_tokens[0], "(", 2)

	queryRoot := strings.TrimSpace(query_tokens[0])

	return defRoot + "." + queryRoot

}

//
// write out csv file
//
func printCSV(w io.Writer, results []json2csv.KeyValue, headerStyle json2csv.KeyStyle, transpose bool) error {
	csv := json2csv.NewCSVWriter(w)
	csv.HeaderStyle = headerStyle
	csv.Transpose = transpose
	if err := csv.WriteCSV(results); err != nil {
		return err
	}
	return nil
}

//
// retrieves templates for reporting queries
// in effect these are just the graphQL queries saved to files
//
func getTemplates(templatesPath string) map[string]string {

	files := make([]string, 0)
	templates := make(map[string]string)

	// get the list of query templates
	gqlFiles, err := filepath.Glob(templatesPath + "*.gql")
	if err != nil {
		log.Fatalln("Unable to read system report template files: ", err)
	}
	files = append(files, gqlFiles...)
	if len(files) == 0 {
		log.Fatalln("No template (*.gql) files found in input folder " + templatesPath)
	}

	// store template against filename
	for _, file := range files {
		query, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalln("Unable to load template query from file: ", err)
		}
		templates[file] = string(query)
	}

	return templates

}
