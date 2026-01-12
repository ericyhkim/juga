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

var aliasCmd = &cobra.Command{
	Use:     "alias",
	Aliases: []string{"a"},
	Short:   "Manage stock aliases",
	Long: `Aliases allow you to link a short nickname to a stock code.
Once set, you can use the nickname in place of the code or full name.`,
	Example: `  juga alias set sam 삼성전자
  juga alias list
  juga alias remove sam`,
}

var aliasSetCmd = &cobra.Command{
	Use:   "set <nickname> <target>",
	Short: "Create or update an alias",
	Long: `Links a nickname to a stock code.
The target can be a 6-digit code or a stock name (which will be auto-resolved).`,
	Example: `  juga alias set sam 삼성전자     -> Finds '삼성전자' code
  juga alias set my-stock 005930 -> Directly links to code`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Usage: "juga alias set <nickname> <stock_target>",
				Examples: []string{
					"juga alias set sam 삼성전자     # Link 'sam' to '삼성전자'",
					"juga alias set my-chip 000660  # Link 'my-chip' to SK하이닉스",
				},
				ErrorMessage: "Missing arguments. Expected <nickname> and <target>.",
			}))
			return
		}

		nick := args[0]
		target := args[1]

		aliasRepo := core.NewAliasRepository()
		if err := aliasRepo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading aliases: %v\n", err)
			os.Exit(1)
		}

		code := target
		resolutionSource := "code"

		if !core.IsValidCode(target) {
			if resolved := aliasRepo.Resolve(target); resolved != "" {
				code = resolved
				resolutionSource = fmt.Sprintf("existing alias '%s'", target)
			} else {
				tickerRepo := core.NewTickerRepository()
				if err := tickerRepo.Load(); err != nil {
					fmt.Fprintf(os.Stderr, "Error loading tickers: %v\n", err)
					os.Exit(1)
				}

				results := core.FindTickers(tickerRepo.GetAll(), target)
				if len(results) == 0 {
					fmt.Printf("Could not resolve '%s' to any stock.\n", target)
					return
				}

				best := results[0]
				code = best.Code
				resolutionSource = fmt.Sprintf("stock name '%s' (%s)", best.Name, best.Code)
			}
		}

		if err := aliasRepo.Add(nick, code); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving alias: %v\n", err)
			os.Exit(1)
		}

		if resolutionSource == "code" {
			fmt.Printf("Alias set: %s -> %s (direct code)\n", nick, code)
		} else {
			fmt.Printf("Alias set: %s -> %s (resolved via %s)\n", nick, code, resolutionSource)
		}
	},
}

var aliasRemoveCmd = &cobra.Command{
	Use:     "remove <nickname>",
	Aliases: []string{"rm"},
	Short:   "Remove an existing alias",
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(ui.RenderContextualHelp(ui.ContextualHelp{
				Usage: "juga alias remove <nickname>",
				Examples: []string{
					"juga alias remove sam",
				},
				Tip:          "Run 'juga alias list' to see all your current nicknames.",
				ErrorMessage: "Please specify the alias nickname to remove.",
			}))
			return
		}

		nick := args[0]

		repo := core.NewAliasRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading aliases: %v\n", err)
			os.Exit(1)
		}

		if repo.Resolve(nick) == "" {
			fmt.Printf("Alias '%s' not found.\n", nick)
			return
		}

		if err := repo.Remove(nick); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing alias: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alias '%s' removed.\n", nick)
	},
}

var aliasListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all registered aliases",
	Run: func(cmd *cobra.Command, args []string) {
		repo := core.NewAliasRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading aliases: %v\n", err)
			os.Exit(1)
		}

		all := repo.GetAll()
		if len(all) == 0 {
			fmt.Println("No aliases defined.")
			return
		}

		// Sort keys for consistent output
		keys := make([]string, 0, len(all))
		for k := range all {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var items []ui.ListItem
		for _, k := range keys {
			items = append(items, ui.ListItem{
				Key:   k,
				Value: all[k],
			})
		}

		fmt.Println(ui.RenderListTable(items))
	},
}

var aliasEditCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit all aliases in your text editor",
	Long: `Opens all your aliases in the default editor ($EDITOR or nano/vi).
Modify the mappings in 'nickname: code' format. Lines starting with # are ignored.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := core.NewAliasRepository()
		if err := repo.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading aliases: %v\n", err)
			os.Exit(1)
		}

		all := repo.GetAll()
		var keys []string
		for k := range all {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var sb strings.Builder
		sb.WriteString("# Editing all aliases\n")
		sb.WriteString("# Format: nickname: code\n")
		sb.WriteString("# Example: sam: 005930\n")
		sb.WriteString("# Lines starting with # are ignored.\n")
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, all[k]))
		}

		newContent, err := sys.OpenEditor(sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
			os.Exit(1)
		}

		newAliases := make(map[string]string)
		lines := strings.Split(newContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue // Skip malformed lines
			}

			nick := strings.TrimSpace(parts[0])
			code := strings.TrimSpace(parts[1])
			if nick != "" && code != "" {
				newAliases[nick] = code
			}
		}

		if err := repo.SetAll(newAliases); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving aliases: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated %d aliases.\n", len(newAliases))
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(aliasSetCmd)
	aliasCmd.AddCommand(aliasRemoveCmd)
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasEditCmd)
}
