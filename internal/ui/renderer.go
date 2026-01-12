package ui

import (
	"fmt"
	"strings"

	"github.com/ericyhkim/juga/internal/core"

	"github.com/charmbracelet/lipgloss"
)

func RenderIndices(indices []core.Stock) string {
	if len(indices) == 0 {
		return ""
	}

	var parts []string
	for _, idx := range indices {
		name := StyleNameActive.Render(idx.Name)
		price := formatNumber(idx.Price)

		changeText := fmt.Sprintf("%s %s (%.2f%%)",
			getDirectionSymbol(idx),
			formatNumber(idx.Change),
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

func RenderMarketDetails(indices []core.Stock) string {
	if len(indices) == 0 {
		return ""
	}

	var blocks []string
	for _, idx := range indices {
		name := StyleNameActive.Render(idx.Name)
		price := formatNumber(idx.Price)
		changeText := fmt.Sprintf("%s %s (%.2f%%)",
			getDirectionSymbol(idx),
			formatNumber(idx.Change),
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

		header := fmt.Sprintf("%-8s %s", name, statsLine)

		details := StyleNameInactive.Render(fmt.Sprintf("         High: %s   Low: %s   Val: %s",
			formatNumber(idx.High),
			formatNumber(idx.Low),
			formatLargeValue(idx.TradingValue),
		))

		blocks = append(blocks, header+"\n"+details)
	}

	return strings.Join(blocks, "\n\n")
}

func formatLargeValue(v float64) string {
	if v >= 1000000 {
		return fmt.Sprintf("%.1fT", v/1000000)
	}
	if v >= 1000 {
		return fmt.Sprintf("%.1fB", v/1000)
	}
	return fmt.Sprintf("%.1fM", v)
}

func RenderStockTable(stocks []core.Stock) string {
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
		priceWidth := lipgloss.Width(formatNumber(s.Price))
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
			Render(formatNumber(s.Price))

		changeText := fmt.Sprintf("%s %s (%.2f%%)",
			getDirectionSymbol(s),
			formatNumber(s.Change),
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

func formatNumber(n float64) string {
	var s string
	if n == float64(int64(n)) {
		s = fmt.Sprintf("%.0f", n)
	} else {
		s = fmt.Sprintf("%.2f", n)
	}

	parts := strings.Split(s, ".")
	intPart := parts[0]

	var res strings.Builder
	l := len(intPart)
	for i, r := range intPart {
		if i > 0 && (l-i)%3 == 0 && intPart[i-1] != '-' {
			res.WriteRune(',')
		}
		res.WriteRune(r)
	}

	if len(parts) > 1 {
		res.WriteRune('.')
		res.WriteString(parts[1])
	}

	return res.String()
}

func getDirectionSymbol(s core.Stock) string {
	if s.IsRising {
		return "▲"
	}
	if s.IsFalling {
		return "▼"
	}
	return "-"
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
