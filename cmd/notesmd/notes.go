package main

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// handleNoteModalKey handles keyboard input for the note creation modal
func (m *model) handleNoteModalKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {
	case "esc":
		// Close modal without saving
		m.showNoteModal = false
		m.noteModal = newNoteModal()
		return true, nil

	case "ctrl+s", "ctrl+enter":
		// Save note
		path, err := m.noteModal.CreateNote(m.currentDir)
		if err != nil {
			// Could show error in status bar, for now just close
			m.showNoteModal = false
			m.noteModal = newNoteModal()
			return true, nil
		}

		// Refresh list and select the new note
		m.baseItems = readDir(m.currentDir)
		m.list.SetItems(m.baseItems)
		for i, item := range m.baseItems {
			if fi, ok := item.(fileItem); ok && fi.path == path {
				m.list.Select(i)
				break
			}
		}

		// Close modal
		m.showNoteModal = false
		m.noteModal = newNoteModal()
		return true, nil
	}

	// Let modal handle the key
	var modalCmd tea.Cmd
	m.noteModal, modalCmd = m.noteModal.Update(msg)
	return true, modalCmd
}

// handleConfirmModalKey handles keyboard input for delete confirmation modal
func (m *model) handleConfirmModalKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {
	case "y":
		// Confirm deletion
		err := os.RemoveAll(m.confirmModal.path)
		if err == nil {
			// Refresh list
			m.baseItems = readDir(m.currentDir)
			m.list.SetItems(m.baseItems)
		}
		m.showConfirmModal = false
		return true, nil

	case "n", "esc":
		// Cancel deletion
		m.showConfirmModal = false
		return true, nil
	}

	return true, nil
}

// handleRenameModalKey handles keyboard input for rename modal
func (m *model) handleRenameModalKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {
	case "esc":
		// Cancel rename
		m.showRenameModal = false
		m.renameModal = renameModal{}
		return true, nil

	case "enter":
		// Perform rename
		err := m.renameModal.Rename()
		if err == nil {
			// Refresh list
			m.baseItems = readDir(m.currentDir)
			m.list.SetItems(m.baseItems)
		}
		m.showRenameModal = false
		m.renameModal = renameModal{}
		return true, nil
	}

	// Let modal handle the key
	var modalCmd tea.Cmd
	m.renameModal, modalCmd = m.renameModal.Update(msg)
	return true, modalCmd
}

func (m *model) handleCreateDirModalKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {
	case "esc":
		m.showCreateDirModal = false
		m.createDirModal = createDirModal{}
		return true, nil

	case "enter":
		_, err := m.createDirModal.CreateDir()
		if err == nil {
			m.baseItems = readDir(m.currentDir)
			m.applyFilters()
		}
		m.showCreateDirModal = false
		m.createDirModal = createDirModal{}
		return true, nil
	}

	var modalCmd tea.Cmd
	m.createDirModal, modalCmd = m.createDirModal.Update(msg)
	return true, modalCmd
}

func (m *model) handleRecentModalKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {
	case "esc":
		m.showRecentModal = false
		return true, nil

	case "enter":
		if it, ok := m.recentModal.list.SelectedItem().(fileItem); ok {
			targetDir := filepath.Dir(it.path)
			m.setDir(targetDir)

			for i, item := range m.baseItems {
				if fi, ok := item.(fileItem); ok && fi.path == it.path {
					m.list.Select(i)
					break
				}
			}
			m.showRecentModal = false
			return true, nil
		}
	}

	var modalCmd tea.Cmd
	m.recentModal, modalCmd = m.recentModal.Update(msg)
	return true, modalCmd
}

func (m *model) handleBookmarksModalKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {
	case "esc":
		m.showBookmarksModal = false
		return true, nil

	case "enter":
		if it, ok := m.bookmarksModal.list.SelectedItem().(fileItem); ok {
			targetDir := filepath.Dir(it.path)
			m.setDir(targetDir)

			for i, item := range m.baseItems {
				if fi, ok := item.(fileItem); ok && fi.path == it.path {
					m.list.Select(i)
					break
				}
			}
			m.showBookmarksModal = false
			return true, nil
		}

	case "D":
		if it, ok := m.bookmarksModal.list.SelectedItem().(fileItem); ok {
			m.toggleBookmark(it.path)
			m.bookmarksModal = newBookmarksModal(m.bookmarks, m.width, m.height)
			return true, nil
		}
	}

	var modalCmd tea.Cmd
	m.bookmarksModal, modalCmd = m.bookmarksModal.Update(msg)
	return true, modalCmd
}
