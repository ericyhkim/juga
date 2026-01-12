package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ericyhkim/juga/internal/core"
	"github.com/ericyhkim/juga/internal/sys"
	"github.com/ericyhkim/juga/internal/ui"
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

		repo := core.NewPortfolioRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading portfolios: %v\n", err)
			os.Exit(1)
		}

		if err := repo.Add(name, items); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving portfolio: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Portfolio '%s' saved with %d items.\n", name, len(items))
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

		repo := core.NewPortfolioRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading portfolios: %v\n", err)
			os.Exit(1)
		}

		items, ok := repo.Get(name)
		if !ok {
			fmt.Printf("Portfolio '%s' does not exist.\n", name)
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
			fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
			os.Exit(1)
		}

		var newItems []string
		lines := strings.Split(newContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			newItems = append(newItems, line)
		}

		if err := repo.Add(name, newItems); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving portfolio: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Portfolio '%s' updated. Now has %d items.\n", name, len(newItems))
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

		repo := core.NewPortfolioRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading portfolios: %v\n", err)
			os.Exit(1)
		}

		if _, ok := repo.Get(name); !ok {
			fmt.Printf("Portfolio '%s' not found.\n", name)
			return
		}

		if err := repo.Remove(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing portfolio: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Portfolio '%s' removed.\n", name)
	},
}

var portListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all portfolios",
	Run: func(cmd *cobra.Command, args []string) {
		repo := core.NewPortfolioRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading portfolios: %v\n", err)
			os.Exit(1)
		}

		all := repo.GetAll()
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
