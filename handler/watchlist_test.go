package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-screener/model"
)

type controllerMock struct{}

func (m *controllerMock) Create(request *model.WatchlistRequest) (*model.Watchlist, error) {
	return nil, nil
}

func TestWatchlistCreateHandler(t *testing.T) {
	t.Run("sends bad request when name is missing", func(t *testing.T) {
		router := mux.NewRouter().Path("/watchlist").Subrouter()
		WatchlistCreateHandler(router, &controllerMock{}, func(r *http.Request) string {
			return "asd@asd.com"
		})

		mcPostBody := map[string]interface{}{
			"question_text": "Is this a test post for MutliQuestion?",
		}
		body, _ := json.Marshal(mcPostBody)

		req := httptest.NewRequest(http.MethodPost, "/watchlist", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		res := rec.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Logf("Expected [%d], got [%d]", http.StatusBadRequest, res.StatusCode)
			t.FailNow()
		}
	})
}
