// nap-writing-sanitiser.go

package main

import (
	"log"

	nws "github.com/nsip/nias2/nap-writing-sanitiser"
)

func main() {

	log.Println("starting html writer...")

	// clear directories
	nws.ClearReportsDirectory()

	// find file(s)
	files := nws.ParseResultsFileDirectory()

	// send them through the processing pipeline
	for _, fileName := range files {
		err := nws.RunPrintingPipeline(fileName)
		if err != nil {
			log.Fatal("error creating printing pipeline: ", err)
		}
	}

	// backup input file
	nws.BackupInputFiles()

	log.Println("...all CSV files written.")
}
