package ui

import (
	"fmt"
	"strings"

	"github.com/ericyhkim/juga/pkg/models"
)

type Presenter struct{}

func NewPresenter() *Presenter {
	return &Presenter{}
}

func (p *Presenter) PrepareList(stocks []models.Stock) []StockViewModel {
	var vms []StockViewModel
	for _, s := range stocks {
		vms = append(vms, p.PrepareStock(s))
	}
	return vms
}

func (p *Presenter) PrepareStock(s models.Stock) StockViewModel {
	nameStyle := StyleInactive
	if isMarketOpen(s.MarketStatus) {
		nameStyle = StyleActive
	}

	changeText := fmt.Sprintf("%s %s (%s)",
		getDirectionSymbol(s),
		formatNumber(s.Change),
		fmt.Sprintf("%.2f%%", s.ChangePercent),
	)

	changeStyle := StyleNeutral
	if s.IsRising {
		changeStyle = StyleRise
	} else if s.IsFalling {
		changeStyle = StyleFall
	}

	return StockViewModel{
		Name:         s.Name,
		NameStyle:    nameStyle,
		Price:        formatNumber(s.Price),
		ChangeInfo:   changeText,
		ChangeStyle:  changeStyle,
		High:         formatNumber(s.High),
		Low:          formatNumber(s.Low),
		TradingValue: formatLargeValue(s.TradingValue),
	}
}

// isMarketOpen checks if the market status indicates active trading.
// This is a simplified check; Naver returns various strings like "OPEN", "CLOSE", "DELAY".
func isMarketOpen(status string) bool {
	status = strings.ToUpper(status)
	return status == "OPEN" || status == "장중" // "장중" is Korean for "During Market"
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

func formatLargeValue(v float64) string {
	if v >= 1000000 {
		return fmt.Sprintf("%.1fT", v/1000000)
	}
	if v >= 1000 {
		return fmt.Sprintf("%.1fB", v/1000)
	}
	return fmt.Sprintf("%.1fM", v)
}

func getDirectionSymbol(s models.Stock) string {
	if s.IsRising {
		return "▲"
	}
	if s.IsFalling {
		return "▼"
	}
	return "-"
}