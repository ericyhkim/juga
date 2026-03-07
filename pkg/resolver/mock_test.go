package resolver

import "github.com/ericyhkim/juga/pkg/models"

type MockPortfolioProvider struct {
	Data map[string][]string
}

func (m *MockPortfolioProvider) Get(name string) ([]string, bool) {
	items, ok := m.Data[name]
	return items, ok
}

type MockAliasProvider struct {
	Data map[string]string
}

func (m *MockAliasProvider) Resolve(name string) string {
	return m.Data[name]
}

type MockCacheProvider struct {
	Data map[string]string
}

func (m *MockCacheProvider) Get(term string) (string, bool) {
	code, ok := m.Data[term]
	return code, ok
}

func (m *MockCacheProvider) Set(term, code string) {
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	m.Data[term] = code
}

type MockTickerProvider struct {
	Tickers []models.Ticker
}

func (m *MockTickerProvider) GetAll() []models.Ticker {
	return m.Tickers
}

func (m *MockTickerProvider) Count() int {
	return len(m.Tickers)
}

func (m *MockTickerProvider) Load() error {
	return nil
}
