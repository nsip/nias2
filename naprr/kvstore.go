package naprr

// kvstore.go
//
// key-value backing store for naprr
// introduced to try and resolve trade-off issues on win32
// between mem usage and disk access, both of which have limits
// which cause issues.

import (
	"log"
	"os"

	"github.com/dgraph-io/badger"
)

var KV *badger.KV

func init() {
	KV = openDB()
}

func openDB() *badger.KV {
	log.Println("opening database...")

	// make sure the working directory exists
	err = os.Mkdir("kvs", os.ModePerm)
	if err != nil {
		log.Println("Error trying to create kvs working directory")
	}

	opt := badger.DefaultOptions
	opt.Dir = "kvs"
	opt.ValueDir = "kvs"
	kv, dbErr := badger.NewKV(&opt)
	if dbErr != nil {
		log.Fatalln("KVS create error: ", dbErr)
	}
	return kv
}
