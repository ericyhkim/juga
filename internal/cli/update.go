package cli

import (
	"fmt"

	"github.com/ericyhkim/juga/pkg/config"
	"github.com/ericyhkim/juga/pkg/naver"
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

		deps := GetDeps(cmd)

		scraper := naver.NewScraper(config.DefaultScraperTimeout, deps.Logger)
		tickers, err := scraper.ScrapeAll()
		if err != nil {
			deps.Logger.Error("Error scraping data: %v", err)
			return
		}

		if err := deps.Tickers.Save(tickers); err != nil {
			deps.Logger.Error("Error saving database: %v", err)
			return
		}

		fmt.Printf("✅ Successfully updated %d tickers.\n", len(tickers))
	},
}
