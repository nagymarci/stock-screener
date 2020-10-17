package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/routes"
	"github.com/nagymarci/stock-screener/service"
	"github.com/robfig/cron/v3"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	config.Init()
	router := routes.Route()
	database.Connect(config.Config.DatabaseConnectionString)

	c := cron.New()
	_, err := c.AddFunc("CRON_TZ=America/New_York * 9-17 * * MON-FRI", service.UpdateStocks)
	log.Println(err)
	_, err = c.AddFunc("CRON_TZ=America/New_York 0 8-18 * * MON-FRI", service.NotifyChanges)
	log.Println(err)

	c.Start()

	log.Println(len(c.Entries()))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), router))
}
