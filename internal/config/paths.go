package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DirName               = ".juga"
	AliasesFileName       = "aliases.json"
	CacheFileName         = "cache.json"
	MasterTickersFileName = "master_tickers.csv"
	PortfoliosFileName    = "portfolios.json"
)

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, DirName), nil
}

func GetAliasesPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AliasesFileName), nil
}

func GetCachePath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, CacheFileName), nil
}

func GetPortfoliosPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, PortfoliosFileName), nil
}

func GetMasterTickersPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, MasterTickersFileName), nil
}

func EnsureConfigDir() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}
	return nil
}
