package core

// NaverResponse represents the top-level structure of the Naver Finance polling API.
type NaverResponse struct {
	Datas []NaverStockData `json:"datas"`
	Time  string           `json:"time"`
}

// NaverStockData represents the individual stock data returned by Naver.
// We use string types for numeric values because Naver returns them formatted with commas.
type NaverStockData struct {
	ItemCode                    string                 `json:"itemCode"`
	StockName                   string                 `json:"stockName"`
	ClosePrice                  string                 `json:"closePrice"`
	CompareToPreviousClosePrice string                 `json:"compareToPreviousClosePrice"`
	FluctuationsRatio           string                 `json:"fluctuationsRatio"`
	HighPrice                   string                 `json:"highPrice"`
	LowPrice                    string                 `json:"lowPrice"`
	AccumulatedTradingValue     string                 `json:"accumulatedTradingValue"`
	MarketStatus                string                 `json:"marketStatus"`
	CompareToPreviousPrice      CompareToPreviousPrice `json:"compareToPreviousPrice"`
}

// CompareToPreviousPrice details the price movement (FALLING, RISING, STABLE).
type CompareToPreviousPrice struct {
	Code string `json:"code"`
	Text string `json:"text"`
	Name string `json:"name"`
}
