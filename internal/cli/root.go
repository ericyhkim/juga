package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ericyhkim/juga/internal/ui"
	"github.com/ericyhkim/juga/pkg/config"
	"github.com/ericyhkim/juga/pkg/diag"
	"github.com/ericyhkim/juga/pkg/models"
	"github.com/ericyhkim/juga/pkg/resolver"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "juga [names...]",
	Short:   "A minimalist CLI for real-time Korean stock prices",
	Version: Version,
	Long: `juga (주가) 📈

A simple terminal tool that bypasses complex APIs and official codes,
letting you check KOSPI/KOSDAQ market data instantly using aliases and fuzzy search.

Deterministic Prefixes:
  @<name>   - Force Portfolio resolution
  :<name>   - Force Alias resolution
  #<code>   - Force Stock Code resolution
  /<query>  - Force Fuzzy Search (bypasses cache/aliases)

Example:
  juga 삼성전자 :sam #005930 @my-tech
  juga /카카오`,
	Args: cobra.ArbitraryArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger := diag.NewStdLogger()
		deps, err := NewDependencies(logger)
		if err != nil {
			return err
		}
		cmd.SetContext(SetDeps(cmd.Context(), deps))
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Description: "juga (주가) 📈 - Minimalist CLI for Korean stock prices.",
				Usage:       "juga <stock_name_or_code> [more_stocks...]",
				Examples: []string{
					"juga 삼성전자       # Check price by name",
					"juga :sam           # Force alias 'sam'",
					"juga #005930         # Force code '005930'",
					"juga /카카오        # Search with picker",
					"juga @tech          # Force portfolio 'tech'",
				},
				Tip: "Run 'juga help' for the full list of commands.",
			}))
			return
		}

		deps := GetDeps(cmd)

		results := deps.Resolver.ResolveAll(args)

		finalResults := make([]resolver.ResolutionResult, 0, len(results))
		for _, res := range results {
			if res.IsAmbiguous {
				listItems := make([]ui.ListItem, 0, len(res.Candidates))
				for _, c := range res.Candidates {
					listItems = append(listItems, ui.ListItem{
						Key:   c.Name,
						Value: c.Code,
					})
				}

				title := fmt.Sprintf("Multiple matches for '%s'. Select one:", res.Input)
				selectedCode, err := ui.RunPicker(title, listItems)
				if err == nil {
					for _, c := range res.Candidates {
						if c.Code == selectedCode {
							res.Code = c.Code
							res.Name = c.Name
							res.IsAmbiguous = false
							
							prefix := "Search"
							if strings.HasPrefix(res.Input, models.PrefixSearch) {
								prefix = models.PrefixSearch
							}
							res.Trace = fmt.Sprintf("[%s] %s → %s (%s)", prefix, strings.TrimPrefix(res.Input, models.PrefixSearch), res.Code, res.Name)
							break
						}
					}
				}
			}
			finalResults = append(finalResults, res)
		}

		for _, res := range finalResults {
			if res.Status == resolver.StatusSuccess && res.Trace != "" {
				fmt.Println(ui.StyleNameInactive.Render(res.Trace))
			} else if res.Status == resolver.StatusNotFound {
				fmt.Printf("⚠️  Could not find stock for '%s'\n", res.Input)
			}
		}

		if cacheErr := deps.Cache.Save(); cacheErr != nil {
			deps.Logger.Error("Failed to save cache: %v", cacheErr)
		}

		fetchRes, err := deps.StockService.FetchStocks(finalResults)
		if err != nil {
			deps.Logger.Error("Error fetching data: %v", err)
			return
		}

		if len(fetchRes.Stocks) == 0 {
			return
		}

		presenter := ui.NewPresenter()
		stockVMs := presenter.PrepareList(fetchRes.Stocks)

		fmt.Println("") 
		fmt.Println(ui.RenderStockTable(stockVMs))

		if fetchRes.IsTruncated {
			fmt.Fprintf(os.Stderr, "\n⚠️  Display limited to some stocks. %d items were ignored.\n", fetchRes.IgnoredCount)
		}
	},
}

func Execute() {
	if err := config.EnsureAppDirs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(marketCmd)
}
