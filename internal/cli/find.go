package cli

import (
	"fmt"
	"os"

	"github.com/ericyhkim/juga/internal/core"
	"github.com/ericyhkim/juga/internal/ui"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:     "find [query]",
	Aliases: []string{"f", "search"},
	Short:   "Fuzzy search for stocks by name",
	Long: `Search the master ticker list using fuzzy matching.
Useful for finding the 6-digit code for a company when you only know its name.

Example:
  juga find 삼전   -> Matches '삼성전자' (005930), '삼성전기' (009150), etc.
  juga find 카카오  -> Matches '카카오' (035720), '카카오뱅크' (323410), etc.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Usage: "juga find <search_term>",
				Examples: []string{
					"juga find 삼전        # Matches '삼성전자', '삼성전기', etc.",
					"juga find NAVER       # Matches 'NAVER'",
				},
				ErrorMessage: "Please provide a search term.",
			}))
			return
		}

		query := args[0]

		repo := core.NewTickerRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading ticker database: %v\n", err)
			os.Exit(1)
		}

		results := core.FindTickers(repo.GetAll(), query)

		if len(results) == 0 {
			fmt.Printf("No matches found for '%s'.\n", query)
			return
		}

		limit := 10
		displayCount := limit
		if len(results) < limit {
			displayCount = len(results)
		}

		var items []ui.ListItem
		for _, t := range results[:displayCount] {
			items = append(items, ui.ListItem{
				Key:   t.Name,
				Value: fmt.Sprintf("%s [%s]", t.Code, t.Market),
			})
		}

		fmt.Println(ui.RenderListTable(items))

		if len(results) > limit {
			fmt.Printf("...and %d more.\n", len(results)-limit)
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
