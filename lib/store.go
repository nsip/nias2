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
const STORE_AND_FORWARD_PREFIX = "ssf:"
const SIF_MEMORY_STORE_PREFIX = "sms:"

// MessageStore listens for messages on the store topic and captures them
// in ledis as lists (persistent qs in effect), messages can be stored on transaction
// and use case basis - use case being the superset of all transactions of a given type
type MessageStore struct {
	C    *goredis.Client
	trkC *goredis.Client
}

func NewMessageStore() *MessageStore {
	ms := MessageStore{
		C: CreateLedisConnection(1024, 1024),
	}
	return &ms
}

// track status in simple hash counter; key is transactionid, value is no. records processed
var status = make(map[string]int)

// mutex to protect status hash for concurrent updates
var mutex = &sync.Mutex{}

// put message into ledis store
// endcode converts nias message to byte array for storage
func (ms *MessageStore) StoreMessage(m *NiasMessage) {

	// store for txaction
	tx_key := m.Target + m.TxID
	_, err := ms.C.Do("rpush", tx_key, EncodeNiasMessage(m))
	if err != nil {
		log.Println("error saving message:tx: - ", err)
	}
	//log.Printf("Storing under %s\n", tx_key)

	// store for use case - disabled for now - in config
	//store_usecase := false
	store_usecase := true
	if store_usecase {
		uc_key := m.Target
		_, err := ms.C.Do("rpush", uc_key, EncodeNiasMessage(m))
		if err != nil {
			log.Println("error saving message:uc - ", err)
		}
	}

}

// update the progress of the validation transaction
func (ms *MessageStore) IncrementTracker(txid string) {

	mutex.Lock()
	status[txid]++
	mutex.Unlock()

}

// Retrieve the data for this transaction - txid
// fulldata if true returns all data in the transaction, if false then return
// is capped at 10,000 records
// If txid is empty, retrieves all records in the use case stream
func GetTxData(txid string, prefix string, fulldata bool) ([]interface{}, error) {

	data := make([]interface{}, 0)
	c := CreateLedisConnection(1024, 1024)
	defer c.Close()

	var endpoint int
	if fulldata {
		endpoint = -1
	} else {
		endpoint = (10000 - 1)
	}

	//log.Printf("retrieving from %s\n", prefix+txid)
	results, err := goredis.Values(c.Do("lrange", prefix+txid, 0, endpoint))
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

// binary encding for messages going to internal q/store.
func EncodeNiasMessage(msg *NiasMessage) []byte {

	encBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(encBuf)
	err := encoder.Encode(msg)
	if err != nil {
		log.Printf("Encoder unable to binary encode message for: %#v\n", msg)
	}
	return encBuf.Bytes()

}

// binary decoding for messages coming from internal q/store.
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
