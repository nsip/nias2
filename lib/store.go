// store.go
package nias2

import (
	"bytes"
	"encoding/gob"
	"github.com/siddontang/goredis"
	"log"
	"strconv"
	"sync"
)

const VALIDATION_PREFIX = "nvr:"

// MessageStore listens for messages on the store topic and captures them
// in ledis as lists (persistent qs in effect), messages can be stored on transaction
// and use case basis - use case being the superset of all transactions of a given type
type MessageStore struct{}

// MessageTracker listens for progress updates from services and
// updates progress monitor in ledis
type MessageTracker struct{}

// track status in simple hash counter; key is transactionid, value is no. records processed
var status = make(map[string]int)

// mutex to protect status hash for concurrent updates
var mutex = &sync.Mutex{}

// launch a set of storage agents that will listen for messages to store
func (ms *MessageStore) Run(poolsize int) {

	ms_ec := CreateNATSConnection()
	for i := 0; i < poolsize; i++ {
		c := CreateLedisConnection()
		// defer c.Close()
		ms_ec.QueueSubscribe(STORE_TOPIC, "msg_store", func(m *NiasMessage) {
			// don't store if no content
			if m.Body != nil {

				// store for txaction
				tx_key := m.Target + m.TxID
				_, err := c.Do("rpush", tx_key, EncodeNiasMessage(m))
				if err != nil {
					log.Println("error saving message:tx: - ", err)
				}

				// store for use case - disabled for now
				store_usecase := false
				if store_usecase {
					uc_key := m.Target
					_, err := c.Do("rpush", uc_key, EncodeNiasMessage(m))
					if err != nil {
						log.Println("error saving message:uc - ", err)
					}
				}

			}

		})
	}
}

// create a set of progress trackers to monitor services' progress
func (mt *MessageTracker) Run(poolsize int) {

	mt_ec := CreateNATSConnection()
	for i := 0; i < poolsize; i++ {

		mt_ec.QueueSubscribe(TRACK_TOPIC, "msg_tracker", func(txid string) {

			mutex.Lock()
			status[txid]++
			mutex.Unlock()

		})
	}

}

// Retrieve the data for this transaction - txid
// fulldata if true returns all data in the transaction, if false then return
// is capped at 10,000 records
func GetTxData(txid string, fulldata bool) ([]interface{}, error) {

	data := make([]interface{}, 0)
	c := CreateLedisConnection()
	defer c.Close()

	var endpoint int
	if fulldata {
		endpoint = -1
	} else {
		endpoint = (10000 - 1)
	}

	results, err := goredis.Values(c.Do("lrange", VALIDATION_PREFIX+txid, 0, endpoint))
	if err != nil {
		log.Println("Error fetching tx data for: ", txid)
		return data, err
	}

	for _, result := range results {
		r := result.([]uint8)
		msg := DecodeNiasMessage(r)
		data = append(data, msg.NiasData.Body)
	}

	return data, nil

}

// return a json string containing the tracking data
func GetTrackingData(txid string) map[string]string {

	trackmap := make(map[string]string)

	mutex.Lock()
	tx_size := strconv.Itoa(status[txid])
	mutex.Unlock()

	trackmap["Total"] = tx_size
	return trackmap
}

// binary encding for messages going to internal q/store, in nats qs this is
// handled automatically by the use of gob encoder on connection
func EncodeNiasMessage(msg *NiasMessage) []byte {

	encBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(encBuf)
	err := encoder.Encode(msg)
	if err != nil {
		log.Printf("Encoder unable to binary encode message for: %#v\n", msg)
	}
	return encBuf.Bytes()

}

// binary encding for messages coming from internal q/store, in nats qs this is
// handled automatically by the use of gob encoder on connection
func DecodeNiasMessage(bytemsg []uint8) *NiasMessage {

	decBuf := bytes.NewBuffer(bytemsg)
	decoder := gob.NewDecoder(decBuf)
	var msgOut NiasMessage
	err := decoder.Decode(&msgOut)
	if err != nil {
		log.Println("Error decoding message from q/store(internal):", err)
	}
	return &msgOut
}
