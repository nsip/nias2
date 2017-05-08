package main

import (
	"flag"
	"github.com/nats-io/nats-streaming-server/server"
	"github.com/nsip/nias2/naprr"
	"log"
	"net"
	"os"
	// "os/exec"
	// "os/signal"
	//"runtime"
	// "time"
)

func main() {

	flag.Parse()

	wd, _ := os.Getwd()
	log.Println("working directory:", wd)

	log.Println("removing old files...")
	clearNSSWorkingDirectory()

	log.Println("Launching stream server...")
	ss := launchNatsStreamingServer()
	defer ss.Shutdown()

	log.Println("Starting data ingest...")

	di := naprr.NewDataIngest()
	di.Run()
	di.RunYr3Writing()
	student_ids := di.Yr3StudentIds
	naprr_config := di.NaprrConfig
	di.Close()

	rb := naprr.NewReportBuilder()
	// must run Year 3 Writing ingest before full XML: Full XML ingest generates map to reconcile student identities between the two
	log.Println("Generating report data, Year 3 Writing...")
	rb.RunYr3W(false, student_ids, naprr_config)

	log.Println("Writing report files...")
	rw := naprr.NewReportWriter()
	rw.WriteYr3WReports()
	log.Println("Report files Done")

	//runtime.Goexit()
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

	ss, _ := server.RunServerWithOpts(stanOpts, &nOpts)

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
