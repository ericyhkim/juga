package ui

// StyleType defines semantic styles for stock data presentation.
type StyleType int

const (
	StyleNeutral StyleType = iota
	StyleRise
	StyleFall
	StyleActive
	StyleInactive
)

// StockViewModel represents the fully formatted data ready for terminal rendering.
type StockViewModel struct {
	Name        string
	NameStyle   StyleType
	Price       string
	ChangeInfo  string
	ChangeStyle StyleType

	High         string
	Low          string
	TradingValue string
}