// system-report-pipeline.go

package naprrql

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

//
// printing pipeline for system-level reports
//

//
// helper type to capture query parameters
//
type systemQueryParams struct {
	schoolAcaraID string
}

//
// create and run a system-level report printing pipeline.
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
func RunSystemReportPipeline(queryFileName string, query string, schools []string) error {

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
	outFileDir := "./out/system_reports"
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

	log.Println("System report file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

func systemParametersValidate(schools []string) error {
	{ //check input variables, handle errors before goroutine starts
		if len(schools) == 0 {
			return errors.Errorf("no schools provided")
		}
	}
	for schoolIndex, school := range schools {
		if school == "" {
			return errors.Errorf("school %v is empty string", schoolIndex+1)
		}
	}
	return nil
}

//
// acts as input feed to the pipeline, sends parameters to retrieve data for
// each school in turn, for the given year level.
//
func systemParametersSource(ctx context.Context, schools ...string) (<-chan systemQueryParams, <-chan error, error) {
	err := systemParametersValidate(schools)
	if err != nil {
		return nil, nil, err
	}

	out := make(chan systemQueryParams)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for _, school := range schools {
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
// query executor transform stage takes query params in, executes gql query
// and writes results to output channel
//
func systemQueryExecutor(ctx context.Context, query, url string, in <-chan systemQueryParams) (<-chan gjson.Result, <-chan error, error) {

	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for params := range in {
			vars := map[string]interface{}{"acaraIDs": []string{params.schoolAcaraID}}
			json, err := gql.DoQuery(url, query, vars)
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				return
			}
			for _, result := range json.Array() {
				// Send the data to the output channel but return early
				// if the context has been cancelled.
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
			// handle edge cases where data is not returned as an array
			if len(json.Array()) == 0 {
				select {
				case out <- json:
				case <-ctx.Done():
					return
				}
			}

		}
	}()
	return out, errc, nil
}

// As above, but do not break up acaraID arguments
func systemQueryExecutorNoBreakdown(ctx context.Context, query, url string, vars map[string]interface{}) (<-chan gjson.Result, <-chan error, error) {

	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		json, err := gql.DoQuery(url, query, vars)
		if err != nil {
			// Handle an error that occurs during the goroutine.
			errc <- err
			return
		}
		for _, result := range json.Array() {
			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- result:
			case <-ctx.Done():
				return
			}
		}
		// handle edge cases where data is not returned as an array
		if len(json.Array()) == 0 {
			select {
			case out <- json:
			case <-ctx.Done():
				return
			}
		}

	}()
	return out, errc, nil
}
