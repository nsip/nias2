// report-pipeline-helpers.go
package naprrql

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//
// utiltiy methods to support reporting pipelines
//

//
// generic print pipeline terminator that writes to a csv file
// mapFile is the path/name of a *_map.csv file which contains the
// display names for each csv column, and the name of the json data fields
// to put in those colums.
// If no map is found or is unreadable the data format will be derived
// from the first record the processor recieves.
//
func csvFileSink(ctx context.Context, csvFileName string, mapFileName string, in <-chan gjson.Result) (<-chan error, error) {

	file, err := os.Create(csvFileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(file)
	rowFormat, dataKeys := readRowFormat(mapFileName)
	headerWritten := false

	errc := make(chan error, 1)
	go func() {
		i := 0
		defer close(errc)
		defer file.Close()

		for record := range in {

			// if no row template established derive one from the data
			if len(rowFormat) == 0 {
				rowFormat, dataKeys = deriveRowFormat(record)
			}

			resultRow := make([]string, 0)
			for _, key := range dataKeys {
				resultRow = append(resultRow, strings.Replace(record.Get(key).String(), "\n", "\\n", -1))
			}

			// if the header hasn't been written then create it
			if !headerWritten {
				headerRow := make([][]string, 0)
				headerRow = append(headerRow, rowFormat, resultRow)
				err := w.WriteAll(headerRow)
				if err != nil {
					// Handle an error that occurs during the goroutine.
					errc <- err
					log.Printf("%+v\n", err)
					return
				}
				headerWritten = true
			} else {
				i++
				err := w.Write(resultRow)
				if err != nil {
					// Handle an error that occurs during the goroutine.
					errc <- err
					log.Printf("%+v\n", err)
					return
				}

			}
			w.Flush()
		}
	}()
	return errc, nil
}

//
// if a formatting file provided to fix layout of
// columns in csv, then attempt to read now
// errors are not propagated, so pipeline will function with
// a derived layout rather than quitting
//
func readRowFormat(fileName string) (displayNames []string, dataKeys []string) {

	displayNames = make([]string, 0)
	dataKeys = make([]string, 0)

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Printf("No format map: %s: found: layout will be auto-derived from data.", fileName)
		return displayNames, dataKeys
	}

	r := csv.NewReader(file)
	formatLines, err := r.ReadAll()
	if err != nil {
		log.Println("Error reading format map: ", err)
		return displayNames, dataKeys
	}

	if len(formatLines) == 2 {
		displayNames = formatLines[0]
		dataKeys = formatLines[1]
	}

	return displayNames, dataKeys

}

//
// if no header provided, then derive an arbitrary order from
// the first result that arrives
//
func deriveRowFormat(result gjson.Result) (displayNames []string, dataKeys []string) {
	displayNames = make([]string, 0)
	dataKeys = make([]string, 0)

	resultMap := result.Map()

	for key, _ := range resultMap {
		displayNames = append(displayNames, key)
		dataKeys = append(dataKeys, key)
	}
	return displayNames, dataKeys
}

// WaitForPipeline waits for results from all error channels.
// It returns early on the first error.
func WaitForPipeline(errs ...<-chan error) error {

	errc := MergeErrors(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

// MergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
func MergeErrors(cs ...<-chan error) <-chan error {

	var wg sync.WaitGroup
	out := make(chan error)
	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls
	// wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	// Start a goroutine to close out once all the output goroutines
	// are done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

//
// from a given report file name, create a filename for the
// csv output to write to.
//
func deriveCSVFileName(queryFileName string) string {

	reportFilename := filepath.Base(queryFileName)
	csvFilename := strings.Replace(reportFilename, ".gql", ".csv", 1)

	return csvFilename

}

//
// from a given report file name, derive the associated map file name
//
func deriveMapFileName(queryFileName string) string {

	mapFilename := strings.Replace(queryFileName, ".gql", "_map.csv", 1)

	return mapFilename

}

//
// from a given report file name, create a base filename stripping suffix
// and directory
//
func deriveQueryFileName(queryFileName string) string {

	reportFilename := filepath.Base(queryFileName)
	csvFilename := strings.Replace(reportFilename, ".gql", "", 1)

	return csvFilename

}

// splits a channel of GJSON into an array of channels of GJSON
func splitter(ctx context.Context, in <-chan gjson.Result, size int) (
	[]chan gjson.Result, <-chan error, error) {
	jsonc1 := make([]chan gjson.Result, size)
	for i := range jsonc1 {
		jsonc1[i] = make(chan gjson.Result)
	}
	errc := make(chan error, 1)
	/*
		var wg sync.WaitGroup
		wg.Add(size)
	*/
	go func() {
		for i := range jsonc1 {
			defer close(jsonc1[i])
		}
		defer close(errc)
		for n := range in {
			// Send the data to the output channel 1 but return early
			// if the context has been cancelled.
			for _, c := range jsonc1 {
				select {
				case c <- n:
				case <-ctx.Done():
					return
				}
			}
		}

	}()
	/*
		go func() {
			defer close(errc)
			for n := range in {
				for _, c := range jsonc1 {
					select {
					case c <- n:
					case <-ctx.Done():
						return
					}
				}
			}
			for _, c := range jsonc1 {
				close(c)
			}
		}()
	*/
	return jsonc1, errc, nil

}

// Specific to the XML report: takes in JSON output, marshals to XML, prints out. Does not insert
// container elements (e.g. no StudentPersonals container), does not segregate objects by class
// into different files. Applies filtering from nappqrl.toml (which is expressed in JSON dot notation)
func xmlFileSink(ctx context.Context, xmlFileName string, in <-chan []byte) (<-chan error, error) {
	config := LoadNAPLANConfig()
	// filter format: SJSON dot notation
	filter := config.XMLFilter

	file, err := os.Create(xmlFileName)
	if err != nil {
		return nil, err
	}

	errc := make(chan error, 1)
	go func() {
		i := 0
		defer close(errc)
		defer file.Close()

		for record := range in {

			i++
			var to TypedObject
			var out []byte
			for _, rule := range filter {
				record1, err := sjson.DeleteBytes(record, rule)
				if err == nil {
					record = record1
					// ignore errors, such as the filter not applying to this object
				}
			}
			json.Unmarshal(record, &to)
			if to.NAPTestScoreSummary != nil {
				out, err = xml.MarshalIndent(to.NAPTestScoreSummary, "", "  ")
			} else if to.SchoolInfo != nil {
				out, err = xml.MarshalIndent(to.SchoolInfo, "", "  ")
			} else if to.StudentPersonal != nil {
				out, err = xml.MarshalIndent(to.StudentPersonal, "", "  ")
			} else if to.NAPEventStudentLink != nil {
				out, err = xml.MarshalIndent(to.NAPEventStudentLink, "", "  ")
			} else if to.NAPTest != nil {
				out, err = xml.MarshalIndent(to.NAPTest, "", "  ")
			} else if to.NAPTestlet != nil {
				out, err = xml.MarshalIndent(to.NAPTestlet, "", "  ")
			} else if to.NAPTestItem != nil {
				out, err = xml.MarshalIndent(to.NAPTestItem, "", "  ")
			} else if to.NAPStudentResponseSet != nil {
				out, err = xml.MarshalIndent(to.NAPStudentResponseSet, "", "  ")
			} else if to.NAPCodeFrame != nil {
				out, err = xml.MarshalIndent(to.NAPCodeFrame, "", "  ")
			} else {
				err = fmt.Errorf("%+v: no type selected for record", to)
			}
			if err != nil {
				// Handle an error that occurs during the goroutine.
				errc <- err
				log.Printf("%+v\n", err)
				return
			}
			_, err = file.Write(out)
			if err != nil {
				errc <- err
				log.Printf("%+v\n", err)
				return
			}
		}
	}()
	return errc, nil
}
