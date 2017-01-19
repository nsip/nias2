// validation_distributor.go

// distributor creates parallel service access on demand
// creates a pool of listeners on a nats q to manage requests in parallel
// requests are routed sequentially through the service chain provided
// in NiasMessage.Route
package napval

import (
	"github.com/nsip/nias2/lib"
)

// creates a pool of message handlers which process the
// routing slip of each message thru the listed services
type ValidationDistributor struct{}

// uses single process for all streaming updates
// given the write-heavy nature of storage activities
// this is significantly faster than shared read/write
// services in parallel
func (vd *ValidationDistributor) Run(poolsize int) {

	ec := lib.CreateNATSConnection()
	vs := NewValidationStore()
	tt := lib.NewTransactionTracker()

	for i := 0; i < poolsize; i++ {

		// create service handler
		go func(vs *ValidationStore, tt *lib.TransactionTracker) {
			sr := NewServiceRegister()
			ec.QueueSubscribe(lib.REQUEST_TOPIC, "distributor", func(m *lib.NiasMessage) {
				responses := sr.ProcessByRoute(m)
				for _, response := range responses {
					r := response
					ec.Publish(lib.STORE_TOPIC, r)
				}
				tt.IncrementTracker(m.TxID)
				// get status of transaction and add message to stream
				// if a notable status change has occurred
				sigChange, msg := tt.GetStatusReport(m.TxID)
				if sigChange {
					ec.Publish(lib.STORE_TOPIC, msg)
				}
			})

		}(vs, tt)

	}

	// create storage handler
	go func(vs *ValidationStore) {

		ec.Subscribe(lib.STORE_TOPIC, func(m *lib.NiasMessage) {
			vs.StoreMessage(m)

		})

	}(vs)

}

//
//
//
//
//
//
//
//
//
//
//
