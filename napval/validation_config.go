// configmanager.go
package napval

import (
	"github.com/BurntSushi/toml"
	"github.com/nsip/nias2/lib"
	"log"
)

// utility object to manage configurable parameters

var nAPLANConfig lib.NIASConfig = lib.NIASConfig{}

func LoadNAPLANConfig() lib.NIASConfig {
	if !nAPLANConfig.Loaded {
		ncfg := lib.NIASConfig{}
		if _, err := toml.DecodeFile("napval.toml", &ncfg); err != nil {
			log.Fatalln("Unable to read NAPLAN config, aborting.", err)
		}
		nAPLANConfig = ncfg
		nAPLANConfig.Loaded = true
	}
	return nAPLANConfig
}
