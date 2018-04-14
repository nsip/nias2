// nap-writing-print.go

package nap_writing_print

import (
	"log"
	"os"
	"path/filepath"
)

//
// support functions for the nap-writing-print app
// for all file-system activities; cerating folders
// setting filenames etc.
//

//
// look for writing results data files
//
func ParseResultsFileDirectory() []string {

	files := make([]string, 0)

	csvFiles, err := filepath.Glob("./in/*.csv")

	files = append(files, csvFiles...)
	if len(files) == 0 {
		log.Fatalln("No writng results data *.csv files found in input folder /in.", err)
	}

	return files

}

//
// remove reports working directory
//
func ClearReportsDirectory() {
	// remove existing logs and recreate the directory
	err := os.RemoveAll("out")
	if err != nil {
		log.Println("Error trying to reset reports directory: ", err)
	}
	createReportsDirectory()

}

//
// create folder for .csv reports
//
func createReportsDirectory() {
	err := os.Mkdir("out", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create reports directory: ", err)
	}

}
