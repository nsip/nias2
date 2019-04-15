// naprrqlhp.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/nsip/nias2/naprrql"
	"github.com/shomali11/parallelizer"
	"github.com/tidwall/gjson"
)

// simple test of worker pool to run all pipelines
func main() {

	// shutdown handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		closeDB()
		os.Exit(1)
	}()

	startWebServer(true)
	runReports()

	// close down
	closeDB()
	os.Exit(1)

}

//
// run all report pipelines as parallel jobs
//
func runReports() {

	schools, err := getSchoolsList()
	if err != nil {
		log.Fatalln("Cannot connect to naprrql server: ", err)
	}

	schoolTemplates := getTemplates("./reporting_templates/school/")
	systemTemplates := getTemplates("./reporting_templates/system/")

	log.Println(schools)
	log.Println(len(schoolTemplates))
	log.Println(len(systemTemplates))

	group := parallelizer.NewGroup(parallelizer.WithPoolSize(1000), parallelizer.WithJobQueueSize(500))
	defer group.Close()

	// school reports
	for _, school := range schools {
		for filename, query := range schoolTemplates {
			f := filename
			q := query
			sch := school
			group.Add(func() {
				naprrql.RunSchoolReportPipeline(f, q, sch)
			})
		}
	}

	// system reports
	for filename, query := range systemTemplates {
		f := filename
		q := query
		group.Add(func() {
			naprrql.RunSystemReportPipeline(f, q, schools)
		})
	}

	// item print reports
	group.Add(func() {
		naprrql.RunItemPipeline(schools)
	})

	// writing extract reports
	group.Add(func() {
		var pipelineError error
		naprrql.Run	WritingExtractPipeline(schools, pipelineError)
	})

	// xml pipeline
	// system-level
	group.Add(func() {
		naprrql.RunXMLPipeline(schools)
	})
	// per school
	for _, school := range schools {
		sc := school
		group.Add(func() { naprrql.RunXMLPipelinePerSchool(sc) })
	}

	// qa reports
	group.Add(func() { naprrql.RunQAErdsReportPipeline(schools) })
	group.Add(func() { naprrql.RunQAItemRespReportPipeline(schools) })
	group.Add(func() { naprrql.RunQAMiscReportPipeline(schools) })
	group.Add(func() { naprrql.RunQAOrphanPipeline(schools) })
	group.Add(func() {
		naprrql.RunQASchoolSummaryPipeline(schools, "./out/qa", "./reporting_templates/qa/qaSchools_map.csv")
	})

	err = group.Wait()

	fmt.Println()
	fmt.Println("Done")
	fmt.Printf("Error: %v\n\n", err)
}

//
// connect to the server & retrieve list of known acara-ids for schools
//
func getSchoolsList() ([]string, error) {

	schoolsList := make([]string, 0)
	gqlc := naprrql.NewGqlClient()

	// create gql query parameters
	query := `query schoolDetails{school_details{ACARAId}}`
	variables := make(map[string]interface{})

	json, err := gqlc.DoQuery(naprrql.DEF_GQL_URL, query, variables)
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
// ensure clean shutdown of data store
//
func closeDB() {
	log.Println("Closing datastore...")
	naprrql.GetDB().Close()
	log.Println("Datastore closed.")
}

//
// launch the webserver
//
func startWebServer(silent bool) {
	go naprrql.RunQLServer()
	if !silent {
		fmt.Printf("\n\nBrowse to following locations:\n\n")
		fmt.Printf("\n\tFor qa reporting user interface:\n\n\t\thttp://localhost:1329/ui\n")
		fmt.Printf("\n\tFor data query explorer:\n\n\t\thttp://localhost:1329/sifql\n")
		fmt.Printf("\n\tFor data model viewer:\n\n\t\thttp://localhost:1329/datamodel\n\n ")
	}

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
