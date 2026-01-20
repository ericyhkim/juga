package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type AliasRepository struct {
	filePath string
	aliases  map[string]string
}

func NewAliasRepository(filePath string) *AliasRepository {
	return &AliasRepository{
		filePath: filePath,
		aliases:  make(map[string]string),
	}
}

func (r *AliasRepository) Load() error {
	data, err := os.ReadFile(r.filePath)
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
	data, err := json.MarshalIndent(r.aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
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
