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

func GetStyle(t StyleType) lipgloss.Style {
	switch t {
	case StyleRise:
		return StyleChangeRise
	case StyleFall:
		return StyleChangeFall
	case StyleActive:
		return StyleNameActive
	case StyleInactive:
		return StyleNameInactive
	case StyleNeutral:
		return StyleChangeNeutral
	default:
		return lipgloss.NewStyle()
	}
}