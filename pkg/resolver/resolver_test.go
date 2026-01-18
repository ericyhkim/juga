package resolver

import (
	"os"
	"testing"

	"github.com/ericyhkim/juga/pkg/storage"
)

func setupTestResolver(t *testing.T) (*Resolver, func()) {
	aliasFile, _ := os.CreateTemp("", "alias_*.json")
	portFile, _ := os.CreateTemp("", "port_*.json")
	cacheFile, _ := os.CreateTemp("", "cache_*.json")
	
	aRepo := storage.NewAliasRepository()
	pRepo := storage.NewPortfolioRepository()
	cRepo := storage.NewCacheRepository(10)
	tRepo := storage.NewTickerRepository()

	aRepo.Add("sam", "005930")
	pRepo.Add("tech", []string{"sam", "kakao"})
	
	r := NewResolver(pRepo, aRepo, cRepo, tRepo)

	cleanup := func() {
		os.Remove(aliasFile.Name())
		os.Remove(portFile.Name())
		os.Remove(cacheFile.Name())
	}

	return r, cleanup
}

func TestResolve_Alias(t *testing.T) {
	r, cleanup := setupTestResolver(t)
	defer cleanup()

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
	r, cleanup := setupTestResolver(t)
	defer cleanup()

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
	r, cleanup := setupTestResolver(t)
	defer cleanup()

	results := r.ResolveAll([]string{"tech"})

	if len(results) != 2 {
		t.Errorf("Expected 2 results from portfolio expansion, got %d", len(results))
	}
	
	if results[0].Input != "sam" || results[0].Code != "005930" {
		t.Errorf("Expected first item to be resolved alias 'sam', got %v", results[0])
	}
}
