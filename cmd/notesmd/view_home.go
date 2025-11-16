package main

import (
	"github.com/charmbracelet/lipgloss"
)

// viewHome renders the home/welcome screen
func (m model) viewHome() string {
	shortcuts := `
Bienvenue dans NotesMD

  ENTER      → Ouvrir l'explorateur de notes
  ?          → Afficher tous les raccourcis
  t          → Changer le thème
  q          → Quitter
`

	// Dynamic accent color
	ts := titleStyle
	accent := lipgloss.Color("208")
	if len(titlePalette) > 0 {
		colorCode := titlePalette[m.themeIndex]
		accent = lipgloss.Color(colorCode)
		ts = ts.Foreground(accent)
	}

	content := ts.Render(asciiTitle) + helpStyle.Render(shortcuts)
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
