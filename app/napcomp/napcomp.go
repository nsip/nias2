package main

// napcomp takes input file from student registration
// and output file from results reporting
// and does a comparison to find students not accounted for
// in both files.

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/nsip/nias2/napcomp"
)

func main() {

	// shutdown handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		napcomp.CloseDB()
		os.Exit(1)
	}()

	napcomp.IngestData()
	napcomp.WriteReports()
	// shut down
	napcomp.CloseDB()
	os.Exit(1)

	// wait for shutdown
	// should only be needed if processing hangs for some reason
	for {
		runtime.Gosched()
	}

}
