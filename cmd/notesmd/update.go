package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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

	// Clipboard copied
	case clipboardCopiedMsg:
		cmd := m.statusBar.SetMessage(msg.message, 2*time.Second)
		return m, cmd

	// Paste completed
	case pasteCompletedMsg:
		if msg.success {
			m.setDir(m.currentDir)
		}
		cmd := m.statusBar.SetMessage(msg.message, 2*time.Second)
		return m, cmd

	// Content search completed (deprecated - using in-note search now)
	// case searchCompletedMsg:
	// 	m.searchResults = msg.results
	// 	m.contentSearchActive = false
	// 	if len(msg.results) > 0 {
	// 		items := make([]blist.Item, len(msg.results))
	// 		for i, result := range msg.results {
	// 			items[i] = searchResultItem{result: result}
	// 		}
	// 		m.list.SetItems(items)
	// 	}
	// 	return m, nil

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
	// Handle help modal
	if m.showHelpModal {
		if msg.String() == "esc" || msg.String() == "?" {
			m.showHelpModal = false
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "?":
		m.showHelpModal = true
		m.helpModal = newHelpModal()
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

	if m.showCreateDirModal {
		handled, cmd := m.handleCreateDirModalKey(msg)
		if handled {
			return m, cmd
		}
	}

	if m.showRecentModal {
		handled, cmd := m.handleRecentModalKey(msg)
		if handled {
			return m, cmd
		}
	}

	if m.showBookmarksModal {
		handled, cmd := m.handleBookmarksModalKey(msg)
		if handled {
			return m, cmd
		}
	}

	if m.showLinksModal {
		handled, cmd := m.handleLinksModalKey(msg)
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

	if m.searchInNoteActive {
		handled, cmd := m.handleNoteSearchKey(msg)
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
			} else {
				m.trackRecentFile(it.path)
				m.currentNotePath = it.path
			}
		}

	case "left", "h":
		if m.searchActive {
			break
		}
		parent := filepath.Dir(m.currentDir)
		if parent != m.currentDir {
			m.setDir(parent)
		}

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

	case "ctrl+o":
		// Navigate back in history
		m.navBack()
		m.lastKey = ""

	case "ctrl+i":
		// Navigate forward in history
		m.navForward()
		m.lastKey = ""

	// Preview scrolling
	case "u":
		// Scroll preview up
		if m.showPreview {
			m.viewport.LineUp(3)
		}
		m.lastKey = ""

	case "d":
		// Scroll preview down
		if m.showPreview {
			m.viewport.LineDown(3)
		}
		m.lastKey = ""

	// File operations
	case "D":
		// Delete file/folder with confirmation
		if it, ok := m.list.SelectedItem().(fileItem); ok {
			m.showConfirmModal = true
			m.confirmModal = newConfirmDeleteModal(it.path, it.name)
		}
		m.lastKey = ""

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

	case ".":
		// Toggle hidden files
		m.showHidden = !m.showHidden
		m.applyFilters()
		m.lastKey = ""

	case "s":
		// Cycle through sort modes: name -> date -> size -> name
		m.sortMode = (m.sortMode + 1) % 3
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
		}
		m.lastKey = ""

	case "y":
		// Copy file path to clipboard
		if it, ok := m.list.SelectedItem().(fileItem); ok {
			msg := copyFilePath(it.path)
			cmd := m.statusBar.SetMessage(msg.message, 2*time.Second)
			return m, cmd
		}
		m.lastKey = ""

	case "Y":
		// Copy file content to clipboard
		if it, ok := m.list.SelectedItem().(fileItem); ok && !it.isDir {
			msg := copyFileContent(it.path)
			cmd := m.statusBar.SetMessage(msg.message, 2*time.Second)
			return m, cmd
		}
		m.lastKey = ""

	case "N":
		// Create new directory
		m.showCreateDirModal = true
		m.createDirModal = newCreateDirModal(m.currentDir)
		m.lastKey = ""
		return m, nil

	case "c":
		// Copy file (to internal clipboard)
		if it, ok := m.list.SelectedItem().(fileItem); ok {
			m.clipboard = &FileClipboard{path: it.path, mode: "copy"}
			cmd := m.statusBar.SetMessage("Copied: "+it.name, 2*time.Second)
			return m, cmd
		}
		m.lastKey = ""

	case "p":
		// Paste file
		if m.clipboard != nil {
			return m, pasteFile(m.clipboard, m.currentDir)
		}
		m.lastKey = ""

	case "b":
		// Toggle bookmark on current file
		if it, ok := m.list.SelectedItem().(fileItem); ok && !it.isDir {
			added := m.toggleBookmark(it.path)
			message := "Bookmark removed"
			if added {
				message = "Bookmark added"
			}
			cmd := m.statusBar.SetMessage(message, 2*time.Second)
			return m, cmd
		}
		m.lastKey = ""

	case "B":
		// Show all bookmarks
		m.showBookmarksModal = true
		m.bookmarksModal = newBookmarksModal(m.bookmarks, m.width, m.height)
		m.lastKey = ""
		return m, nil

	case "L":
		// Show links in current note
		if m.currentNotePath != "" {
			// Parse links from raw content
			rawContent := loadMarkdownRaw(m.currentNotePath)
			links := parseWikiLinks(rawContent)

			if len(links) > 0 {
				m.showLinksModal = true
				m.linksModal = newLinksModal(links, m.rootDir, m.width, m.height)
			} else {
				cmd := m.statusBar.SetMessage("Aucun lien trouvé dans cette note", 2*time.Second)
				return m, cmd
			}
		}
		m.lastKey = ""
		return m, nil

	case "ctrl+r":
		// Show recent files
		m.showRecentModal = true
		m.recentModal = newRecentFilesModal(m.recentFiles, m.width, m.height)
		m.lastKey = ""
		return m, nil

	case "F":
		// Start in-note search (only if a note is open)
		if m.showPreview && m.currentNotePath != "" {
			m.searchInNoteActive = true
			m.noteSearchQuery = ""
			m.lastKey = ""
		}
		return m, nil

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

	// Auto-preview on selection change
	if m.autoPreview {
		currentIndex := m.list.Index()
		if currentIndex != m.lastSelectedIndex {
			m.lastSelectedIndex = currentIndex
			if it, ok := m.list.SelectedItem().(fileItem); ok && !it.isDir {
				content := loadMarkdownWithLinks(it.path, m.rootDir)
				m.viewport.SetContent(content)
				m.showPreview = true
				m.currentNotePath = it.path
			}
		}
	}

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

// handleNoteSearchKey handles keyboard input during in-note search
func (m *model) handleNoteSearchKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {

	case "esc":
		// Cancel search and restore original content
		m.searchInNoteActive = false
		if m.currentNotePath != "" {
			content := loadMarkdownWithLinks(m.currentNotePath, m.rootDir)
			m.viewport.SetContent(content)
		}
		m.noteSearchQuery = ""
		return true, nil

	case "backspace", "backspace2":
		if len(m.noteSearchQuery) > 0 {
			m.noteSearchQuery = m.noteSearchQuery[:len(m.noteSearchQuery)-1]
			// Live update on backspace
			if m.currentNotePath != "" {
				content := loadMarkdownWithHighlight(m.currentNotePath, m.noteSearchQuery)
				m.viewport.SetContent(content)
			}
		}
		return true, nil

	// Allow viewport navigation during search (arrows and Ctrl keys only)
	case "up", "down", "pgup", "pgdown", "ctrl+u", "ctrl+d":
		var viewportCmd tea.Cmd
		m.viewport, viewportCmd = m.viewport.Update(msg)
		return true, viewportCmd
	}

	// Add character to search query
	if len(s) == 1 && s >= " " && s <= "~" {
		m.noteSearchQuery += s
		// Live update as user types
		if m.currentNotePath != "" {
			content := loadMarkdownWithHighlight(m.currentNotePath, m.noteSearchQuery)
			m.viewport.SetContent(content)
		}
		return true, nil
	}

	return false, nil
}
