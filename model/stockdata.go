package model

import (
	"time"

	"github.com/nagymarci/stock-screener/config"
)

type pERatioInfo struct {
	Avg        float32   `json:"avg" bson:"avg"`
	Min        float32   `json:"min" bson:"min"`
	NextUpdate time.Time `json:"-" bson:"nextUpdate"`
}

type dividendYieldInfo struct {
	Avg        float32   `json:"avg" bson:"avg"`
	Max        float32   `json:"max" bson:"max"`
	NextUpdate time.Time `json:"-" bson:"nextUpdate"`
}

//StockDataInfo holds the information for one stock
type StockDataInfo struct {
	Ticker           string            `json:"ticker" bson:"ticker"`
	Price            float32           `json:"price" bson:"price"`
	Eps              float32           `json:"eps" bson:"eps"`
	Dividend         float32           `json:"dividend" bson:"dividend"`
	PeRatio5yr       pERatioInfo       `json:"peRatio5yr" bson:"peRatio5yr"`
	DividendYield5yr dividendYieldInfo `json:"dividendYield5yr" bson:"dividendYield5yr"`
	NextUpdate       time.Time         `json:"-" bson:"nextUpdate"`
}

//NextUpdateTimes calculates the next update times based on the configuration
func NextUpdateTimes() (time.Time, time.Time, time.Time) {
	stockUpdateInterval, _ := time.ParseDuration(config.Config.StockUpdateInterval)
	peUpdateInterval, _ := time.ParseDuration(config.Config.PeUpdateInterval)
	divYieldUpdateInterval, _ := time.ParseDuration(config.Config.DivYieldUpdateInterval)

	return time.Now().Add(stockUpdateInterval), time.Now().Add(peUpdateInterval), time.Now().Add(divYieldUpdateInterval)
}
