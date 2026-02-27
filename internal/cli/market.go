package cli

import (
	"fmt"

	"github.com/ericyhkim/juga/internal/ui"
	"github.com/spf13/cobra"
)

var marketCmd = &cobra.Command{
	Use:     "market",
	Aliases: []string{"m"},
	Short:   "Show detailed market index information",
	Long:    `Display detailed statistics for KOSPI and KOSDAQ, including high/low prices and trading volume/value.`,
	Run: func(cmd *cobra.Command, args []string) {
		deps := GetDeps(cmd)

		indices, err := deps.Client.FetchIndices()
		if err != nil {
			deps.Logger.Error("Error fetching market data: %v", err)
			return
		}

		presenter := ui.NewPresenter()
		marketVMs := presenter.PrepareList(indices)

		fmt.Println(ui.RenderMarketDetails(marketVMs))
	},
}
