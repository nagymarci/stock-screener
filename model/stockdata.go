package model

type pERatioInfo struct {
	Avg float32 `json:"avg"`
	Min float32 `json:"min"`
}

type dividendYieldInfo struct {
	Avg float32 `json:"avg"`
	Max float32 `json:"max"`
}

//StockDataInfo holds the information for one stock
type StockDataInfo struct {
	Ticker           string            `json:"ticker"`
	Price            float32           `json:"price"`
	Eps              float32           `json:"eps"`
	PeRatio5yr       pERatioInfo       `json:"peRatio5yr"`
	DividendYield5yr dividendYieldInfo `json:"dividendYield5yr"`
}
