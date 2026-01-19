package ui

// StockViewModel represents the fully formatted data ready for terminal rendering.
type StockViewModel struct {
	Name          string
	Price         string
	Change        string
	ChangePercent string
	TrendIcon     string
	IsRising      bool
	IsFalling     bool
	High          string
	Low           string
	TradingValue  string
	MarketStatus  string
}
