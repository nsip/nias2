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
			// doc.AddPage()

			//
			// print header id
			doc.SetColor("black").
				SetFont("Helvetica-Bold", 24).
				SetXY(3.5, 2).DrawText(rmap["Anonymised Id"])

			colText := splitParas(rmap["Item Response"])
			doc.SetColor("black").
				SetFont("Helvetica", 10).
				// DrawUnitGrid().
				DrawTextInBox(1.5, 2.5, 17, 28, "LT", colText)

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

//
// use the paragraph markers that have been left in the text
// to insert spacing characters.
//
func splitParas(fullText string) string {

	return strings.Replace(fullText, "<p>", "\n", -1)

}
