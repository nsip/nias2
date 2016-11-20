// harness.go

// harness runs the various services and web server
package main

import (
	"github.com/nsip/nias2/lib"
	// "github.com/pkg/profile"
	"log"
	"runtime"
)

func main() {

	// uncomment line (and import) below to invoke code profiler when running
	// note: macs must be running osx El Capitan for this to report correctly
	// defer profile.Start(profile.ProfilePath(".")).Stop()

	log.Println("Loading default config")
	log.Println("Config values are: ", nias2.NiasConfig)

	poolsize := nias2.NiasConfig.PoolSize

	log.Println("Starting ledis storage server...")
	go nias2.LaunchStorageServer()
	log.Println("...ledis storage service running")

	log.Println("Starting ledis lookup server...")
	go nias2.LaunchLookupServer()
	log.Println("...ledis lookup service running")

	log.Println("Loading ASL Lookup data")
	nias2.LoadASLLookupData()

	log.Println("Starting distributor....")
	dist := nias2.Distributor{}
	switch nias2.NiasConfig.MsgTransport {
	case "MEM":
		go dist.RunMemBus(poolsize)
	case "NATS":
		go dist.RunNATSBus2(poolsize)
	case "STAN":
		go dist.RunSTANBus(poolsize)
	default:
		go dist.RunMemBus(poolsize)
	}
	log.Println("...Distributor running")

	log.Println("Starting web services...")
	ws := &nias2.NIASWebServer{}
	go ws.Run()
	log.Println("...web services running")

	runtime.Goexit()

}
