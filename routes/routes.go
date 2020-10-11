package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/controllers"
)

//Route configures the routing
func Route(router *mux.Router) {
	router.Use(corsMiddleware)
	router.HandleFunc("/", welcome)
	router.HandleFunc("/stocks/update", controllers.UpdateAll).Methods("POST")
	router.HandleFunc("/stocks/recommended", controllers.GetAllRecommendedStock).Methods("GET")
	router.HandleFunc("/stocks/calculated", controllers.GetAllCalculatedStocks).Methods("GET")
	router.HandleFunc("/stocks/{symbol}", controllers.RegisterStock).Methods("POST")
	router.HandleFunc("/stocks/{symbol}", controllers.GetStockInfo).Methods("GET")
	router.HandleFunc("/stocks/{symbol}", controllers.DeleteStock).Methods("DELETE")
	router.HandleFunc("/stocks/{symbol}/calculated", controllers.GetCalculatedStockInfo).Methods("GET")
	router.HandleFunc("/stocks", controllers.GetAllStocks).Methods("GET")
	router.HandleFunc("/profiles/{name}", controllers.SaveProfile).Methods("POST")
	router.HandleFunc("/profiles/{name}", controllers.DeleteProfile).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/profiles/{name}/stocks/recommended", controllers.GetRecommendedStocksInProfile).Methods("GET")
	router.HandleFunc("/profiles/{name}/stocks/calculated", controllers.GetCalculatedStocksInProfile).Methods("GET")
	router.HandleFunc("/profiles/{name}/stocks", controllers.GetStocksInProfile).Methods("GET")
	router.HandleFunc("/profiles", controllers.ListProfiles).Methods("GET")

	router.HandleFunc("/notifyTest", controllers.NotifyTest).Methods("GET")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
