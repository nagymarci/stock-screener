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
	router.HandleFunc("/stocks/calculated", controllers.GetAllRecommendedStock).Methods("GET")
	router.HandleFunc("/stocks/{symbol}", controllers.RegisterStock).Methods("POST")
	router.HandleFunc("/stocks/{symbol}", controllers.GetStockInfo).Methods("GET")
	router.HandleFunc("/stocks/{symbol}", controllers.DeleteStock).Methods("DELETE")
	router.HandleFunc("/stocks/{symbol}/calculated", controllers.GetCalculatedStockInfo).Methods("GET")
	router.HandleFunc("/stocks", controllers.GetAllStocks).Methods("GET")
	router.HandleFunc("/profiles/{name}", controllers.SaveProfile).Methods("POST")
	router.HandleFunc("/profiles/{name}", controllers.DeleteProfile).Methods("DELETE")
	router.HandleFunc("/profiles/{name}/stocks/recommended", controllers.GetRecommendedStocksInProfile).Methods("GET")
	router.HandleFunc("/profiles/{name}/stocks/calculated", controllers.GetCalculatedStocksInProfile).Methods("GET")
	router.HandleFunc("/profiles/{name}/stocks", controllers.GetStocksInProfile).Methods("GET")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}
