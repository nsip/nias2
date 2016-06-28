// niasmessage.go

// Message wrapper types to embed metadata for nias
// service handlers
package nias2

// meta-data for all messages
type NiasContext struct {
	MsgID  string   // guid for this message (if required)
	SeqNo  string   // orignal sequence no. of this message e.g. if read from a file
	TxID   string   // the transactionID associated with this message
	Target string   // the use case/url to be used to retrieve the processed message or results of processing the message
	Route  []string // slice/array of named services to pass this message through
	Source string   // service that produced this event
}

// the content of the message
type NiasData struct {
	Body interface{} // the payload of the message
}

// type to bundle both aspects together
type NiasMessage struct {
	NiasContext
	NiasData
}
