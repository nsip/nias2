// distributor.go

// distributor creates parallel service access on demand
// creates a pool of listeners on a nats q to manage requests in parallel
// requests are routed sequentially through the service chain provided
// in NiasMessage.Route
package nias2

import (
	"github.com/nats-io/nats"
	"log"
)

type Distributor struct{}

func createServiceRegister() *ServiceRegister {

	log.Println("Creating services & register")
	sr := NewServiceRegister()

	schema1, err := NewCoreSchemaService()
	if err != nil {
		log.Fatal("Unable to create schema service ", err)
	}

	schema2, err := NewCustomSchemaService("local.json")
	if err != nil {
		log.Fatal("Unable to create schema service ", err)
	}

	id1, err := NewIDService()
	if err != nil {
		log.Fatal("Unable to create id service ", err)
	}

	dob1, err := NewDOBService(NiasConfig.TestYear)
	if err != nil {
		log.Fatal("Unable to create dob service ", err)
	}

	asl1, err := NewASLService()
	if err != nil {
		log.Fatal("Unable to create asl service ", err)
	}

	sr.AddService("schema", schema1)
	sr.AddService("local", schema2)
	sr.AddService("id", id1)
	sr.AddService("dob", dob1)
	sr.AddService("asl", asl1)

	log.Println("services created & installed in register")

	return sr

}

func serviceHandler(m *NiasMessage, sr *ServiceRegister, dist_ec *nats.EncodedConn) {

	route := m.Route

	for _, sname := range route {

		// retrieve service from registry & execute
		srvc := sr.FindService(sname)
		responses, err := srvc.HandleMessage(m)
		if err != nil {
			log.Println("\t *** got an error on service handler " + sname + " ***")
			log.Println("\t", err)
		} else {
			// pass the responses to the message store
			for _, r := range responses {
				response := r
				response.Source = sname
				err := dist_ec.Publish(STORE_TOPIC, response)
				if err != nil {
					log.Println("Error saving response to message store: ", err)
				}
			}
		}
	}
	// update the progress tracker
	err := dist_ec.Publish(TRACK_TOPIC, m.TxID)
	if err != nil {
		log.Println("Error saving tracking data: ", err)
	}

}

// creates a pool of message handlers which process the
// routing slip of each message thru the listed services
func (d *Distributor) Run(poolsize int) {

	for i := 0; i < poolsize; i++ {

		sr := createServiceRegister()
		dist_ec := CreateNATSConnection()
		dist_ec.QueueSubscribe(REQUEST_TOPIC, "distributor", func(m *NiasMessage) {
			serviceHandler(m, sr, dist_ec)
		})
	}
}
