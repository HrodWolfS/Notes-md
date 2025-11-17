package main

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderBreadcrumb creates a breadcrumb path display
func renderBreadcrumb(currentPath string, accent lipgloss.Color) string {
	parts := strings.Split(currentPath, string(filepath.Separator))
	var crumbs []string

	for i, part := range parts {
		if part == "" {
			continue
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		if i == len(parts)-1 {
			style = lipgloss.NewStyle().
				Foreground(accent).
				Bold(true)
		}

		crumbs = append(crumbs, style.Render(part))
	}

	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(" → ")

	return strings.Join(crumbs, separator)
}

// viewBrowser renders the file browser view
func (m model) viewBrowser() string {
	m.updateStatusBar()

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
	header += renderBreadcrumb(m.currentDir, accent) + "\n\n"

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
		rightContent = "Preview masqué. Appuie sur 'o' pour afficher."
	}

	right := borderStyle.
		BorderForeground(accent).
		Width(rightWidth).
		Render(rightContent)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	var footer string
	if m.searchInNoteActive {
		footer = helpStyle.Render(
			"\n" +
				sl.Render("Rechercher dans la note: ") +
				searchQueryStyle.Render(m.noteSearchQuery) +
				"\nMises à jour en temps réel • ↑/↓ ou Ctrl+u/d pour naviguer • ESC pour annuler\n",
		)
	} else if m.searchActive {
		footer = helpStyle.Render(
			"\n" +
				sl.Render("Search: ") +
				searchQueryStyle.Render(m.searchQuery) +
				"\nENTER pour ouvrir le résultat sélectionné — ESC pour annuler\n",
		)
	} else {
		footer = helpStyle.Render(
			"\n? aide • ↑/↓ naviguer • D supprimer • r renommer • m filtre .md • n nouvelle note • / rechercher • F recherche note • q quitter\n",
		)
	}

	// Status bar using custom component
	statusBar := m.statusBar.View()

	baseView := header + layout + "\n" + statusBar + footer

	// If any modal is open, overlay it on top
	var modalView string
	if m.showNoteModal {
		modalView = m.noteModal.View()
	} else if m.showConfirmModal {
		modalView = m.confirmModal.View()
	} else if m.showRenameModal {
		modalView = m.renameModal.View()
	} else if m.showCreateDirModal {
		modalView = m.createDirModal.View()
	} else if m.showRecentModal {
		modalView = m.recentModal.View()
	} else if m.showBookmarksModal {
		modalView = m.bookmarksModal.View()
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
