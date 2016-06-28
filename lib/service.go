// service.go

// service interface to handle message requests
package nias2

import ()

// interface for all handlers,
// req -
// is a niasmessage containing metadata and payload
// request can be changed by a handler - allows for changes to payload such as
// creatiion of a guid, reformatting etc.
//
// resp - slice of nias messages, can be empty
// is the result of the processing activity, such as a validation item error reports.
// handlers do not have to emit responses.
//
// errros -
// will be returned if any system error occurs during the execution of the handler process
type NiasService interface {
	HandleMessage(req *NiasMessage) ([]NiasMessage, error)
}
