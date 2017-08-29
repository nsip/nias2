// datastore.go

package naprrql

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var db *leveldb.DB
var dbOpen bool = false

var ge = GobEncoder{}

func GetDB() *leveldb.DB {
	if !dbOpen {
		log.Println("DB not initialised. Opening...")
		openDB()
	}
	return db
}

//
// open the kv store, this must be called before any access is attempted
//
func openDB() {

	workingDir := "kvs"

	var dbErr error
	db, dbErr = leveldb.OpenFile(workingDir, nil)
	if dbErr != nil {
		log.Fatalln("DB Create error: ", dbErr)
	}
	dbOpen = true
}

//
// Given a key or key-prefix, returns the reference ids that
// can be used in a Get operation to retreive the
// desired object
//
func getIdentifiers(keyPrefix string) []string {

	db = GetDB()
	objIDs := make([]string, 0)

	searchKey := []byte(keyPrefix + ":")
	// log.Printf("search_key: %s\n\n", searchKey)
	iter := db.NewIterator(util.BytesPrefix(searchKey), nil)
	for iter.Next() {
		id := fmt.Sprintf("%s", iter.Value())
		objIDs = append(objIDs, id)
		// break
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Println("Iterator error: ", err)
	}

	return objIDs
}

//
// get objects stored against the list of keys
//
func getObjects(objIDs []string) ([]interface{}, error) {

	db = GetDB()
	objects := []interface{}{}

	for _, objID := range objIDs {
		var object interface{}
		data, err := db.Get([]byte(objID), nil)
		if err != nil {
			log.Println("Cannot find object with key: ", string(objID))
			return objects, err
		}
		err = ge.Decode(data, &object)
		if err != nil {
			log.Println("Cannot decode object with key: ", objID, err)
			return objects, err
		}
		objects = append(objects, object)
	}

	return objects, nil

}

//
//
//
