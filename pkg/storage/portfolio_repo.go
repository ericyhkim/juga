package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ericyhkim/juga/pkg/diag"
)

type PortfolioRepository struct {
	filePath   string
	portfolios map[string][]string
	logger     diag.Logger
}

func NewPortfolioRepository(filePath string, logger diag.Logger) *PortfolioRepository {
	return &PortfolioRepository{
		filePath:   filePath,
		portfolios: make(map[string][]string),
		logger:     logger,
	}
}

func (r *PortfolioRepository) Load() error {
	data, err := os.ReadFile(r.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read portfolios file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &r.portfolios); err != nil {
		return fmt.Errorf("failed to parse portfolios JSON: %w", err)
	}

	return nil
}

func (r *PortfolioRepository) Save() error {
	data, err := json.MarshalIndent(r.portfolios, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal portfolios: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write portfolios file: %w", err)
	}
	return nil
}

func (r *PortfolioRepository) Add(name string, items []string) error {
	r.portfolios[name] = items
	return r.Save()
}

func (r *PortfolioRepository) Remove(name string) error {
	delete(r.portfolios, name)
	return r.Save()
}

func (r *PortfolioRepository) Get(name string) ([]string, bool) {
	items, ok := r.portfolios[name]
	return items, ok
}

func (r *PortfolioRepository) GetAll() map[string][]string {
	copy := make(map[string][]string, len(r.portfolios))
	for k, v := range r.portfolios {
		itemsCopy := make([]string, len(v))
		for i, item := range v {
			itemsCopy[i] = item
		}
		copy[k] = itemsCopy
	}
	return copy
}
