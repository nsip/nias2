// nats.go

package lib

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nuid"
	"log"
)

// create nats connections, and useful constants for topics etc.
//

const REQUEST_TOPIC = "requests"
const STORE_TOPIC = "store"
const TRACK_TOPIC = "track"
const NIAS_CLUSTER_ID = "nias"
const NAP_VAL_CID = "nap-val"

var req_chan = make(chan NiasMessage, 1)

/*
ProcessChains provide conceptual links between the major
components of nias:

inbound (distributor/dist) connections
	represetnitng entry points for data such as a web gateway or file reader

service (srvc) connections
	multiplexing inbound data to service handlers

storage (store) connections
	feeding processed data to storage; streams or database

*/

// ProcessChain when running with NATS Streaming Server, known as STAN
type STANProcessChain struct {
	dist_in_conn     stan.Conn
	dist_in_subject  string
	dist_out_conn    stan.Conn
	dist_out_subject string

	srvc_in_conn     stan.Conn
	srvc_in_subject  string
	srvc_out_conn    stan.Conn
	srvc_out_subject string

	store_in_conn    stan.Conn
	store_in_subject string
}

// ProcessChain for running with standard NATS
type NATSProcessChain struct {
	dist_in_conn     *nats.EncodedConn
	dist_in_subject  string
	dist_out_conn    *nats.EncodedConn
	dist_out_subject string

	srvc_in_conn     *nats.EncodedConn
	srvc_in_subject  string
	srvc_out_conn    *nats.EncodedConn
	srvc_out_subject string

	store_in_conn    *nats.EncodedConn
	store_in_subject string
}

type NATSConfig struct {
	Port string
}

// ProcessChain for running in memory - useful for standaloe use
// and for resource constrained enviornments, use of blocking channels accross
// the chain means solution will balance processing accross the chain based on
// speed of host
// Channel of struct rather than of pointers to struct: pointer channel was not
// updating properly when read out
type MemProcessChain struct {
	req_chan   chan NiasMessage
	srvc_chan  chan NiasMessage
	store_chan chan NiasMessage
}

func NewMemProcessChain() MemProcessChain {
	mpc, _ := createMemProcessChain()
	return mpc
}

func createMemProcessChain() (MemProcessChain, error) {

	pc := MemProcessChain{}

	pc.req_chan = req_chan
	pc.srvc_chan = make(chan NiasMessage, 1)
	pc.store_chan = make(chan NiasMessage, 1)

	return pc, nil

}

func NewNATSProcessChain(cfg NATSConfig) NATSProcessChain {
	npc, _ := createNATSProcessChain(cfg)
	return npc
}

func createNATSProcessChain(cfg NATSConfig) (NATSProcessChain, error) {

	pc := NATSProcessChain{}

	ec := CreateNATSConnection(cfg)

	distID := nuid.Next()
	pc.dist_in_conn = ec
	pc.dist_out_conn = ec
	pc.dist_in_subject = REQUEST_TOPIC
	pc.dist_out_subject = distID

	srvcID := nuid.Next()
	pc.srvc_in_conn = ec
	pc.srvc_out_conn = ec
	pc.srvc_in_subject = distID
	pc.srvc_out_subject = srvcID

	pc.store_in_conn = ec
	pc.store_in_subject = srvcID

	return pc, nil
}

// helper function to provide encoded connections for standard NATA
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

func NewSTANProcessChain() STANProcessChain {
	spc, _ := createSTANProcessChain()
	return spc
}

func createSTANProcessChain() (STANProcessChain, error) {

	pc := STANProcessChain{}

	distID := nuid.Next()
	pc.dist_in_conn, _ = stan.Connect(NIAS_CLUSTER_ID, nuid.Next())
	pc.dist_out_conn, _ = stan.Connect(NIAS_CLUSTER_ID, nuid.Next())
	pc.dist_in_subject = REQUEST_TOPIC
	pc.dist_out_subject = distID

	srvcID := nuid.Next()
	pc.srvc_in_conn, _ = stan.Connect(NIAS_CLUSTER_ID, nuid.Next())
	pc.srvc_out_conn, _ = stan.Connect(NIAS_CLUSTER_ID, nuid.Next())
	pc.srvc_in_subject = distID
	pc.srvc_out_subject = srvcID

	pc.store_in_conn, _ = stan.Connect(NIAS_CLUSTER_ID, nuid.Next())
	pc.store_in_subject = srvcID

	return pc, nil
}
