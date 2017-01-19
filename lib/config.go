// configmanager.go
package lib

import (
	"github.com/BurntSushi/toml"
	"log"
)

// utility object to manage configurable parameters

type NIASConfig struct {
	TestYear         string
	WebServerPort    string
	ValidationRoute  []string
	SSFRoute         []string
	SMSRoute         []string
	PoolSize         int // number of service processors
	MsgTransport     string
	TxReportInterval int // progress report after every n records
	UIMessageLimit   int //how many messages to send to web ui
	TxStorageLimit   int
}

var DefaultConfig = loadDefaultConfig()

func loadDefaultConfig() NIASConfig {

	ncfg := NIASConfig{}
	if _, err := toml.DecodeFile("napval.toml", &ncfg); err != nil {
		log.Fatalln("Unable to read default config, aborting.", err)
	}
	return ncfg

}
