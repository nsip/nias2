// harness runs the validation services and web server
package main

import (
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/napval"
	"log"
	"runtime"
)

func main() {

	config := napval.LoadNAPLANConfig()
	NAPLAN_NATS_CFG := lib.NATSConfig{Port: config.NATSPort}
	log.Println("NAPVAL: Loading default config")
	log.Println("NAPVAL: Config values are: ", config)

	poolsize := config.PoolSize

	log.Println("NAPVAL: Loading ASL Lookup data")
	napval.LoadASLLookupData()

	log.Println("NAPVAL: Starting distributor....")
	dist := &napval.ValidationDistributor{}
	go dist.Run(poolsize, NAPLAN_NATS_CFG)
	log.Println("...Distributor running")

	log.Println("NAPVAL: Starting web services...")
	ws := &napval.ValidationWebServer{}
	go ws.Run(NAPLAN_NATS_CFG)
	log.Println("...web services running")

	runtime.Goexit()

}
