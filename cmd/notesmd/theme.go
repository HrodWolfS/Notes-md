package main

import (
	"github.com/charmbracelet/lipgloss"
)

// ASCII art title for home screen
const asciiTitle = `
███╗   ██╗ ██████╗ ████████╗███████╗███████╗
████╗  ██║██╔═══██╗╚══██╔══╝██╔════╝██╔════╝
██╔██╗ ██║██║   ██║   ██║   █████╗  ███████╗
██║╚██╗██║██║   ██║   ██║   ██╔══╝  ╚════██║
██║ ╚████║╚██████╔╝   ██║   ███████╗███████║
╚═╝  ╚═══╝ ╚═════╝    ╚═╝   ╚══════╝╚══════╝

                N O T E S . m d
`

// Color palette for theme cycling
var titlePalette = []string{"208", "196", "46", "51", "201"}

// Markdown theme for glamour rendering
var markdownTheme = "dark"

// Lipgloss styles
var (
	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")). // orange-ish
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Faint(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")). // light blue
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("208"))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("208")). // orange-ish, same as logo
			Padding(1, 4).
			Align(lipgloss.Center)

	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("208")).
			Foreground(lipgloss.Color("0")).
			Bold(true)

	searchLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("213")).
				Bold(true)

	searchQueryStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229"))
)
