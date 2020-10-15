package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/controllers"
)

//Route configures the routing
func Route(router *mux.Router) {
	router.Use(corsMiddleware)

	stocks := router.PathPrefix("/stocks").Subrouter()
	stocks.HandleFunc("/update", controllers.UpdateAll).Methods("POST")
	stocks.HandleFunc("/recommended", controllers.GetAllRecommendedStock).Methods("GET")
	stocks.HandleFunc("/calculated", controllers.GetAllCalculatedStocks).Methods("GET")
	stocks.HandleFunc("/{symbol}", controllers.RegisterStock).Methods("POST")
	stocks.HandleFunc("/{symbol}", controllers.GetStockInfo).Methods("GET")
	stocks.HandleFunc("/{symbol}", controllers.DeleteStock).Methods("DELETE")
	stocks.HandleFunc("/{symbol}/calculated", controllers.GetCalculatedStockInfo).Methods("GET")
	stocks.HandleFunc("/", controllers.GetAllStocks).Methods("GET")

	profiles := router.PathPrefix("/profiles").Subrouter()
	profiles.HandleFunc("/{name}", controllers.SaveProfile).Methods("POST")
	profiles.HandleFunc("/{name}", controllers.DeleteProfile).Methods(http.MethodDelete, http.MethodOptions)
	profiles.HandleFunc("/{name}/stocks/recommended", controllers.GetRecommendedStocksInProfile).Methods("GET")
	profiles.HandleFunc("/{name}/stocks/calculated", controllers.GetCalculatedStocksInProfile).Methods("GET")
	profiles.HandleFunc("/{name}/stocks", controllers.GetStocksInProfile).Methods("GET")
	profiles.HandleFunc("/", controllers.ListProfiles).Methods("GET")

	router.HandleFunc("/notifyTest", controllers.NotifyTest).Methods("GET")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
