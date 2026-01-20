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
		Name:        "Samsung Electronics",
		NameStyle:   StyleActive,
		Price:       "75,000",
		ChangeInfo:  "▲ 1,500 (2.04%)",
		ChangeStyle: StyleRise,
		High:        "76,000",
		Low:         "74,000",
		TradingValue: "1.5T",
	}

	result := p.PrepareStock(input)

	if result.Name != expected.Name {
		t.Errorf("Expected Name %q, got %q", expected.Name, result.Name)
	}
	if result.NameStyle != expected.NameStyle {
		t.Errorf("Expected NameStyle %v, got %v", expected.NameStyle, result.NameStyle)
	}
	if result.Price != expected.Price {
		t.Errorf("Expected Price %q, got %q", expected.Price, result.Price)
	}
	if result.ChangeInfo != expected.ChangeInfo {
		t.Errorf("Expected ChangeInfo %q, got %q", expected.ChangeInfo, result.ChangeInfo)
	}
	if result.ChangeStyle != expected.ChangeStyle {
		t.Errorf("Expected ChangeStyle %v, got %v", expected.ChangeStyle, result.ChangeStyle)
	}
	if result.TradingValue != expected.TradingValue {
		t.Errorf("Expected TradingValue %q, got %q", expected.TradingValue, result.TradingValue)
	}
}