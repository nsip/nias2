// configmanager.go
package nias2

import (
	"github.com/BurntSushi/toml"
	"log"
)

// utility object to manage configurable parameters

type NIASConfig struct {
	TestYear        string
	WebServerPort   string
	ValidationRoute []string
	PoolSize        int
}

var NiasConfig = loadDefaultConfig()

func loadDefaultConfig() NIASConfig {

	ncfg := NIASConfig{}
	if _, err := toml.DecodeFile("nias.toml", &ncfg); err != nil {
		log.Fatalln("Unable to read default config, aborting.", err)
	}
	return ncfg

}
