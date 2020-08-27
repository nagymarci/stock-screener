package service

import (
	"log"
	"time"

	"github.com/nagymarci/stock-screener/database"
)

//UpdateStocks checks NextUpdate attribute of the stock and updates it if the time passed
func UpdateStocks() {
	stocks := database.GetAllExpired()
	now := time.Now()

	for _, stockInfo := range stocks {
		fields := []string{}
		if stockInfo.NextUpdate.Before(now) {
			fields = append(fields, "price", "eps")
		}
		if stockInfo.DividendYield5yr.NextUpdate.Before(now) {
			fields = append(fields, "div")
		}
		if stockInfo.PeRatio5yr.NextUpdate.Before(now) {
			fields = append(fields, "pe")
		}

		newStockInfo, err := GetWithFields(stockInfo.Ticker, fields)
		if err != nil {
			log.Println(err)
			continue
		}

		duration, _ := time.ParseDuration("1h")
		newStockInfo.NextUpdate = time.Now().Add(duration)
		newStockInfo.DividendYield5yr.NextUpdate = time.Now().Add(duration)
		newStockInfo.PeRatio5yr.NextUpdate = time.Now().Add(duration)

		database.Update(newStockInfo)
	}

}
