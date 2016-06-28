// serviceregister.go
package nias2

// simple hashtable of service handlers stored by name
// will be matched against required tasks in NasMessage.Route meta-data

import (
	"sync"
)

// simple thread-safe container for group of services that will be available
// to process messages passed from a distributor node
type ServiceRegister struct {
	sync.RWMutex
	registry map[string]NiasService
}

// creates a ServiceRegister with properly initilaised internal map
// processing services are stored with a name and the referenced NiasService
func NewServiceRegister() *ServiceRegister {
	reg := ServiceRegister{}
	reg.registry = make(map[string]NiasService)
	return &reg
}

// add a service to the registry with a name
func (sr *ServiceRegister) AddService(servicename string, service NiasService) {
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
func (sr *ServiceRegister) FindService(servicename string) NiasService {
	sr.RLock()
	defer sr.RUnlock()
	return sr.registry[servicename]
}
