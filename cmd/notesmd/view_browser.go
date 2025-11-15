package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// viewBrowser renders the file browser view
func (m model) viewBrowser() string {
	// Dynamic styles based on current accent color
	ts := titleStyle
	sl := searchLabelStyle
	accent := lipgloss.Color("208")
	if len(titlePalette) > 0 {
		colorCode := titlePalette[m.themeIndex]
		accent = lipgloss.Color(colorCode)
		ts = ts.Foreground(accent)
		sl = sl.Foreground(accent)
	}

	header := ts.Render("\nExplorateur de notes") + "\n\n"
	header += fmt.Sprintf("Dossier : %s\n\n", m.currentDir)

	// Total width known by the model
	totalWidth := m.width
	if totalWidth <= 0 {
		// Fallback in case we haven't received WindowSizeMsg yet
		totalWidth = 120
	}

	// ~30% for list, 70% for preview
	leftWidth := int(float64(totalWidth) * 0.30)
	if leftWidth < 20 {
		leftWidth = 20
	}

	// Keep some margin for borders
	rightWidth := totalWidth - leftWidth - 4
	if rightWidth < 20 {
		rightWidth = 20
	}

	left := borderStyle.
		BorderForeground(accent).
		Width(leftWidth).
		Render(m.list.View())

	var rightContent string
	if m.showPreview {
		rightContent = m.viewport.View()
	} else {
		rightContent = "Appuie sur 'o' pour prévisualiser un fichier Markdown."
	}

	right := borderStyle.
		BorderForeground(accent).
		Width(rightWidth).
		Render(rightContent)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	var footer string
	if m.searchActive {
		footer = helpStyle.Render(
			"\n" +
				sl.Render("Search: ") +
				searchQueryStyle.Render(m.searchQuery) +
				"\nENTER pour ouvrir le résultat sélectionné — ESC pour annuler\n",
		)
	} else {
		footer = helpStyle.Render(
			"\n? aide • ↑/↓ naviguer • dd supprimer • r renommer • m filtre .md • n nouvelle note • / rechercher • q quitter\n",
		)
	}

	// Status bar: current folder + item count + mode
	modeLabel := "browse"
	if m.searchActive {
		modeLabel = "search"
	}

	itemCount := len(m.list.Items())
	status := statusBarStyle.
		Background(accent).
		Foreground(lipgloss.Color("0")).
		Render(fmt.Sprintf(" %s | %d éléments | mode: %s ", m.currentDir, itemCount, modeLabel))

	baseView := header + status + "\n" + layout + footer

	// If any modal is open, overlay it on top
	var modalView string
	if m.showNoteModal {
		modalView = m.noteModal.View()
	} else if m.showConfirmModal {
		modalView = m.confirmModal.View()
	} else if m.showRenameModal {
		modalView = m.renameModal.View()
	} else if m.showHelpModal {
		modalView = m.helpModal.View()
	}

	if modalView != "" {
		// Center the modal on screen
		w := m.width
		h := m.height
		if w <= 0 {
			w = 120
		}
		if h <= 0 {
			h = 40
		}

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
