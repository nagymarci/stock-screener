package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/nagymarci/stock-screener/config"
	"github.com/nagymarci/stock-screener/model"
)

var defaultDividendPerYear float64 = 4
var minOptInYieldWeight float64 = 0.4

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

//Calculate returns the dynamically computed data from the latest information
func Calculate(stockInfo *model.StockDataInfo) model.CalculatedStockInfo {
	var result model.CalculatedStockInfo

	now := time.Now()
	if model.Sp500DivYield.NextUpdate.Before(now) {
		model.Sp500DivYield.Mux.Lock()

		if model.Sp500DivYield.NextUpdate.Before(now) {
			yield, err := getSp500DivYield()
			if err != nil {
				model.Sp500DivYield.Mux.Unlock()
				log.Printf("Failed to update sp500 dividend yield: [%v]", err)
				log.Println("Using old sp500 dividend yield")

			} else {
				nextUpdateInterval, _ := time.ParseDuration("24h")
				model.Sp500DivYield.Yield = yield
				model.Sp500DivYield.NextUpdate = now.Add(nextUpdateInterval)
				model.Sp500DivYield.Mux.Unlock()
				log.Println("SP500 dividend yield updated")
			}
		}
	}

	optInYield, minOptInYield := calculateOptInYield(stockInfo.DividendYield5yr.Max, stockInfo.DividendYield5yr.Avg, model.Sp500DivYield.Yield)

	result.Ticker = stockInfo.Ticker
	result.AnnualDividend = stockInfo.Dividend * defaultDividendPerYear
	result.Price = stockInfo.Price
	result.DividendYield = result.AnnualDividend / result.Price * 100
	result.CurrentPe = result.Price / stockInfo.Eps
	result.OptInYield = optInYield
	result.DividendColor = calculateDividendColor(result.DividendYield, minOptInYield, stockInfo.DividendYield5yr.Avg)

	return result
}

func calculateDividendColor(dividendYield float64, minOptInYield float64, avg float64) string {
	if dividendYield > minOptInYield {
		return "green"
	}
	if dividendYield > avg {
		return "yellow"
	}

	return "blank"
}

//TODO use expected dividend raise for the calcualtion
func calculateOptInYield(max float64, avg float64, sp float64) (float64, float64) {
	minOptInYield := calculateMinOptInYield(max, avg)
	return math.Max(minOptInYield, sp), minOptInYieldWeight
}

func calculateMinOptInYield(max float64, avg float64) float64 {
	return (max-avg)*0.4 + avg
}

func getSp500DivYield() (float64, error) {
	resp, err := http.Get(config.Config.StockInfoProviderURL + "sp500/divYield")

	if err != nil {
		return 0, fmt.Errorf("Failed to get SP500 div yield: [%v]", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 299 {
		var response string
		fmt.Fscan(resp.Body, &response)
		return 0, fmt.Errorf("Failed to get SP500 div yield:, status code [%d], response [%v]", resp.StatusCode, response)
	}

	var response float64
	fmt.Fscan(resp.Body, &response)

	return response, nil
}
