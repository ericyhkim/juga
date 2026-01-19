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
	return StockViewModel{
		Name:          s.Name,
		Price:         formatNumber(s.Price),
		Change:        formatNumber(s.Change),
		ChangePercent: fmt.Sprintf("%.2f%%", s.ChangePercent),
		TrendIcon:     getDirectionSymbol(s),
		IsRising:      s.IsRising,
		IsFalling:     s.IsFalling,
		High:          formatNumber(s.High),
		Low:           formatNumber(s.Low),
		TradingValue:  formatLargeValue(s.TradingValue),
		MarketStatus:  s.MarketStatus,
	}
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
