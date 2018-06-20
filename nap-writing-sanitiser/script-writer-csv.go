// script-writer-html.go

package nap_writing_sanitiser

import (
	"context"
	"encoding/csv"
	//"fmt"
	"log"
	"os"
	"strings"
)

//
// creates the full html file output of the writing response
//
func createScriptWriterCSV(ctx context.Context, csvFileName string, in <-chan map[string]string) (<-chan error, error) {

	errc := make(chan error, 1)
	file, err := os.Create("out/" + csvFileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(file)
	headerWritten := false

	go func() {
		defer close(errc)
		defer file.Close()
		var rowFormat []string
		var dataKeys []string

		for record := range in {

			resultRow := make([]string, 0)

			// if the header hasn't been written then create it
			if !headerWritten {
				rowFormat, dataKeys = deriveRowFormat(record)
				headerRow := make([][]string, 0)
				headerRow = append(headerRow, rowFormat)
				err := w.WriteAll(headerRow)
				if err != nil {
					// Handle an error that occurs during the goroutine.
					errc <- err
					log.Printf("%+v\n", err)
					return
				}
				headerWritten = true
			}
			for _, key := range dataKeys {
				resultRow = append(resultRow, strings.Replace(record[key], "\n", "\\n", -1))
			}
			if !emptyCsvRow(resultRow) {
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

func emptyCsvRow(row []string) bool {
	for _, r := range row {
		if len(r) > 0 {
			return false
		}
	}
	return true
}

func deriveRowFormat(resultMap map[string]string) (displayNames []string, dataKeys []string) {
	displayNames = make([]string, 0)
	dataKeys = make([]string, 0)

	for key, _ := range resultMap {
		displayNames = append(displayNames, key)
		dataKeys = append(dataKeys, key)
	}
	return displayNames, dataKeys
}
