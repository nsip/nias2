// qa-misc-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"
)

//
// printing pipeline for QA reports on ACARA ID orphan records
// (records that do not link to a school mentioned in the list of
// ACARA IDs that reports iterate over)
//

func runQAOrphanPipeline(schools []string) error {
	var pipelineError error
	systemTemplates := getTemplates("./reporting_templates/qa/")
	orphan_queries := make(map[string]string)
	for filename, query := range systemTemplates {
		// query filenames prefixed with "orphan" need to be run once with their entire acaraIDs argument list,
		// rather than once per acaraID instance
		matched, _ := regexp.MatchString(`([/\\]orphan|^orphan)[^/\\]*$`, filename)
		if matched {
			orphan_queries[filename] = query
		}
	}
	for filename, query := range orphan_queries {
		pipelineError = runQASystemReportAllSchoolsPipeline(filename, query, schools)
	}
	return pipelineError
}

// Run the query once with the argument schools
func runQASystemReportAllSchoolsPipeline(queryFileName string, query string, schools []string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	// check if system query needs to iterate schools, if not
	// pass single dummy school to fire the query
	if !strings.Contains(query, "$acaraIDs") {
		schools = []string{"no-op"}
	}

	if err := systemParametersValidate(schools); err != nil {
		return err
	}

	// transform stage
	jsonc, errc, err := systemQueryExecutorNoBreakdown(ctx, query, DEF_GQL_URL, map[string]interface{}{"acaraIDs": schools})
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/qa"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	csvFileName := deriveCSVFileName(queryFileName)
	outFileName := outFileDir + "/" + csvFileName
	mapFileName := deriveMapFileName(queryFileName)
	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("QA file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}
