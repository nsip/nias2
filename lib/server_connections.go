// server_connections.go

//
// utility routines to standardise nats/stan server access
//

package lib

import (
	"fmt"
	"log"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nuid"
)

//
// config wrapper
//
type NATSConfig struct {
	Port string
}

//
// helper function to provide encoded connections for standard NATA
//
func CreateNATSConnection(cfg NATSConfig) *nats.EncodedConn {

	// var servers = "nats://localhost:4222, nats://localhost:5222, nats://localhost:6222" //cluster
	var servers = fmt.Sprintf("nats://localhost:%s", cfg.Port) // standalone

	nc, err := nats.Connect(servers)
	if err != nil {
		log.Fatalln("Unable to connect to gnatsd, services aborting...\n", err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Fatalln("Unable to connect to gnatsd, services aborting...\n", err)
	}
	return ec

}

//
// helper function for STAN conenctions
//
func CreateSTANConnection(cfg NATSConfig) stan.Conn {

	server_url := fmt.Sprintf("nats://localhost:%s", cfg.Port)
	sc, err := stan.Connect(NAP_VAL_CID, nuid.Next(), stan.NatsURL(server_url))
	if err != nil {
		log.Fatalln("Unable to establish connection to STAN server, aborting.", err)
	}
	return sc
}
