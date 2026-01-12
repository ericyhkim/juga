package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/ericyhkim/juga/internal/core"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"up"},
	Short:   "Refresh the master ticker database",
	Long: `Scrapes the latest stock list (KOSPI & KOSDAQ) from Naver Finance.
This is useful if new companies are listed or names change.
Process takes about 10-20 seconds.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating ticker database... (this may take a moment)")

		scraper := core.NewScraper()
		tickers, err := scraper.ScrapeAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scraping data: %v\n", err)
			os.Exit(1)
		}

		repo := core.NewTickerRepository()
		if err := repo.Load(); err != nil {
		}

		if err := repo.Save(tickers); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving tickers: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated %d tickers.\n", len(tickers))

		last, _ := repo.LastUpdated()
		fmt.Printf("Last updated: %s\n", last.Format(time.RFC822))
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
