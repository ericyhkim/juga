package cli

import (
	"fmt"
	"os"

	"github.com/ericyhkim/juga/internal/core"
	"github.com/ericyhkim/juga/internal/ui"

	"github.com/spf13/cobra"
)

var marketCmd = &cobra.Command{
	Use:     "market",
	Aliases: []string{"m"},
	Short:   "Show detailed market index information",
	Long:    `Display detailed statistics for KOSPI and KOSDAQ, including high/low prices and trading volume/value.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := core.NewClient()
		indices, err := client.FetchIndices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching market data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(ui.RenderMarketDetails(indices))
	},
}
