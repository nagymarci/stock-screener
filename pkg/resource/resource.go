package resource

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterStock registers a stock symbol to the watchlist to evaluate it
func RegisterStock(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	fmt.Fprintf(w, "Register %s!", symbol)
}
