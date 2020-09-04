package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/routes"
	"github.com/nagymarci/stock-screener/service"
	"github.com/robfig/cron/v3"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	router := mux.NewRouter().StrictSlash(true)
	config.Init()
	routes.Route(router)
	database.Connect(config.Config.DatabaseConnectionString)

	c := cron.New()
	c.AddFunc("CRON_TZ=America/New_York 0/30,1 9-10 * * MON-FRI", service.UpdateStocks)
	c.AddFunc("CRON_TZ=America/New_York * 10-16 * * MON-FRI", service.UpdateStocks)

	//c.Start()

	log.Println(len(c.Entries()))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), router))
}
