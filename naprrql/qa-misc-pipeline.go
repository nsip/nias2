// qa-misc-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"strings"
)

//
// printing pipeline for miscellaneous QA reports
//

// Run a single query for each school, then run a series of reports off that query
func runQASystemSingleQueryReportPipeline(query string, queryFileNames []string, schools []string) error {

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
	varsc, errc, err := systemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	jsonc, errc, err := systemQueryExecutor(ctx, query, DEF_GQL_URL, varsc)
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
	jsonc1, errc, err := splitter(ctx, jsonc, len(queryFileNames))
	errcList = append(errcList, errc)
	for i, queryFileName := range queryFileNames {
		csvFileName := deriveCSVFileName(queryFileName)
		outFileName := outFileDir + "/" + csvFileName
		mapFileName := deriveMapFileName(queryFileName)
		errc1, err := csvFileSink(ctx, outFileName, mapFileName, jsonc1[i])
		if err != nil {
			return err
		}
		errcList = append(errcList, errc1)

		log.Println("QA file writing... " + outFileName)
	}
	return WaitForPipeline(errcList...)
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
