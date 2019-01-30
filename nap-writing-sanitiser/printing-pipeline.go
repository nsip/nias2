// printing-pipeline.go

package nap_writing_sanitiser

import (
	"context"
	"sync"
)

//
// builds and executes the prinitng pipeline
//

func RunPrintingPipeline(csvFileName string) error {

	// create a context to manage pipeline
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error

	// create csv reader input stage
	csv_map_out, errc, err := createCSVReader(ctx, csvFileName)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	// validate structure
	// inject missing/empty field defaults
	valid_line_out, errc, err := createRecordValidator(ctx, csv_map_out)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// sanitize provided writing response html
	//
	sanitize_out, sanitize_report_out, errc, err := createHtmlSanitizer(ctx, valid_line_out)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// write script CSV file
	//
	errc, err = createScriptWriterCSV(ctx, "writing_extract_sanitised.csv", sanitize_out)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// write sanitise report CSV file
	//
	errc, err = createScriptWriterCSV(ctx, "sanitiser_report.csv", sanitize_report_out)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	//
	// wait for pipeline to complete
	//
	return WaitForPipeline(errcList...)

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
	// We must ensure that the output channel has the capacity to
	// hold as many errors
	// as there are error channels.
	// This will ensure that it never blocks, even
	// if WaitForPipeline returns early.
	out := make(chan error, len(cs))
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
