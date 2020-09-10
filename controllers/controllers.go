package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/nagymarci/stock-screener/model"

	"github.com/nagymarci/stock-screener/service"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/database"
)

// RegisterStock registers a stock symbol to the watchlist to evaluate it
func RegisterStock(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	_, err := database.Get(symbol)

	if err == nil {
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

	stockData.CalculateNextUpdateTimes()

	err = database.Save(stockData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(stockData)
}

// GetStockInfo returns the information of a stock symbol
func GetStockInfo(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	result, err := database.Get(symbol)

	if err != nil {
		log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
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

//GetCalculatedStockInfo returns the calculated informatin of a stock
func GetCalculatedStockInfo(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	stockInfo, err := database.Get(symbol)

	if err != nil {
		log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}

	calculatedStockInfo := service.Calculate(&stockInfo)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calculatedStockInfo)
}

//UpdateAll updates all stocks in the database
func UpdateAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating all stocks")

	go service.UpdateStocks()

	w.WriteHeader(http.StatusOK)
}

//DeleteStock deletes the given stock from the database
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	log.Printf("Delete [%s]", symbol)

	err := database.Delete(symbol)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

//GetAllRecommendedStock calculates the data for all stocks and returns the recommended ones
func GetAllRecommendedStock(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllCalculatedStockInfo")
	min := r.FormValue("min")

	if min == "" {
		min = "3"
	}

	numReqs, err := strconv.Atoi(min)

	if err != nil || numReqs < 1 || numReqs > 3 {
		log.Println("Invalid parameter, changing to 3", err)
		numReqs = 3
	}

	log.Println(numReqs, err)

	stocks := database.GetAll()

	result := service.GetAllRecommendedStock(stocks, numReqs)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}

func SaveProfile(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	var stocks model.Stocks

	err := json.NewDecoder(r.Body).Decode(&stocks)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	database.DeleteProfile(name)

	profile := model.Profile{Name: name}

	for _, symbol := range stocks.Values {
		_, err := database.Get(symbol)

		if err == nil {
			profile.Stocks = append(profile.Stocks, symbol)
			continue
		}

		stockData, err := service.Get(symbol)

		if err != nil {
			log.Println(err)
			continue
		}

		err = database.Save(stockData)

		if err != nil {
			log.Println(err)
			continue
		}

		profile.Stocks = append(profile.Stocks, symbol)
	}

	err = database.SaveProfile(profile)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(profile)
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	err := database.DeleteProfile(name)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetStocksInProfile(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	profile, err := database.GetProfile(name)

	if err != nil {
		log.Printf("Failed to get profile [%s]: [%v]\n", name, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}

	var stockInfos []model.StockDataInfo

	for _, symbol := range profile.Stocks {
		result, err := database.Get(symbol)

		if err != nil {
			log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
			continue
		}

		stockInfos = append(stockInfos, result)
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(stockInfos)

}

func GetCalculatedStocksInProfile(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	profile, err := database.GetProfile(name)

	if err != nil {
		log.Printf("Failed to get profile [%s]: [%v]\n", name, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}

	var stockInfos []model.CalculatedStockInfo

	for _, symbol := range profile.Stocks {
		result, err := database.Get(symbol)

		if err != nil {
			log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
			continue
		}

		calculatedStockInfo := service.Calculate(&result)

		stockInfos = append(stockInfos, calculatedStockInfo)
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(stockInfos)

}

func GetRecommendedStocksInProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllCalculatedStockInfo")
	name := mux.Vars(r)["name"]
	min := r.FormValue("min")

	if min == "" {
		min = "3"
	}

	numReqs, err := strconv.Atoi(min)

	if err != nil || numReqs < 1 || numReqs > 3 {
		log.Println("Invalid parameter, changing to 3", err)
		numReqs = 3
	}

	log.Println(numReqs, err)

	profile, err := database.GetProfile(name)

	if err != nil {
		log.Printf("Failed to get profile [%s]: [%v]\n", name, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}

	var stockInfos []model.StockDataInfo

	for _, symbol := range profile.Stocks {
		result, err := database.Get(symbol)

		if err != nil {
			log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
			continue
		}

		stockInfos = append(stockInfos, result)
	}

	result := service.GetAllRecommendedStock(stockInfos, numReqs)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)

}
