// validationstore.go
package napval

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nuid"
	"github.com/nsip/nias2/lib"
	"log"
)

// amount of error reports to store for any given input file
var STORE_LIMIT = config.TxStorageLimit

//var STORE_LIMIT = DefaultValidationConfig.TxStorageLimit

// ValidationStore assigns messages (validation results) to
// output streams for retrieval by clients
type ValidationStore struct {
	C         stan.Conn
	TxCounter map[string]int
	TxLimit   map[string]bool
}

// Returns a ValidationStore with an active connection
// to the stan server.
func NewValidationStore() *ValidationStore {
	sc, err := stan.Connect(lib.NAP_VAL_CID, nuid.Next())
	if err != nil {
		log.Fatalln("Unable to establish storage connection to STAN server, aborting.", err)
	}
	vs := ValidationStore{C: sc, TxCounter: make(map[string]int), TxLimit: make(map[string]bool)}
	return &vs
}

// put message into stan store/stream
// endcode converts nias message to byte array for storage
func (ms *ValidationStore) StoreMessage(msg *lib.NiasMessage) {

	// store for txaction
	// respecting limits of no. error reports for any
	// transaction

	// var storage_limit_reached bool

	switch t := msg.Body.(type) {
	case ValidationError:
		ms.TxCounter[msg.TxID]++
		if !ms.TxLimit[msg.TxID] {
			err := ms.C.Publish(msg.TxID, lib.EncodeNiasMessage(msg))
			if err != nil {
				log.Println("publish to store error: ", err)
			}
		}
		if ms.TxCounter[msg.TxID] >= STORE_LIMIT {
			ms.TxLimit[msg.TxID] = true
		}
	case lib.TxStatusUpdate:
		err := ms.C.Publish(msg.TxID, lib.EncodeNiasMessage(msg))
		if err != nil {
			log.Println("publish to store error: ", err)
		}
	default:
		_ = t
		log.Printf("unknown message type in storage handler: %v", msg)
	}

	// err := ms.C.Publish(m.TxID, lib.EncodeNiasMessage(m))
	// if err != nil {
	// 	log.Println("publish to store error: ", err)
	// }

}
