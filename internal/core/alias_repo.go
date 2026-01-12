package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ericyhkim/juga/internal/config"
)

type AliasRepository struct {
	aliases map[string]string
}

func NewAliasRepository() *AliasRepository {
	return &AliasRepository{
		aliases: make(map[string]string),
	}
}

func (r *AliasRepository) Load() error {
	path, err := config.GetAliasesPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read aliases file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &r.aliases); err != nil {
		return fmt.Errorf("failed to parse aliases JSON: %w", err)
	}

	return nil
}

func (r *AliasRepository) Save() error {
	path, err := config.GetAliasesPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(r.aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write aliases file: %w", err)
	}
	return nil
}

func (r *AliasRepository) Add(nick, code string) error {
	r.aliases[nick] = code
	return r.Save()
}

func (r *AliasRepository) Remove(nick string) error {
	delete(r.aliases, nick)
	return r.Save()
}

func (r *AliasRepository) Resolve(nick string) string {
	return r.aliases[nick]
}

func (r *AliasRepository) GetAll() map[string]string {
	copy := make(map[string]string, len(r.aliases))
	for k, v := range r.aliases {
		copy[k] = v
	}
	return copy
}

func (r *AliasRepository) SetAll(aliases map[string]string) error {
	r.aliases = aliases
	return r.Save()
}
