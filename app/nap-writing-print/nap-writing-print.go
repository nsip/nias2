// nap-writing-print.go

package main

import (
	"flag"
	"log"

	nwp "github.com/nsip/nias2/nap-writing-print"
)

func main() {

	// check any cli params
	flag.Parse()

	// clear directories
	nwp.ClearReportsDirectory()

	// find file(s)
	files := nwp.ParseResultsFileDirectory()

	for _, fileName := range files {
		err := nwp.RunPrintingPipeline(fileName)
		if err != nil {
			log.Fatal("error creating printing pipeline: ", err)
		}
	}

	// backup input file

	// report files written?
}
