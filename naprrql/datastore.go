// datastore.go

package naprrql

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Read Only optimised
var db *leveldb.DB
var dbOpen bool = false

// Read Write for ingest only
var dbReadWrite *leveldb.DB
var dbOpenReadWrite bool = false

var ge = GobEncoder{}

// if allow_empty is false, abort if database is empty: enforces
// requirement to run ingest first
func GetDB() *leveldb.DB {
	if !dbOpen {
		log.Println("DB not initialised. Opening...")
		openDB()
	}
	return db
}

// if allow_empty is false, abort if database is empty: enforces
// requirement to run ingest first
func GetDBReadWrite() *leveldb.DB {
	if !dbOpenReadWrite {
		log.Println("DB not initialised. Opening ReadWrite...")
		openDBReadWrite()
	}
	return dbReadWrite
}

//
// Check if the databse contains any entries
//
func DatabaseIsEmpty() bool {

	iter := GetDB().NewIterator(nil, nil)
	iter.Next()
	if len(iter.Key()) == 0 {
		return true
	}
	return false
}

//
// run compaction on db, typically post-ingest
//
func CompactDatastore() {
	err := GetDB().CompactRange(util.Range{Limit: nil, Start: nil})
	if err != nil {
		log.Println("Error compacting datastore: ", err)
	}
}

//
// open the kv store, this must be called before any access is attempted
//
func openDB() {

	workingDir := "kvs"

	o := &opt.Options{

		// Original
		// // Filter:             filter.NewBloomFilter(10),
		// // BlockCacheCapacity: 128 * 1024 * 1024,
		// // NoSync: true,
		// // OpenFilesCacheCapacity: 1024,
		// // ReadOnly:            true,
		// CompactionTableSize: (4 * opt.MiB),

		// Experiment from Matt 2019-04-12
		Filter:                        filter.NewBloomFilter(10),
		BlockCacheCapacity:            2 * opt.GiB,
		CompactionTotalSizeMultiplier: 20,
		OpenFilesCacheCapacity:        1024,
		ReadOnly:                      true,
		CompactionTableSize:           (4 * opt.MiB),   //default 2
		CompactionTotalSize:           (128 * opt.MiB), //default 10 mb
	}
	var dbErr error
	db, dbErr = leveldb.OpenFile(workingDir, o)
	if dbErr != nil {
		log.Fatalln("DB Create error: ", dbErr)
	}

	dbOpen = true
}

//
// open the kv store, this must be called before any access is attempted (ReadWrite)
//
func openDBReadWrite() {

	workingDir := "kvs"

	o := &opt.Options{

		// Original
		// // Filter:             filter.NewBloomFilter(10),
		// // BlockCacheCapacity: 128 * 1024 * 1024,
		// // NoSync: true,
		// // OpenFilesCacheCapacity: 1024,
		// // ReadOnly:            true,
		// CompactionTableSize: (4 * opt.MiB),

		// Experiment from Matt 2019-04-12
		Filter:                        filter.NewBloomFilter(10),
		BlockCacheCapacity:            2 * opt.GiB,
		CompactionTotalSizeMultiplier: 20,
		OpenFilesCacheCapacity:        1024,
		ReadOnly:                      false,
		CompactionTableSize:           (4 * opt.MiB),   //default 2
		CompactionTotalSize:           (128 * opt.MiB), //default 10 mb
	}
	var dbErr error
	dbReadWrite, dbErr = leveldb.OpenFile(workingDir, o)
	if dbErr != nil {
		log.Fatalln("DB Create error: ", dbErr)
	}

	dbOpenReadWrite = true
}

//
// Given a key or key-prefix, returns the reference ids that
// can be used in a Get operation to retreive the
// desired object
//
func getIdentifiers(keyPrefix string) []string {

	db = GetDB()
	objIDs := make([]string, 0)

	searchKey := []byte(keyPrefix)
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
