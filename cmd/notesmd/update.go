package main

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all state updates
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Dynamic resize handling
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Responsive calculation
		leftWidth := int(float64(m.width) * 0.30)
		rightWidth := m.width - leftWidth - 4 // compensate for borders
		viewHeight := m.height - 10           // header/footer (adjustable)

		// Apply to components
		m.list.SetWidth(leftWidth)
		m.list.SetHeight(viewHeight)

		m.viewport.Width = rightWidth
		m.viewport.Height = viewHeight

		return m, nil

	// Editor finished
	case editorDoneMsg:
		if m.mode == modeBrowser && m.showPreview {
			if it, ok := m.list.SelectedItem().(fileItem); ok && !it.isDir {
				m.viewport.SetContent(loadMarkdown(it.path))
			}
		}
		return m, nil

	// Editor error
	case editorErrorMsg:
		m.showPreview = true
		m.viewport.SetContent(fmt.Sprintf("Erreur lors de l'ouverture de l'éditeur :\n\n%s", string(msg)))
		return m, nil

	// Key handling
	case tea.KeyMsg:
		switch m.mode {

		case modeHome:
			return m.updateHome(msg)

		case modeBrowser:
			return m.updateBrowser(msg)
		}
	}

	return m, nil
}

// updateHome handles updates for home screen
func (m model) updateHome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.mode = modeBrowser
		m.setDir(m.rootDir)
	case "t":
		m.toggleTheme()
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

// updateBrowser handles updates for browser mode
func (m model) updateBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If modal is open, handle that first
	if m.showModal {
		handled, cmd := m.handleModalKey(msg)
		if handled {
			return m, cmd
		}
	}

	if m.searchActive {
		handled, cmd := m.handleSearchKey(msg)
		if handled {
			return m, cmd
		}
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "right", "l", "enter":
		if m.searchActive {
			break
		}
		if it, ok := m.list.SelectedItem().(fileItem); ok {
			if it.isDir {
				m.setDir(it.path)
				m.showPreview = false
				m.viewport.SetContent("Appuie sur 'o' pour prévisualiser un fichier Markdown.")
			}
		}

	case "left", "h":
		if m.searchActive {
			break
		}
		parent := filepath.Dir(m.currentDir)
		if parent != m.currentDir {
			m.setDir(parent)
			m.showPreview = false
			m.viewport.SetContent("Appuie sur 'o' pour prévisualiser un fichier Markdown.")
		}

	case "o":
		if it, ok := m.list.SelectedItem().(fileItem); ok && !it.isDir {
			if m.showPreview {
				m.showPreview = false
				m.viewport.SetContent("Appuie sur 'o' pour prévisualiser un fichier Markdown.")
			} else {
				content := loadMarkdown(it.path)
				m.viewport.SetContent(content)
				m.showPreview = true
			}
		}

	case "k":
		m.viewport.LineUp(1)
	case "j":
		m.viewport.LineDown(1)
	case "u":
		m.viewport.LineUp(5)
	case "d":
		m.viewport.LineDown(5)

	case "e":
		if it, ok := m.list.SelectedItem().(fileItem); ok && !it.isDir {
			return m, openInEditor(it.path)
		}

	case "t":
		m.toggleTheme()
		return m, nil

	case "n":
		// Open note creation modal
		m.showModal = true
		m.modal = newNoteModal()
		return m, nil

	case "/":
		m.searchActive = true
		m.searchQuery = ""
		m.ensureAllFilesScanned()
		m.buildSearchResults()
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// handleSearchKey handles keyboard input during search
func (m *model) handleSearchKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {

	case "enter":
		if it, ok := m.list.SelectedItem().(fileItem); ok {
			targetDir := filepath.Dir(it.path)
			m.setDir(targetDir)
			for i, item := range m.baseItems {
				if fi, ok := item.(fileItem); ok && fi.path == it.path {
					m.list.Select(i)
					break
				}
			}
		}
		m.searchActive = false
		m.searchQuery = ""
		return true, nil

	case "esc":
		m.searchActive = false
		m.searchQuery = ""
		m.list.SetItems(m.baseItems)
		return true, nil

	case "backspace", "backspace2":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.buildSearchResults()
		}
		return true, nil
	}

	if len(s) == 1 && s >= " " && s <= "~" {
		m.searchQuery += s
		m.buildSearchResults()
		return true, nil
	}

	return false, nil
}
