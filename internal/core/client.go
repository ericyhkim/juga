package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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

func (c *Client) FetchStocks(codes []string) ([]Stock, error) {
	if len(codes) == 0 {
		return []Stock{}, nil
	}

	joinedCodes := strings.Join(codes, ",")
	url := naverPollingURL + joinedCodes

	return c.fetchData(url)
}

func (c *Client) FetchIndices() ([]Stock, error) {
	return c.fetchData(naverIndexURL)
}

func (c *Client) fetchData(url string) ([]Stock, error) {
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

	var stocks []Stock
	for _, data := range naverResp.Datas {
		stocks = append(stocks, MapToStock(data))
	}

	return stocks, nil
}
