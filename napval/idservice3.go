// idservice3.go
package napval

import (
	"fmt"
	"log"
	"strconv"

	//"github.com/nats-io/go-nats"
	"github.com/nats-io/nats.go"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
	"github.com/orcaman/concurrent-map"
	"gopkg.in/fatih/set.v0"
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

// second simple duplicate check: have we seen this PSI for this school before
type ID3SimpleKey2 struct {
	PSI         string
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
	Locations          cmap.ConcurrentMap
	SimpleKeysLocalId  *set.Set
	SimpleKeysPSI      *set.Set
	ExtendedKeys       *set.Set
	CrossSchoolMatches *set.Set
	CrossSchoolFTEs    cmap.ConcurrentMap
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
		/*
			// and now issue report on FTEs of data matched students
			tdata, ok := ids.Transactions.Get(txID)
			tids, ok1 := tdata.(TransactionIDs)
			if ok && ok1 {
				tids.CrossSchoolMatches.Each(func(item interface{}) bool {
					crossSchoolKey := item.(string)
					fte, ok := tids.CrossSchoolFTEs.Get(crossSchoolKey)
					if ok && fte != 1.0 {
						loc, _ := tids.Locations.Get(crossSchoolKey)
						ol, _ := loc.(string)
						desc := fmt.Sprintf("Student with key %s at %s matched across schools has an FTE of %0.2f across all enrolments\n",
							crossSchoolKey, ol, fte)
						ve := ValidationError{
							Description:  desc,
							Field:        "FTE",
							OriginalLine: ol,
							Vtype:        "identity",
							Severity:     "warning",
						}
						r := lib.NiasMessage{}
						r.TxID = txID
						r.SeqNo = ol
						r.Body = ve
						r.Source = "id"
						responses = append(responses, r)
						log.Printf("%v\n", responses)
					}
					return true
				})
				for _, r := range responses {
					ids.C.Publish(lib.STORE_TOPIC, r)
				}
			}
		*/
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
			SimpleKeysLocalId:  set.New(set.ThreadSafe).(*set.Set),
			SimpleKeysPSI:      set.New(set.ThreadSafe).(*set.Set),
			ExtendedKeys:       set.New(set.ThreadSafe).(*set.Set),
			CrossSchoolMatches: set.New(set.ThreadSafe).(*set.Set),
			CrossSchoolFTEs:    cmap.New(),
		})

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

	platformid := rr.GetOtherId("NAPPlatformStudentId")

	// perform lookups
	k11 := ID3SimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	k12 := ID3SimpleKey2{
		PSI:         platformid,
		ASLSchoolId: rr.ASLSchoolId,
	}
	k2 := ID3ExtendedKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
		FamilyName:  rr.FamilyName,
		GivenName:   rr.GivenName,
		BirthDate:   rr.BirthDate,
	}

	simpleKey1 := fmt.Sprintf("%v", k11)
	simpleKey2 := fmt.Sprintf("%v", k12)
	complexKey := fmt.Sprintf("%v", k2)
	crossSchoolKey := rr.FieldsKey(config.StudentMatch)

	// record all FTEs, preemptively; we will report on those corresponding to cross-school collisions
	fte, ok := tids.CrossSchoolFTEs.Get(crossSchoolKey)
	fte_num := 0.0
	if ok {
		fte_num = fte.(float64)
	}
	fte_new, err := strconv.ParseFloat(rr.FTE, 64)
	if err != nil {
		fte_new = 0.0
	}
	tids.CrossSchoolFTEs.Set(crossSchoolKey, fte_num+fte_new)

	//log.Printf("simplekey: %s\nsimplekey2: %s\ncompllexkey: %s", simpleKey1, simpleKey2, complexKey)
	var simpleRecordExists1, simpleRecordExists2, complexRecordExists, crossSchoolRecordExists bool

	if simpleRecordExists1 = (tids.SimpleKeysLocalId.Has(simpleKey1) && len(rr.LocalId) > 0); !simpleRecordExists1 {
		tids.SimpleKeysLocalId.Add(simpleKey1)
	}
	if simpleRecordExists2 = (tids.SimpleKeysPSI.Has(simpleKey2) && len(platformid) > 0); !simpleRecordExists2 {
		tids.SimpleKeysPSI.Add(simpleKey2)
	}

	if complexRecordExists = (tids.ExtendedKeys.Has(complexKey) && len(rr.LocalId) > 0); !complexRecordExists {
		tids.ExtendedKeys.Add(complexKey)
	}
	if crossSchoolRecordExists = (tids.CrossSchoolMatches.Has(crossSchoolKey) && len(crossSchoolKey) > 0); !crossSchoolRecordExists {
		tids.CrossSchoolMatches.Add(crossSchoolKey)
	}
	tids.Locations.SetIfAbsent(simpleKey1, req.SeqNo)
	tids.Locations.SetIfAbsent(simpleKey2, req.SeqNo)
	tids.Locations.SetIfAbsent(complexKey, req.SeqNo)
	tids.Locations.SetIfAbsent(crossSchoolKey, req.SeqNo)

	// if record is new then just return
	if !complexRecordExists && !simpleRecordExists1 && !simpleRecordExists2 && !crossSchoolRecordExists {
		return responses, nil
	}

	// if we have seen it before then construct validation error
	if crossSchoolRecordExists {
		loc, _ := tids.Locations.Get(crossSchoolKey)
		ol, _ := loc.(string)
		desc := "Potential duplicate of record: " + ol + "\n" +
			fmt.Sprintf("based on matching: %v", config.StudentMatch)
		ve := ValidationError{
			Description:  desc,
			Field:        "Multiple (see description)",
			OriginalLine: req.SeqNo,
			Vtype:        "identity",
			Severity:     "warning",
		}
		r := lib.NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		// r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

		if fte_num+fte_new != 1.0 {
			// FTEs for colliding student do not add up to 1.0
			// This may be a false alarm if there is a third enrolment; but we can only address that if we have STAN queues, and the ability to rewind a queue, so that we can issue a check at the end of a transaction.

			desc1 := fmt.Sprintf("Student with key %s at %s matched across schools has an FTE of %0.2f across all enrolments so far\n",
				crossSchoolKey, ol, fte_num+fte_new)
			ve1 := ValidationError{
				Description:  desc1,
				Field:        "FTE",
				OriginalLine: req.SeqNo,
				Vtype:        "identity",
				Severity:     "warning",
			}
			r1 := lib.NiasMessage{}
			r1.TxID = req.TxID
			r1.SeqNo = req.SeqNo
			r1.Body = ve1
			r1.Source = "id"
			responses = append(responses, r1)
		}

	} else if complexRecordExists {
		loc, _ := tids.Locations.Get(complexKey)
		ol, _ := loc.(string)
		desc := "Potential duplicate of record: " + ol + "\n" +
			"based on matching: student local id, school asl id, family & given names and birthdate"
		ve := ValidationError{
			Description:  desc,
			Field:        "Multiple (see description)",
			OriginalLine: req.SeqNo,
			Vtype:        "identity",
			Severity:     "warning",
		}
		r := lib.NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		// r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

	} else if simpleRecordExists1 {
		loc, _ := tids.Locations.Get(simpleKey1)
		ol, _ := loc.(string)
		desc := "LocalID (Student) and ASL ID (School) are potential duplicate of record: " + ol
		ve := ValidationError{
			Description:  desc,
			Field:        "LocalID/ASL ID",
			OriginalLine: req.SeqNo,
			Vtype:        "identity",
			Severity:     "warning",
		}
		r := lib.NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		// r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)

	} else if simpleRecordExists2 {
		loc, _ := tids.Locations.Get(simpleKey2)
		ol, _ := loc.(string)
		desc := "Platform Student ID (Student) and ASL ID (School) are potential duplicate of record: " + ol
		ve := ValidationError{
			Description:  desc,
			Field:        "PSI/ASL ID",
			OriginalLine: req.SeqNo,
			Vtype:        "identity",
			Severity:     "warning",
		}
		r := lib.NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		// r.Target = VALIDATION_PREFIX
		r.Body = ve
		responses = append(responses, r)
		// log.Printf("simplekey: %s\nsimplekey2: %s\ncompllexkey: %s", simpleKey1, simpleKey2, complexKey)

	}

	return responses, nil
}
