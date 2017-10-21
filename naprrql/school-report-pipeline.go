// school-report-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"strings"
)

//
// report pipleines to generate per-school reports
//
//
// create and run a school-level report printing pipeline.
//
// Pipeline streams requests a given report, school by school
// feeding results to the ouput csv file.
//
// This means the server & parser never have to deal with query data
// volumes larger than a single school at a time.
//
// Overall round-trip latency is less than querying for all data at once
// and ensures we can't run out of memory
//
func runSchoolReportPipeline(queryFileName string, query string, school string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// wrap single school in array to reuse system report pipeline transformers
	schools := []string{school}

	// input stage
	// check if system query needs to iterate schools, if not
	// pass single dummy school to fire the query
	if !strings.Contains(query, "$acaraIDs") {
		schools = []string{"no-op"}
	}
	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("systemQueryExecutor()")
	// transform stage
	jsonc, errc, err := systemQueryExecutor(ctx, query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("send to sink")
	// sink stage
	// create working directory if not there
	outFileDir := "./out/school_reports/" + school
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

	log.Println("School report file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}
