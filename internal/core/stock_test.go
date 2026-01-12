package core

import "testing"

func TestMapToStock(t *testing.T) {
	raw := NaverStockData{
		ItemCode:                    "005930",
		StockName:                   "Samsung",
		ClosePrice:                  "75,200",
		CompareToPreviousClosePrice: "-200",
		FluctuationsRatio:           "-0.26",
		HighPrice:                   "76,000",
		LowPrice:                    "75,000",
		AccumulatedTradingValue:     "1,234,567",
		CompareToPreviousPrice: CompareToPreviousPrice{
			Name: "FALLING",
		},
		MarketStatus: "OPEN",
	}

	stock := MapToStock(raw)

	if stock.Price != 75200 {
		t.Errorf("Expected Price 75200, got %f", stock.Price)
	}
	if stock.Change != -200 {
		t.Errorf("Expected Change -200, got %f", stock.Change)
	}
	if stock.ChangePercent != -0.26 {
		t.Errorf("Expected ChangePercent -0.26, got %f", stock.ChangePercent)
	}
	if stock.High != 76000 {
		t.Errorf("Expected High 76000, got %f", stock.High)
	}
	if stock.Low != 75000 {
		t.Errorf("Expected Low 75000, got %f", stock.Low)
	}
	if stock.TradingValue != 1234567 {
		t.Errorf("Expected TradingValue 1234567, got %f", stock.TradingValue)
	}
	if !stock.IsFalling {
		t.Error("Expected IsFalling to be true")
	}
	if stock.IsRising {
		t.Error("Expected IsRising to be false")
	}
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1,000", 1000},
		{"-500", -500},
		{"0", 0},
		{"invalid", 0},
		{"1,234,567", 1234567},
	}

	for _, test := range tests {
		if res := parsePrice(test.input); res != test.expected {
			t.Errorf("parsePrice(%q) = %f; want %f", test.input, res, test.expected)
		}
	}
}
