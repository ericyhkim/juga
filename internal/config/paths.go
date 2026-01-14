package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	AliasesFileName       = "aliases.json"
	CacheFileName         = "cache.json"
	MasterTickersFileName = "master_tickers.csv"
	PortfoliosFileName    = "portfolios.json"

	// Environment variable overrides
	EnvConfigHome = "JUGA_CONFIG_HOME"
	EnvDataHome   = "JUGA_DATA_HOME"
	EnvCacheHome  = "JUGA_CACHE_HOME"

	// XDG standard variables
	XDGConfigHome = "XDG_CONFIG_HOME"
	XDGDataHome   = "XDG_DATA_HOME"
	XDGCacheHome  = "XDG_CACHE_HOME"
)

// getConfigHome returns the directory for configuration files (aliases, portfolios).
// Precedence: JUGA_CONFIG_HOME > XDG_CONFIG_HOME > Default (~/.config/juga or %APPDATA%/juga)
func getConfigHome() (string, error) {
	if path := os.Getenv(EnvConfigHome); path != "" {
		return path, nil
	}
	if path := os.Getenv(XDGConfigHome); path != "" {
		return filepath.Join(path, "juga"), nil
	}

	if runtime.GOOS == "windows" {
		configDir, err := os.UserConfigDir() // returns %APPDATA%
		if err != nil {
			return "", err
		}
		return filepath.Join(configDir, "juga"), nil
	}

	// Unix-like default: ~/.config/juga
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "juga"), nil
}

// getDataHome returns the directory for data files (master_tickers.csv).
// Precedence: JUGA_DATA_HOME > XDG_DATA_HOME > Default (~/.local/share/juga or %LOCALAPPDATA%/juga)
func getDataHome() (string, error) {
	if path := os.Getenv(EnvDataHome); path != "" {
		return path, nil
	}
	if path := os.Getenv(XDGDataHome); path != "" {
		return filepath.Join(path, "juga"), nil
	}

	if runtime.GOOS == "windows" {
		// On Windows, local data usually goes to %LOCALAPPDATA%
		cacheDir, err := os.UserCacheDir() // returns %LOCALAPPDATA%
		if err != nil {
			return "", err
		}
		return filepath.Join(cacheDir, "juga"), nil
	}

	// Unix-like default: ~/.local/share/juga
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "juga"), nil
}

// getCacheHome returns the directory for cache files (cache.json).
// Precedence: JUGA_CACHE_HOME > XDG_CACHE_HOME > Default (~/.cache/juga or %LOCALAPPDATA%/juga/cache)
func getCacheHome() (string, error) {
	if path := os.Getenv(EnvCacheHome); path != "" {
		return path, nil
	}
	if path := os.Getenv(XDGCacheHome); path != "" {
		return filepath.Join(path, "juga"), nil
	}

	if runtime.GOOS == "windows" {
		cacheDir, err := os.UserCacheDir() // returns %LOCALAPPDATA%
		if err != nil {
			return "", err
		}
		return filepath.Join(cacheDir, "juga", "cache"), nil
	}

	// Unix-like default: ~/.cache/juga
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "juga"), nil
}

func GetAliasesPath() (string, error) {
	dir, err := getConfigHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AliasesFileName), nil
}

func GetCachePath() (string, error) {
	dir, err := getCacheHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, CacheFileName), nil
}

func GetPortfoliosPath() (string, error) {
	dir, err := getConfigHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, PortfoliosFileName), nil
}

func GetMasterTickersPath() (string, error) {
	dir, err := getDataHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, MasterTickersFileName), nil
}

func EnsureAppDirs() error {
	configDir, err := getConfigHome()
	if err != nil {
		return fmt.Errorf("failed to resolve config dir: %w", err)
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	dataDir, err := getDataHome()
	if err != nil {
		return fmt.Errorf("failed to resolve data dir: %w", err)
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data dir: %w", err)
	}

	cacheDir, err := getCacheHome()
	if err != nil {
		return fmt.Errorf("failed to resolve cache dir: %w", err)
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache dir: %w", err)
	}

	return nil
}
