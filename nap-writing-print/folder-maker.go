// folder-maker.go

package nap_writing_print

import (
	"context"
	"os"
)

//
// checks all referenced output folders in the record are
// created before any file writes are attempted
//
//
func createFolderMaker(ctx context.Context, in <-chan map[string]string) (
	<-chan map[string]string, <-chan error, error) {

	out := make(chan map[string]string)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for rMap := range in {
			rmap := rMap

			os.MkdirAll(rmap["yr_level_script_path"], os.ModePerm)
			os.MkdirAll(rmap["yr_level_audit_path"], os.ModePerm)
			os.MkdirAll(rmap["school_level_script_path"], os.ModePerm)
			os.MkdirAll(rmap["school_level_audit_path"], os.ModePerm)

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
