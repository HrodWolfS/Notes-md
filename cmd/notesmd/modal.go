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

// ========== Delete Confirmation Modal ==========

type confirmModal struct {
	message string
	path    string
}

func newConfirmDeleteModal(path, itemName string) confirmModal {
	return confirmModal{
		message: fmt.Sprintf("Supprimer '%s' ?", itemName),
		path:    path,
	}
}

func (m confirmModal) View() string {
	// Modal title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Red for danger
		Bold(true).
		Render("⚠️  Confirmation")

	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Render(m.message)

	helpText := helpStyle.Render("y: confirmer • n/Esc: annuler")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		message,
		"",
		helpText,
	)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")). // Red border
		Padding(1, 4).
		Width(50)

	return modalStyle.Render(content)
}

// ========== Rename Modal ==========

type renameModal struct {
	input    textinput.Model
	oldPath  string
	itemName string
}

func newRenameModal(path, currentName string) renameModal {
	ti := textinput.New()
	ti.Placeholder = "nouveau-nom"
	ti.SetValue(currentName)
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// Select all text for easy replacement
	ti.CursorEnd()

	return renameModal{
		input:    ti,
		oldPath:  path,
		itemName: currentName,
	}
}

func (m renameModal) Update(msg tea.Msg) (renameModal, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m renameModal) View() string {
	title := titleStyle.Render("Renommer")

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Nouveau nom:")

	inputField := m.input.View()

	helpText := helpStyle.Render("Enter: confirmer • Esc: annuler")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		label,
		inputField,
		"",
		helpText,
	)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(1, 2).
		Width(60)

	return modalStyle.Render(content)
}

func (m renameModal) GetNewName() string {
	return strings.TrimSpace(m.input.Value())
}

func (m renameModal) Rename() error {
	newName := m.GetNewName()
	if newName == "" || newName == m.itemName {
		return fmt.Errorf("nom invalide ou inchangé")
	}

	dir := filepath.Dir(m.oldPath)
	newPath := filepath.Join(dir, newName)

	return os.Rename(m.oldPath, newPath)
}

// ========== Help Modal ==========

type helpModal struct{}

func newHelpModal() helpModal {
	return helpModal{}
}

func (m helpModal) View() string {
	title := titleStyle.Render("📖 Guide des raccourcis")

	// Navigation section
	navTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Navigation Vim:")
	navContent := `  gg         → Aller au début
  G          → Aller à la fin
  Ctrl+d     → Page down (½ page)
  Ctrl+u     → Page up (½ page)
  -          → Dossier parent
  ~          → Dossier home
  ↑/↓ ou j/k → Naviguer dans la liste
  ←/h        → Remonter au parent
  →/l/Enter  → Entrer dans un dossier`

	// File operations section
	fileTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Actions fichiers:")
	fileContent := `  n          → Nouvelle note (modal)
  dd         → Supprimer (avec confirmation)
  r          → Renommer
  e          → Éditer dans $EDITOR`

	// Filters section
	filterTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Filtres & Recherche:")
	filterContent := `  m          → Toggle .md uniquement
  /          → Recherche globale fuzzy`

	// Interface section
	uiTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Interface:")
	uiContent := `  o          → Toggle preview
  j/k        → Scroll preview (↓/↑)
  u/d        → Scroll rapide preview
  t          → Changer thème
  ?          → Afficher cette aide
  q          → Quitter`

	helpText := helpStyle.Render("Appuyez sur Esc ou ? pour fermer")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		navTitle,
		navContent,
		"",
		fileTitle,
		fileContent,
		"",
		filterTitle,
		filterContent,
		"",
		uiTitle,
		uiContent,
		"",
		helpText,
	)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")). // Light blue
		Padding(1, 2).
		Width(70)

	return modalStyle.Render(content)
}
