package model

import (
	"time"
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

//Stocks represent list of stocks
type Stocks struct {
	Values []string `json:"values"`
}
