// distributor.go

// distributor creates parallel service access on demand
// creates a pool of listeners on a nats q to manage requests in parallel
// requests are routed sequentially through the service chain provided
// in NiasMessage.Route
package nias2

import (
	"github.com/nats-io/go-nats-streaming"
	"log"
)

type Distributor struct{}

// creates a pool of message handlers which process the
// routing slip of each message thru the listed services

// Use STAN as message bus
func (d *Distributor) RunSTANBus(poolsize int) {

	for i := 0; i < poolsize; i++ {

		sr := NewServiceRegister()
		pc := NewSTANProcessChain()
		ms := NewMessageStore()

		// create storage handler
		go func(pc STANProcessChain, ms *MessageStore, id int) {

			pc.store_in_conn.Subscribe(pc.store_in_subject, func(m *stan.Msg) {
				msg := DecodeNiasMessage(m.Data)
				if msg.Target == SIF_MEMORY_STORE_PREFIX {
					ms.StoreGraph(msg)
				} else {
					ms.StoreMessage(msg)
				}
			})

		}(pc, ms, i) //drop ids

		// create service handler
		go func(pc STANProcessChain, sr *ServiceRegister, ms *MessageStore, id int) {

			pc.srvc_in_conn.Subscribe(pc.srvc_in_subject, func(m *stan.Msg) {
				msg := DecodeNiasMessage(m.Data)
				if msg.Target == STORE_AND_FORWARD_PREFIX {
					responses1 := sr.ProcessByPrivacy(msg)
					for _, response1 := range responses1 {
						pc.srvc_out_conn.Publish(pc.store_in_subject, EncodeNiasMessage(&response1))
					}
				}
				responses := sr.ProcessByRoute(msg)
				for _, response := range responses {
					pc.srvc_out_conn.Publish(pc.store_in_subject, EncodeNiasMessage(&response))
				}
				ms.IncrementTracker(msg.TxID)
			})

		}(pc, sr, ms, i)

		// create an inbound handler to multiplex validation requests, create last or will drop messages
		go func(pc STANProcessChain, id int) {

			pc.dist_in_conn.QueueSubscribe(pc.dist_in_subject, "distributor", func(m *stan.Msg) {
				pc.dist_out_conn.Publish(pc.dist_out_subject, m.Data)
			})

		}(pc, i)

	}
}

// Use regular NATS as message bus
func (d *Distributor) RunNATSBus(poolsize int) {

	for i := 0; i < poolsize; i++ {

		i := i

		sr := NewServiceRegister()
		pc := NewNATSProcessChain()
		ms := NewMessageStore()

		log.Printf("Process Chain %d\n%#v", i, pc)

		// create storage handler
		go func(pc NATSProcessChain, ms *MessageStore) {

			pc.store_in_conn.Subscribe(pc.store_in_subject, func(m *NiasMessage) {
				if m.Target == SIF_MEMORY_STORE_PREFIX {
					ms.StoreGraph(m)
				} else {
					ms.StoreMessage(m)
				}
			})

		}(pc, ms)

		// create service handler
		go func(pc NATSProcessChain, sr *ServiceRegister, ms *MessageStore) {

			pc.srvc_in_conn.Subscribe(pc.srvc_in_subject, func(m *NiasMessage) {
				if m.Target == STORE_AND_FORWARD_PREFIX {
					responses1 := sr.ProcessByPrivacy(m)
					for _, response1 := range responses1 {
						pc.srvc_out_conn.Publish(pc.store_in_subject, response1)
					}
				}
				responses := sr.ProcessByRoute(m)
				for _, response := range responses {
					pc.srvc_out_conn.Publish(pc.store_in_subject, response)
				}
				ms.IncrementTracker(m.TxID)
			})

		}(pc, sr, ms)

		// create an inbound handler to multiplex validation requests, create last or will drop messages
		go func(pc NATSProcessChain) {

			pc.dist_in_conn.QueueSubscribe(pc.dist_in_subject, "distributor", func(m *NiasMessage) {
				pc.dist_out_conn.Publish(pc.dist_out_subject, m)
			})

		}(pc)

	}
}

// use internal channels as message bus
func (d *Distributor) RunMemBus(poolsize int) {

	for i := 0; i < poolsize; i++ {

		i := i

		sr := NewServiceRegister()
		pc := NewMemProcessChain()
		ms := NewMessageStore()

		log.Printf("Process Chain %d\n%#v", i, pc)

		// create storage handler
		go func(pc MemProcessChain, ms *MessageStore) {

			var msg NiasMessage
			for {
				msg = <-pc.store_chan
				if msg.Target == SIF_MEMORY_STORE_PREFIX {
					ms.StoreGraph(&msg)
				} else {
					ms.StoreMessage(&msg)
				}
				//log.Printf("\t>%v %s %s\n", msg.Target, msg.MsgID, msg.SeqNo)
			}

		}(pc, ms)

		// create service handler
		go func(pc MemProcessChain, sr *ServiceRegister, ms *MessageStore) {

			for {
				msg := <-pc.req_chan
				if msg.Target == STORE_AND_FORWARD_PREFIX {
					//log.Printf("#%v %s %s\n", msg.Target, msg.MsgID, msg.SeqNo)
					responses1 := sr.ProcessByPrivacy(&msg)
					for _, response1 := range responses1 {
						pc.store_chan <- response1
						//log.Printf("<%v %d %s %s\n", response1.Target, i, response1.MsgID, response1.SeqNo)
					}
				}
				responses := sr.ProcessByRoute(&msg)
				for _, response := range responses {
					pc.store_chan <- response
				}
				ms.IncrementTracker(msg.TxID)
			}

		}(pc, sr, ms)

	}

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
