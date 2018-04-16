// audit-writer.go

package nap_writing_print

import (
	"context"
	"fmt"

	pdf "github.com/balacode/one-file-pdf"
)

//
// creates the meta-data audit file to accompany the writing response
//
func createAuditWriter(ctx context.Context, in <-chan map[string]string) (<-chan error, error) {

	// out := make(chan map[string]string)
	errc := make(chan error, 1)

	go func() {
		defer close(errc)

		for rMap := range in {
			rmap := rMap

			// create the writing script pdfs (one in each output folder structure)
			outputFile1 := fmt.Sprintf("%s%s", rmap["yr_level_audit_path"], rmap["pdf_filename"])
			outputFile2 := fmt.Sprintf("%s%s", rmap["school_level_audit_path"], rmap["pdf_filename"])

			var doc = pdf.NewPDF("A4") // create a new PDF using 'A4' page size
			doc.SetUnits("cm")

			//
			// print header id
			doc.SetColor("black").
				SetFont("Helvetica-Bold", 24).
				SetXY(3.5, 2.7).DrawText(rmap["Anonymised Id"])
			//
			// draw the metadata text
			doc.SetColor("black").
				// DrawUnitGrid().
				SetFont("Helvetica-Bold", 18).
				SetXY(3.5, 5).
				DrawText(fmt.Sprintf("Participation Code:   %s", rmap["Participation Code"])).
				SetXY(3.5, 6).
				DrawText(fmt.Sprintf("TAA Student ID:   %s", rmap["TAA student ID"])).
				SetXY(3.5, 7).
				DrawText(fmt.Sprintf("PSI:   %s", rmap["PSI"])).
				SetXY(3.5, 8).
				DrawText(fmt.Sprintf("Local School Student ID:   %s", rmap["Local school student ID"])).
				SetXY(3.5, 9).
				DrawText(fmt.Sprintf("ACARA School ID:   %s", rmap["ACARA ID"])).
				SetXY(3.5, 10).
				DrawText(fmt.Sprintf("Test Year:   %s", rmap["Test Year"])).
				SetXY(3.5, 11).
				DrawText(fmt.Sprintf("Test Level:   %s", rmap["Test level"])).
				SetXY(3.5, 12).
				DrawText(fmt.Sprintf("Jurisdiction Identifier:   %s", rmap["Jurisdiction Id"]))
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
