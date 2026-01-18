package cli

import (
	"fmt"
	"os"

	"github.com/ericyhkim/juga/internal/ui"
	"github.com/ericyhkim/juga/pkg/config"
	"github.com/ericyhkim/juga/pkg/models"
	"github.com/ericyhkim/juga/pkg/naver"
	"github.com/ericyhkim/juga/pkg/search"
	"github.com/ericyhkim/juga/pkg/storage"

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

Example:
  juga 삼성전자 kakao
  juga 005930`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Description: "juga (주가) 📈 - Minimalist CLI for Korean stock prices.",
				Usage:       "juga <stock_name_or_code> [more_stocks...]",
				Examples: []string{
					"juga 삼성전자       # Check price by name",
					"juga 005930         # Check price by code",
					"juga find 카카오    # Find a stock code",
					"juga market         # Check KOSPI/KOSDAQ",
				},
				Tip: "Run 'juga help' for the full list of commands.",
			}))
			return
		}

		portRepo := storage.NewPortfolioRepository()
		if err := portRepo.Load(); err != nil {
		}

		var expandedArgs []string
		for _, arg := range args {
			if items, ok := portRepo.Get(arg); ok {
				expandedArgs = append(expandedArgs, items...)
			} else {
				expandedArgs = append(expandedArgs, arg)
			}
		}

		isTruncated := false
		ignoredCount := 0
		if len(expandedArgs) > config.DefaultMaxStocks {
			isTruncated = true
			ignoredCount = len(expandedArgs) - config.DefaultMaxStocks
			expandedArgs = expandedArgs[:config.DefaultMaxStocks]
		}

		aliasRepo := storage.NewAliasRepository()
		if err := aliasRepo.Load(); err != nil {
		}

		cacheRepo := storage.NewCacheRepository(config.DefaultCacheSize)
		if err := cacheRepo.Load(); err != nil {
		}

		tickerRepo := storage.NewTickerRepository()
		tickerLoaded := false

		var targetCodes []string
		seen := make(map[string]bool)

		for _, arg := range expandedArgs {
			var code string

			if resolved := aliasRepo.Resolve(arg); resolved != "" {
				code = resolved
			} else if models.IsValidCode(arg) {
				code = arg
			} else if cached, ok := cacheRepo.Get(arg); ok {
				code = cached
			} else {
				if !tickerLoaded {
					if err := tickerRepo.Load(); err != nil {
						fmt.Fprintf(os.Stderr, "⚠️  Could not load ticker database: %v\n", err)
						continue
					}
					tickerLoaded = true
				}

				results := search.FindTickers(tickerRepo.GetAll(), arg)
				if len(results) > 0 {
					code = results[0].Code
					cacheRepo.Set(arg, code)
				} else {
					fmt.Printf("⚠️  Could not find stock for '%s'\n", arg)
					continue
				}
			}

			if code != "" && !seen[code] {
				targetCodes = append(targetCodes, code)
				seen[code] = true
			}
		}

		if cacheErr := cacheRepo.Save(); cacheErr != nil {
		}

		if len(targetCodes) == 0 {
			return
		}

		client := naver.NewClient(naver.WithTimeout(config.DefaultClientTimeout))
		stockResult, stockErr := client.FetchStocks(targetCodes)

		if stockErr != nil {
			fmt.Fprintf(os.Stderr, "Error fetching data: %v\n", stockErr)
			return
		}

		fmt.Println(ui.RenderStockTable(stockResult))

		if isTruncated {
			fmt.Fprintf(os.Stderr, "\n⚠️  Display limited to %d stocks. %d items were ignored.\n", config.DefaultMaxStocks, ignoredCount)
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
