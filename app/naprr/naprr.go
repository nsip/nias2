package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/nats-io/nats-streaming-server/server"
	"github.com/nsip/nias2/naprr"
	// "os/exec"
	// "os/signal"
	"runtime"
	// "time"
)

var rewrite = flag.Bool("rewrite", false, "rewrite regenerates all reports without re-loading data")
var webonly = flag.Bool("webonly", false, "just launch web data explorer")
var yr3writing = flag.Bool("yr3writing", false, "read external yr3 writing files")

func main() {

	flag.Parse()

	wd, _ := os.Getwd()
	log.Println("working directory:", wd)

	if !*rewrite {
		log.Println("removing old files...")
		clearNSSWorkingDirectory()
	}

	log.Println("Launching stream server...")
	ss := launchNatsStreamingServer()
	defer ss.Shutdown()

	if !*rewrite {
		log.Println("Starting data ingest...")

		di := naprr.NewDataIngest()
		di.Run()

		//
		// only parse external yr3 if requested
		//
		// if *yr3writing {
		// 	log.Println("attempting to read yr3 writing files...")
		// 	di.RunYr3Writing()
		// }

		di.Close()

		rb := naprr.NewReportBuilder()
		log.Println("Generating report data...")
		rb.Run()

		//
		// report on yr3 only if requested
		//
		// if *yr3writing {
		// 	log.Println("Generating report data, Year 3 Writing...")
		// 	rb.RunYr3W(false)
		// }
	}

	log.Println("Writing report files...")
	rw := naprr.NewReportWriter()
	rw.Run()

	log.Println("Ingest and report writing complete.")

	log.Println("Starting web data browser server.")

	rrs := naprr.NewResultsReportingServer()
	rrs.Run()

	runtime.Goexit()
}

func clearNSSWorkingDirectory() {

	// remove existing logs and recreate the directory
	err := os.RemoveAll("nss")
	err = os.Mkdir("nss", os.ModePerm)
	if err != nil {
		log.Println("Error trying to remove nss working directory")
	}
}

func launchNatsStreamingServer() *server.StanServer {

	stanOpts := server.GetDefaultOptions()

	stanOpts.ID = "nap-rr"
	stanOpts.MaxChannels = 30000
	stanOpts.MaxMsgs = 2000000
	stanOpts.MaxBytes = 0 //unlimited
	stanOpts.MaxSubscriptions = 10000

	stanOpts.StoreType = "FILE"
	stanOpts.FilestoreDir = "nss"
	// stanOpts.Debug = true

	nOpts := server.DefaultNatsServerOptions
	nOpts.Port = 5222

	ss := server.RunServerWithOpts(stanOpts, &nOpts)

	return ss

}

// getAvailPort asks the OS for an unused port.
// There's a race here, where the port could be grabbed by someone else
// before the caller gets to Listen on it, but in practice such races
// are rare. Uses net.Listen("tcp", ":0") to determine a free port, then
// releases it back to the OS with Listener.Close().
func getAvailPort() int {
	l, _ := net.Listen("tcp", ":0")
	r := l.Addr()
	l.Close()
	return r.(*net.TCPAddr).Port
}
