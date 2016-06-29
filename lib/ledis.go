// ledis.go
// sets up the ledis server & provides useful constants etc.

package nias2

import (
	"github.com/siddontang/goredis"
	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/server"
	"log"
)

const DEF_ADDRESS = "127.0.0.1:6380"

func CreateLedisConnection() *goredis.Client {

	log.Println("Creating ledis client")
	c := goredis.NewClient(DEF_ADDRESS, "")
	return c

}

func LaunchLedisServer() {

	cfg := config.NewConfigDefault()
	cfg.LevelDB.CacheSize = 524288000
	cfg.LevelDB.WriteBufferSize = 67108864

	// cfg.LevelDB.MaxOpenFiles = 10240
	// log.Println("\tLDB MxOF is: ", cfg.LevelDB.MaxOpenFiles)

	cfg.ConnReadBufferSize = 1024
	cfg.ConnWriteBufferSize = 1024

	// log.Println("\tRead buffer is", cfg.ConnReadBufferSize)
	// log.Println("\tWrite buffer is", cfg.ConnWriteBufferSize)

	cfg.Addr = DEF_ADDRESS
	cfg.DataDir = "./var"
	app, err := server.NewApp(cfg)
	if err != nil {
		log.Fatalln("Cannot launch database server, aborting....", err)
	}
	go app.Run()

}
