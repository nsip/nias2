// ledis.go
// sets up the ledis server & provides useful constants etc.

package sms

import (
	"github.com/siddontang/goredis"
	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/server"
	"log"
)

const DEF_STORAGE_DB_ADDRESS = "127.0.0.1:6397"
const DEF_STORAGE_DB_PROXY_ADDRESS = "127.0.0.1:6397"
const DEF_LOOKUP_DB_ADDRESS = "127.0.0.1:6399"
const DEF_LOOKUP_DB_PROXY_ADDRESS = "127.0.0.1:6399"

// alternative ports for db clustering; twemproxy, xcodis etc.
// const DEF_DB_ADDRESS = "127.0.0.1:6379"
// const DEF_PROXY_ADDRESS = "127.0.0.1:22121"
// const DEF_PROXY_ADDRESS = "127.0.0.1:6379"
// const DEF_PROXY_ADDRESS = "127.0.0.1:19000"
// const DEF_ADDRESS = "127.0.0.1:19000"

// returns a connection to the storage database
func CreateStorageConnection() *goredis.Client {

	c := goredis.NewClient(DEF_STORAGE_DB_PROXY_ADDRESS, "")
	return c

}

// returns a connection to the lookup service database
func CreateLookupConnection() *goredis.Client {

	c := goredis.NewClient(DEF_LOOKUP_DB_PROXY_ADDRESS, "")
	return c

}

// launch ledis server (local) for message store
func LaunchStorageServer() {

	cfg := config.NewConfigDefault()

	cfg.LevelDB.CacheSize = 524288000
	// cfg.LevelDB.WriteBufferSize = (1024 * 1024 * 250) //67108864
	cfg.LevelDB.WriteBufferSize = 67108864

	cfg.LevelDB.MaxOpenFiles = 10240
	// log.Println("\tLDB MxOF is: ", cfg.LevelDB.MaxOpenFiles)

	// commented from default
	// cfg.ConnReadBufferSize = (1024 * 1024 * 5)
	// cfg.ConnWriteBufferSize = (1024 * 1024 * 50)

	// log.Println("\tRead buffer is", cfg.ConnReadBufferSize)
	// log.Println("\tWrite buffer is", cfg.ConnWriteBufferSize)

	cfg.Addr = DEF_STORAGE_DB_ADDRESS
	cfg.DataDir = "./var/storage"
	app, err := server.NewApp(cfg)
	if err != nil {
		log.Fatalln("Cannot launch database server, aborting....", err)
	}
	go app.Run()

}

// launch ledis server (local) for lookup duties
func LaunchLookupServer() {

	cfg := config.NewConfigDefault()

	// cfg.LevelDB.CacheSize = 524288000
	cfg.LevelDB.WriteBufferSize = (1024 * 1024 * 250) //67108864

	cfg.LevelDB.MaxOpenFiles = 10240
	// log.Println("\tLDB MxOF is: ", cfg.LevelDB.MaxOpenFiles)

	cfg.ConnReadBufferSize = (1024 * 1024 * 5)
	cfg.ConnWriteBufferSize = (1024 * 1024 * 50)

	// log.Println("\tRead buffer is", cfg.ConnReadBufferSize)
	// log.Println("\tWrite buffer is", cfg.ConnWriteBufferSize)

	cfg.Addr = DEF_LOOKUP_DB_ADDRESS
	cfg.DataDir = "./var/lookup"
	app, err := server.NewApp(cfg)
	if err != nil {
		log.Fatalln("Cannot launch database server, aborting....", err)
	}
	go app.Run()

}
