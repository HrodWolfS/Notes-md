package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleModalKey handles keyboard input for the note creation modal
func (m *model) handleModalKey(msg tea.KeyMsg) (handled bool, cmd tea.Cmd) {
	s := msg.String()

	switch s {
	case "esc":
		// Close modal without saving
		m.showModal = false
		m.modal = newNoteModal()
		return true, nil

	case "ctrl+s", "ctrl+enter":
		// Save note
		path, err := m.modal.CreateNote(m.currentDir)
		if err != nil {
			// Could show error in status bar, for now just close
			m.showModal = false
			m.modal = newNoteModal()
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
		m.showModal = false
		m.modal = newNoteModal()
		return true, nil
	}

	// Let modal handle the key
	var modalCmd tea.Cmd
	m.modal, modalCmd = m.modal.Update(msg)
	return true, modalCmd
}
