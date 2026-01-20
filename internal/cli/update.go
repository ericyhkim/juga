package cli

import (
	"fmt"
	"os"

	"github.com/ericyhkim/juga/pkg/config"
	"github.com/ericyhkim/juga/pkg/naver"
	"github.com/ericyhkim/juga/pkg/storage"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"up"},
	Short:   "Update the local ticker database",
	Long: `Scrapes Naver Finance to update the master list of stock codes (KOSPI/KOSDAQ/ETF).
Process takes about 10-20 seconds.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating ticker database... (this may take a moment)")

		scraper := naver.NewScraper(config.DefaultScraperTimeout)
		tickers, err := scraper.ScrapeAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scraping data: %v\n", err)
			os.Exit(1)
		}

		tickerPath, _ := config.GetMasterTickersPath()
		repo := storage.NewTickerRepository(tickerPath)
		if err := repo.Load(); err != nil {
		}

		if err := repo.Save(tickers); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving database: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully updated %d tickers.\n", len(tickers))
	},
}
