package model

import (
	"sync"
	"time"

	"github.com/nagymarci/stock-screener/config"
)

type pERatioInfo struct {
	Avg        float64   `json:"avg" bson:"avg"`
	Min        float64   `json:"min" bson:"min"`
	NextUpdate time.Time `json:"-" bson:"nextUpdate"`
}

type dividendYieldInfo struct {
	Avg        float64   `json:"avg" bson:"avg"`
	Max        float64   `json:"max" bson:"max"`
	NextUpdate time.Time `json:"-" bson:"nextUpdate"`
}

//StockDataInfo holds the information for one stock
type StockDataInfo struct {
	Ticker           string            `json:"ticker" bson:"ticker"`
	Price            float64           `json:"price" bson:"price"`
	Eps              float64           `json:"eps" bson:"eps"`
	Dividend         float64           `json:"dividend" bson:"dividend"`
	PeRatio5yr       pERatioInfo       `json:"peRatio5yr" bson:"peRatio5yr"`
	DividendYield5yr dividendYieldInfo `json:"dividendYield5yr" bson:"dividendYield5yr"`
	NextUpdate       time.Time         `json:"-" bson:"nextUpdate"`
}

//CalculatedStockInfo holds the data calculated for investment suggestions
type CalculatedStockInfo struct {
	Ticker         string  `json:"ticker"`
	Price          float64 `json:"price"`
	OptInPrice     float64 `json:"optInPrice"`
	PriceColor     string  `json:"priceColor"`
	AnnualDividend float64 `json:"dividend"`
	DividendYield  float64 `json:"dividendYield"`
	OptInYield     float64 `json:"optInYield"`
	DividendColor  string  `json:"dividendColor"`
	CurrentPe      float64 `json:"currentPe"`
	OptInPe        float64 `json:"optInPe"`
	PeColor        string  `json:"pecolor"`
}

type sp500DivYield struct {
	Yield      float64
	NextUpdate time.Time
	Mux        sync.Mutex
}

//Sp500DivYield stores information of the S&P500 dividend yield, and when we should update it next
var Sp500DivYield sp500DivYield

//NextUpdateTimes calculates the next update times based on the configuration
func NextUpdateTimes() (time.Time, time.Time, time.Time) {
	stockUpdateInterval, _ := time.ParseDuration(config.Config.StockUpdateInterval)
	peUpdateInterval, _ := time.ParseDuration(config.Config.PeUpdateInterval)
	divYieldUpdateInterval, _ := time.ParseDuration(config.Config.DivYieldUpdateInterval)

	return time.Now().Add(stockUpdateInterval), time.Now().Add(peUpdateInterval), time.Now().Add(divYieldUpdateInterval)
}
