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
	Long: `Remove the search cache (cache.json) and the local ticker database (master_tickers.csv).

This cleans files from your XDG Cache and Data directories.
User-defined aliases and portfolios in your Config directory are preserved.`,
	Run: func(cmd *cobra.Command, args []string) {
		cachePath, _ := config.GetCachePath()
		tickersPath, _ := config.GetMasterTickersPath()

		filesToRemove := []struct {
			path  string
			label string
		}{
			{cachePath, "search cache"},
			{tickersPath, "ticker database"},
		}

		removedCount := 0

		for _, item := range filesToRemove {
			if _, err := os.Stat(item.path); err == nil {
				if err := os.Remove(item.path); err == nil {
					fmt.Printf("Removed %s: %s\n", item.label, item.path)
					removedCount++
				} else {
					fmt.Fprintf(os.Stderr, "Failed to remove %s: %v\n", item.path, err)
				}
			}
		}

		if removedCount > 0 {
			fmt.Println("✨ Cleanup complete.")
		} else {
			fmt.Println("Nothing to clean.")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
