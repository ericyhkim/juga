package resolver

import (
	"fmt"
	"strings"

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
		if strings.HasPrefix(input, models.PrefixPortfolio) {
			name := strings.TrimPrefix(input, models.PrefixPortfolio)
			if items, ok := r.portfolios.Get(name); ok {
				expandedInputs = append(expandedInputs, items...)
				continue
			}
		}

		// Fallback to legacy behavior: check if it's a portfolio without prefix
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
	if strings.HasPrefix(input, models.PrefixAlias) {
		nick := strings.TrimPrefix(input, models.PrefixAlias)
		if resolved := r.aliases.Resolve(nick); resolved != "" {
			return ResolutionResult{
				Input:  input,
				Code:   resolved,
				Source: SourceAlias,
				Status: StatusSuccess,
				Trace:  fmt.Sprintf("[%s] %s → %s", models.PrefixAlias, nick, resolved),
			}
		}
		return ResolutionResult{
			Input:  input,
			Status: StatusNotFound,
			Error:  fmt.Errorf("%w: alias '%s'", ErrNotFound, nick),
		}
	}

	if strings.HasPrefix(input, models.PrefixCode) {
		code := strings.TrimPrefix(input, models.PrefixCode)
		if models.IsValidCode(code) {
			return ResolutionResult{
				Input:  input,
				Code:   code,
				Source: SourceCode,
				Status: StatusSuccess,
				Trace:  fmt.Sprintf("[%s] %s", models.PrefixCode, code),
			}
		}
		return ResolutionResult{
			Input:  input,
			Status: StatusNotFound,
			Error:  fmt.Errorf("invalid stock code: %s", code),
		}
	}

	if strings.HasPrefix(input, models.PrefixSearch) {
		query := strings.TrimPrefix(input, models.PrefixSearch)
		return r.resolveSearch(input, query, true)
	}

	if resolved := r.aliases.Resolve(input); resolved != "" {
		return ResolutionResult{
			Input:  input,
			Code:   resolved,
			Source: SourceAlias,
			Status: StatusSuccess,
			Trace:  fmt.Sprintf("[Alias] %s → %s", input, resolved),
		}
	}

	if models.IsValidCode(input) {
		return ResolutionResult{
			Input:  input,
			Code:   input,
			Source: SourceCode,
			Status: StatusSuccess,
			Trace:  fmt.Sprintf("[Code] %s", input),
		}
	}

	if cached, ok := r.cache.Get(input); ok {
		return ResolutionResult{
			Input:  input,
			Code:   cached,
			Source: SourceCache,
			Status: StatusSuccess,
			Trace:  fmt.Sprintf("[Cache] %s → %s", input, cached),
		}
	}

	return r.resolveSearch(input, input, false)
}

func (r *Resolver) resolveSearch(input, query string, isExplicit bool) ResolutionResult {
	if r.tickers.Count() == 0 {
		if err := r.tickers.Load(); err != nil {
			r.logger.Error("Failed to load ticker list: %v", err)
		}
	}

	results := search.FindTickers(r.tickers.GetAll(), query)
	if len(results) > 0 {
		bestMatch := results[0]

		prefix := "Search"
		if isExplicit {
			prefix = models.PrefixSearch
		} else {
			r.cache.Set(input, bestMatch.Code)
		}

		isAmbiguous := len(results) > 1

		return ResolutionResult{
			Input:       input,
			Code:        bestMatch.Code,
			Name:        bestMatch.Name,
			Source:      SourceSearch,
			Status:      StatusSuccess,
			IsAmbiguous: isAmbiguous,
			Candidates:  results,
			Trace:       fmt.Sprintf("[%s] %s → %s (%s)", prefix, query, bestMatch.Code, bestMatch.Name),
		}
	}

	return ResolutionResult{
		Input:  input,
		Status: StatusNotFound,
		Error:  fmt.Errorf("%w: %s", ErrNotFound, query),
	}
}
