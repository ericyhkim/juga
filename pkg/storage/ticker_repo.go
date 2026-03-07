package storage

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/ericyhkim/juga/pkg/diag"
	"github.com/ericyhkim/juga/pkg/models"
)

//go:embed master_tickers.csv
var defaultTickersCSV []byte

type TickerRepository struct {
	filePath string
	tickers  []models.Ticker
	logger   diag.Logger
}

func NewTickerRepository(filePath string, logger diag.Logger) *TickerRepository {
	return &TickerRepository{
		filePath: filePath,
		tickers:  []models.Ticker{},
		logger:   logger,
	}
}

func (r *TickerRepository) LastUpdated() (time.Time, error) {
	info, err := os.Stat(r.filePath)
	if os.IsNotExist(err) {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, err
	}

	return info.ModTime(), nil
}

func (r *TickerRepository) IsFresh(d time.Duration) bool {
	last, err := r.LastUpdated()
	if err != nil || last.IsZero() {
		return false
	}
	return time.Since(last) < d
}

func (r *TickerRepository) Load() error {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		if err := os.WriteFile(r.filePath, defaultTickersCSV, 0644); err != nil {
			return fmt.Errorf("failed to create default tickers file: %w", err)
		}
	}

	f, err := os.Open(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to open tickers file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse tickers CSV: %w", err)
	}

	var loaded []models.Ticker
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		loaded = append(loaded, models.Ticker{
			Code:   record[0],
			Name:   record[1],
			Market: record[2],
		})
	}

	embeddedCount := bytes.Count(defaultTickersCSV, []byte{'\n'})
	if len(loaded) < embeddedCount {
		reader := csv.NewReader(bytes.NewReader(defaultTickersCSV))
		embeddedRecords, err := reader.ReadAll()
		if err == nil {
			var embeddedLoaded []models.Ticker
			for _, record := range embeddedRecords {
				if len(record) < 3 {
					continue
				}
				embeddedLoaded = append(embeddedLoaded, models.Ticker{
					Code:   record[0],
					Name:   record[1],
					Market: record[2],
				})
			}

			if len(embeddedLoaded) > len(loaded) {
				loaded = embeddedLoaded
				if saveErr := r.Save(loaded); saveErr != nil {
					r.logger.Warn("Warning: failed to update local tickers: %v\n", saveErr)
				}
			}
		}
	}

	r.tickers = loaded
	return nil
}

func (r *TickerRepository) Save(tickers []models.Ticker) error {
	f, err := os.Create(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to create tickers file: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, t := range tickers {
		if err := writer.Write([]string{t.Code, t.Name, t.Market}); err != nil {
			return err
		}
	}

	r.tickers = tickers
	return nil
}

func (r *TickerRepository) GetAll() []models.Ticker {
	return r.tickers
}

func (r *TickerRepository) Count() int {
	return len(r.tickers)
}
