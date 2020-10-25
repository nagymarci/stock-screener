package service

import (
	"log"

	"github.com/nagymarci/stock-screener/database"
	"github.com/nagymarci/stock-screener/model"
)

type Notifier struct {
	recommendations database.RecommendationCollection
	watchlists      database.WatchlistCollection
}

func NewNotifier(r database.RecommendationCollection, w database.WatchlistCollection) *Notifier {
	return &Notifier{
		recommendations: r,
		watchlists:      w,
	}
}

func (n *Notifier) NotifyChanges() {
	watchlists, err := n.watchlists.List()

	if err != nil {
		log.Printf("Failed to get watchlists [%v]", err)
		return
	}

	for _, watchlist := range watchlists {
		previouStocks, _ := n.recommendations.Get(watchlist.ID)

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

		n.recommendations.Update(watchlist.ID, currentStocks)
	}
}
