package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderIndices(indices []StockViewModel) string {
	if len(indices) == 0 {
		return ""
	}

	var parts []string
	for _, idx := range indices {
		name := StyleNameActive.Render(idx.Name)
		price := idx.Price

		changeText := fmt.Sprintf("%s %s (%s)",
			idx.TrendIcon,
			idx.Change,
			idx.ChangePercent,
		)

		var stats string
		if idx.IsRising {
			stats = StyleChangeRise.Render(fmt.Sprintf("%s %s", price, changeText))
		} else if idx.IsFalling {
			stats = StyleChangeFall.Render(fmt.Sprintf("%s %s", price, changeText))
		} else {
			stats = StyleChangeNeutral.Render(fmt.Sprintf("%s %s", price, changeText))
		}

		parts = append(parts, fmt.Sprintf("%s %s", name, stats))
	}

	return strings.Join(parts, "   |   ")
}

func RenderMarketDetails(indices []StockViewModel) string {
	if len(indices) == 0 {
		return ""
	}

	var blocks []string
	for _, idx := range indices {
		name := StyleNameActive.Render(idx.Name)
		price := idx.Price
		changeText := fmt.Sprintf("%s %s (%s)",
			idx.TrendIcon,
			idx.Change,
			idx.ChangePercent,
		)

		var statsLine string
		if idx.IsRising {
			statsLine = StyleChangeRise.Render(fmt.Sprintf("%s %s", price, changeText))
		} else if idx.IsFalling {
			statsLine = StyleChangeFall.Render(fmt.Sprintf("%s %s", price, changeText))
		} else {
			statsLine = StyleChangeNeutral.Render(fmt.Sprintf("%s %s", price, changeText))
		}

		header := fmt.Sprintf("% -8s %s", name, statsLine)

		details := StyleNameInactive.Render(fmt.Sprintf("         High: %s   Low: %s   Val: %s",
			idx.High,
			idx.Low,
			idx.TradingValue,
		))

		blocks = append(blocks, header+"\n"+details)
	}

	return strings.Join(blocks, "\n\n")
}

func RenderStockTable(stocks []StockViewModel) string {
	if len(stocks) == 0 {
		return ""
	}

	maxNameWidth := 0
	maxPriceWidth := 0

	for _, s := range stocks {
		nameWidth := lipgloss.Width(s.Name)
		if nameWidth > maxNameWidth {
			maxNameWidth = nameWidth
		}
		priceWidth := lipgloss.Width(s.Price)
		if priceWidth > maxPriceWidth {
			maxPriceWidth = priceWidth
		}
	}

	var rows []string
	for _, s := range stocks {
		var name string
		if isMarketOpen(s.MarketStatus) {
			name = StyleNameActive.Copy().Width(maxNameWidth).Render(s.Name)
		} else {
			name = StyleNameInactive.Copy().Width(maxNameWidth).Render(s.Name)
		}

		price := StylePrice.Copy().
			Width(maxPriceWidth).
			Align(lipgloss.Right).
			Render(s.Price)

		changeText := fmt.Sprintf("%s %s (%s)",
			s.TrendIcon,
			s.Change,
			s.ChangePercent,
		)

		var change string
		if s.IsRising {
			change = StyleChangeRise.Render(changeText)
		} else if s.IsFalling {
			change = StyleChangeFall.Render(changeText)
		} else {
			change = StyleChangeNeutral.Render(changeText)
		}

		rows = append(rows, fmt.Sprintf("%s  %s  %s", name, price, change))
	}

	return strings.Join(rows, "\n")
}

// isMarketOpen checks if the market status indicates active trading.
// This is a simplified check; Naver returns various strings like "OPEN", "CLOSE", "DELAY".
func isMarketOpen(status string) bool {
	status = strings.ToUpper(status)
	return status == "OPEN" || status == "장중" // "장중" is Korean for "During Market"
}

// ListItem represents a single row in a key-value list (e.g., Alias -> Code)
type ListItem struct {
	Key   string
	Value string
}

// RenderListTable renders a list of key-value pairs with perfect alignment.
// It mimics the style of the main stock table: Bold Key, 2-space gutter, Grey Value.
func RenderListTable(items []ListItem) string {
	if len(items) == 0 {
		return StyleNameInactive.Render("No items found.")
	}

	maxKeyWidth := 0
	for _, item := range items {
		w := lipgloss.Width(item.Key)
		if w > maxKeyWidth {
			maxKeyWidth = w
		}
	}

	var rows []string
	for _, item := range items {
		key := StyleNameActive.Copy().Width(maxKeyWidth).Render(item.Key)
		val := StyleNameInactive.Render(item.Value)
		rows = append(rows, fmt.Sprintf("%s  %s", key, val))
	}

	return strings.Join(rows, "\n")
}
