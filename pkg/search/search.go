package search

import (
	"sort"

	"github.com/ericyhkim/juga/pkg/models"
	"github.com/sahilm/fuzzy"
)

type TickerSource []models.Ticker

func (t TickerSource) String(i int) string {
	return t[i].Name
}

func (t TickerSource) Len() int {
	return len(t)
}

func FindTickers(tickers []models.Ticker, query string) []models.Ticker {
	if query == "" {
		return nil
	}

	source := TickerSource(tickers)
	matches := fuzzy.FindFrom(query, source)

	sort.Sort(matches)

	var results []models.Ticker
	for _, match := range matches {
		results = append(results, tickers[match.Index])
	}

	return results
}
