// lookup.go
package nias2

import (
	"errors"
	"fmt"
	"github.com/siddontang/goredis"
	//"log"
)

//
// encapsulation of all hash-based ledis lookups
//

var id_c = CreateLedisConnection(64, 64)
var asl_c = CreateLedisConnection(1024, 1024)

const ID_PREFIX = "id:"
const ASL_KEY = "asl:lookup"

// seconds until key expires, for expirable keys
const KEY_EXPIRATION = 100000

// asl lookup
func ASLKeyExists(msg *NiasMessage) bool {

	rr := msg.Body.(RegistrationRecord)

	if resp, _ := asl_c.Do("hget", ASL_KEY, rr.ASLSchoolId); resp != nil {
		return true
	}

	if rr.ASLSchoolId == "" {
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

	if resp, _ := id_c.Do("hget", ID_PREFIX+msg.TxID, k); resp != nil {
		return true
	}

	return false

}

// SetNx simple key
func SimpleIDKeySetnx(msg *NiasMessage) bool {

	rr := msg.Body.(RegistrationRecord)
	var resp int
	var err error

	k := IDSimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	//log.Printf("SimpleIDKeySetnx %s %v", ID_PREFIX+msg.TxID, k)
	//log.Println("queried")

	resp, err = goredis.Int(id_c.Do("set", fmt.Sprintf("%s%s::%v NX EX %d", ID_PREFIX, msg.TxID, k, KEY_EXPIRATION), msg.SeqNo))
	//log.Printf(":: %v %v %v", resp, err, msg.SeqNo)
	if err == nil && resp == 0 {
		//log.Println("set")
		return true
	}

	//log.Printf("found:%v", resp)
	return false

}

// Seen simple key: its value is a msg.SeqNo other than the current one
func SimpleIDKeySeen(msg *NiasMessage) (string, error) {

	rr := msg.Body.(RegistrationRecord)

	k := IDSimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	//log.Printf("SimpleIDKeySeen %s %v", ID_PREFIX+msg.TxID, k)
	//log.Println("queried")

	resp, err := goredis.String(id_c.Do("get", fmt.Sprintf("%s%s::%v", ID_PREFIX, msg.TxID, k)))
	if err != nil {
		//log.Println("err")
		return "", err
	}
	if resp == msg.SeqNo {
		//log.Println("not found")
		return "", nil
	}
	//log.Println("found")
	return fmt.Sprint(resp), errors.New("ID already seen")

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

	if resp, _ := id_c.Do("hget", ID_PREFIX+msg.TxID, ek); resp != nil {
		return true
	}

	return false
}

// SetNx complex key
func ComplexIDKeySetnx(msg *NiasMessage) bool {

	rr := msg.Body.(RegistrationRecord)

	ek := IDExtendedKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
		FamilyName:  rr.FamilyName,
		GivenName:   rr.GivenName,
		BirthDate:   rr.BirthDate,
	}

	//log.Printf("ComplexIDKeySetnx %s %v", ID_PREFIX+msg.TxID, ek)
	//log.Println("queried")
	//if resp, err := goredis.Int(id_c.Do("setnx", fmt.Sprintf("%s%s::%v", ID_PREFIX, msg.TxID, ek))); err != nil && resp == 0 {
	if resp, err := goredis.Int(id_c.Do("set", fmt.Sprintf("%s%s::%v NX EX %d", ID_PREFIX, msg.TxID, ek, KEY_EXPIRATION))); err != nil && resp == 0 {
		//log.Println("set")
		return true
	}

	//log.Println("found")
	return false
}

// SetNx simple and complex keys; return false if either found
// SetNx complex key
func SimpleAndComplexIDKeySetnx(msg *NiasMessage) bool {

	rr := msg.Body.(RegistrationRecord)

	k := IDSimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	ek := IDExtendedKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
		FamilyName:  rr.FamilyName,
		GivenName:   rr.GivenName,
		BirthDate:   rr.BirthDate,
	}

	//log.Printf("ComplexIDKeySetnx %s %v: %v", ID_PREFIX+msg.TxID, ek, msg.SeqNo)

	resp, err := goredis.Int(id_c.Do("setnx", fmt.Sprintf("%s%s::%v", ID_PREFIX, msg.TxID, ek), msg.SeqNo))
	ret := (err == nil && resp == 1)
	resp, err = goredis.Int(id_c.Do("setnx", fmt.Sprintf("%s%s::%v", ID_PREFIX, msg.TxID, k), msg.SeqNo))
	//log.Printf(":: %v %v %v", resp, err, msg.SeqNo)
	//log.Printf("%v %v %v\n", ret, err, resp)
	ret = ret && (err == nil && resp == 1)
	return ret
}

// Seen complex key: its value is a msg.SeqNo other than the current one
func ComplexIDKeySeen(msg *NiasMessage) (string, error) {

	rr := msg.Body.(RegistrationRecord)

	ek := IDExtendedKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
		FamilyName:  rr.FamilyName,
		GivenName:   rr.GivenName,
		BirthDate:   rr.BirthDate,
	}

	//log.Printf("ComplexIDKeySeen %s %v", ID_PREFIX+msg.TxID, ek)
	//log.Println("queried")
	resp, err := goredis.String(id_c.Do("get", fmt.Sprintf("%s%s::%v", ID_PREFIX, msg.TxID, ek)))
	if err != nil {
		//log.Println("err")
		return "", err
	}
	if resp == msg.SeqNo {
		//log.Println("not found")
		return "", nil
	}
	//log.Println("found: "+resp)
	return fmt.Sprint(resp), errors.New("ID already seen")
}

// return the original line (ol) number that holds this identity if a duplicate
func GetIDValue(msg *NiasMessage) (string, error) {

	rr := msg.Body.(RegistrationRecord)

	k := IDSimpleKey{
		LocalId:     rr.LocalId,
		ASLSchoolId: rr.ASLSchoolId,
	}
	//log.Printf("GetIDValue %s %v", ID_PREFIX+msg.TxID, k)
	//log.Println("queried")

	if ol, err := goredis.String(id_c.Do("hget", ID_PREFIX+msg.TxID, k)); err != nil {
		//log.Println("not found")
		return "", err
	} else {
		//log.Println("found")
		return ol, nil
	}

}

// store an identity
func SetIDValue(msg *NiasMessage) error {
	/*
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
	*/
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
