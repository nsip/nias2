// html-sanitizer.go

package nap_writing_sanitiser

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"strings"
)

//
// strips html tags from the provided writing response
//
func createHtmlSanitizer(ctx context.Context, in <-chan map[string]string) (
	<-chan map[string]string, <-chan map[string]string, <-chan error, error) {

	out := make(chan map[string]string)
	report := make(chan map[string]string)
	errc := make(chan error, 1)
	p := bluemonday.NewPolicy()
	p.AllowElements("strong", "em", "span", "p", "ol", "ul", "li", "br", "u", "font", "h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowAttrs("style").Matching(regexp.MustCompile("^((text-decoration:underline|text-decoration-line:underline|font-size:16px|font-size:18px|font-size:large|text-align:left|text-align:center|text-align:start|background-color:rgba\\(255, 255, 255, 0\\));)+$")).Globally()
	p.AllowAttrs("size").Matching(regexp.MustCompile("^\\d+$")).OnElements("font")

	go func() {
		defer close(out)
		defer close(report)
		defer close(errc)

		for rMap := range in {
			rmap := rMap

			fragmentHtml := rmap["Item Response"]
			sanitised := unescape(p.Sanitize(fragmentHtml))
			rmap["Item Response"] = sanitised

			select {
			case out <- rmap:
				// log.Printf("\nrecord read:\n\n%v\n\n", recordMap)
			case <-ctx.Done():
				return
			}

			if sanitised != fragmentHtml {
				rpt := make(map[string]string)
				for k, v := range rmap {
					rpt[k] = v
				}
				rpt["Original Item Response"] = fragmentHtml
				select {
				case report <- rpt:
					// log.Printf("\nrecord read:\n\n%v\n\n", recordMap)
				case <-ctx.Done():
					return
				}

			}
		}
	}()
	return out, report, errc, nil

}

func unescape(x string) string {
	x = strings.Replace(x, "&#39;", "'", -1)
	x = strings.Replace(x, "&#34;", "\"", -1)
	x = strings.Replace(x, "\u201C", "&ldquo;", -1)
	x = strings.Replace(x, "\u201D", "&rdquo;", -1)
	x = strings.Replace(x, "\u2018", "&lsquo;", -1)
	x = strings.Replace(x, "\u2019", "&rsquo;", -1)
	x = strings.Replace(x, "\u2026", "&hellip;", -1)
	x = strings.Replace(x, "\u00A0", "&nbsp;", -1)
	x = strings.Replace(x, "\u00B4", "&acute;", -1)
	x = strings.Replace(x, "\u00A8", "&uml;", -1)
	x = strings.Replace(x, "\u00B0", "&deg;", -1)
	x = strings.Replace(x, "\u2013", "&ndash;", -1)
	x = strings.Replace(x, "\n", "", -1)
	return x
}
