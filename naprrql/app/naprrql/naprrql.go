package main

// Alternative naprr server offering graphql data feeds and KV store backend
// to try and elimnate issues with tcp, memory and disk access
// on Win32 that hamper all other deployments so far.

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/matt-farmer/naprrql"
)

var ingest = flag.Bool("ingest", false, "forces re-load of data from results file.")

func main() {

	flag.Parse()

	// shutdown handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		closeDB()
		os.Exit(1)
	}()

	// ingest results data, and exit to save memory
	if *ingest {
		clearDBWorkingDirectory()
		naprrql.IngestResultsFile("master_nap.xml.zip")
		closeDB()
		os.Exit(1)
	} else {
		// launch sif-ql service
		go naprrql.RunQLServer()
		// launch ui service
	}

	// wait for shutdown
	for {
		runtime.Gosched()
	}

}

//
// ensure clean shutdown of data store
//
func closeDB() {
	log.Println("Closing datastore...")
	naprrql.GetDB().Close()
	log.Println("Datastore closed.")
}

func clearDBWorkingDirectory() {

	// remove existing logs and recreate the directory
	err := os.RemoveAll("kvs")
	if err != nil {
		log.Println("Error trying to reset datastore working directory: ", err)
	}
	createDBWorkingDirectory()
}

func createDBWorkingDirectory() {
	err := os.Mkdir("kvs", os.ModePerm)
	if !os.IsExist(err) && err != nil {
		log.Fatalln("Error trying to create datastore working directory: ", err)
	}

}
