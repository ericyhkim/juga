package naver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericyhkim/juga/pkg/models"
)

const (
	naverPollingURL = "https://polling.finance.naver.com/api/realtime/domestic/stock/"
	naverIndexURL   = "https://polling.finance.naver.com/api/realtime/domestic/index/KOSPI,KOSDAQ"
	defaultTimeout  = 2 * time.Second
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *Client) FetchStocks(codes []string) ([]models.Stock, error) {
	if len(codes) == 0 {
		return []models.Stock{}, nil
	}

	joinedCodes := strings.Join(codes, ",")
	url := naverPollingURL + joinedCodes

	return c.fetchData(url)
}

func (c *Client) FetchIndices() ([]models.Stock, error) {
	return c.fetchData(naverIndexURL)
}

func (c *Client) fetchData(url string) ([]models.Stock, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("naver api returned status: %d", resp.StatusCode)
	}

	var naverResp NaverResponse
	if err := json.NewDecoder(resp.Body).Decode(&naverResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var stocks []models.Stock
	for _, data := range naverResp.Datas {
		stocks = append(stocks, MapToStock(data))
	}

	return stocks, nil
}

func MapToStock(apiData NaverStockData) models.Stock {
	price := parsePrice(apiData.ClosePrice)
	change := parsePrice(apiData.CompareToPreviousClosePrice)
	changePercent := parsePercent(apiData.FluctuationsRatio)
	high := parsePrice(apiData.HighPrice)
	low := parsePrice(apiData.LowPrice)
	tradingValue := parsePrice(apiData.AccumulatedTradingValue)

	return models.Stock{
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
