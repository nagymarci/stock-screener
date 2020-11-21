package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/nagymarci/stock-screener/api"
	"github.com/nagymarci/stock-screener/controllers"

	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/routes"
	"github.com/nagymarci/stock-screener/service"
	"github.com/robfig/cron/v3"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	rand.Seed(time.Now().UnixNano())

	db := database.New(os.Getenv("DB_CONNECTION_URI"))
	stockInfo := database.NewStockinfos(db)

	stockscraper := api.New(os.Getenv("STOCKINFO_PROVIDER_URL"))

	controller := controllers.New(stockInfo, stockscraper)

	router := routes.Route(controller)

	updater := service.New(stockInfo, stockscraper, os.Getenv("STOCK_UPDATE_INTERVAL"), os.Getenv("PE_UPDATE_INTERVAL"), os.Getenv("DIV_UPDATE_INTERVAL"))

	c := cron.New()
	_, err := c.AddFunc("CRON_TZ=America/New_York * 9-17 * * MON-FRI", updater.UpdateStocks)

	if err != nil {
		log.Errorln(err)
	}

	c.Start()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}
