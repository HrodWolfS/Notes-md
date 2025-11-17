package main

import (
	"github.com/charmbracelet/lipgloss"
)

// viewHome renders the home/welcome screen
func (m model) viewHome() string {
	// Dynamic accent color
	ts := titleStyle
	accent := lipgloss.Color("208")
	if len(titlePalette) > 0 {
		colorCode := titlePalette[m.themeIndex]
		accent = lipgloss.Color(colorCode)
		ts = ts.Foreground(accent)
	}

	// Style for NOTES.md in blue
	notesStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("81")). // blue
		Bold(true)

	// Build content with proper alignment
	title := ts.Render(asciiTitle)
	notesText := notesStyle.Render("NOTES.md")
	helpText := helpStyle.Render("Press any key to continue")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		notesText,
		helpText,
	)

	card := cardStyle.
		BorderForeground(accent).
		Render(content)

	// Fallback if size not yet known
	w := m.width
	h := m.height
	if w <= 0 {
		w = 120
	}
	if h <= 0 {
		h = 40
	}

	// Base view
	baseView := lipgloss.Place(
		w,
		h,
		lipgloss.Center,
		lipgloss.Center,
		card,
	)

	// Handle modals (overlay on top of base view)
	var modalView string
	if m.showHelpModal {
		modalView = m.helpModal.View()
	}

	if modalView != "" {
		overlayedModal := lipgloss.Place(
			w,
			h,
			lipgloss.Center,
			lipgloss.Center,
			modalView,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
		)
		return overlayedModal
	}

	return baseView
}
