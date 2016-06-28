// nats.go

package nias2

import (
	"github.com/nats-io/nats"
	"log"
)

// create nats connections, and useful constants for topics etc.
//

const REQUEST_TOPIC = "requests"
const STORE_TOPIC = "store"
const TRACK_TOPIC = "track"

func CreateNATSConnection() *nats.EncodedConn {
	nc, err := nats.Connect(nats.DefaultURL)
	ec, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Fatalln("Unable to connect to gnatsd, services aborting...\n")
	}
	return ec

}
