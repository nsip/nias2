// configmanager.go
package napval

import (
	"github.com/BurntSushi/toml"
	"log"
)

// utility object to manage configurable parameters

type NIASConfig struct {
	TestYear         string
	WebServerPort    string
	NATSPort         string
	ValidationRoute  []string
	SSFRoute         []string
	SMSRoute         []string
	PoolSize         int // number of service processors
	MsgTransport     string
	TxReportInterval int // progress report after every n records
	UIMessageLimit   int //how many messages to send to web ui
	TxStorageLimit   int
	Loaded           bool
}

var defaultConfig NIASConfig = NIASConfig{}
var nAPLANConfig NIASConfig = NIASConfig{}

func LoadNAPLANConfig() NIASConfig {
	if !nAPLANConfig.Loaded {
		ncfg := NIASConfig{}
		if _, err := toml.DecodeFile("napval.toml", &ncfg); err != nil {
			log.Fatalln("Unable to read NAPLAN config, aborting.", err)
		}
		nAPLANConfig = ncfg
		nAPLANConfig.Loaded = true
	}
	return nAPLANConfig
}

func LoadDefaultConfig() NIASConfig {
	if !defaultConfig.Loaded {
		ncfg := NIASConfig{}
		if _, err := toml.DecodeFile("nias.toml", &ncfg); err != nil {
			log.Fatalln("Unable to read default config, aborting.", err)
		}
		defaultConfig = ncfg
		defaultConfig.Loaded = true
	}
	return defaultConfig
}
