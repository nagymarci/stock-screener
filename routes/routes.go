package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/controller"
)

//Route configures the routing
func Route(router *mux.Router) {
	router.HandleFunc("/", welcome)
	router.HandleFunc("/register/{symbol}", controller.RegisterStock)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}
