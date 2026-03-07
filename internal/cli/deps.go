package cli

import (
	"context"
	"fmt"

	"github.com/ericyhkim/juga/pkg/config"
	"github.com/ericyhkim/juga/pkg/diag"
	"github.com/ericyhkim/juga/pkg/naver"
	"github.com/ericyhkim/juga/pkg/resolver"
	"github.com/ericyhkim/juga/pkg/service"
	"github.com/ericyhkim/juga/pkg/storage"
	"github.com/spf13/cobra"
)

type depsKey struct{}

// Dependencies holds all initialized services and repositories for the application.
type Dependencies struct {
	Logger     diag.Logger
	Aliases    *storage.AliasRepository
	Portfolios *storage.PortfolioRepository
	Cache      *storage.CacheRepository
	Tickers    *storage.TickerRepository
	Resolver   *resolver.Resolver
	Client     *naver.Client

	// Services
	AliasService     *service.AliasService
	PortfolioService *service.PortfolioService
	StockService     *service.StockService
}

// NewDependencies initializes all core application components.
func NewDependencies(logger diag.Logger) (*Dependencies, error) {
	aliasPath, err := config.GetAliasesPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get alias path: %w", err)
	}
	aliasRepo := storage.NewAliasRepository(aliasPath, logger)
	if err := aliasRepo.Load(); err != nil {
		logger.Error("Failed to load aliases: %v", err)
	}

	portPath, err := config.GetPortfoliosPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio path: %w", err)
	}
	portRepo := storage.NewPortfolioRepository(portPath, logger)
	if err := portRepo.Load(); err != nil {
		logger.Error("Failed to load portfolios: %v", err)
	}

	cachePath, err := config.GetCachePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache path: %w", err)
	}
	cacheRepo := storage.NewCacheRepository(cachePath, config.DefaultCacheSize, logger)
	if err := cacheRepo.Load(); err != nil {
		logger.Error("Failed to load cache: %v", err)
	}

	tickerPath, err := config.GetMasterTickersPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker path: %w", err)
	}
	tickerRepo := storage.NewTickerRepository(tickerPath, logger)

	resSvc := resolver.NewResolver(portRepo, aliasRepo, cacheRepo, tickerRepo, logger)
	client := naver.NewClient(logger, naver.WithTimeout(config.DefaultClientTimeout))

	aliasService := service.NewAliasService(aliasRepo, resSvc)
	portfolioService := service.NewPortfolioService(portRepo)
	stockService := service.NewStockService(
		tickerRepo,
		client,
		logger,
		config.DefaultScraperTimeout,
		config.DefaultMaxStocks,
	)

	return &Dependencies{
		Logger:           logger,
		Aliases:          aliasRepo,
		Portfolios:       portRepo,
		Cache:            cacheRepo,
		Tickers:          tickerRepo,
		Resolver:         resSvc,
		Client:           client,
		AliasService:     aliasService,
		PortfolioService: portfolioService,
		StockService:     stockService,
	}, nil
}

// GetDeps retrieves the Dependencies from the command context.
func GetDeps(cmd *cobra.Command) *Dependencies {
	if deps, ok := cmd.Context().Value(depsKey{}).(*Dependencies); ok {
		return deps
	}
	return nil
}

// SetDeps attaches Dependencies to the context.
func SetDeps(ctx context.Context, deps *Dependencies) context.Context {
	return context.WithValue(ctx, depsKey{}, deps)
}
