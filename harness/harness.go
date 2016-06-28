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

	log.Println("Starting storage engine...")
	go nias2.LaunchLedisServer()
	log.Println("...Storage service running")

	log.Println("Loading ASL Lookup data")
	nias2.LoadASLLookupData()

	log.Println("Starting distributor....")
	dist := nias2.Distributor{}
	go dist.Run(poolsize)
	log.Println("...Distributor running")

	log.Println("Starting storage agent...")
	msg_store := &nias2.MessageStore{}
	go msg_store.Run(poolsize)
	log.Println("...Storage engine running")

	log.Println("Starting progress tracker...")
	msg_trk := &nias2.MessageTracker{}
	go msg_trk.Run(poolsize)
	log.Println("...progress tracker running")

	log.Println("Starting web services...")
	ws := &nias2.NIASWebServer{}
	go ws.Run()
	log.Println("...web services running")

	runtime.Goexit()

}
