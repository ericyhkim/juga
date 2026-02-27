package resolver

import (
	"fmt"

	"github.com/ericyhkim/juga/pkg/diag"
	"github.com/ericyhkim/juga/pkg/models"
	"github.com/ericyhkim/juga/pkg/search"
)

type Resolver struct {
	portfolios PortfolioProvider
	aliases    AliasProvider
	cache      CacheProvider
	tickers    TickerProvider
	logger     diag.Logger
}

func NewResolver(
	p PortfolioProvider,
	a AliasProvider,
	c CacheProvider,
	t TickerProvider,
	logger diag.Logger,
) *Resolver {
	return &Resolver{
		portfolios: p,
		aliases:    a,
		cache:      c,
		tickers:    t,
		logger:     logger,
	}
}

func (r *Resolver) ResolveAll(inputs []string) []ResolutionResult {
	var expandedInputs []string
	for _, input := range inputs {
		if items, ok := r.portfolios.Get(input); ok {
			expandedInputs = append(expandedInputs, items...)
		} else {
			expandedInputs = append(expandedInputs, input)
		}
	}

	results := make([]ResolutionResult, 0, len(expandedInputs))
	seen := make(map[string]bool)

	for _, input := range expandedInputs {
		res := r.Resolve(input)
		
		if res.Status == StatusSuccess {
			if seen[res.Code] {
				continue
			}
			seen[res.Code] = true
		}
		
		results = append(results, res)
	}

	return results
}

func (r *Resolver) Resolve(input string) ResolutionResult {
	if resolved := r.aliases.Resolve(input); resolved != "" {
		return ResolutionResult{
			Input:  input,
			Code:   resolved,
			Source: SourceAlias,
			Status: StatusSuccess,
		}
	}

	if models.IsValidCode(input) {
		return ResolutionResult{
			Input:  input,
			Code:   input,
			Source: SourceCode,
			Status: StatusSuccess,
		}
	}

	if cached, ok := r.cache.Get(input); ok {
		return ResolutionResult{
			Input:  input,
			Code:   cached,
			Source: SourceCache,
			Status: StatusSuccess,
		}
	}

	if r.tickers.Count() == 0 {
		if err := r.tickers.Load(); err != nil {
			r.logger.Error("Failed to load ticker list: %v", err)
		}
	}

	results := search.FindTickers(r.tickers.GetAll(), input)
	if len(results) > 0 {
		bestMatch := results[0]
		
		r.cache.Set(input, bestMatch.Code)
		
		return ResolutionResult{
			Input:  input,
			Code:   bestMatch.Code,
			Name:   bestMatch.Name,
			Source: SourceSearch,
			Status: StatusSuccess,
		}
	}

	return ResolutionResult{
		Input:  input,
		Status: StatusNotFound,
		Error:  fmt.Errorf("%w: %s", ErrNotFound, input),
	}
}
