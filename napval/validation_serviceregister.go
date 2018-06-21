// val_serviceregister.go
package napval

// simple hashtable of service handlers stored by name
// will be matched against required tasks in NasMessage.Route meta-data

import (
	"log"
	"sync"

	"github.com/nsip/nias2/lib"
)

var config = LoadNAPLANConfig()

// simple thread-safe container for group of services that will be available
// to process messages passed from a distributor node
type ServiceRegister struct {
	sync.RWMutex
	registry map[string]lib.NiasService
}

// creates a ServiceRegister with properly initilaised internal map
// processing services are stored with a name and the referenced NiasService
func NewServiceRegister(nats_cfg lib.NATSConfig) *ServiceRegister {
	return createDefaultServiceRegister(nats_cfg)
}

// add a service to the registry with a name
func (sr *ServiceRegister) AddService(servicename string, service lib.NiasService) {
	sr.Lock()
	sr.registry[servicename] = service
	sr.Unlock()
}

// remove a service from the registry by name
func (sr *ServiceRegister) RemoveService(servicename string) {
	sr.Lock()
	delete(sr.registry, servicename)
	sr.Unlock()
}

// return a service by providing the name
func (sr *ServiceRegister) FindService(servicename string) lib.NiasService {
	sr.RLock()
	defer sr.RUnlock()
	return sr.registry[servicename]
}

// build register with default set of services
func createDefaultServiceRegister(nats_cfg lib.NATSConfig) *ServiceRegister {

	log.Println("Creating services & register")
	sr := ServiceRegister{}
	sr.registry = make(map[string]lib.NiasService)

	schema1, err := NewCoreSchemaService()
	if err != nil {
		log.Fatal("Unable to create schema service ", err)
	}

	schema11, err := NewCustomSchemaService("core_parent2.json")
	if err != nil {
		log.Fatal("Unable to create schema service ", err)
	}

	schema2, err := NewCustomSchemaService("local.json")
	if err != nil {
		log.Fatal("Unable to create schema service ", err)
	}

	id1, err := NewIDService3(nats_cfg)
	if err != nil {
		log.Fatal("Unable to create id service ", err)
	}

	dob1, err := NewDOBService(config.TestYear)
	if err != nil {
		log.Fatal("Unable to create dob service ", err)
	}

	asl1, err := NewASLService()
	if err != nil {
		log.Fatal("Unable to create asl service ", err)
	}

	psi1, err := NewPsiService()
	if err != nil {
		log.Fatal("Unable to create psi service ", err)
	}

	num, err := NewNumericValidService()
	if err != nil {
		log.Fatal("Unable to create numeric validation service ", err)
	}

	nam, err := NewNameValidService()
	if err != nil {
		log.Fatal("Unable to create name validation service ", err)
	}

	sr.AddService("schema", schema1)
	sr.AddService("schema2", schema11)
	sr.AddService("local", schema2)
	sr.AddService("id", id1)
	sr.AddService("dob", dob1)
	sr.AddService("asl", asl1)
	sr.AddService("psi", psi1)
	sr.AddService("numericvalid", num)
	sr.AddService("namevalid", nam)

	log.Println("services created & installed in register")

	return &sr

}

func (sr *ServiceRegister) ProcessByRoute(m *lib.NiasMessage) []lib.NiasMessage {

	response_msgs := make([]lib.NiasMessage, 0)

	route := m.Route

	// log.Printf("\t\tservice register recieved msg: %+v", m)

	for _, sname := range route {

		// retrieve service from registry & execute
		srvc := sr.FindService(sname)
		responses, err := srvc.HandleMessage(m)
		if err != nil {
			log.Println("\t *** got an error on service handler " + sname + " ***")
			log.Println("\t", err)
		} else {
			// pass the responses to the message store
			// log.Printf("\t\tservice %s returned %d responses: %+v", sname, len(responses), responses)
			for _, r := range responses {
				response := r
				response.Source = sname
				response_msgs = append(response_msgs, response)
			}
		}

		//
		// enable this for low-cpu machines, drastically increases
		// processing time, but reduces cpu load to around 30%
		// may be a worthwhile tradeoff for some users
		//
		// time.Sleep(1 * time.Millisecond)
	}

	//log.Printf("\t\tresponse messages: %+v", response_msgs)
	return response_msgs

}
