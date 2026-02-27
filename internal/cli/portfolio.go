package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/ericyhkim/juga/internal/sys"
	"github.com/ericyhkim/juga/internal/ui"
	"github.com/ericyhkim/juga/pkg/service"

	"github.com/spf13/cobra"
)

var portfolioCmd = &cobra.Command{
	Use:     "portfolio",
	Aliases: []string{"p", "port"},
	Short:   "Manage stock portfolios (groups)",
	Long: `Portfolios allow you to group multiple stocks under a single name.
When you run 'juga <portfolio_name>', it expands into all stocks in that group.
Items in a portfolio can be aliases, stock codes, or company names.`,
	Example: `  juga portfolio set my-tech 삼전 카카오 035420
  juga portfolio edit my-tech
  juga my-tech`,
}

var portSetCmd = &cobra.Command{
	Use:   "set <name> [stocks...]",
	Short: "Create or overwrite a portfolio",
	Long: `Creates a new portfolio or overwrites an existing one with the given stocks.
The stocks can be provided as names, codes, or existing aliases.`,
	Example: `  juga portfolio set favorites 삼전 카카오 NAVER
  juga portfolio set tech 005930 000660`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Usage: "juga portfolio set <name> <item1> [item2...]",
				Examples: []string{
					"juga portfolio set tech 005930 035420	# Create/Overwrite 'tech' with 2 stocks",
					"juga portfolio set faves 삼전 카카오	 # Use aliases or names",
				},
				ErrorMessage: "Please provide a portfolio name and at least one stock.",
			}))
			return
		}

		name := args[0]
		items := args[1:]

		deps := GetDeps(cmd)

		res, err := deps.PortfolioService.CreatePortfolio(name, items)
		if err != nil {
			deps.Logger.Error("Error saving portfolio: %v", err)
			return
		}

		fmt.Printf("Portfolio '%s' saved with %d items.\n", res.Name, res.Count)
	},
}

var portEditCmd = &cobra.Command{
	Use:     "edit <name>",
	Aliases: []string{"e"},
	Short:   "Edit a portfolio in your text editor",
	Long: `Opens the portfolio list in your default editor ($EDITOR or nano/vi).
Add or remove stocks line by line. Lines starting with # are ignored.`,
	Example: `  juga portfolio edit my-tech`,
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Usage: "juga portfolio edit <name>",
				Examples: []string{
					"juga portfolio edit my-tech",
				},
				Description:  "Action: Opens your default text editor to modify the portfolio list.",
				ErrorMessage: "Please specify the portfolio name to edit.",
			}))
			return
		}

		name := args[0]

		deps := GetDeps(cmd)

		items, err := deps.PortfolioService.GetPortfolio(name)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				fmt.Printf("Portfolio '%s' does not exist.\n", name)
			} else {
				deps.Logger.Error("Error: %v", err)
			}
			return
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("# Editing portfolio: %s\n", name))
		sb.WriteString("# Add one stock per line (name, code, or alias).\n")
		sb.WriteString("# Lines starting with # are ignored.\n")
		for _, item := range items {
			sb.WriteString(item + "\n")
		}

		newContent, err := sys.OpenEditor(sb.String())
		if err != nil {
			deps.Logger.Error("Error opening editor: %v", err)
			return
		}

		res, err := deps.PortfolioService.ParseAndSave(name, newContent)
		if err != nil {
			deps.Logger.Error("Error saving portfolio: %v", err)
			return
		}

		fmt.Printf("Portfolio '%s' updated. Now has %d items.\n", res.Name, res.Count)
	},
}

var portRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a portfolio",
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Usage: "juga portfolio remove <name>",
				Examples: []string{
					"juga portfolio remove my-tech",
				},
				Tip:          "Run 'juga portfolio list' to see your portfolios.",
				ErrorMessage: "Please specify the portfolio name to remove.",
			}))
			return
		}

		name := args[0]

		deps := GetDeps(cmd)

		if err := deps.PortfolioService.RemovePortfolio(name); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				fmt.Printf("Portfolio '%s' not found.\n", name)
			} else {
				deps.Logger.Error("Error removing portfolio: %v", err)
			}
			return
		}

		fmt.Printf("Portfolio '%s' removed.\n", name)
	},
}

var portListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all portfolios",
	Run: func(cmd *cobra.Command, args []string) {
		deps := GetDeps(cmd)
		all := deps.PortfolioService.ListPortfolios()
		if len(all) == 0 {
			fmt.Println("No portfolios defined.")
			return
		}

		keys := make([]string, 0, len(all))
		for k := range all {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var items []ui.ListItem
		for _, k := range keys {
			items = append(items, ui.ListItem{
				Key:   k,
				Value: strings.Join(all[k], ", "),
			})
		}

		fmt.Println(ui.RenderListTable(items))
	},
}

func init() {
	rootCmd.AddCommand(portfolioCmd)
	portfolioCmd.AddCommand(portSetCmd)
	portfolioCmd.AddCommand(portEditCmd)
	portfolioCmd.AddCommand(portRemoveCmd)
	portfolioCmd.AddCommand(portListCmd)
}
