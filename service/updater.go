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
	log.Printf("Updating [%d] stocks\n", len(stocks))

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

		newStockInfo.CalculateNextUpdateTimes()

		database.Update(newStockInfo)
	}
	mux.Unlock()
}

func NotifyChanges() {
	profiles, err := database.GetAllProfileName()

	if err != nil {
		log.Printf("Failed to get profiles [%v]", err)
		return
	}

	for _, profileName := range profiles {
		previouStocks, _ := database.GetPreviouslyRecommendedStocks(profileName)

		profile, err := database.GetProfile(profileName)

		if err != nil {
			log.Printf("Failed to get profile [%s]: [%v]\n", profileName, err)
			continue
		}

		var stockInfos []model.StockDataInfo

		for _, symbol := range profile.Stocks {
			result, err := database.Get(symbol)

			if err != nil {
				log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
				continue
			}

			stockInfos = append(stockInfos, result)
		}

		calculatedStockData := GetAllRecommendedStock(stockInfos, 2)

		currentStocks := filterGreenPrices(calculatedStockData)

		removed, added := getChanges(previouStocks, currentStocks)

		if len(removed) == 0 && len(added) == 0 {
			continue
		}

		err = sendNotification(profileName, removed, added, currentStocks)

		if err != nil {
			log.Printf("Failed to send notification for profile [%v], [%v]", profileName, err)
			continue
		}

		database.SaveRecommendation(profileName, currentStocks)
	}
}

func NotifyChanges_watchlist() {
	watchlists, err := database.WatchLists.List()

	if err != nil {
		log.Printf("Failed to get watchlists [%v]", err)
		return
	}

	for _, watchlist := range watchlists {
		previouStocks, _ := database.Recommendations.Get(watchlist.ID)

		var stockInfos []model.StockDataInfo

		for _, symbol := range watchlist.Stocks {
			result, err := database.Get(symbol)

			if err != nil {
				log.Printf("Failed to get stock [%s]: [%v]\n", symbol, err)
				continue
			}

			stockInfos = append(stockInfos, result)
		}

		calculatedStockData := GetAllRecommendedStock(stockInfos, 2)

		currentStocks := filterGreenPrices(calculatedStockData)

		removed, added := getChanges(previouStocks, currentStocks)

		if len(removed) == 0 && len(added) == 0 {
			continue
		}

		err = sendNotification_watchlist(watchlist.Name, removed, added, currentStocks, watchlist.Email)

		if err != nil {
			log.Printf("Failed to send notification for profile [%v], [%v]", watchlist.ID, err)
			continue
		}

		database.Recommendations.Update(watchlist.ID, currentStocks)
	}
}

func filterGreenPrices(stockInfos []model.CalculatedStockInfo) []string {
	var result []string

	for _, calc := range stockInfos {
		if calc.PriceColor != "green" {
			continue
		}

		result = append(result, calc.Ticker)
	}
	return result
}

func getChanges(old, new []string) ([]string, []string) {
	if len(old) == 0 || len(new) == 0 {
		return old, new
	}

	var removed []string

	for _, symbol := range old {
		if !contains(new, symbol) {
			removed = append(removed, symbol)
		}
	}

	var added []string

	for _, symbol := range new {
		if !contains(old, symbol) {
			added = append(added, symbol)
		}
	}

	return removed, added
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
