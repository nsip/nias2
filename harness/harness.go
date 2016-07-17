// harness.go

// harness runs the various services and web server
package main

import (
	"github.com/nsip/nias2/lib"
	"log"
	"runtime"
)

func main() {

	log.Println("Loading default config")
	log.Println("Config values are: ", nias2.NiasConfig)

	poolsize := nias2.NiasConfig.PoolSize

	log.Println("Starting ledis server...")
	go nias2.LaunchLedisServer()
	log.Println("...ledis service running")

	log.Println("Loading ASL Lookup data")
	nias2.LoadASLLookupData()

	log.Println("Starting distributor....")
	dist := nias2.Distributor{}
	switch nias2.NiasConfig.MsgTransport {
	case "MEM":
		go dist.RunMemBus(poolsize)
	case "NATS":
		go dist.RunNATSBus(poolsize)
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
