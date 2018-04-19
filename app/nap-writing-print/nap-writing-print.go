// nap-writing-print.go

package main

import (
	"log"

	nwp "github.com/nsip/nias2/nap-writing-print"
)

func main() {

	log.Println("starting html writer...")

	// clear directories
	nwp.ClearReportsDirectory()

	// find file(s)
	files := nwp.ParseResultsFileDirectory()

	// send them through the processing pipeline
	for _, fileName := range files {
		err := nwp.RunPrintingPipeline(fileName)
		if err != nil {
			log.Fatal("error creating printing pipeline: ", err)
		}
	}

	// backup input file
	nwp.BackupInputFiles()

	log.Println("...all html files written.")
}
