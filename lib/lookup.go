// lookup.go
package nias2

import (
	"github.com/siddontang/goredis"
	"log"
)

//
// encapsulation of all hash-based ledis lookups
//

var id_c = CreateLedisConnection(64, 64)
var asl_c = CreateLedisConnection(1024, 1024)

const ID_PREFIX = "id:"
const ASL_KEY = "asl:lookup"

// asl lookup
func ASLKeyExists(msg *NiasMessage) bool {

	rr := msg.Body.(RegistrationRecord)

	if resp, _ := asl_c.Do("hget", ASL_KEY, rr.ASLSchoolId); resp != nil {
		return true
	}

	return false
}

//
// ID Lookup Handler - simple key
//
func SimpleIDKeyExists(msg *NiasMessage) bool {

	rr := msg.Body.(RegistrationRecord)

	k := IDSimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	log.Printf("%s %v", ID_PREFIX+msg.TxID, k)
	log.Println("queried")

	if resp, _ := id_c.Do("hget", ID_PREFIX+msg.TxID, k); resp != nil {
		log.Println("found")
		return true
	}

	log.Println("not found")
	return false

}

//
// ID Lookup Handler - complex key
//
func ComplexIDKeyExists(msg *NiasMessage) bool {

	rr := msg.Body.(RegistrationRecord)

	ek := IDExtendedKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
		FamilyName:  rr.FamilyName,
		GivenName:   rr.GivenName,
		BirthDate:   rr.BirthDate,
	}

	log.Println(ek)
	log.Println("queried")
	if resp, _ := id_c.Do("hget", ID_PREFIX+msg.TxID, ek); resp != nil {
		log.Println("found")
		return true
	}

	log.Println("not found")
	return false
}

// return the original line (ol) number that holds this identity if a duplicate
func GetIDValue(msg *NiasMessage) (string, error) {

	rr := msg.Body.(RegistrationRecord)

	k := IDSimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	log.Printf("%s %v", ID_PREFIX+msg.TxID, k)
	log.Println("queried")

	if ol, err := goredis.String(id_c.Do("hget", ID_PREFIX+msg.TxID, k)); err != nil {
		log.Println("not found")
		return "", err
	} else {
		log.Println("found")
		return ol, nil
	}

}

// store an identity
func SetIDValue(msg *NiasMessage) error {

	rr := msg.Body.(RegistrationRecord)

	// first set a simple key type
	k := IDSimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	if _, err := id_c.Do("hset", ID_PREFIX+msg.TxID, k, msg.SeqNo); err != nil {
		log.Println("err1")
		return err
	}
	log.Printf("%s %v", ID_PREFIX+msg.TxID, k)
	log.Println("set")

	// then the complex type
	ek := IDExtendedKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
		FamilyName:  rr.FamilyName,
		GivenName:   rr.GivenName,
		BirthDate:   rr.BirthDate,
	}
	if _, err := id_c.Do("hset", ID_PREFIX+msg.TxID, ek, msg.SeqNo); err != nil {
		log.Println("err2")
		return err
	}
	log.Println(ek)
	log.Println("set")

	return nil

}

// for a given message returns the state is associated with the given asl id.
func GetASLValue(msg *NiasMessage) (string, error) {

	rr := msg.Body.(RegistrationRecord)

	if stateid, err := goredis.String(asl_c.Do("hget", ASL_KEY, rr.ASLSchoolId)); err != nil {
		return "", err
	} else {
		return stateid, nil
	}

}

// store an ASL record
// only minimal data - key of acaraid to determine if record exists
// value is the state associated with the acara id for extended checking
func SetASLValue(st ASLSearchTerms) error {

	if _, err := asl_c.Do("hset", ASL_KEY, st.AcaraID, st.State); err != nil {
		return err
	}
	return nil

}
