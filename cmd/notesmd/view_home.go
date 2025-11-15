package main

import (
	"github.com/charmbracelet/lipgloss"
)

// viewHome renders the home/welcome screen
func (m model) viewHome() string {
	shortcuts := `
Raccourcis :

  q          → quitter
  ENTER      → ouvrir l'explorateur de notes
  t          → changer la couleur d'accent
  n          → créer une nouvelle note rapide

  ↑ / ↓      → naviguer dans la liste
  ← / h      → remonter au dossier parent
  → / l / ↵  → entrer dans un dossier
  o          → ouvrir / fermer la prévisualisation
  j / k      → scroller la prévisualisation (bas / haut)
  u / d      → scroll rapide (haut / bas)
  e          → éditer la note dans $EDITOR (Neovim)
  /          → recherche globale fuzzy (style Finder) depuis le dossier de départ
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

	// Center the card on screen
	return lipgloss.Place(
		w,
		h,
		lipgloss.Center,
		lipgloss.Center,
		card,
	)
}
