package ui

import (
	"strings"
	"testing"

	"github.com/ericyhkim/juga/internal/core"
)

func TestRenderStockTable(t *testing.T) {
	stocks := []core.Stock{
		{
			Name:          "Samsung",
			Price:         75000,
			Change:        1000,
			ChangePercent: 1.35,
			IsRising:      true,
			MarketStatus:  "OPEN",
		},
		{
			Name:          "SK Hynix",
			Price:         130000,
			Change:        -2500,
			ChangePercent: -1.89,
			IsFalling:     true,
			MarketStatus:  "OPEN",
		},
	}

	result := RenderStockTable(stocks)
	lines := strings.Split(result, "\n")

	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(lines))
	}

	// Verify both stocks are present
	for _, s := range stocks {
		if !strings.Contains(result, s.Name) {
			t.Errorf("RenderStockTable() missing stock name: %s", s.Name)
		}
	}

	// Check for formatting/alignment (basic check)
	if !strings.Contains(result, " 75,000") {
		t.Errorf("RenderStockTable() missing formatted price 75,000")
	}
	if !strings.Contains(result, "130,000") {
		t.Errorf("RenderStockTable() missing formatted price 130,000")
	}
}

func TestRenderMarketDetails(t *testing.T) {
	indices := []core.Stock{
		{
			Name:          "KOSPI",
			Price:         2586.32,
			Change:        33.95,
			ChangePercent: 1.33,
			High:          2590.00,
			Low:           2550.00,
			TradingValue:  23478838, // 23.4T
			IsRising:      true,
		},
	}

	result := RenderMarketDetails(indices)
	if !strings.Contains(result, "KOSPI") {
		t.Error("Missing KOSPI")
	}
	if !strings.Contains(result, "High: 2,590") {
		t.Error("Missing High price")
	}
	if !strings.Contains(result, "Val: 23.5T") { // 23.47... -> 23.5T
		t.Errorf("Expected 23.5T, got formatted value in: %s", result)
	}
}

func TestFormatLargeValue(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{23478838, "23.5T"},
		{9033037, "9.0T"},
		{123456, "123.5B"},
		{1234, "1.2B"},
		{123, "123.0M"},
	}

	for _, tt := range tests {
		if got := formatLargeValue(tt.input); got != tt.expected {
			t.Errorf("formatLargeValue(%f) = %s, want %s", tt.input, got, tt.expected)
		}
	}
}

func TestRenderIndices(t *testing.T) {
	indices := []core.Stock{
		{
			Name:          "KOSPI",
			Price:         2586.32,
			Change:        33.95,
			ChangePercent: 1.33,
			IsRising:      true,
		},
		{
			Name:          "KOSDAQ",
			Price:         847.92,
			Change:        -3.86,
			ChangePercent: -0.45,
			IsFalling:     true,
		},
	}

	result := RenderIndices(indices)

	if !strings.Contains(result, "KOSPI") || !strings.Contains(result, "KOSDAQ") {
		t.Errorf("RenderIndices() missing index names")
	}
	if !strings.Contains(result, "2,586.32") || !strings.Contains(result, "847.92") {
		t.Errorf("RenderIndices() missing index prices")
	}
	if !strings.Contains(result, " | ") {
		t.Errorf("RenderIndices() missing separator")
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{-1234, "-1,234"},
		{2586.32, "2,586.32"},
		{947.92, "947.92"},
		{0.41, "0.41"},
		{-3.86, "-3.86"},
	}

	for _, tt := range tests {
		if got := formatNumber(tt.input); got != tt.expected {
			t.Errorf("formatNumber(%f) = %s, want %s", tt.input, got, tt.expected)
		}
	}
}
