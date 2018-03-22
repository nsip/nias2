package main

// napcomp takes input file from student registration
// and output file from results reporting
// and does a comparison to find students not accounted for
// in both files.

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/nsip/nias2/napcomp"
	"github.com/nsip/nias2/version"
)

var vers = flag.Bool("version", false, "Reports version of NIAS distribution")

func main() {

	flag.Parse()

	// shutdown handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		napcomp.CloseDB()
		os.Exit(1)
	}()

	if *vers {
		fmt.Printf("NIAS: Version %s\n", version.TagName)
		os.Exit(1)
	}

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
