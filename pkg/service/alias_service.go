package service

import (
	"fmt"

	"github.com/ericyhkim/juga/pkg/models"
	"github.com/ericyhkim/juga/pkg/resolver"
)

type AliasRepository interface {
	Add(nick, code string) error
	Remove(nick string) error
	Resolve(nick string) string
	GetAll() map[string]string
	SetAll(aliases map[string]string) error
}

type AliasService struct {
	repo     AliasRepository
	resolver *resolver.Resolver
}

func NewAliasService(repo AliasRepository, res *resolver.Resolver) *AliasService {
	return &AliasService{
		repo:     repo,
		resolver: res,
	}
}

func (s *AliasService) SetAlias(nick, target string) (*AliasOpResult, error) {
	if models.IsValidCode(nick) {
		return nil, ErrReservedName
	}

	res := s.resolver.Resolve(target)
	if res.Status != resolver.StatusSuccess {
		return nil, ErrInvalidTarget
	}

	if err := s.repo.Add(nick, res.Code); err != nil {
		return nil, fmt.Errorf("failed to save alias: %w", err)
	}

	return &AliasOpResult{
		Nickname: nick,
		Code:     res.Code,
		Name:     res.Name,
		Source:   res.Source,
	}, nil
}

func (s *AliasService) RemoveAlias(nick string) error {
	if s.repo.Resolve(nick) == "" {
		return ErrNotFound
	}
	return s.repo.Remove(nick)
}

func (s *AliasService) ListAliases() map[string]string {
	return s.repo.GetAll()
}

func (s *AliasService) BulkUpdate(aliases map[string]string) error {
	return s.repo.SetAll(aliases)
}
