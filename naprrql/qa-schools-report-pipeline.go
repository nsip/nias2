// qa-schools-report-pipeline.go

package naprrql

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/nsip/nias2/xml"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

//
// reporting pipeline to produce qa reports
// summarising school attributes; system type, no. students, general info
//

//
// helper type to consolidate data
//
type schoolQASummary struct {
	xml.SchoolInfo
	TestAttempts         map[string]int // internal map structure where key is composite "YrLevel:Domain" eg. "3:Numeracy"
	AttemptParticipation map[string]int // map of participation codes with summary totals for each
	TotalAttempts        int
	AttemptDisruptions   int
	Yr3registered        int
	Yr5registered        int
	Yr7registered        int
	Yr9registered        int
	YrUnknowRegistered   int
	TotalStudents        int
	DerivedSector        string
	DerivedSystem        string
}

//
// create and run an qa school summary report printing pipeline.
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
func runQASchoolSummaryPipeline(schools []string) error {

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
	qasummc3, errc, err := qaYrLevelQueryExecutor(ctx, query, DEF_GQL_URL, qasummc2)
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
	qasummc4, errc, err := qaAttemptsQueryExecutor(ctx, query, DEF_GQL_URL, qasummc3)
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
	qasummc5, errc, err := qaParticipationQueryExecutor(ctx, query, DEF_GQL_URL, qasummc4)
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
	outFileDir := "./out/qa"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	csvFileName := "qaSchools.csv"
	outFileName := outFileDir + "/" + csvFileName
	mapFileName := "./reporting_templates/qa/qaSchools_map.csv"
	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("QA Schools Summary report file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

//
// acts as input feed to the pipeline, sends parameters to retrieve data for
// each school in turn.
//
func qaSchoolParametersSource(ctx context.Context, schools ...string) (<-chan systemQueryParams, <-chan error, error) {

	{ //check input variables, handle errors before goroutine starts
		if len(schools) == 0 {
			return nil, nil, errors.Errorf("no schools provided")
		}
	}

	out := make(chan systemQueryParams)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for schoolIndex, school := range schools {
			if school == "" {
				// Handle an error that occurs during the goroutine.
				errc <- errors.Errorf("school %v is empty string", schoolIndex+1)
				return
			}
			vars := systemQueryParams{schoolAcaraID: school}
			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- vars:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}

//
// query executor transform stage takes query params in, excutes gql query
// and writes results to output chaneel
//
func qaSchoolQueryExecutor(ctx context.Context, query string, url string, in <-chan systemQueryParams) (<-chan schoolQASummary, <-chan error, error) {

	out := make(chan schoolQASummary)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for params := range in {
			vars := map[string]interface{}{"acaraIDs": []string{params.schoolAcaraID}}
			jsonResp, err := gql.DoQuery(url, query, vars)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			for _, result := range jsonResp.Array() {

				sqaSumm := schoolQASummary{}
				if err := json.Unmarshal([]byte(result.Raw), &sqaSumm); err != nil {
					errc <- err
					return
				}

				// Send the data to the output channel but return early
				// if the context has been cancelled.
				select {
				case out <- sqaSumm:
				case <-ctx.Done():
					return
				}

			}
		}
	}()
	return out, errc, nil
}

//
// receives json encoded schoolinfo, uses to construct classification of school
//
func qaSchoolClassifier(ctx context.Context, in <-chan schoolQASummary) (<-chan schoolQASummary, <-chan error, error) {

	out := make(chan schoolQASummary)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for school := range in {
			deriveClassifications(&school)
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
// derive the school system calssification from the school info
//
func deriveClassifications(school *schoolQASummary) {

	// derive sector
	sector := school.SchoolSector
	switch {
	case strings.EqualFold(sector, "ng"):
		school.DerivedSector = "Non-Government"
	default:
		school.DerivedSector = "Government"
		return
	}

	// non gov systemic?
	// no further checks required if non-systemic
	ngsystemic := school.NonGovSystemicStatus
	if !strings.EqualFold(ngsystemic, "s") {
		return
	}

	// if non-gov systemic...
	// derive system
	system := school.System
	switch system {
	case "0001":
		school.DerivedSystem = "Catholic"
	case "0002":
		school.DerivedSystem = "Anglican"
	case "0003":
		school.DerivedSystem = "Lutheran"
	case "0004":
		school.DerivedSystem = "Seventh Day Adventist"
	default:
		school.DerivedSystem = "Other"
	}

}

//
// receives school qa summary, uses to retrieve student yr level counts
//
func qaYrLevelQueryExecutor(ctx context.Context, query string, url string, in <-chan schoolQASummary) (<-chan schoolQASummary, <-chan error, error) {

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
				case "7":
					school.Yr7registered++
				case "9":
					school.Yr9registered++
				default:
					school.YrUnknowRegistered++
				}
				school.TotalStudents++
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
// receives school qa summary & collects test attempts by yr level & domain
//
func qaAttemptsQueryExecutor(ctx context.Context, query string, url string, in <-chan schoolQASummary) (<-chan schoolQASummary, <-chan error, error) {

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
				key := result.Get("Test.TestContent.TestLevel").String() + ":" +
					result.Get("Test.TestContent.TestDomain").String()
				school.TestAttempts[key]++
				school.TotalAttempts++
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
// receives qa summary as input, summaries participation codes for test attempts, and disruptions
//
func qaParticipationQueryExecutor(ctx context.Context, query string, url string, in <-chan schoolQASummary) (<-chan schoolQASummary, <-chan error, error) {

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
func qaTransformSummary(ctx context.Context, in <-chan schoolQASummary) (<-chan gjson.Result, <-chan error, error) {

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
