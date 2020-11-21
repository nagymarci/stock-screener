package routes

import (
	"net/http"

	"github.com/nagymarci/stock-screener/handler"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/controllers"
	"github.com/urfave/negroni"
)

//Route configures the routing
func Route(controller *controllers.Controller) http.Handler {
	router := mux.NewRouter()

	stocks := router.PathPrefix("/stocks").Subrouter()
	handler.RegisterStockHandler(stocks, controller)
	handler.GetStockInfoHandler(stocks, controller)
	handler.DeleteStockHandler(stocks, controller)
	handler.GetAllStocksHandler(stocks, controller)

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	n := negroni.New(recovery, negroni.NewLogger())
	n.UseHandler(router)
	return n
}
