package watchlist

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/nagymarci/stock-screener/model"
	"github.com/nagymarci/stock-screener/service"
)

type WatchlistController struct {
	watchlists database.WatchlistCollection
}

func NewWatchlistController(w database.WatchlistCollection) *WatchlistController {
	return &WatchlistController{
		watchlists: w,
	}
}

//Create creates a new watchlist
func (wl *WatchlistController) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")

	email := user.(*jwt.Token).Claims.(jwt.MapClaims)[config.Config.EmailClaim].(string)
	log.Printf("User email: %s", email)

	var watchlistRequest model.WatchlistRequest

	err := json.NewDecoder(r.Body).Decode(&watchlistRequest)

	if err != nil {
		message := "Failed to deserialize payload: " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	if watchlistRequest.Stocks == nil || len(watchlistRequest.Stocks) < 1 {
		message := "Required value 'stocks' is missing"
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	if len(watchlistRequest.Name) < 1 || watchlistRequest.Name == " " {
		message := "Required value 'name' is missing"
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	watchlistRequest.Email = email
	var addedStocks []string

	for _, symbol := range watchlistRequest.Stocks {
		err = saveStock(symbol)

		if err != nil {
			log.Println(err)
			continue
		}

		addedStocks = append(addedStocks, symbol)
	}

	watchlistRequest.Stocks = addedStocks
	id, err := wl.watchlists.Create(watchlistRequest)

	if err != nil {
		message := "Watchlist creation failed: " + err.Error()
		handleError(message, w, http.StatusInternalServerError)
		log.Println(message)
		return
	}

	watchlistResponse := model.Watchlist{
		ID:     id,
		Name:   watchlistRequest.Name,
		Stocks: watchlistRequest.Stocks,
		Email:  watchlistRequest.Email}

	handleJSONResponse(watchlistResponse, w, http.StatusCreated)
}

//Delete deletes the specified watchlist if that belongs to the authorized user
func (wl *WatchlistController) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user := r.Context().Value("user")
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)[config.Config.EmailClaim].(string)
	log.Printf("User email: %s", email)

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		message := "Invalid watchlist id: " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	_, err = wl.getAndValidateUserAuthorization(objectID, email)

	if err != nil {
		message := "Cannot remove watchlist " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	result, err := wl.watchlists.Delete(objectID)

	if result != 1 || err != nil {
		errorText := ""
		if err != nil {
			errorText = err.Error()
		}
		message := "Failed to delete watchlist: " + errorText
		handleError(message, w, http.StatusInternalServerError)
		log.Println(message)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (wl *WatchlistController) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user := r.Context().Value("user")
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)[config.Config.EmailClaim].(string)
	log.Printf("User email: %s", email)

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		message := "Invalid watchlist id: " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	watchlist, err := wl.getAndValidateUserAuthorization(objectID, email)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	handleJSONResponse(watchlist, w, http.StatusOK)
}

func (wl *WatchlistController) GetAll(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)[config.Config.EmailClaim].(string)
	log.Printf("User email: %s", email)

	watchlists, err := wl.watchlists.GetAll(email)

	if err != nil {
		message := "Unable to list watchlists " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	handleJSONResponse(watchlists, w, http.StatusOK)
}

func (wl *WatchlistController) GetCalculated(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user := r.Context().Value("user")
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)[config.Config.EmailClaim].(string)
	log.Printf("User email: %s", email)

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		message := "Invalid watchlist id: " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	watchlist, err := wl.getAndValidateUserAuthorization(objectID, email)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		handleError(message, w, http.StatusBadRequest)
		log.Println(message)
		return
	}

	var stockInfos []model.CalculatedStockInfo

	for _, symbol := range watchlist.Stocks {
		result, err := database.Get(symbol)

		if err != nil {
			log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
			continue
		}

		calculatedStockInfo := service.Calculate(&result)

		stockInfos = append(stockInfos, calculatedStockInfo)
	}

	handleJSONResponse(stockInfos, w, http.StatusOK)
}

func (w *WatchlistController) getAndValidateUserAuthorization(id primitive.ObjectID, email string) (model.Watchlist, error) {
	watchlist, err := w.watchlists.Get(id)
	if err != nil {
		return watchlist, err
	}

	if watchlist.Email != email {
		return watchlist, errors.New("Watchlist does not belong to user")
	}

	return watchlist, err
}

func saveStock(stock string) error {
	_, err := database.Get(stock)

	if err == nil {
		return err
	}

	stockData, err := service.Get(stock)

	if err != nil {
		return err
	}

	err = database.Save(stockData)

	return err
}

func handleError(msg string, w http.ResponseWriter, status int) {
	response := model.ErrorResponse{Message: msg}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, model.UnknownError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}

func handleJSONResponse(object interface{}, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(object)
}
