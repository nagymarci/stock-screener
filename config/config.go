package config

import (
	"fmt"
	"os"

	"github.com/tkanos/gonfig"
)

//Configuration holds the config of the app
type Configuration struct {
	DatabaseConnectionString string
	Port                     int
	StockInfoProviderURL     string
}

//Get returns the current configuration
func Get() Configuration {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)

	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}

	return configuration
}
