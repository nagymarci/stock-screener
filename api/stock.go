package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nagymarci/stock-screener/model"
)

type StockScraper struct {
	host string
}

func New(h string) *StockScraper {
	return &StockScraper{
		host: h,
	}
}

//Get returns the requested stock from the provider
func (ss *StockScraper) Get(symbol string) (model.StockDataInfo, error) {
	return ss.GetWithFields(symbol, []string{})
}

//GetWithFields returns the stock from the provider with the requested fields filled
func (ss *StockScraper) GetWithFields(symbol string, fields []string) (model.StockDataInfo, error) {
	resp, err := http.Get(ss.host + symbol + "?fields=" + strings.Join(fields, ","))

	if err != nil {
		return model.StockDataInfo{}, fmt.Errorf("Failed to get [%s] with error [%v]", symbol, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 299 {
		var response string
		fmt.Fscan(resp.Body, &response)
		return model.StockDataInfo{}, fmt.Errorf("Failed to get [%s], status code [%d], response [%v]", symbol, resp.StatusCode, response)
	}

	stockData := model.StockDataInfo{}

	err = json.NewDecoder(resp.Body).Decode(&stockData)

	if err != nil {
		return model.StockDataInfo{}, fmt.Errorf("Failed to deserialize data for [%s], error: [%v]", symbol, err)
	}

	return stockData, nil
}
