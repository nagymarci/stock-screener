package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/model"
)

//Get returns the requested stock from the provider
func Get(symbol string) (model.StockDataInfo, error) {
	return GetWithFields(symbol, []string{})
}

//GetWithFields returns the stock from the provider with the requested fields filled
func GetWithFields(symbol string, fields []string) (model.StockDataInfo, error) {
	resp, err := http.Get(config.Config.StockInfoProviderURL + symbol + "?fields=" + strings.Join(fields, ","))

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
