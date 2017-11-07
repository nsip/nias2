package naprrql

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/tidwall/gjson"
)

// creates reports based on gql queries
// by school & aggregate reports
//
func GenerateReports() {

	schools, err := getSchoolsList()
	if err != nil {
		log.Fatalln("Cannot connect to naprrql server: ", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = runSchoolReports(schools)
		if err != nil {
			log.Printf("Error creating school reports: ", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = runSystemReports(schools)
		if err != nil {
			log.Printf("Error creating system reports: ", err)
		}
	}()

	wg.Wait()
}

//
// generates a specific 'report' which is the input
// file for isr printing processes
//
func GenerateISRPrintReports() {

	schools, err := getSchoolsList()
	if err != nil {
		log.Fatalln("Cannot connect to naprrql server: ", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = runISRPrintReports(schools)
		if err != nil {
			log.Println("Error creating isr printing reports: ", err)
		}
	}()

	wg.Wait()
}

//
// generates a specific 'report' which is the input
// file for item printing processes
//
func GenerateItemPrintReports() {

	schools, err := getSchoolsList()
	if err != nil {
		log.Fatalln("Cannot connect to naprrql server: ", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = runItemPrintReports(schools)
		if err != nil {
			log.Println("Error creating isr printing reports: ", err)
		}
	}()

	wg.Wait()
}

// generates a specific 'report' which is the input
// file for item printing processes
//
func GenerateQAReports() {

	schools, err := getSchoolsList()
	if err != nil {
		log.Fatalln("Cannot connect to naprrql server: ", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = runQAReports(schools)
		if err != nil {
			log.Println("Error creating QA reports: ", err)
		}
	}()

	wg.Wait()
}

func runISRPrintReports(schools []string) error {

	var pipelineError error

	years := []string{"3", "5", "7", "9"}
	for _, year := range years {
		pipelineError = runISRPipeline(year, schools)

	}

	return pipelineError

}

func runItemPrintReports(schools []string) error {

	var pipelineError error
	pipelineError = runItemPipeline(schools)
	return pipelineError

}

func runQAReports(schools []string) error {

	var pipelineError error
	systemTemplates := getTemplates("./reporting_templates/qa/")
	querymap := make(map[string][]string)
	orphan_queries := make(map[string]string)
	erds_queries := make(map[string]string)
	itemresp_queries := make(map[string]string)
	for filename, query := range systemTemplates {
		// query filenames prefixed with "orphan" need to be run once with their entire acaraIDs argument list,
		// rather than once per acaraID instance
		matched, _ := regexp.MatchString("(/orphan|^orphan)[^/]*$", filename)
		if matched {
			orphan_queries[filename] = query
		} else {
			re := regexp.MustCompile("^.*/")
			filename1 := re.ReplaceAllString(filename, "")
			if filename1 == "systemTestAttempts.gql" ||
				filename1 == "systemParticipationCodeImpacts.gql" ||
				filename1 == "systemTestTypeImpacts.gql" ||
				filename1 == "systemObjectFrequency.gql" ||
				filename1 == "systemTestCompleteness.gql" ||
				filename1 == "systemTestIncidents.gql" {
				erds_queries[filename] = query
			} else if filename1 == "systemTestTypeItemImpacts.gql" ||
				filename1 == "systemItemCounts.gql" ||
				filename1 == "systemParticipationCodeItemImpacts.gql" {
				itemresp_queries[filename] = query
			} else {

				if _, ok := querymap[query]; !ok {
					querymap[query] = make([]string, 0)
				}
				querymap[query] = append(querymap[query], filename)
			}
		}
	}
	if len(erds_queries) > 0 {
		pipelineError = runQAErdsReportPipeline(erds_queries, schools)
	}
	if len(itemresp_queries) > 0 {
		pipelineError = runQAItemRespReportPipeline(itemresp_queries, schools)
	}
	for query := range querymap {
		// log.Printf("Running reports %+v from the query: %+v\n", querymap[query], query)
		pipelineError = runQASystemSingleQueryReportPipeline(query, querymap[query], schools)
	}
	for filename, query := range orphan_queries {
		// log.Printf("Running report %s across all ACARA IDs\n", filename)
		pipelineError = runQASystemReportAllSchoolsPipeline(filename, query, schools)
	}
	return pipelineError

}

//
// produces system (pan-school) reports
//
func runSystemReports(schools []string) error {

	var pipelineError error

	systemTemplates := getTemplates("./reporting_templates/system/")

	for filename, query := range systemTemplates {
		pipelineError = runSystemReportPipeline(filename, query, schools)
	}

	return pipelineError

}

//
// produces per-school reports
//
func runSchoolReports(schools []string) error {

	var pipelineError error

	schoolTemplates := getTemplates("./reporting_templates/school/")

	for _, school := range schools {
		for filename, query := range schoolTemplates {
			pipelineError = runSchoolReportPipeline(filename, query, school)
		}
	}

	return pipelineError

}

//
// connect to the server & retrieve list of known acara-ids for schools
//
func getSchoolsList() ([]string, error) {

	schoolsList := make([]string, 0)
	gqlc := NewGqlClient()

	// create gql query parameters
	query := `query schoolDetails{school_details{ACARAId}}`
	variables := make(map[string]interface{})

	json, err := gqlc.DoQuery(DEF_GQL_URL, query, variables)
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
