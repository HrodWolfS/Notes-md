package main

import (
	"fmt"
	"os"
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
	// If any modal is open, handle that first
	if m.showNoteModal {
		handled, cmd := m.handleNoteModalKey(msg)
		if handled {
			return m, cmd
		}
	}

	if m.showConfirmModal {
		handled, cmd := m.handleConfirmModalKey(msg)
		if handled {
			return m, cmd
		}
	}

	if m.showRenameModal {
		handled, cmd := m.handleRenameModalKey(msg)
		if handled {
			return m, cmd
		}
	}

	if m.showHelpModal {
		// Close help modal on Esc or ?
		if msg.String() == "esc" || msg.String() == "?" {
			m.showHelpModal = false
			return m, nil
		}
		return m, nil
	}

	if m.searchActive {
		handled, cmd := m.handleSearchKey(msg)
		if handled {
			return m, cmd
		}
	}

	// Reset lastKey for non-continuation keys
	key := msg.String()
	if key != "g" && m.lastKey == "g" {
		m.lastKey = ""
	}
	if key != "d" && m.lastKey == "d" {
		m.lastKey = ""
	}

	switch key {

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

	// Vim-style navigation
	case "g":
		// Handle 'gg' to go to top
		if m.lastKey == "g" {
			m.list.Select(0)
			m.lastKey = ""
			return m, nil
		}
		m.lastKey = "g"
		return m, nil

	case "G":
		// Go to bottom of list
		itemCount := len(m.list.Items())
		if itemCount > 0 {
			m.list.Select(itemCount - 1)
		}
		m.lastKey = ""

	case "ctrl+d":
		// Page down (half page)
		current := m.list.Index()
		pageSize := m.list.Height() / 2
		newIndex := current + pageSize
		itemCount := len(m.list.Items())
		if newIndex >= itemCount {
			newIndex = itemCount - 1
		}
		if newIndex >= 0 {
			m.list.Select(newIndex)
		}
		m.lastKey = ""

	case "ctrl+u":
		// Page up (half page)
		current := m.list.Index()
		pageSize := m.list.Height() / 2
		newIndex := current - pageSize
		if newIndex < 0 {
			newIndex = 0
		}
		m.list.Select(newIndex)
		m.lastKey = ""

	// File operations
	case "d":
		// Handle 'dd' to delete
		if m.lastKey == "d" {
			m.lastKey = ""
			if it, ok := m.list.SelectedItem().(fileItem); ok {
				m.showConfirmModal = true
				m.confirmModal = newConfirmDeleteModal(it.path, it.name)
			}
			return m, nil
		}
		m.lastKey = "d"
		return m, nil

	case "r":
		// Rename file/folder
		if it, ok := m.list.SelectedItem().(fileItem); ok {
			m.showRenameModal = true
			m.renameModal = newRenameModal(it.path, it.name)
		}
		m.lastKey = ""

	case "e":
		if it, ok := m.list.SelectedItem().(fileItem); ok && !it.isDir {
			return m, openInEditor(it.path)
		}

	case "t":
		m.toggleTheme()
		m.lastKey = ""
		return m, nil

	case "?":
		// Open help modal
		m.showHelpModal = true
		m.helpModal = newHelpModal()
		m.lastKey = ""
		return m, nil

	case "n":
		// Open note creation modal
		m.showNoteModal = true
		m.noteModal = newNoteModal()
		m.lastKey = ""
		return m, nil

	case "m":
		// Toggle .md filter
		m.mdOnly = !m.mdOnly
		m.applyFilters()
		m.lastKey = ""

	case "-":
		// Navigate to parent directory
		if m.searchActive {
			break
		}
		parent := filepath.Dir(m.currentDir)
		if parent != m.currentDir {
			m.setDir(parent)
			m.showPreview = false
			m.viewport.SetContent("Appuie sur 'o' pour prévisualiser un fichier Markdown.")
		}
		m.lastKey = ""

	case "~":
		// Navigate to home directory
		if m.searchActive {
			break
		}
		home, err := os.UserHomeDir()
		if err == nil {
			m.setDir(home)
			m.rootDir = home
			m.showPreview = false
			m.viewport.SetContent("Appuie sur 'o' pour prévisualiser un fichier Markdown.")
		}
		m.lastKey = ""

	case "/":
		m.searchActive = true
		m.searchQuery = ""
		m.ensureAllFilesScanned()
		m.buildSearchResults()
		m.lastKey = ""
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
