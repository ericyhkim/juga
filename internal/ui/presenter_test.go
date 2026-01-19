package ui

import (
	"testing"

	"github.com/ericyhkim/juga/pkg/models"
)

func TestPresenter_PrepareStock(t *testing.T) {
	p := NewPresenter()

	input := models.Stock{
		Name:          "Samsung Electronics",
		Price:         75000,
		Change:        1500,
		ChangePercent: 2.04,
		IsRising:      true,
		IsFalling:     false,
		High:          76000,
		Low:           74000,
		TradingValue:  1500000,
		MarketStatus:  "OPEN",
	}

	expected := StockViewModel{
		Name:          "Samsung Electronics",
		Price:         "75,000",
		Change:        "1,500",
		ChangePercent: "2.04%",
		TrendIcon:     "▲",
		IsRising:      true,
		IsFalling:     false,
		High:          "76,000",
		Low:           "74,000",
		TradingValue:  "1.5T",
		MarketStatus:  "OPEN",
	}

	result := p.PrepareStock(input)

	if result.Name != expected.Name {
		t.Errorf("Expected Name %q, got %q", expected.Name, result.Name)
	}
	if result.Price != expected.Price {
		t.Errorf("Expected Price %q, got %q", expected.Price, result.Price)
	}
	if result.Change != expected.Change {
		t.Errorf("Expected Change %q, got %q", expected.Change, result.Change)
	}
	if result.ChangePercent != expected.ChangePercent {
		t.Errorf("Expected ChangePercent %q, got %q", expected.ChangePercent, result.ChangePercent)
	}
	if result.TrendIcon != expected.TrendIcon {
		t.Errorf("Expected TrendIcon %q, got %q", expected.TrendIcon, result.TrendIcon)
	}
	if result.TradingValue != expected.TradingValue {
		t.Errorf("Expected TradingValue %q, got %q", expected.TradingValue, result.TradingValue)
	}
}
