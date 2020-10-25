package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/controllers/watchlist"
	"github.com/nagymarci/stock-screener/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WatchlistCreateHandler(mux *mux.Router, watchlist *watchlist.WatchlistController) {
	mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		email := extractEmail(r)

		var watchlistRequest *model.WatchlistRequest

		err := json.NewDecoder(r.Body).Decode(&watchlistRequest)

		if err != nil {
			message := "Failed to deserialize payload: " + err.Error()
			handleErrorResponse(message, w, http.StatusBadRequest)
			log.Println(message)
			return
		}

		if watchlistRequest.Stocks == nil || len(watchlistRequest.Stocks) < 1 {
			message := "Required value 'stocks' is missing"
			handleErrorResponse(message, w, http.StatusBadRequest)
			log.Println(message)
			return
		}

		if len(watchlistRequest.Name) < 1 || watchlistRequest.Name == " " {
			message := "Required value 'name' is missing"
			handleErrorResponse(message, w, http.StatusBadRequest)
			log.Println(message)
			return
		}

		watchlistRequest.Email = email

		result, err := watchlist.Create(watchlistRequest)

		if err != nil {
			message := "Watchlist creation failed: " + err.Error()
			handleErrorResponse(message, w, http.StatusInternalServerError)
			log.Println(message)
			return
		}

		handleJSONResponse(result, w, http.StatusCreated)

	}).Methods(http.MethodPost, http.MethodOptions)
}

func WatchlistDeleteHandler(router *mux.Router, watchlist *watchlist.WatchlistController) {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		email, id, err := extractEmailAndId(r)

		if err != nil {
			handleError(err, w)
			return
		}

		err = watchlist.Delete(id, email)

		if err != nil {
			handleError(err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete, http.MethodOptions)
}

func WatchlistGetAllHandler(router *mux.Router, watchlist *watchlist.WatchlistController) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		email := extractEmail(r)

		result, err := watchlist.GetAll(email)

		if err != nil {
			handleError(err, w)
			return
		}

		handleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func WatchlistGetHandler(router *mux.Router, watchlist *watchlist.WatchlistController) {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		email, id, err := extractEmailAndId(r)

		if err != nil {
			handleError(err, w)
			return
		}

		result, err := watchlist.Get(id, email)

		if err != nil {
			handleError(err, w)
			return
		}

		handleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func WatchlistGetCalculatedHandler(router *mux.Router, watchlist *watchlist.WatchlistController) {
	router.HandleFunc("/{id}/calculated", func(w http.ResponseWriter, r *http.Request) {
		email, id, err := extractEmailAndId(r)

		if err != nil {
			handleError(err, w)
			return
		}

		result, err := watchlist.GetCalculated(id, email)

		if err != nil {
			handleError(err, w)
			return
		}

		handleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func handleError(err error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	if err, ok := interface{}(&err).(model.HttpError); ok {
		statusCode = err.Status()
	}
	message := "Failed to process request: " + err.Error()
	handleErrorResponse(message, w, statusCode)
	log.Println(message)
}

func handleErrorResponse(msg string, w http.ResponseWriter, status int) {
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

func extractEmail(r *http.Request) string {
	user := r.Context().Value("user")
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)[config.Config.EmailClaim].(string)
	log.Printf("User email: %s", email)
	return email
}

func extractEmailAndId(r *http.Request) (string, primitive.ObjectID, error) {
	email := extractEmail(r)
	id, err := extractId(r)

	return email, id, err
}

func extractId(r *http.Request) (primitive.ObjectID, error) {
	id := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		message := "Invalid watchlist id: " + err.Error()
		log.Println(message)
		return primitive.NilObjectID, model.NewBadRequestError(message)
	}

	return objectID, nil
}
