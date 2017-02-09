// harness runs the validation services and web server
package main

import (
	"github.com/nsip/nias2/napval"
	"log"
	"runtime"
)

func main() {

	log.Println("Loading default config")
	log.Printf("Config values are: \n%v\n\n", napval.DefaultValidationConfig)

	poolsize := napval.DefaultValidationConfig.PoolSize

	log.Println("Loading ASL Lookup data")
	napval.LoadASLLookupData()

	log.Println("Starting distributor....")
	dist := &napval.ValidationDistributor{}
	go dist.Run(poolsize)
	log.Println("...Distributor running")

	log.Println("Starting web services...")
	ws := &napval.ValidationWebServer{}
	go ws.Run()
	log.Println("...web services running")

	runtime.Goexit()

}
