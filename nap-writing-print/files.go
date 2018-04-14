// nap-writing-print.go

package nap_writing_print

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
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
// Creates a timestamped backup of the input files
// used to create a run of pdf files.
//
// This so that if the files are distributed, and then the tool is re-run
// the output files will be over-written.
// The original files can be used to recerate the pdf and audit files
// again.
//
func BackupInputFiles() {

	paths := ParseResultsFileDirectory()

	now := time.Now()
	timestamp := fmt.Sprintf("%d-%d-%d-%02d%02d%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())

	backupPath := "./backup/" + timestamp

	for _, sourceFilePath := range paths {
		fileName := filepath.Base(sourceFilePath)
		dstFilePath := backupPath + "/" + fileName
		err := os.MkdirAll(backupPath, os.ModePerm)
		if err != nil {
			log.Println("Error creating input backup directory: ", err)
		}
		err = copyFileContents(sourceFilePath, dstFilePath)
		if err != nil {
			log.Println("Error trying to create backup of input files: ", err)
		}
	}

	log.Println("backup of input files created...")

}

//
// underlying file copy
//
func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
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
