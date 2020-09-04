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
var maxOptInPeWeight float64 = 0.5
var lowerDividendYieldGuardScore float64 = 1.5

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

	optInPe := calculateOptInPe(stockInfo.PeRatio5yr.Min, stockInfo.PeRatio5yr.Avg)

	result.Ticker = stockInfo.Ticker
	result.AnnualDividend = stockInfo.Dividend * defaultDividendPerYear
	result.Price = stockInfo.Price
	result.DividendYield = result.AnnualDividend / result.Price * 100
	result.CurrentPe = result.Price / stockInfo.Eps
	result.OptInYield = optInYield
	result.DividendColor = calculateDividendColor(result.DividendYield, minOptInYield, stockInfo.DividendYield5yr.Avg)
	result.OptInPe = optInPe
	result.PeColor = calculatePeColor(result.CurrentPe, optInPe, stockInfo.PeRatio5yr.Avg)

	optInPrice := calculateOptInPrice(optInYield, result.AnnualDividend, model.Sp500DivYield.Yield)

	result.OptInPrice = optInPrice
	result.PriceColor = calculatePriceColor(result.Price, optInPrice)

	return result
}

func calculatePriceColor(price float64, optInPrice float64) string {
	if price < optInPrice {
		return "green"
	}
	if price < optInPrice*1.05 {
		return "yellow"
	}

	return "red"
}

func calculateOptInPrice(optInYield float64, annualDividend float64, sp float64) float64 {
	spOptInPrice := annualDividend / (sp * lowerDividendYieldGuardScore) * 100
	minOptInPrice := annualDividend / optInYield * 100

	return math.Min(spOptInPrice, minOptInPrice)
}

func calculatePeColor(currentPe float64, optInPe float64, avg float64) string {
	if currentPe < optInPe {
		return "green"
	}

	if currentPe < avg {
		return "yellow"
	}

	return "blank"
}

func calculateOptInPe(min float64, avg float64) float64 {
	return (avg-min)*maxOptInPeWeight + min
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
	return math.Max(minOptInYield, sp*lowerDividendYieldGuardScore), minOptInYield
}

func calculateMinOptInYield(max float64, avg float64) float64 {
	return (max-avg)*minOptInYieldWeight + avg
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

//GetAllRecommendedStock returns all the recommended stocks based on the requirements
func GetAllRecommendedStock(stocks []model.StockDataInfo, numReqs int) []model.CalculatedStockInfo {
	var result []model.CalculatedStockInfo

	for _, stockInfo := range stocks {
		calculated := Calculate(&stockInfo)

		reqsFulfilled := calculateReqsFulfilled(&calculated)

		if reqsFulfilled >= numReqs {
			result = append(result, calculated)
		}
	}

	return result
}

func calculateReqsFulfilled(stock *model.CalculatedStockInfo) int {
	result := 0
	if stock.DividendColor == "green" {
		result++
	}

	if stock.PriceColor == "green" {
		result++
	}

	if stock.PeColor == "green" {
		result++
	}

	return result
}
