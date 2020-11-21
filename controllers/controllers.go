package controllers

import (
	"github.com/nagymarci/stock-screener/api"
	"github.com/nagymarci/stock-screener/model"
	"github.com/sirupsen/logrus"

	"github.com/nagymarci/stock-screener/database"

	stockHttp "github.com/nagymarci/stock-commons/http"
)

type Controller struct {
	database *database.Stockinfos
	client   *api.StockScraper
}

func New(db *database.Stockinfos, cl *api.StockScraper) *Controller {
	return &Controller{
		database: db,
		client:   cl,
	}
}

// RegisterStock registers a stock symbol to the watchlist to evaluate it
func (c *Controller) RegisterStock(symbol string) error {
	_, err := c.database.Get(symbol)

	if err == nil {
		return nil
	}

	stockData, err := c.client.Get(symbol)

	if err != nil {
		return stockHttp.NewFailedDependencyError(err.Error())
	}

	err = c.database.Save(stockData)

	if err != nil {
		return stockHttp.NewInternalServerError(err.Error())
	}

	return nil
}

// GetStockInfo returns the information of a stock symbol
func (c *Controller) GetStockInfo(symbol string) (model.StockDataInfo, error) {
	result, err := c.database.Get(symbol)

	if err != nil {
		return result, stockHttp.NewNotFoundError(err.Error())
	}

	return result, nil
}

// GetAllStocks returns the information of all of the stocks
func (c *Controller) GetAllStocks() []model.StockDataInfo {
	result, err := c.database.GetAll()

	if err != nil {
		logrus.Warnln(err)
	}

	return result
}

/*
//UpdateAll updates all stocks in the database
func UpdateAll() {

	go service.UpdateStocks()

	w.WriteHeader(http.StatusOK)
}*/

//DeleteStock deletes the given stock from the database
func (c *Controller) DeleteStock(symbol string) error {
	err := c.database.Delete(symbol)

	if err != nil {
		return stockHttp.NewInternalServerError(err.Error())
	}

	return nil
}
