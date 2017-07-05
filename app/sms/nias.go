// harness runs the validation services and web server
package main

import (
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/sms"
	"log"
	"runtime"
)

func main() {

	config := lib.LoadDefaultConfig()
	SMS_NATS_CFG := lib.NATSConfig{Port: config.NATSPort}
	log.Println("SMS: Loading default config")
	log.Println("SMS: Config values are: ", config)
	log.Println("SMS: Starting ledis storage server...")
	go sms.LaunchStorageServer()
	log.Println("...ledis storage service running")

	log.Println("SMS: Starting ledis lookup server...")
	//go sms.LaunchLookupServer()
	log.Println("...ledis lookup service running")
	poolsize := config.PoolSize

	log.Println("SMS: Starting distributor....")
	dist := &sms.Distributor{}
	go dist.Run(poolsize, SMS_NATS_CFG)
	log.Println("...Distributor running")

	log.Println("SMS: Starting web services...")
	ws := &sms.NIASWebServer{}
	go ws.Run(SMS_NATS_CFG)
	log.Println("...web services running")

	runtime.Goexit()

}
