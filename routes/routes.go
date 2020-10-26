package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/controllers"
	watchlistControllers "github.com/nagymarci/stock-screener/controllers/watchlist"
	"github.com/nagymarci/stock-screener/handler"
	"github.com/urfave/negroni"
)

//Route configures the routing
func Route(watchlistController *watchlistControllers.WatchlistController) http.Handler {
	router := mux.NewRouter()

	router.Use(corsMiddleware)

	stocks := router.PathPrefix("/stocks").Subrouter()
	stocks.HandleFunc("/update", controllers.UpdateAll).Methods("POST")
	stocks.HandleFunc("/recommended", controllers.GetAllRecommendedStock).Methods("GET")
	stocks.HandleFunc("/calculated", controllers.GetAllCalculatedStocks).Methods("GET")
	stocks.HandleFunc("/{symbol}", controllers.RegisterStock).Methods("POST")
	stocks.HandleFunc("/{symbol}", controllers.GetStockInfo).Methods("GET")
	stocks.HandleFunc("/{symbol}", controllers.DeleteStock).Methods("DELETE")
	stocks.HandleFunc("/{symbol}/calculated", controllers.GetCalculatedStockInfo).Methods("GET")
	stocks.HandleFunc("", controllers.GetAllStocks).Methods("GET")

	profiles := router.PathPrefix("/profiles").Subrouter()
	auth := negroni.New(
		negroni.HandlerFunc(AuthorizationMiddleware().HandlerWithNext),
		negroni.HandlerFunc(ScopeMiddleware))
	profiles.Handle("/{name}", auth.With(negroni.Wrap(http.HandlerFunc(controllers.SaveProfile)))).Methods("POST")
	profiles.Handle("/{name}", auth.With(negroni.Wrap(http.HandlerFunc(controllers.DeleteProfile)))).Methods(http.MethodDelete, http.MethodOptions)
	profiles.HandleFunc("/{name}/stocks/recommended", controllers.GetRecommendedStocksInProfile).Methods("GET")
	profiles.HandleFunc("/{name}/stocks/calculated", controllers.GetCalculatedStocksInProfile).Methods("GET")
	profiles.HandleFunc("/{name}/stocks", controllers.GetStocksInProfile).Methods("GET")
	profiles.HandleFunc("", controllers.ListProfiles).Methods("GET")

	router.HandleFunc("/crashTest", func(w http.ResponseWriter, req *http.Request) {
		panic("shit")
	}).Methods(http.MethodGet)

	router.HandleFunc("/logTest", func(w http.ResponseWriter, req *http.Request) {
		return
	}).Methods(http.MethodGet)

	watchlist := mux.NewRouter().PathPrefix("/watchlist").Subrouter()
	handler.WatchlistCreateHandler(watchlist, watchlistController, handler.DefaultExtractEmail)
	handler.WatchlistGetAllHandler(watchlist, watchlistController, handler.DefaultExtractEmail)
	handler.WatchlistDeleteHandler(watchlist, watchlistController, handler.DefaultExtractEmail)
	handler.WatchlistGetHandler(watchlist, watchlistController, handler.DefaultExtractEmail)
	handler.WatchlistGetCalculatedHandler(watchlist, watchlistController, handler.DefaultExtractEmail)

	router.PathPrefix("/watchlist").Handler(auth.With(negroni.Wrap(watchlist)))

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	n := negroni.New(recovery, negroni.NewLogger())
	n.UseHandler(router)
	return n
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
