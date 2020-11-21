package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"github.com/nagymarci/stock-screener/controllers"

	stockHttp "github.com/nagymarci/stock-commons/http"
)

// RegisterStock registers a stock symbol to the watchlist to evaluate it
func RegisterStockHandler(router *mux.Router, controller *controllers.Controller) {
	router.HandleFunc("/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		symbol := mux.Vars(r)["symbol"]

		log := logrus.WithField("symbol", symbol)

		err := controller.RegisterStock(symbol)

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodPost, http.MethodOptions)
}

// GetStockInfo returns the information of a stock symbol
func GetStockInfoHandler(router *mux.Router, controller *controllers.Controller) {
	router.HandleFunc("/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		symbol := mux.Vars(r)["symbol"]

		log := logrus.WithField("symbol", symbol)

		result, err := controller.GetStockInfo(symbol)

		if err != nil {
			log.Errorln(err)
			stockHttp.HandleError(err, w)
			return
		}

		stockHttp.HandleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

// GetAllStocks returns the information of all of the stocks
func GetAllStocksHandler(router *mux.Router, controller *controllers.Controller) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {

		result := controller.GetAllStocks()

		stockHttp.HandleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

/*
//UpdateAll updates all stocks in the database
func UpdateAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating all stocks")

	go service.UpdateStocks()

	w.WriteHeader(http.StatusOK)
}
*/

//DeleteStock deletes the given stock from the database
func DeleteStockHandler(router *mux.Router, controller *controllers.Controller) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		symbol := mux.Vars(r)["symbol"]

		log := logrus.WithField("symbol", symbol)

		err := controller.DeleteStock(symbol)

		if err != nil {
			log.Println(err)
			stockHttp.HandleError(err, w)
		}

		w.WriteHeader(http.StatusNoContent)

	}).Methods(http.MethodDelete)
}
