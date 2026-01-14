package core

import (
	"strconv"
	"strings"
)

type Stock struct {
	Code          string
	Name          string
	Price         float64
	Change        float64
	ChangePercent float64
	High          float64
	Low           float64
	TradingValue  float64
	IsRising      bool
	IsFalling     bool
	MarketStatus  string
}

type Ticker struct {
	Code   string
	Name   string
	Market string
}

func IsValidCode(s string) bool {
	if len(s) != 6 {
		return false
	}
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return false
		}
	}
	return true
}

func MapToStock(apiData NaverStockData) Stock {
	price := parsePrice(apiData.ClosePrice)
	change := parsePrice(apiData.CompareToPreviousClosePrice)
	changePercent := parsePercent(apiData.FluctuationsRatio)
	high := parsePrice(apiData.HighPrice)
	low := parsePrice(apiData.LowPrice)
	tradingValue := parsePrice(apiData.AccumulatedTradingValue)

	return Stock{
		Code:          apiData.ItemCode,
		Name:          apiData.StockName,
		Price:         price,
		Change:        change,
		ChangePercent: changePercent,
		High:          high,
		Low:           low,
		TradingValue:  tradingValue,
		IsRising:      apiData.CompareToPreviousPrice.Name == "RISING",
		IsFalling:     apiData.CompareToPreviousPrice.Name == "FALLING",
		MarketStatus:  apiData.MarketStatus,
	}
}

func parsePrice(s string) float64 {
	clean := strings.ReplaceAll(s, ",", "")
	clean = strings.ReplaceAll(clean, "백만", "")
	clean = strings.ReplaceAll(clean, "천주", "")
	val, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0.0
	}
	return val
}

func parsePercent(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return val
}
