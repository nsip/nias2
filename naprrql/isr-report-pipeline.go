// query-report-writer.go

package naprrql

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

//
// helper type to capture query parameters
//
type isrQueryParams struct {
	year   string
	school string
}

//
// create and run an isr printing pipeline.
//
// Pipeline streams requests at a given year level, school by school
// feeding results to the ouput csv file.
//
// This means the server & parser never have to deal with query data
// volumes larger than a single school at a time.
//
// Overall round-trip latency is less than querying for all data at once
// and ensures we can't run out of memory
//
func runISRPipeline(year string, schools []string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	varsc, errc, err := isrParametersSource(ctx, year, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	queryTemplates := getTemplates("./reporting_templates/isr_printing/")
	var query string
	for queryFileName, queryText := range queryTemplates {
		query = queryText
		jsonc, errc, err := isrQueryExecutor(ctx, query, DEF_ISR_URL, varsc)
		if err != nil {
			return err
		}
		errcList = append(errcList, errc)

		// sink stage
		// create working directory if not there
		outFileDir := "./out/isr_printing"
		err = os.MkdirAll(outFileDir, os.ModePerm)
		if err != nil {
			return err
		}
		csvFileName := deriveQueryFileName(queryFileName) + "_yr" + year + ".csv"
		outFileName := outFileDir + "/" + csvFileName
		mapFileName := deriveMapFileName(queryFileName)
		errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc)
		if err != nil {
			return err
		}
		errcList = append(errcList, errc)

		log.Println("ISR print file writing... " + outFileName)
	}
	return WaitForPipeline(errcList...)
}

//
// acts as input feed to the pipeline, sends parameters to retrieve data for
// each school in turn, for the given year level.
//
func isrParametersSource(ctx context.Context, year string, schools ...string) (<-chan isrQueryParams, <-chan error, error) {

	{ //check input variables, handle errors before goroutine starts
		if len(schools) == 0 {
			return nil, nil, errors.Errorf("no schools provided")
		}

		if year == "" {
			return nil, nil, errors.Errorf("no year provided")
		}
	}

	out := make(chan isrQueryParams)
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
			vars := isrQueryParams{school: school, year: year}
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
func isrQueryExecutor(ctx context.Context, query string, url string, in <-chan isrQueryParams) (<-chan gjson.Result, <-chan error, error) {

	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for params := range in {
			vars := map[string]interface{}{"acaraID": params.school, "yrLevel": params.year}
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
		}
	}()
	return out, errc, nil
}
