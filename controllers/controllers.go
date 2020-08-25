package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/model"
)

// RegisterStock registers a stock symbol to the watchlist to evaluate it
func RegisterStock(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	result := database.Get(symbol)

	if result.Ticker != "" {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	resp, err := http.Get(config.Get().StockInfoProviderURL + symbol)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	stockData := model.StockDataInfo{}

	err = json.NewDecoder(resp.Body).Decode(&stockData)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(stockData)

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
