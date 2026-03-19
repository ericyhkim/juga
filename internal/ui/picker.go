package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
)

// RunPicker displays an interactive selection list using the huh library.
func RunPicker(title string, items []ListItem) (string, error) {
	if !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		if len(items) > 0 {
			return items[0].Value, nil
		}
		return "", fmt.Errorf("no items to pick from")
	}

	options := make([]huh.Option[string], 0, len(items))
	for _, item := range items {
		options = append(options, huh.NewOption(
			fmt.Sprintf("%s (%s)", item.Key, item.Value),
			item.Value,
		))
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(options...).
				Inline(true).
				Value(&selected),
		),
	)

	theme := huh.ThemeCharm()
	theme.Focused.Title = StyleNameActive.Copy()
	theme.Focused.SelectedOption = StylePrice.Copy()
	theme.Focused.UnselectedOption = StyleNameInactive.Copy()
	theme.Blurred.Title = StyleNameInactive.Copy()
	
	form.WithTheme(theme)

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}
