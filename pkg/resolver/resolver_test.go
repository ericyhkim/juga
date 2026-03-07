package resolver

import (
	"testing"

	"github.com/ericyhkim/juga/pkg/diag"
	"github.com/ericyhkim/juga/pkg/models"
)

func setupTestResolver(t *testing.T) *Resolver {
	aMock := &MockAliasProvider{
		Data: map[string]string{
			"sam": "005930",
		},
	}
	pMock := &MockPortfolioProvider{
		Data: map[string][]string{
			"tech": {"sam", "kakao"},
		},
	}
	cMock := &MockCacheProvider{
		Data: make(map[string]string),
	}
	tMock := &MockTickerProvider{
		Tickers: []models.Ticker{},
	}

	return NewResolver(pMock, aMock, cMock, tMock, diag.NewNopLogger())
}

func TestResolve_Alias(t *testing.T) {
	r := setupTestResolver(t)

	res := r.Resolve("sam")

	if res.Status != StatusSuccess {
		t.Errorf("Expected success, got %s", res.Status)
	}
	if res.Source != SourceAlias {
		t.Errorf("Expected SourceAlias, got %s", res.Source)
	}
	if res.Code != "005930" {
		t.Errorf("Expected code 005930, got %s", res.Code)
	}
}

func TestResolve_DirectCode(t *testing.T) {
	r := setupTestResolver(t)

	res := r.Resolve("000660")

	if res.Status != StatusSuccess {
		t.Errorf("Expected success, got %s", res.Status)
	}
	if res.Source != SourceCode {
		t.Errorf("Expected SourceCode, got %s", res.Source)
	}
	if res.Code != "000660" {
		t.Errorf("Expected code 000660, got %s", res.Code)
	}
}

func TestResolveAll_Portfolio(t *testing.T) {
	r := setupTestResolver(t)

	results := r.ResolveAll([]string{"tech"})

	if len(results) != 2 {
		t.Errorf("Expected 2 results from portfolio expansion, got %d", len(results))
	}

	if results[0].Input != "sam" || results[0].Code != "005930" {
		t.Errorf("Expected first item to be resolved alias 'sam', got %v", results[0])
	}
}