package cli

import (
	"fmt"

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

		res, err := deps.StockService.UpdateTickerDatabase()
		if err != nil {
			deps.Logger.Error("Error updating ticker database: %v", err)
			return
		}

		fmt.Printf("✅ Successfully updated %d tickers.\n", res.Count)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
