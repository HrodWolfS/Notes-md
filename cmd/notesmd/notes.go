package main

import (
	"os"

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
