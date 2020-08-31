package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/controllers"
)

//Route configures the routing
func Route(router *mux.Router) {
	router.HandleFunc("/", welcome)
	router.HandleFunc("/stocks/update", controllers.UpdateAll).Methods("POST")
	router.HandleFunc("/stocks/{symbol}", controllers.RegisterStock).Methods("POST")
	router.HandleFunc("/stocks/{symbol}", controllers.GetStockInfo).Methods("GET")
	router.HandleFunc("/stocks/{symbol}/calculated", controllers.GetCalculatedStockInfo).Methods("GET")
	router.HandleFunc("/stocks", controllers.GetAllStocks).Methods("GET")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}
