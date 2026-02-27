package service

import (
	"errors"

	"github.com/ericyhkim/juga/pkg/models"
	"github.com/ericyhkim/juga/pkg/resolver"
)

var (
	ErrInvalidTarget = errors.New("could not resolve target to a valid stock")
	ErrReservedName  = errors.New("nickname cannot be a valid stock code")
	ErrNotFound      = errors.New("resource not found")
)

// AliasOpResult represents the outcome of an alias management operation.
type AliasOpResult struct {
	Nickname string
	Code     string
	Name     string
	Source   resolver.ResolutionSource
}

// PortfolioOpResult represents the outcome of a portfolio management operation.
type PortfolioOpResult struct {
	Name  string
	Count int
}

// StockFetchResult represents the outcome of a batch stock fetching operation.
type StockFetchResult struct {
	Stocks       []models.Stock
	IsTruncated  bool
	IgnoredCount int
}

// TickerUpdateResult represents the outcome of a ticker database update.
type TickerUpdateResult struct {
	Count int
}
