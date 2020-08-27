package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nagymarci/stock-screener/service"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/database"
)

// RegisterStock registers a stock symbol to the watchlist to evaluate it
func RegisterStock(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	result := database.Get(symbol)

	if result.Ticker != "" {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	stockData, err := service.Get(symbol)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusFailedDependency)
		fmt.Fprint(w, err)
		return
	}

	duration, _ := time.ParseDuration("1h")
	stockData.NextUpdate = time.Now().Add(duration)
	stockData.DividendYield5yr.NextUpdate = time.Now().Add(duration)
	stockData.PeRatio5yr.NextUpdate = time.Now().Add(duration)

	database.Save(stockData)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(stockData)
}

// GetStockInfo returns the information of a stock symbol
func GetStockInfo(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	result := database.Get(symbol)

	if result.Ticker == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{}")
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}

// GetAllStocks returns the information of all of the stocks
func GetAllStocks(w http.ResponseWriter, r *http.Request) {

	result := database.GetAll()

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}
