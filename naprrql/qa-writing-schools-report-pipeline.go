// qa-writing-schools-report-pipeline.go

package naprrql

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	//"github.com/nsip/nias2/xml"
	//"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

//
// reporting pipeline to produce qa reports
// summarising school attributes; system type, no. students, general info
//

//
// create and run an qa school summary report printing pipeline specific to writing.
//
// Pipeline streams requests school by school
// feeding results to the ouput processors and eventaully csv file.
//
// This means the server & parser never have to deal with query data
// volumes larger than a single school at a time.
//
// Overall round-trip latency is less than querying for all data at once
// and ensures we can't run out of memory
//
func runQAWritingSchoolSummaryPipeline(schools []string, outFileDir string, mapFileName string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage - emits schools one at a time into pipeline
	varsc, errc, err := qaSchoolParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// load query templates from disk
	var query string
	queryTemplates := getTemplates("./reporting_templates/qa/")

	//
	// transform query stage 1 : get the school info
	//
	for name, queryText := range queryTemplates {
		if strings.Contains(name, "qaSchoolInfo.gql") {
			query = queryText
		}
	}
	qasummc1, errc, err := qaSchoolQueryExecutor(ctx, query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	//
	// transform stage 2: do school system classification
	//
	qasummc2, errc, err := qaSchoolClassifier(ctx, qasummc1)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// transform stage 3: get student registration numbers [by year]
	//
	for name, queryText := range queryTemplates {
		if strings.Contains(name, "qaStudentsByYrLevel.gql") {
			query = queryText
		}
	}
	qasummc3, errc, err := qaWritingYrLevelQueryExecutor(ctx, query, DEF_GQL_URL, qasummc2)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// transform stge 4: get actual test event numbers [by year]
	//
	for name, queryText := range queryTemplates {
		if strings.Contains(name, "qaTestAttempts.gql") {
			query = queryText
		}
	}
	qasummc4, errc, err := qaWritingAttemptsQueryExecutor(ctx, query, DEF_GQL_URL, qasummc3)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// transform stage 5: get disruptions, adjustments, exemptions etc.
	//
	for name, queryText := range queryTemplates {
		if strings.Contains(name, "qaParticipation.gql") {
			query = queryText
		}
	}
	// this is the only difference from the qa report: filters by Writing
	qasummc5, errc, err := qaWritingParticipationQueryExecutor(ctx, query, DEF_GQL_URL, qasummc4)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// transform stage 6: flatten summary & convert to gjson for writing to csv
	//
	jsonc, errc, err := qaTransformSummary(ctx, qasummc5)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// sink stage - write to csv file
	//
	// create working directory if not there
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	csvFileName := "qaSchools.csv"
	outFileName := outFileDir + "/" + csvFileName
	// mapFileName := "./reporting_templates/qa/qaSchools_map.csv"
	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("QA Schools Summary report file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

//
// receives qa summary as input, summaries participation codes for test attempts, and disruptions
//
func qaWritingParticipationQueryExecutor(ctx context.Context, query string, url string, in <-chan schoolQASummary) (<-chan schoolQASummary, <-chan error, error) {

	out := make(chan schoolQASummary)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for school := range in {
			vars := map[string]interface{}{"acaraIDs": []string{school.ACARAId}}
			json, err := gql.DoQuery(url, query, vars)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}

			school.AttemptParticipation = make(map[string]int)
			for _, result := range json.Array() {
				// get participation codes
				iter := result.Get("Summary")
				iter.ForEach(func(_, value gjson.Result) bool {
					domain := value.Get("Domain").String()
					if domain != "Writing" {
						return true // skip this value
					}
					key := value.Get("ParticipationCode").String()
					school.AttemptParticipation[key]++
					return true // keep iterating
				})
				// get disruptions
				iter = result.Get("EventInfos")
				iter.ForEach(func(_, value gjson.Result) bool {
					disr := value.Get("Event.TestDisruptionList.TestDisruption").Array()
					school.AttemptDisruptions += len(disr)
					return true // keep iterating
				})
			}

			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- school:
			case <-ctx.Done():
				return
			}

		}
	}()
	return out, errc, nil

}

//
// receives qa summary, tranforms into gjson object csv writer can work with
//
func qaWritingTransformSummary(ctx context.Context, in <-chan schoolQASummary) (<-chan gjson.Result, <-chan error, error) {

	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)

		for school := range in {

			// comvert summary to gjson
			jsonSchool, err := json.Marshal(school)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}

			result := gjson.ParseBytes(jsonSchool)

			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- result:
			case <-ctx.Done():
				return
			}

		}
	}()
	return out, errc, nil

}

//
// receives school qa summary & collects test attempts by yr level & domain
//
func qaWritingAttemptsQueryExecutor(ctx context.Context, query string, url string, in <-chan schoolQASummary) (<-chan schoolQASummary, <-chan error, error) {

	out := make(chan schoolQASummary)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for school := range in {
			vars := map[string]interface{}{"acaraIDs": []string{school.ACARAId}}
			json, err := gql.DoQuery(url, query, vars)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}

			school.TestAttempts = make(map[string]int)
			for _, result := range json.Array() {
				participation := result.Get("Event.ParticipationCode").String()
				domain := result.Get("Test.TestContent.TestDomain").String()
				if domain == "Writing" && participation == "P" {
					key := result.Get("Test.TestContent.TestLevel").String() + ":" + domain
					school.TestAttempts[key]++
					school.TotalAttempts++
				}
			}

			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- school:
			case <-ctx.Done():
				return
			}

		}
	}()
	return out, errc, nil

}

//
// receives school qa summary, uses to retrieve student yr level counts
//
func qaWritingYrLevelQueryExecutor(ctx context.Context, query string, url string, in <-chan schoolQASummary) (<-chan schoolQASummary, <-chan error, error) {

	out := make(chan schoolQASummary)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for school := range in {
			vars := map[string]interface{}{"acaraIDs": []string{school.ACARAId}}
			json, err := gql.DoQuery(url, query, vars)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}

			for _, result := range json.Array() {
				switch result.Get("YearLevel").String() {
				case "3":
					school.Yr3registered++
				case "5":
					school.Yr5registered++
					school.TotalStudents++
				case "7":
					school.Yr7registered++
					school.TotalStudents++
				case "9":
					school.Yr9registered++
					school.TotalStudents++
				default:
					school.YrUnknowRegistered++
					school.TotalStudents++
				}
				switch result.Get("TestLevel").String() {
				case "3":
					school.TestLvl3registered++
				case "5":
					school.TestLvl5registered++
				case "7":
					school.TestLvl7registered++
				case "9":
					school.TestLvl9registered++
				default:
					school.TestLvlUnknowRegistered++
				}
			}

			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- school:
			case <-ctx.Done():
				return
			}

		}
	}()
	return out, errc, nil

}
