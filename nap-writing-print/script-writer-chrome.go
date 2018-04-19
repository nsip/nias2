// script-writer-chrome.go

package nap_writing_print

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
)

//
// creates the simple pdf output of the writing response
//
func createScriptWriterChrome(ctx context.Context, in <-chan map[string]string) (<-chan error, error) {

	var chromePath string
	switch runtime.GOOS {
	case "windows":
		chromePath = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"
	case "darwin":
		chromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	case "linux":
		chromePath = "/usr/bin/chromium-browser"
	}

	cmd := chromePath

	errc := make(chan error, 1)

	buffers := make(chan bool, 3)

	go func() {
		defer close(errc)
		var wg sync.WaitGroup

		for rMap := range in {
			rmap := rMap

			wg.Add(1)

			go func(rmap map[string]string) {
				// This channel will block until a job finishes
				buffers <- true

				yrlvlPdf := fmt.Sprintf("%s%s", rmap["yr_level_script_path"], rmap["pdf_filename"])
				schoolPdf := fmt.Sprintf("%s%s", rmap["school_level_script_path"], rmap["pdf_filename"])

				output := "--print-to-pdf=" + schoolPdf
				input := "./in/test.html"
				args := []string{"--headless", "--disable-gpu", output, input}

				if err := exec.Command(cmd, args...).Run(); err != nil {
					log.Println("writing error", err)
					// os.Exit(1)
				}

				err := copyFileContents(schoolPdf, yrlvlPdf)
				if err != nil {
					log.Println("unable to copy pdf to schools folder")
				}

				// output = "--print-to-pdf=" + outputFile2
				// args = []string{"--headless", "--disable-gpu", output, input}
				// if err := exec.Command(cmd, args...).Run(); err != nil {
				// 	log.Println("writing error", err)
				// 	// os.Exit(1)
				// }

				// As soon as we've finished our work, we discard the value
				<-buffers

				// runtime.Gosched()

				wg.Done()
			}(rmap)

			// wg.Wait()

			// log.Println("printed pdf")

			// log.Println(rmap["Item Response"])

			// // create the writing script pdfs (one in each output folder structure)
			// outputFile1 := fmt.Sprintf("%s%s", rmap["yr_level_script_path"], rmap["pdf_filename"])
			// outputFile2 := fmt.Sprintf("%s%s", rmap["school_level_script_path"], rmap["pdf_filename"])

			// var doc = pdf.NewPDF("A4") // create a new PDF using 'A4' page size
			// doc.SetUnits("cm")
			// // doc.AddPage()

			// //
			// // print header id
			// doc.SetColor("black").
			// 	SetFont("Helvetica-Bold", 24).
			// 	SetXY(3.5, 2).DrawText(rmap["Anonymised Id"])

			// colText := splitParas(rmap["Item Response"])
			// doc.SetColor("black").
			// 	SetFont("Helvetica", 10).
			// 	// DrawUnitGrid().
			// 	DrawTextInBox(1.5, 2.5, 17, 28, "LT", colText)

			// //
			// // save the files
			// doc.SaveFile(outputFile1)
			// doc.SaveFile(outputFile2)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
		wg.Wait()
	}()
	return errc, nil

}

//
// use the paragraph markers that have been left in the text
// to insert spacing characters.
//
// func splitParas(fullText string) string {

// 	return strings.Replace(fullText, "<p>", "\n", -1)

// }
