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
	PeUpdateInterval         string
	DivYieldUpdateInterval   string
	StockUpdateInterval      string
	SmptServerHost           string
	SmptServerPort           string
	SmptSenderUsername       string
	SmptSenderPassword       string
	NotificationRecipient    string
	AuthorizationServer      string
	AuthorizationAudience    string
	RequiredScopes           string
	EmailClaim               string
}

//Config holds the current configuration
var Config Configuration

//Init initializes the config object with the current configuration
//This should be called only once, at service startup
func Init() {
	err := gonfig.GetConf("config.json", &Config)

	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
}
