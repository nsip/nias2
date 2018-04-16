// output-path-injector.go

package nap_writing_print

import (
	"context"
	"fmt"
)

var basePath = "./out"

//
// accepts map[string]string of a csv record-row, adds the
// vrious output file-paths to which the eventual
// print artefacts will be written
//
func createOutputPathInjector(ctx context.Context, in <-chan map[string]string) (
	<-chan map[string]string, <-chan error, error) {

	out := make(chan map[string]string)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for rMap := range in {
			recordMap := rMap
			addFileNames(recordMap)
			addSchoolPaths(recordMap)
			addYearLevelPaths(recordMap)
			select {
			case out <- recordMap:
				// log.Printf("\nrecord read:\n\n%v\n\n", recordMap)
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil

}

//
// add the filenames for this record
//
func addFileNames(rmap map[string]string) {

	pdfFileName := fmt.Sprintf("%s_%s.pdf", rmap["Participation Code"], rmap["Anonymised Id"])

	rmap["pdf_filename"] = pdfFileName

}

//
// create output paths for year-level printing
//
func addYearLevelPaths(rmap map[string]string) {

	scriptPath := fmt.Sprintf("%s/yr-level/%s/script/", basePath, rmap["Test level"])
	auditPath := fmt.Sprintf("%s/yr-level/%s/audit/", basePath, rmap["Test level"])

	rmap["yr_level_script_path"] = scriptPath
	rmap["yr_level_audit_path"] = auditPath

}

//
// create output paths for school-level printing
//
func addSchoolPaths(rmap map[string]string) {

	scriptPath := fmt.Sprintf("%s/schools/%s/script/", basePath, rmap["ACARA ID"])
	auditPath := fmt.Sprintf("%s/schools/%s/audit/", basePath, rmap["ACARA ID"])

	rmap["school_level_script_path"] = scriptPath
	rmap["school_level_audit_path"] = auditPath

}
