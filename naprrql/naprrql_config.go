// configmanager.go
package naprrql

import (
	"github.com/BurntSushi/toml"
	"github.com/nsip/nias2/lib"
	"log"
)

// utility object to manage configurable parameters

var nAPLANConfig lib.NaprrqlConfig = lib.NaprrqlConfig{}

func LoadNAPLANConfig() lib.NaprrqlConfig {
	if !nAPLANConfig.Loaded {
		ncfg := lib.NaprrqlConfig{}
		if _, err := toml.DecodeFile("naprrql.toml", &ncfg); err != nil {
			log.Fatalln("Unable to read NAPLAN config, aborting.", err)
		}
		nAPLANConfig = ncfg
		nAPLANConfig.Loaded = true
	}
	return nAPLANConfig
}
