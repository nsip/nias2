// script-writer.go

package nap_writing_print

import (
	"context"
	"fmt"
	"strings"

	pdf "github.com/balacode/one-file-pdf"
)

//
// creates the simple pdf output of the writing response
//
func createScriptWriter(ctx context.Context, in <-chan map[string]string) (<-chan error, error) {

	errc := make(chan error, 1)

	go func() {
		defer close(errc)

		for rMap := range in {
			rmap := rMap

			// create the writing script pdfs (one in each output folder structure)
			outputFile1 := fmt.Sprintf("%s%s", rmap["yr_level_script_path"], rmap["pdf_filename"])
			outputFile2 := fmt.Sprintf("%s%s", rmap["school_level_script_path"], rmap["pdf_filename"])

			var doc = pdf.NewPDF("A4") // create a new PDF using 'A4' page size
			doc.SetUnits("cm")

			//
			// print header id
			doc.SetColor("black").
				SetFont("Helvetica-Bold", 24).
				SetXY(3.5, 2.7).DrawText(rmap["Anonymised Id"])
			//
			// draw the column of text
			var colText = strings.Replace(rmap["Item Response"], "\\n", " ", -1)
			doc.SetColor("black").
				SetFont("Helvetica", 12).
				// DrawUnitGrid().
				DrawTextInBox(3.5, 4, 12, 28, "LT", colText)
			//
			// save the files
			doc.SaveFile(outputFile1)
			doc.SaveFile(outputFile2)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	return errc, nil

}
