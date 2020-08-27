package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/routes"
	"github.com/nagymarci/stock-screener/service"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	conf := config.Get()
	routes.Route(router)
	database.Connect(conf.DatabaseConnectionString)

	go service.UpdateStocks()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), router))
}
