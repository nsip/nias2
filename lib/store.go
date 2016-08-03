// store.go
package nias2

import (
	"bytes"
	"encoding/gob"
	"github.com/siddontang/goredis"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

const VALIDATION_PREFIX = "nvr:"
const STORE_AND_FORWARD_PREFIX = "ssf:"
const SIF_MEMORY_STORE_PREFIX = "sms:"

type MyGoRedisClient struct {
	*goredis.Client
}

// MessageStore listens for messages on the store topic and captures them
// in ledis as lists (persistent qs in effect), messages can be stored on transaction
// and use case basis - use case being the superset of all transactions of a given type
type MessageStore struct {
	//C    *goredis.Client
	//trkC *goredis.Client
	C    MyGoRedisClient
	trkC MyGoRedisClient
}

func NewMessageStore() *MessageStore {
	ms := MessageStore{
		C: MyGoRedisClient{CreateLedisConnection(1024, 1024)},
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

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func SliceDiff(xs []string, ys []string) []string {
	found := make(map[string]bool)
	ret := make([]string, 1)
	for _, y := range ys {
		found[y] = true
	}
	for _, x := range xs {
		if !found[x] {
			ret = append(ret, x)
		}
	}
	return ret
}

// parse GraphStruct, and store sets in SMS
func (ms *MessageStore) StoreGraph(m *NiasMessage) error {
	graphstruct := m.Body.(GraphStruct)
	// get the nodes equivalent to the current node
	prev_equivalents, err := goredis.Strings(ms.C.Do("smembers", "equivalent:ids:"+graphstruct.Guid))
	if err != nil {
		log.Println("error saving message:storegraph:1 - ", err)
		return err
	}
	equivalents := append(prev_equivalents, graphstruct.EquivalentIds...)
	RemoveDuplicates(&equivalents)

	// are there any new equivalences because of this tuple? If so, duplicate the existing links among all equivalences
	new_equivalents := SliceDiff(graphstruct.EquivalentIds, prev_equivalents)
	if len(new_equivalents) > 0 {
		new_equivalents = append(new_equivalents, graphstruct.Guid)
		for i := range new_equivalents {
			if new_equivalents[i] == "" {
				continue
			}
			for j := range new_equivalents {
				if new_equivalents[j] == "" {
					continue
				}
				if new_equivalents[i] != new_equivalents[j] {
					_, err := ms.C.Do("sunionstore", new_equivalents[i], new_equivalents[i], new_equivalents[j])
					if err != nil {
						log.Println("error saving message:storegraph:2 - ", err)
						log.Printf("sunionstore %s %s %s", new_equivalents[i], new_equivalents[i], new_equivalents[j])
						return err
					}
				}
			}
		}
	}

	// no responses needed from redis so pipeline for speed
	// siddontag does not implement pipelining on the client, only on the connection. Will need to supplement
	// his client's .Do() with a .Send(), which uses .Send() not .Do() on the connection, then issues a final .Do("EXEC")
	if graphstruct.Label != "" {
		_, err := ms.C.Do("hset", "labels", graphstruct.Guid, graphstruct.Label)
		if err != nil {
			log.Printf("hset labels %s %s", graphstruct.Guid, graphstruct.Label)
			log.Println("error saving message:storegraph:3 - ", err)
			return err
		}
	}
	if graphstruct.Type != "" {
		_, err := ms.C.Do("sadd", "known:collections", graphstruct.Type)
		if err != nil {
			log.Println("error saving message:storegraph:4 - ", err)
			return err
		}
		_, err = ms.C.Do("sadd", graphstruct.Type, graphstruct.Guid)
		if err != nil {
			log.Println("error saving message:storegraph:5 - ", err)
			return err
		}
	}

	if len(graphstruct.Links) > 0 {
		args := make([]interface{}, 0)
		args = append(args, graphstruct.Guid)
		for _, link := range graphstruct.Links {
			args = append(args, link)
		}

		//_, err := ms.C.Do("sadd", graphstruct.Guid, graphstruct.Links)
		_, err := ms.C.Do("sadd", args...)
		if err != nil {
			log.Println("error saving message:storegraph:6 - ", err)
			return err
		}
		for _, id := range equivalents {
			args = make([]interface{}, 0)
			args = append(args, id)
			for _, link := range graphstruct.Links {
				args = append(args, link)
			}
			//_, err := ms.C.Do("sadd", id, graphstruct.Links)
			_, err := ms.C.Do("sadd", args...)
			if err != nil {
				log.Println("error saving message:storegraph:7 - ", err)
				return err
			}
		}
	}
	for key, value := range graphstruct.OtherIds {
		_, err := ms.C.Do("hset", "oid:"+value, key, graphstruct.Guid)
		if err != nil {
			log.Println("error saving message:storegraph:8 - ", err)
			return err
		}
		_, err = ms.C.Do("sadd", "other:ids", "oid:"+value)
		if err != nil {
			log.Println("error saving message:storegraph:9 - ", err)
			return err
		}
	}

	// extract equivalent ids
	for _, equiv := range graphstruct.EquivalentIds {
		refs := make([]interface{}, 0)
		refs = append(refs, "equivalent:ids:"+equiv)
		for _, equiv2 := range graphstruct.EquivalentIds {
			if equiv != equiv2 {
				refs = append(refs, equiv2)
			}
		}
		refs = append(refs, graphstruct.Guid)
		//_, err := ms.C.Do("sadd", "equivalent:ids:"+equiv, refs...)
		_, err := ms.C.Do("sadd", refs...)
		if err != nil {
			log.Println("error saving message:storegraph:10 - ", err)
			return err
		}
	}

	// then add id to sets for links
	for _, link := range graphstruct.Links {
		refs := make([]interface{}, 0)
		refs = append(refs, link)
		for _, link2 := range graphstruct.Links {
			if link != link2 {
				refs = append(refs, link2)
			}
		}
		refs = append(refs, graphstruct.Guid)
		for _, equiv := range equivalents {
			refs = append(refs, equiv)
		}
		//_, err := ms.C.Do("sadd", link, refs...)
		_, err := ms.C.Do("sadd", refs...)
		if err != nil {
			log.Println("error saving message:storegraph:11 - ", err)
			return err
		}
	}
	//_, err = ms.C.Do("exec")
	if err != nil {
		log.Println("error saving message:storegraph:12 - ", err)
		return err
	}
	log.Printf("Stored graph for %s", graphstruct.Guid)
	return nil
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

// Implement pipeline command on goredis client

func (c *MyGoRedisClient) Send(cmd string, args ...interface{}) error {
	var co *goredis.PoolConn
	var err error
	//var r interface{}

	for i := 0; i < 2; i++ {
		//co, err = c.get()
		co, err = c.Get()
		if err != nil {
			return err
		}

		err = co.Send(cmd, args...)
		if err != nil {
			co.Close()

			if e, ok := err.(*net.OpError); ok && strings.Contains(e.Error(), "use of closed network connection") {
				//send to a closed connection, try again
				continue
			}

			return err
		} else {
			// c.put(co)
			// TODO: can't shove this out of package, for error recovery
		}

		return nil
	}

	return err
}
