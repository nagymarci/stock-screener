package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/routes"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	conf := config.Get()
	routes.Route(router)
	database.Connect(conf.DatabaseConnectionString)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), router))
}
