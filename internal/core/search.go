package core

import (
	"sort"

	"github.com/sahilm/fuzzy"
)

type TickerSource []Ticker

func (t TickerSource) String(i int) string {
	return t[i].Name
}

func (t TickerSource) Len() int {
	return len(t)
}

func FindTickers(tickers []Ticker, query string) []Ticker {
	if query == "" {
		return nil
	}

	source := TickerSource(tickers)
	matches := fuzzy.FindFrom(query, source)

	sort.Sort(matches)

	var results []Ticker
	for _, match := range matches {
		results = append(results, tickers[match.Index])
	}

	return results
}
