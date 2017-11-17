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
// printing pipeline for miscellaneous QA reports
//

func runQAMiscReportPipeline(schools []string) error {
	var pipelineError error
	systemTemplates := getTemplates("./reporting_templates/qa/")
	querymap := make(map[string][]string)
	for filename, query := range systemTemplates {
		re := regexp.MustCompile(`^.*[/\\]`)
		filename1 := re.ReplaceAllString(filename, "")
		matched, _ := regexp.MatchString(`^orphan[^/\\]*$`, filename1)
		if matched || filename1 == "systemTestAttempts.gql" ||
			filename1 == "systemParticipationCodeImpacts.gql" ||
			filename1 == "systemTestTypeImpacts.gql" ||
			filename1 == "systemObjectFrequency.gql" ||
			filename1 == "systemTestCompleteness.gql" ||
			filename1 == "systemTestIncidents.gql" ||
			filename1 == "systemTestTypeItemImpacts.gql" ||
			filename1 == "itemPrinting.gql" ||
			filename1 == "itemWritingPrinting.gql" ||
			filename1 == "systemItemCounts.gql" ||
			filename1 == "systemParticipationCodeItemImpacts.gql" ||
			filename1 == "qaSchoolInfo.gql" ||
			filename1 == "qaStudentsByYrLevel.gql" ||
			filename1 == "qaTestAttempts.gql" ||
			filename1 == "qaParticipation.gql" ||
			filename1 == "systemRubricSubscoreMatches.gql" {
			// ignore, they're addressed by one of the other pipelines
		} else {

			if _, ok := querymap[query]; !ok {
				querymap[query] = make([]string, 0)
			}
			querymap[query] = append(querymap[query], filename)
		}
	}
	for query := range querymap {
		pipelineError = runQASystemSingleQueryReportPipeline(query, querymap[query], schools)
	}
	return pipelineError
}

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
