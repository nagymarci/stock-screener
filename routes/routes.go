package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/controllers"
	"github.com/urfave/negroni"
)

//Route configures the routing
func Route(watchlistController *watchlistControllers.WatchlistController) http.Handler {
	router := mux.NewRouter()

	stocks := router.PathPrefix("/stocks").Subrouter()
	stocks.HandleFunc("/update", controllers.UpdateAll).Methods("POST")
	stocks.HandleFunc("/recommended", controllers.GetAllRecommendedStock).Methods("GET")
	stocks.HandleFunc("/calculated", controllers.GetAllCalculatedStocks).Methods("GET")
	stocks.HandleFunc("/{symbol}", controllers.RegisterStock).Methods("POST")
	stocks.HandleFunc("/{symbol}", controllers.GetStockInfo).Methods("GET")
	stocks.HandleFunc("/{symbol}", controllers.DeleteStock).Methods("DELETE")
	stocks.HandleFunc("/{symbol}/calculated", controllers.GetCalculatedStockInfo).Methods("GET")
	stocks.HandleFunc("", controllers.GetAllStocks).Methods("GET")

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	n := negroni.New(recovery, negroni.NewLogger())
	n.UseHandler(router)
	return n
}
