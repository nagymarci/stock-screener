package service

import (
	"log"
	"sync"
	"time"

	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/model"
)

var mux sync.Mutex

//UpdateStocks checks NextUpdate attribute of the stock and updates it if the time passed
func UpdateStocks() {
	mux.Lock()

	stocks := database.GetAllExpired()
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

		newStockInfo, err := GetWithFields(stockInfo.Ticker, fields)
		if err != nil {
			log.Println(err)
			continue
		}

		newStockInfo.NextUpdate, newStockInfo.DividendYield5yr.NextUpdate, newStockInfo.PeRatio5yr.NextUpdate = model.NextUpdateTimes()

		database.Update(newStockInfo)
	}
	mux.Unlock()
}
