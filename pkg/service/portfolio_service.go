package service

import (
	"fmt"
	"strings"
)

type PortfolioRepository interface {
	Add(name string, items []string) error
	Remove(name string) error
	Get(name string) ([]string, bool)
	GetAll() map[string][]string
}

type PortfolioService struct {
	repo PortfolioRepository
}

func NewPortfolioService(repo PortfolioRepository) *PortfolioService {
	return &PortfolioService{
		repo: repo,
	}
}

func (s *PortfolioService) CreatePortfolio(name string, items []string) (*PortfolioOpResult, error) {
	if err := s.repo.Add(name, items); err != nil {
		return nil, fmt.Errorf("failed to save portfolio: %w", err)
	}
	return &PortfolioOpResult{
		Name:  name,
		Count: len(items),
	}, nil
}

func (s *PortfolioService) GetPortfolio(name string) ([]string, error) {
	items, ok := s.repo.Get(name)
	if !ok {
		return nil, ErrNotFound
	}
	return items, nil
}

func (s *PortfolioService) RemovePortfolio(name string) error {
	if _, ok := s.repo.Get(name); !ok {
		return ErrNotFound
	}
	return s.repo.Remove(name)
}

func (s *PortfolioService) ListPortfolios() map[string][]string {
	return s.repo.GetAll()
}

func (s *PortfolioService) ParseAndSave(name, content string) (*PortfolioOpResult, error) {
	var newItems []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		newItems = append(newItems, line)
	}

	return s.CreatePortfolio(name, newItems)
}
