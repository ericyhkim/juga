package models

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
