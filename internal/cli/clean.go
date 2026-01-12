package cli

import (
	"fmt"
	"os"

	"github.com/ericyhkim/juga/internal/config"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clear the search cache and master ticker list",
	Long:  `Remove the search cache (cache.json) and the local ticker database (master_tickers.csv). User-defined aliases and portfolios are preserved.`,
	Run: func(cmd *cobra.Command, args []string) {
		cachePath, _ := config.GetCachePath()
		tickersPath, _ := config.GetMasterTickersPath()

		filesToRemove := []string{cachePath, tickersPath}
		removedCount := 0

		for _, path := range filesToRemove {
			if _, err := os.Stat(path); err == nil {
				if err := os.Remove(path); err == nil {
					removedCount++
				}
			}
		}

		if removedCount > 0 {
			fmt.Println("✨ Cache and ticker database cleaned.")
		} else {
			fmt.Println("Nothing to clean.")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
