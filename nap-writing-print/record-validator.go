// record-validator.go

package nap_writing_print

import (
	"context"
	"fmt"
	"log"
	"time"
)

//
// accepts map[string]string of a csv record-row, checks for
// missing fields and assigns values to empty fields as
// appropriate
//
func createRecordValidator(ctx context.Context, in <-chan map[string]string) (
	<-chan map[string]string, <-chan error, error) {

	out := make(chan map[string]string)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for rMap := range in {
			recordMap := rMap
			fillMissingFields(recordMap)
			fillEmptyFields(recordMap)
			select {
			case out <- recordMap:
				// log.Printf("\nrecord read:\n\n%v\n\n", recordMap)
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil

}

//
// ensure there are no missing fields in the onward record
//
func fillMissingFields(rmap map[string]string) {

	expectedHeaders := []string{"Test level", "Jurisdiction Id", "Local school student ID",
		"TAA student ID", "Participation Code", "Anonymised Id", "Test Year", "ACARA ID", "PSI"}

	for _, key := range expectedHeaders {
		_, ok := rmap[key]
		if !ok {
			rmap[key] = ""
		}
	}

}

//
// empty fields in the csv map are legitimate, but need printable
// content for audit purposes
//
func fillEmptyFields(rmap map[string]string) {

	key := "Test level"
	if rmap[key] == "" {
		rmap[key] = "Test Level not provided"
	}

	key = "Jurisdiction Id"
	if rmap[key] == "" {
		rmap[key] = "Jurisdiction Id not provided"
	}

	key = "Local school student ID"
	if rmap[key] == "" {
		rmap[key] = "Local-school student ID not provided"
	}

	key = "TAA student ID"
	if rmap[key] == "" {
		rmap[key] = "TAA student ID not provided"
	}

	key = "Participation Code"
	if rmap[key] == "" {
		// single string as may be used as part of output file path
		rmap[key] = "Participation_code_not_provided"
	}

	key = "Anonymised Id"
	if rmap[key] == "" {
		// single string as may be used as part of output file path
		log.Println("WARNING: record received without anonymised ID")
		rmap[key] = "Anonymised_Id_not_provided"
	}

	key = "Test Year"
	if rmap[key] == "" {
		rmap[key] = fmt.Sprintf("%d", time.Now().Year())
	}

	key = "ACARA ID"
	if rmap[key] == "" {
		// single string as may be used as part of output file path
		rmap[key] = "ACARA_School_ID_not_provided"
	}

	key = "PSI"
	if rmap[key] == "" {
		rmap[key] = "PSI not provided"
	}

	key = "Item Response"
	if rmap[key] == "" {
		rmap[key] = "No writing-response text recorded. Check Participation Code."
	}

}
