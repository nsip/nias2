// html-sanitizer.go

package nap_writing_print

import (
	"context"
	"strings"
)

//
// strips html tags from the provided writing response
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
