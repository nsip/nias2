// query-report-writer.go

package naprrql

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//
// helper type to capture query parameters
//
type itemQueryParams struct {
	//year   string
	school string
}

//
// create and run an item printing pipeline.
//
// Pipeline streams requests at a school by school
// feeding results to the ouput csv file.
//
// This means the server & parser never have to deal with query data
// volumes larger than a single school at a time.
//
// Overall round-trip latency is less than querying for all data at once
// and ensures we can't run out of memory
//

// TODO may need to reintroduce year level here to break payload size down further
func RunItemPipeline(schools []string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// input stage
	varsc, errc, err := itemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	queryTemplates := getTemplates("./reporting_templates/item_printing/")
	var query string
	for _, queryText := range queryTemplates {
		query = queryText
	}
	jsonc, errc, err := itemQueryExecutor(ctx, query, DEF_ITEM_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/item_printing"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	csvFileName := "itemResults.csv"
	outFileName := outFileDir + "/" + csvFileName
	mapFileName := "./reporting_templates/item_printing/itemPrinting_map.csv"
	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("Item print file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

// Slight variant of the foregoing
func RunWritingExtractPipeline(schools []string, psi_exceptions []string, blacklist bool, psi2prompt map[string]string) error {

	// setup pipeline cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	//
	// this pipeline call moved up a level to GenerateWritingExtractReports() so that
	// both reports are created in parallel.
	//
	// pipelineError := runQAWritingSchoolSummaryPipeline(schools, "./out/writing_extract", "./reporting_templates/writing_extract/qaSchools_map.csv")
	// if pipelineError != nil {
	// 	return pipelineError
	// }

	// input stage
	varsc, errc, err := itemParametersSource(ctx, schools...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// transform stage
	queryTemplates := getTemplates("./reporting_templates/writing_extract/")
	var query string
	for _, queryText := range queryTemplates {
		query = queryText
	}
	jsonc, errc, err := itemQueryExecutor(ctx, query, DEF_GQL_URL, varsc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	jsonc0, errc, err := filterPSI(ctx, jsonc, psi_exceptions, blacklist)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	jsonc01, errc, err := replacePrompt(ctx, jsonc0, psi2prompt)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	jsonc1, errc, err := splitter(ctx, jsonc01, 2)
	errcList = append(errcList, errc)

	// sink stage
	// create working directory if not there
	outFileDir := "./out/writing_extract"
	err = os.MkdirAll(outFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	csvFileName := "writing_extract.csv"
	outFileName := outFileDir + "/" + csvFileName

	//
	// switch below deprecated as wordcount always included
	//
	// mapFileName := "./reporting_templates/writing_extract/itemWritingPrinting_map.csv"
	// if wordcount {

	mapFileName := "./reporting_templates/writing_extract/itemWritingPrintingWordCount_map.csv"
	// }

	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc1[0])
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	outFileName = outFileDir + "/" + "qaPsi.csv"
	mapFileName = "./reporting_templates/writing_extract/qaPsi_map.csv"
	errc, err = csvFileSink(ctx, outFileName, mapFileName, jsonc1[1])
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	log.Println("Writing extract file writing... " + outFileName)
	return WaitForPipeline(errcList...)
}

//
// acts as input feed to the pipeline, sends parameters to retrieve data for
// each school in turn
//
func itemParametersSource(ctx context.Context, schools ...string) (<-chan systemQueryParams, <-chan error, error) {

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
			//vars := itemQueryParams{school: school}
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
func itemQueryExecutor(ctx context.Context, query string, url string, in <-chan systemQueryParams) (<-chan gjson.Result, <-chan error, error) {

	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		gql := NewGqlClient()
		for params := range in {
			//vars := map[string]interface{}{"acaraID": params.school}
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
		}
	}()
	return out, errc, nil
}

// ignore an empty psi_exceptions array
func filterPSI(ctx context.Context, in <-chan gjson.Result, psi_exceptions []string, blacklist bool) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	exclude := make(map[string]bool)
	ignore := len(psi_exceptions) == 0
	for _, p := range psi_exceptions {
		exclude[p] = true
	}
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			psi := record.Get("Student.OtherIdList.OtherId.#[Type==NAPPlatformStudentId].Value").String()
			if !ignore && (exclude[psi] && blacklist || !exclude[psi] && !blacklist) {
				continue
			}
			select {
			case out <- record:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}

func replacePrompt(ctx context.Context, in <-chan gjson.Result, psi2prompt map[string]string) (<-chan gjson.Result, <-chan error, error) {
	out := make(chan gjson.Result)
	errc := make(chan error, 1)
	ignore := len(psi2prompt) == 0
	go func() {
		defer close(out)
		defer close(errc)
		for record := range in {
			psi := record.Get("Student.OtherIdList.OtherId.#[Type==NAPPlatformStudentId].Value").String()
			if !ignore {
				if prompt, ok := psi2prompt[psi]; ok {
if len(prompt)>0 {
					record1, _ := sjson.Set(record.Raw, "Test.TestContent.LocalId", prompt)
					record = gjson.Parse(record1)
}
				}
			}
			select {
			case out <- record:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}
