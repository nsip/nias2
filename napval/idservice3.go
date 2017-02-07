// idservice3.go
package napval

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
	"github.com/orcaman/concurrent-map"
	"gopkg.in/fatih/set.v0"
	"log"
)

// id service requires too much overhead using ledis.
// This version replaces with in-memory structures that are
// cleaned up on notificaiton from the transaction tracker

// threadsafe store for transaction id data, shared by
// all service instances
var transactionStore = cmap.New()

// set up the data structures used as hash keys
// simple duplicate check is, have we seen this userid for this school before
type ID3SimpleKey struct {
	LocalId     string
	ASLSchoolId string
}

// this checks the user against a set of likely colliding matches
type ID3ExtendedKey struct {
	LocalId     string
	ASLSchoolId string
	FamilyName  string
	GivenName   string
	BirthDate   string
}

// implementation of the id service
type IDService3 struct {
	Transactions cmap.ConcurrentMap
	C            *nats.EncodedConn
}

// transaction id lookups
type TransactionIDs struct {
	Locations    cmap.ConcurrentMap
	SimpleKeys   *set.Set
	ExtendedKeys *set.Set
}

// create a new id service instance
func NewIDService3(nats_config lib.NATSConfig) (*IDService3, error) {
	ids := IDService3{Transactions: transactionStore,
		C: lib.CreateNATSConnection(nats_config)}
	ids.txMonitor()
	return &ids, nil
}

// start a listener process for notifications that
// transactions are complete; when they are
// remove id data from datastore to prevent resource leak over time
func (ids *IDService3) txMonitor() {

	log.Println("tx monitor is listening...")
	go ids.C.QueueSubscribe(lib.TRACK_TOPIC, "id3", func(txID string) {
		log.Println("Transaction complete message received for tx: ", txID)
		transactionStore.Remove(txID)
	})

}

// implement the service interface
func (ids *IDService3) HandleMessage(req *lib.NiasMessage) ([]lib.NiasMessage, error) {

	responses := make([]lib.NiasMessage, 0)

	rr, ok := req.Body.(xml.RegistrationRecord)
	if !ok {
		log.Println("IDService3 received a message that is not a RegistrationRecord, ignoring")
		return responses, nil
	}

	// see if dataset exists for this transaction, create if not
	ids.Transactions.SetIfAbsent(req.TxID,
		TransactionIDs{Locations: cmap.New(),
			SimpleKeys:   set.New(),
			ExtendedKeys: set.New()})

	// retrieve the transaction dataset from the store
	tdata, ok := ids.Transactions.Get(req.TxID)
	if !ok {
		log.Println("Unable to find transaction id dataset for: ", req.TxID)
		return responses, nil
	}
	tids, ok := tdata.(TransactionIDs)
	if !ok {
		log.Printf("Unable to cast tid store type is: %T %T", tids, tdata)
		return responses, nil
	}

	// perform lookups
	k1 := ID3SimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	k2 := ID3ExtendedKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
		FamilyName:  rr.FamilyName,
		GivenName:   rr.GivenName,
		BirthDate:   rr.BirthDate,
	}

	simpleKey := fmt.Sprintf("%v", k1)
	complexKey := fmt.Sprintf("%v", k2)
	// log.Printf("simplekey: %s\ncompllexkey: %s", simpleKey, complexKey)
	var simpleRecordExists, complexRecordExists bool

	if simpleRecordExists = tids.SimpleKeys.Has(simpleKey); !simpleRecordExists {
		tids.SimpleKeys.Add(simpleKey)
	}

	if complexRecordExists = tids.ExtendedKeys.Has(complexKey); !complexRecordExists {
		tids.ExtendedKeys.Add(complexKey)
	}
	tids.Locations.SetIfAbsent(simpleKey, req.SeqNo)
	tids.Locations.SetIfAbsent(complexKey, req.SeqNo)

	// if record is new then just return
	if !complexRecordExists && !simpleRecordExists {
		return responses, nil
	}

	// if we have seen it before then construct validation error
	if complexRecordExists {
		loc, _ := tids.Locations.Get(complexKey)
		ol, _ := loc.(string)
		desc := "Potential duplicate of record: " + ol + "\n" +
			"based on matching: student local id, school asl id, family & given names and birthdate"
		ve := ValidationError{
			Description:  desc,
			Field:        "Multiple (see description)",
			OriginalLine: req.SeqNo,
			Vtype:        "identity",
		}
		r := lib.NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		// r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

	} else if simpleRecordExists {
		loc, _ := tids.Locations.Get(simpleKey)
		ol, _ := loc.(string)
		desc := "LocalID (Student) and ASL ID (School) are potential duplicate of record: " + ol
		ve := ValidationError{
			Description:  desc,
			Field:        "LocalID/ASL ID",
			OriginalLine: req.SeqNo,
			Vtype:        "identity",
		}
		r := lib.NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		// r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

	}

	return responses, nil
}
