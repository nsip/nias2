// script-writer-html.go

package nap_writing_print

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

var scriptHeader = `
<html>
<head>
    <style>
    p,
    li,
    h2 {
        font-family: Verdana, Arial, sans-serif;
    }

    .response-body {
        width: 600px;
        margin-top: 25px;
        margin-bottom: 30px;
        margin-left: auto;
        margin-right: auto;
    }
    </style>
</head>

<body>
    <div class="response-body">`

var scriptFooter = `
    </div>
</body>
</html>
`

//
// creates the full html file output of the writing response
//
func createScriptWriterHtml(ctx context.Context, in <-chan map[string]string) (<-chan error, error) {

	errc := make(chan error, 1)

	go func() {
		defer close(errc)

		for rMap := range in {
			rmap := rMap

			// create the writing script html files (one in each output folder structure)
			yrLevelFileName := fmt.Sprintf("%s%s", rmap["yr_level_script_path"], rmap["html_script_filename"])
			schoolFileName := fmt.Sprintf("%s%s", rmap["school_level_script_path"], rmap["html_script_filename"])

			f, err := os.Create(schoolFileName)
			if err != nil {
				log.Printf("error creating file: %s error is: %s\n", schoolFileName, err)
				continue
			}

			scriptFileName := strings.TrimSuffix(rmap["html_script_filename"], ".html")
			scriptFileNameBanner := fmt.Sprintf("<h2 style=\"text-align: center;\">%s</h2>", scriptFileName)
			topFileName := scriptFileNameBanner
			bottomfileName := scriptFileNameBanner

			doc := []string{scriptHeader, topFileName, rmap["Item Response"], bottomfileName, scriptFooter}
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
