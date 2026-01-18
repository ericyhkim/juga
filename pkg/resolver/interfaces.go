package resolver

import "github.com/ericyhkim/juga/pkg/models"

type PortfolioProvider interface {
	Get(name string) ([]string, bool)
}

type AliasProvider interface {
	Resolve(name string) string
}

type CacheProvider interface {
	Get(term string) (string, bool)
	Set(term, code string)
}

type TickerProvider interface {
	GetAll() []models.Ticker
	Count() int
	Load() error
}
