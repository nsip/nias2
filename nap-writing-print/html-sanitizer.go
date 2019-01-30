// html-sanitizer.go

package nap_writing_print

import (
	"context"
	"strings"
)

//
// strips redundant characters from the provided writing response
//
// Note:
// in theory the html provided from the NAPLAN platform has already been sanitized
// and so should be fit for output.
//
// in reality in 2018 this turned out not to be the case, so a separate tool
// nap-writing-sanitizer was created to deal with the html issues in one go.
//
// if users have badly corrupted inputs we advise them to run the sanitizer first
// and then pass the output through this tool for final printing.
//
func createHtmlSanitizer(ctx context.Context, in <-chan map[string]string) (
	<-chan map[string]string, <-chan error, error) {

	out := make(chan map[string]string)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for rMap := range in {
			rmap := rMap

			fragmentHtml := rmap["Item Response"]
			noLFHtml := strings.Replace(fragmentHtml, "\\n", "", -1)

			rmap["Item Response"] = noLFHtml

			select {
			case out <- rmap:
				// log.Printf("\nrecord read:\n\n%v\n\n", recordMap)
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil

}
