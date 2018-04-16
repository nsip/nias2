// html-sanitizer.go

package nap_writing_print

import (
	"context"
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
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

			p := bluemonday.StrictPolicy()
			// keep paragraph markers to preserve layout in pdf printing
			p.AllowElements("p")
			sanitizedHtml := p.Sanitize(rmap["Item Response"])
			// now unescape the html for speech marks, apostrophe's etc. for printing
			unescapedHtml := html.UnescapeString(sanitizedHtml)

			// remove unnecessary characters
			// remove non-breaking spaces and line feeds & backticks
			noBrHtml := strings.Replace(unescapedHtml, "&nbsp;", " ", -1)
			noLfHtml := strings.Replace(noBrHtml, "\\n", "", -1)

			// remove end-of para markers, not needed for pdf print
			noEndParaHtml := strings.Replace(noLfHtml, "</p>", "", -1)

			rmap["Item Response"] = noEndParaHtml

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
