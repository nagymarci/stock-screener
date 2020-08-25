package model

import "time"

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
	PeRatio5yr       pERatioInfo       `json:"peRatio5yr" bson:"peRatio5yr"`
	DividendYield5yr dividendYieldInfo `json:"dividendYield5yr" bson:"dividendYield5yr"`
	NextUpdate       time.Time         `json:"-" bson:"nextUpdate"`
}
