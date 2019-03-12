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
	"strconv"
	"strings"
	"sync"

	"github.com/beevik/etree"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//
// utiltiy methods to support reporting pipelines
//

func emptyCsvRow(row []string) bool {
	for _, r := range row {
		if len(r) > 0 {
			return false
		}
	}
	return true
}

//
// generic print pipeline terminator that writes to a csv file
// mapFile is the path/name of a *_map.csv file which contains the
// display names for each csv column, and the name of the json data fields
// to put in those colums.
// If no map is found or is unreadable the data format will be derived
// from the first record the processor recieves.
//
// Will also deal with mapFile specifying Fixed Format width: in that case,
// the first line of the mapFile is "FixedFormat", the second is the
// absolute character widths of each colum (if zero padding if prefixed with 0),
// and the third is the name of the json data fields to put in those columns.
// Blanks are printed  if the data field is named as Blank.
//
func csvFileSink(ctx context.Context, csvFileName string, mapFileName string, in <-chan gjson.Result) (<-chan error, error) {

	file, err := os.Create(csvFileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(file)
	rowFormat, dataKeys, dataWidths := readRowFormat(mapFileName)
	headerWritten := false
	if len(dataWidths) > 0 {
		return fixedFormatFileSink(ctx, csvFileName, mapFileName, in)
	}

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
				if !emptyCsvRow(resultRow) {
					i++
					err := w.Write(resultRow)
					if err != nil {
						// Handle an error that occurs during the goroutine.
						errc <- err
						log.Printf("%+v\n", err)
						return
					}
				}
			}
			w.Flush()
		}
	}()
	return errc, nil
}

func fixedFormatFileSink(ctx context.Context, csvFileName string, mapFileName string, in <-chan gjson.Result) (<-chan error, error) {
	file, err := os.Create(csvFileName)
	if err != nil {
		return nil, err
	}
	_, dataKeys, dataWidths := readRowFormat(mapFileName)
	errc := make(chan error, 1)
	if err != nil {
		errc <- err
		log.Printf("%+v\n", err)
		return errc, err

	}
	go func() {
		i := 0
		defer close(errc)
		defer file.Close()
		for record := range in {
			resultRow := make([]string, 0)
			for _, key := range dataKeys {
				if key == "Blank" {
					resultRow = append(resultRow, "")
				} else {
					value := record.Get(key).String()
					if strings.HasSuffix(key, ".BirthDate") {
						value = strings.Replace(value, "-", "", -1)
					}
					resultRow = append(resultRow, strings.Replace(value, "\n", "\\n", -1))
				}
			}
			if !emptyCsvRow(resultRow) {
				//log.Printf("FIXED: %#v", resultRow)
				out := fixedFormat(resultRow, dataWidths)
				//log.Printf("%s", out)
				//log.Printf("%s", strings.Join(out, ""))
				i++
				_, err := file.WriteString(strings.Join(out, "") + "\n")
				if err != nil {
					// Handle an error that occurs during the goroutine.
					errc <- err
					log.Printf("%+v\n", err)
					return
				}
			}
			file.Sync()
		}
	}()
	return errc, nil
}

// convert slice of strings into slice of strings with fixed width given in format.
// if width prefixed with 0, zero pad output
func fixedFormat(data []string, format []string) []string {
	ret := make([]string, 0)
	for i, value := range data {
		zeropad := strings.HasPrefix(format[i], "0")
		length, err := strconv.Atoi(format[i])
		if err != nil || length < 1 {
			continue
		}
		if length > len(value) {
			replchar := " "
			if zeropad {
				replchar = "0"
			}
			value = strings.Repeat(replchar, length-len(value)) + value
		} else if length < len(value) {
			//log.Println("Truncating: " + value)
			value = string(value[:length]) + ""
			//log.Println("Truncated: " + value)
		}
		ret = append(ret, value)
	}
	return ret
}

//
// if a formatting file provided to fix layout of
// columns in csv, then attempt to read now
// errors are not propagated, so pipeline will function with
// a derived layout rather than quitting.

// If first line is "FixedFormat", then return fixed format specification.
//
func readRowFormat(fileName string) (displayNames []string, dataKeys []string, dataWidths []string) {

	displayNames = make([]string, 0)
	dataKeys = make([]string, 0)
	dataWidths = make([]string, 0)

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Printf("No format map: %s: found: layout will be auto-derived from data.", fileName)
		return displayNames, dataKeys, dataWidths
	}

	r := csv.NewReader(file)
	r.FieldsPerRecord = -1 // don't immediately check number of columns
	formatLines, err := r.ReadAll()
	if err != nil {
		if perr, ok := err.(*csv.ParseError); !ok || perr.Err != csv.ErrFieldCount {
			log.Println("Error reading format map: ", err)
			return displayNames, dataKeys, dataWidths
		}
	}

	if len(formatLines) == 2 {
		if len(formatLines[0]) != len(formatLines[1]) {
			log.Println("Error reading format map: unequal count of fields in header and definitions line")
			return displayNames, dataKeys, dataWidths
		}
		displayNames = formatLines[0]
		dataKeys = formatLines[1]
	} else if len(formatLines) == 3 && formatLines[0][0] == "FixedFormat" {
		if len(formatLines[1]) != len(formatLines[2]) {
			log.Println("Error reading format map: unequal count of fields in widths and definitions line")
			return displayNames, dataKeys, dataWidths
		}
		dataWidths = formatLines[1]
		dataKeys = formatLines[2]
	}

	return displayNames, dataKeys, dataWidths

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
	_, _ = file.Write([]byte("<sif xmlns=\"http://www.sifassociation.org/datamodel/au/3.4\">\n"))

	errc := make(chan error, 1)
	go func() {
		i := 0
		defer close(errc)
		defer file.Close()

		for record := range in {
			//log.Println(string(record))
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
				// to_SIF() to get the attributes back in to the flattened structure
				out, err = xml.MarshalIndent(to.StudentPersonal.To_SIF(), "", "  ")
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

			doc := etree.NewDocument()
			err = doc.ReadFromBytes(out)
			if err != nil {
				errc <- err
				log.Printf("%+v\n", err)
				return
			}
			for _, path := range filter {
				elems := doc.Root().FindElements(path)
				for _, elem := range elems {
					elem.Parent().RemoveChildAt(elem.Index())
				}
			}
			doc.Indent(2)
			out1, err := doc.WriteToBytes()
			if err != nil {
				errc <- err
				log.Printf("%+v\n", err)
				return
			}

			_, err = file.Write(out1)
			if err != nil {
				errc <- err
				log.Printf("%+v\n", err)
				return
			}
		}
		_, _ = file.Write([]byte("</sif>\n"))
	}()
	return errc, nil
}
