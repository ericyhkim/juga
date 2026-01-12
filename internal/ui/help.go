package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	StyleHelpHeader  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	StyleHelpCommand = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	StyleHelpComment = lipgloss.NewStyle().Foreground(ColorGrey)
	StyleHelpWarning = lipgloss.NewStyle().Foreground(ColorYellow)
	StyleHelpError   = lipgloss.NewStyle().Foreground(ColorRed)
	StyleHelpTip     = lipgloss.NewStyle().Foreground(ColorBlue)
)

// ContextualHelp defines the structure for a command's help/empty state.
type ContextualHelp struct {
	Usage        string
	Description  string // Optional: Brief explanation
	Examples     []string
	Tip          string // Optional: A helpful hint
	ErrorMessage string // Optional: displayed in yellow at the bottom
}

// RenderContextualHelp generates a styled help message string.
func RenderContextualHelp(h ContextualHelp) string {
	var sb strings.Builder

	if h.Description != "" {
		sb.WriteString(h.Description + "\n\n")
	}

	if h.Usage != "" {
		sb.WriteString(StyleHelpHeader.Render("Usage:") + "\n")
		sb.WriteString("  " + h.Usage + "\n\n")
	}

	if len(h.Examples) > 0 {
		sb.WriteString(StyleHelpHeader.Render("Examples:") + "\n")
		for _, ex := range h.Examples {
			// Split example command and comment if present (using #)
			parts := strings.SplitN(ex, "#", 2)
			cmd := parts[0]
			
			line := "  " + StyleHelpCommand.Render(cmd)
			if len(parts) > 1 {
				comment := strings.TrimSpace(parts[1])
				line += StyleHelpComment.Render("# " + comment)
			}
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
	}

	if h.Tip != "" {
		sb.WriteString(StyleHelpTip.Render("💡 Tip: "+h.Tip) + "\n\n")
	}

	if h.ErrorMessage != "" {
		// User input errors are styled as warnings (Yellow) to be less aggressive
		sb.WriteString(StyleHelpWarning.Render("Error: "+h.ErrorMessage) + "\n")
	}

	return strings.TrimRight(sb.String(), "\n")
}
