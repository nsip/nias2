// csv-reader.go

package nap_writing_sanitiser

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
	csvutils "github.com/wildducktheories/go-csv"
)

//
// opens & reads a csv file, passing each line as a map[string]string
// structure to the next component in the pipe
//
func createCSVReader(ctx context.Context, csvFileName string) (
	<-chan map[string]string, <-chan error, error) {

	out := make(chan map[string]string)
	errc := make(chan error, 1)

	// build the reader, raw csv reader is wrapped to provide
	// channel access & row-reord to map functionality
	f, err := os.Open(csvFileName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to open csv file: "+csvFileName)
	}

	csvr := csv.NewReader(f)
	reader := csvutils.WithCsvReader(csvr, f)

	go func() {
		defer close(out)
		defer close(errc)
		defer f.Close()
		defer reader.Close()

		for record := range reader.C() {
			recordMap := record.AsMap()
			select {
			case out <- recordMap:
				// log.Printf("\nrecord read:\n\n%v\n\n", recordMap)
			case <-ctx.Done():
				return
			}
		}
		// check for any errors that may have closed the stream early
		if reader.Error() != nil {
			errc <- errors.Wrap(reader.Error(), "csv stream unexpectedly closed")
		}
	}()
	return out, errc, nil

}
