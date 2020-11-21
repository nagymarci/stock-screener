package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/nagymarci/stock-screener/model"
	"github.com/sirupsen/logrus"

	"github.com/nagymarci/stock-screener/database"
)

type Updater struct {
	mux                    sync.Mutex
	database               *database.Stockinfos
	stockClient            getStockWithFields
	stockUpdateInterval    string
	peUpdateInterval       string
	divYieldUpdateInterval string
}

type getStockWithFields interface {
	GetWithFields(symbol string, fields []string) (model.StockDataInfo, error)
}

func New(db *database.Stockinfos, sc getStockWithFields, stockInterval, peInterval, divInterval string) *Updater {
	return &Updater{
		database:               db,
		stockClient:            sc,
		stockUpdateInterval:    stockInterval,
		peUpdateInterval:       peInterval,
		divYieldUpdateInterval: divInterval,
	}
}

//UpdateStocks checks NextUpdate attribute of the stock and updates it if the time passed
func (u *Updater) UpdateStocks() {
	log := logrus.WithField("component", "updater")
	u.mux.Lock()
	defer u.mux.Unlock()

	stocks, err := u.database.GetAllExpired()
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infof("Updating [%d] stocks\n", len(stocks))

	now := time.Now()

	for _, stockInfo := range stocks {
		fields := []string{}
		if stockInfo.NextUpdate.Before(now) {
			fields = append(fields, "price", "eps", "div")
		}
		if stockInfo.DividendYield5yr.NextUpdate.Before(now) {
			fields = append(fields, "divHist")
		}
		if stockInfo.PeRatio5yr.NextUpdate.Before(now) {
			fields = append(fields, "pe")
		}

		newStockInfo, err := u.stockClient.GetWithFields(stockInfo.Ticker, fields)
		if err != nil {
			log.WithField("ticker", stockInfo.Ticker).Warningln(err)
			continue
		}

		u.calculateNextUpdateTimes(&newStockInfo)

		u.database.Update(newStockInfo)
	}
}

//CalculateNextUpdateTimes calculates the next update times based on the configuration
func (u *Updater) calculateNextUpdateTimes(stock *model.StockDataInfo) {
	stockUpdateInterval, _ := time.ParseDuration(u.stockUpdateInterval)
	peUpdateInterval, _ := time.ParseDuration(u.peUpdateInterval)
	divYieldUpdateInterval, _ := time.ParseDuration(u.divYieldUpdateInterval)

	randMinutes := rand.Intn(30)
	randMinutesInterval, _ := time.ParseDuration(fmt.Sprintf("%dm", randMinutes))

	randHours := rand.Intn(24)
	randHoursInterval, _ := time.ParseDuration(fmt.Sprintf("%dh", randHours))

	stock.NextUpdate = time.Now().Add(stockUpdateInterval).Add(randMinutesInterval)
	stock.PeRatio5yr.NextUpdate = time.Now().Add(peUpdateInterval).Add(randHoursInterval)
	stock.DividendYield5yr.NextUpdate = time.Now().Add(divYieldUpdateInterval).Add(randHoursInterval)
}
