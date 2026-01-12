package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorRed    = lipgloss.Color("#EF4444")
	ColorBlue   = lipgloss.Color("#3B82F6")
	ColorGrey   = lipgloss.Color("#626262")
	ColorYellow = lipgloss.Color("#EAB308")

	StyleNameActive   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	StyleNameInactive = lipgloss.NewStyle().Foreground(ColorGrey)

	StylePrice = lipgloss.NewStyle().Bold(true)

	StyleChangeRise    = lipgloss.NewStyle().Foreground(ColorRed)
	StyleChangeFall    = lipgloss.NewStyle().Foreground(ColorBlue)
	StyleChangeNeutral = lipgloss.NewStyle().Foreground(ColorGrey)
)
