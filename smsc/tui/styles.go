package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorBorder = lipgloss.Color("238")
	colorTitle  = lipgloss.Color("12")
	colorDim    = lipgloss.Color("243")
	colorBright = lipgloss.Color("15")
	colorGreen  = lipgloss.Color("10")
	colorYellow = lipgloss.Color("11")
	colorCyan   = lipgloss.Color("14")

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(colorTitle)
	dimStyle    = lipgloss.NewStyle().Foreground(colorDim)
	brightStyle = lipgloss.NewStyle().Bold(true).Foreground(colorBright)
	greenStyle  = lipgloss.NewStyle().Foreground(colorGreen)
	yellowStyle = lipgloss.NewStyle().Foreground(colorYellow)
	cyanStyle   = lipgloss.NewStyle().Foreground(colorCyan)
)
