package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// focusField represents which field is currently focused
type focusField int

const (
	focusName focusField = iota
	focusContent
)

// noteModal represents the note creation modal state
type noteModal struct {
	nameInput    textinput.Model
	contentInput textarea.Model
	focused      focusField
	width        int
	height       int
}

// newNoteModal creates a new note creation modal
func newNoteModal() noteModal {
	// Name input
	ti := textinput.New()
	ti.Placeholder = "ma-note"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// Content textarea
	ta := textarea.New()
	ta.Placeholder = "Écrivez votre note en Markdown..."
	ta.ShowLineNumbers = false
	ta.CharLimit = 5000
	ta.SetWidth(60)
	ta.SetHeight(10)

	return noteModal{
		nameInput:    ti,
		contentInput: ta,
		focused:      focusName,
	}
}

// Update updates the modal state
func (m noteModal) Update(msg tea.Msg) (noteModal, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			// Switch focus between fields
			if m.focused == focusName {
				m.focused = focusContent
				m.nameInput.Blur()
				m.contentInput.Focus()
			} else {
				m.focused = focusName
				m.contentInput.Blur()
				m.nameInput.Focus()
			}
			return m, nil
		}
	}

	// Update the focused field
	if m.focused == focusName {
		m.nameInput, cmd = m.nameInput.Update(msg)
	} else {
		m.contentInput, cmd = m.contentInput.Update(msg)
	}

	return m, cmd
}

// View renders the modal
func (m noteModal) View() string {
	// Modal title
	title := titleStyle.Render("Nouvelle note")

	// Name field
	nameLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Nom du fichier:")
	nameField := m.nameInput.View()

	// Content field
	contentLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Contenu:")
	contentField := m.contentInput.View()

	// Help text
	helpText := helpStyle.Render(
		"TAB: changer de champ • Ctrl+S/Ctrl+Enter: sauvegarder • Esc: annuler",
	)

	// Assemble modal content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		nameLabel,
		nameField,
		"",
		contentLabel,
		contentField,
		"",
		helpText,
	)

	// Modal style with border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(1, 2).
		Width(70)

	modal := modalStyle.Render(content)

	return modal
}

// GetName returns the current name input value
func (m noteModal) GetName() string {
	return strings.TrimSpace(m.nameInput.Value())
}

// GetContent returns the current content input value
func (m noteModal) GetContent() string {
	return strings.TrimSpace(m.contentInput.Value())
}

// CreateNote creates the note file and returns the path
func (m noteModal) CreateNote(currentDir string) (string, error) {
	name := m.GetName()
	if name == "" {
		return "", fmt.Errorf("le nom ne peut pas être vide")
	}

	// Add .md extension if not present
	if filepath.Ext(name) == "" {
		name += ".md"
	}

	path := filepath.Join(currentDir, name)

	// Title derived from filename (without extension)
	title := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	content := fmt.Sprintf("# %s\n\n%s\n", title, m.GetContent())

	err := os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		return "", err
	}

	return path, nil
}
