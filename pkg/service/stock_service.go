package service

import (
	"fmt"
	"time"

	"github.com/ericyhkim/juga/pkg/diag"
	"github.com/ericyhkim/juga/pkg/models"
	"github.com/ericyhkim/juga/pkg/naver"
	"github.com/ericyhkim/juga/pkg/resolver"
	"github.com/ericyhkim/juga/pkg/search"
)

type TickerRepository interface {
	Save(tickers []models.Ticker) error
	Load() error
	Count() int
}

type NaverClient interface {
	FetchStocks(codes []string) ([]models.Stock, error)
	FetchIndices() ([]models.Stock, error)
}

type StockService struct {
	tickerRepo     TickerRepository
	client         NaverClient
	logger         diag.Logger
	scraperTimeout time.Duration
	maxStocks      int
}

func NewStockService(
	tickerRepo TickerRepository,
	client NaverClient,
	logger diag.Logger,
	scraperTimeout time.Duration,
	maxStocks int,
) *StockService {
	return &StockService{
		tickerRepo:     tickerRepo,
		client:         client,
		logger:         logger,
		scraperTimeout: scraperTimeout,
		maxStocks:      maxStocks,
	}
}

func (s *StockService) UpdateTickerDatabase() (*TickerUpdateResult, error) {
	scraper := naver.NewScraper(s.scraperTimeout, s.logger)
	tickers, err := scraper.ScrapeAll()
	if err != nil {
		return nil, fmt.Errorf("failed to scrape tickers: %w", err)
	}

	if err := s.tickerRepo.Save(tickers); err != nil {
		return nil, fmt.Errorf("failed to save tickers: %w", err)
	}

	return &TickerUpdateResult{
		Count: len(tickers),
	}, nil
}

func (s *StockService) FetchStocks(results []resolver.ResolutionResult) (*StockFetchResult, error) {
	isTruncated := false
	ignoredCount := 0
	
	toFetch := results
	if len(results) > s.maxStocks {
		isTruncated = true
		ignoredCount = len(results) - s.maxStocks
		toFetch = results[:s.maxStocks]
	}

	var targetCodes []string
	for _, res := range toFetch {
		if res.Status == resolver.StatusSuccess {
			targetCodes = append(targetCodes, res.Code)
		}
	}

	if len(targetCodes) == 0 {
		return &StockFetchResult{
			Stocks:       []models.Stock{},
			IsTruncated:  isTruncated,
			IgnoredCount: ignoredCount,
		}, nil
	}

	stocks, err := s.client.FetchStocks(targetCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock data: %w", err)
	}

	return &StockFetchResult{
		Stocks:       stocks,
		IsTruncated:  isTruncated,
		IgnoredCount: ignoredCount,
	}, nil
}

func (s *StockService) FetchIndices() ([]models.Stock, error) {
	return s.client.FetchIndices()
}

func (s *StockService) SearchTickers(query string) ([]models.Ticker, error) {
	if s.tickerRepo.Count() == 0 {
		if err := s.tickerRepo.Load(); err != nil {
			return nil, fmt.Errorf("failed to load ticker database: %w", err)
		}
	}

	return search.FindTickers(s.tickerRepo.(interface{ GetAll() []models.Ticker }).GetAll(), query), nil
}
