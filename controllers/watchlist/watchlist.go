package watchlist

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nagymarci/stock-screener/database"

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
func (wl *WatchlistController) Create(request *model.WatchlistRequest) (*model.Watchlist, error) {
	var addedStocks []string

	for _, symbol := range request.Stocks {
		err := saveStock(symbol)

		if err != nil {
			log.Println(err)
			continue
		}

		addedStocks = append(addedStocks, symbol)
	}

	request.Stocks = addedStocks
	id, err := wl.watchlists.Create(*request)

	if err != nil {
		return nil, model.NewInternalServerError(err.Error())
	}

	watchlistResponse := model.Watchlist{
		ID:     id,
		Name:   request.Name,
		Stocks: request.Stocks,
		Email:  request.Email}

	return &watchlistResponse, err
}

//Delete deletes the specified watchlist if that belongs to the authorized user
func (wl *WatchlistController) Delete(id primitive.ObjectID, email string) error {
	_, err := wl.getAndValidateUserAuthorization(id, email)

	if err != nil {
		return model.NewBadRequestError(err.Error())
	}

	result, err := wl.watchlists.Delete(id)

	if result != 1 {
		return model.NewInternalServerError("No object were removed from database")
	}

	if err != nil {
		return model.NewInternalServerError(err.Error())
	}

	return nil
}

func (wl *WatchlistController) Get(id primitive.ObjectID, email string) (model.Watchlist, error) {
	watchlist, err := wl.getAndValidateUserAuthorization(id, email)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		log.Println(message)
		return model.Watchlist{}, model.NewBadRequestError(message)
	}

	return watchlist, nil
}

func (wl *WatchlistController) GetAll(email string) ([]model.Watchlist, error) {
	watchlists, err := wl.watchlists.GetAll(email)

	if err != nil {
		message := "Unable to list watchlists " + err.Error()
		log.Println(message)
		return nil, model.NewBadRequestError(message)
	}

	return watchlists, nil
}

func (wl *WatchlistController) GetCalculated(id primitive.ObjectID, email string) ([]model.CalculatedStockInfo, error) {
	watchlist, err := wl.getAndValidateUserAuthorization(id, email)

	if err != nil {
		message := "Cannot read watchlist " + err.Error()
		log.Println(message)
		return nil, model.NewBadRequestError(message)
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

	return stockInfos, nil
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
