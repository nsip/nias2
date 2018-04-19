// audit-writer-html.go

// script-writer-html.go

package nap_writing_print

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

var auditHeader = `
<html>
<head>
    <style>
    p,
    li,
    h2 {
        font-family: Verdana, Arial, sans-serif;
    }

    .audit-body {
        width: 650px;
        margin-top: 45px;
        margin-bottom: 30px;
        margin-left: auto;
        margin-right: auto;
    }
    </style>
</head>

<body>
    <div class="audit-body">`

var auditReport = `
        <p style="text-align: left;">Test Year: %s</p>
        <p style="text-align: left;">Test Level: %s:</p>
        <p style="text-align: left;">TAA Student ID: %s</p>
        <p style="text-align: left;">PSI: %s</p>
        <p style="text-align: left;">Local School Student ID: %s</p>
        <p style="text-align: left;">Jurisdiction ID: %s</p>
        <p style="text-align: left;">Anonymised ID: %s</p>        
        <p style="text-align: left;">Participation Code: %s</p>
        <p style="text-align: left;">School ACARA ID: %s</p>
		<p style="text-align: left;">Script File Name: %s</p>        	
`

var auditFooter = `
    </div>
</body>
</html>
`

//
// creates the audit html file that accompanies the html script file
//
func createAuditWriterHtml(ctx context.Context, in <-chan map[string]string) (<-chan error, error) {

	errc := make(chan error, 1)

	go func() {
		defer close(errc)

		for rMap := range in {
			rmap := rMap

			// create the writing script html files (one in each output folder structure)
			yrLevelFileName := fmt.Sprintf("%s%s", rmap["yr_level_audit_path"], rmap["html_audit_filename"])
			schoolFileName := fmt.Sprintf("%s%s", rmap["school_level_audit_path"], rmap["html_audit_filename"])

			// build the audit report
			auditBody := fmt.Sprintf(auditReport, rmap["Test Year"], rmap["Test level"],
				rmap["TAA student ID"], rmap["PSI"], rmap["Local school student ID"],
				rmap["Jurisdiction Id"], rmap["Anonymised Id"], rmap["Participation Code"], rmap["ACARA ID"],
				rmap["html_script_filename"])

			f, err := os.Create(schoolFileName)
			if err != nil {
				log.Printf("error creating file: %s error is: %s\n", schoolFileName, err)
				continue
			}

			anonID := strings.TrimSuffix(rmap["html_script_filename"], ".html")
			fileNameBanner := fmt.Sprintf("<h2 style=\"text-align: center;\">%s</h2>", anonID)
			topFileName := fileNameBanner

			doc := []string{auditHeader, topFileName, auditBody, auditFooter}
			for _, element := range doc {
				_, err := f.WriteString(element)
				if err != nil {
					log.Printf("unable to print to document %s error is:%s\n ", schoolFileName, err)
				}
			}

			f.Sync()
			f.Close()

			// create copy in other file structure
			err = copyFileContents(schoolFileName, yrLevelFileName)
			if err != nil {
				log.Printf("unable to copy file to yr level folder: %s\n", err)
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	return errc, nil

}
