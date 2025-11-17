package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

// ========== Create Directory Modal ==========

type createDirModal struct {
	input    textinput.Model
	basePath string
}

func newCreateDirModal(basePath string) createDirModal {
	ti := textinput.New()
	ti.Placeholder = "nouveau-dossier"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return createDirModal{
		input:    ti,
		basePath: basePath,
	}
}

func (m createDirModal) Update(msg tea.Msg) (createDirModal, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m createDirModal) View() string {
	title := titleStyle.Render("Nouveau Dossier")

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Nom du dossier:")

	inputField := m.input.View()

	helpText := helpStyle.Render("Enter: créer • Esc: annuler")

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

func (m createDirModal) GetDirName() string {
	return strings.TrimSpace(m.input.Value())
}

func (m createDirModal) CreateDir() (string, error) {
	name := m.GetDirName()
	if name == "" {
		return "", fmt.Errorf("le nom ne peut pas être vide")
	}

	path := filepath.Join(m.basePath, name)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}

	return path, nil
}

// ========== Recent Files Modal ==========

type recentFilesModal struct {
	list blist.Model
}

func newRecentFilesModal(recentFiles []string, width, height int) recentFilesModal {
	items := make([]blist.Item, 0, len(recentFiles))
	for _, path := range recentFiles {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		items = append(items, fileItem{
			name:    filepath.Base(path),
			path:    path,
			isDir:   false,
			size:    info.Size(),
			modTime: info.ModTime().Unix(),
		})
	}

	l := blist.New(items, blist.NewDefaultDelegate(), width-10, height-10)
	l.Title = "Recent Files"
	l.SetShowHelp(false)

	return recentFilesModal{
		list: l,
	}
}

func (m recentFilesModal) Update(msg tea.Msg) (recentFilesModal, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m recentFilesModal) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true).
		Render("📋 Recent Files")

	listView := m.list.View()

	helpText := helpStyle.Render("Enter: open • Esc: close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
		"",
		helpText,
	)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("212")).
		Padding(1, 2).
		Width(70)

	return modalStyle.Render(content)
}

// ========== Bookmarks Modal ==========

type bookmarksModal struct {
	list blist.Model
}

func newBookmarksModal(bookmarks []string, width, height int) bookmarksModal {
	items := make([]blist.Item, 0, len(bookmarks))
	for _, path := range bookmarks {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		item := fileItem{
			name:    filepath.Base(path),
			path:    path,
			isDir:   info.IsDir(),
			size:    info.Size(),
			modTime: info.ModTime().Unix(),
		}

		items = append(items, item)
	}

	l := blist.New(items, blist.NewDefaultDelegate(), width-10, height-10)
	l.Title = "Bookmarks"
	l.SetShowHelp(false)

	return bookmarksModal{
		list: l,
	}
}

func (m bookmarksModal) Update(msg tea.Msg) (bookmarksModal, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m bookmarksModal) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true).
		Render("★ Bookmarks")

	listView := m.list.View()

	helpText := helpStyle.Render("Enter: open • d: remove • Esc: close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
		"",
		helpText,
	)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("212")).
		Padding(1, 2).
		Width(70)

	return modalStyle.Render(content)
}

// ========== Help Modal ==========

type helpModal struct{}

func newHelpModal() helpModal {
	return helpModal{}
}

func (m helpModal) View() string {
	title := titleStyle.Render("📖 Raccourcis")

	// Navigation section
	navTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Navigation:")
	navContent := `j/k/↑/↓: liste | gg: début | G: fin | Ctrl+d/u: page
Ctrl+o/i: historique | -: parent | ~: home | h/l: ←/→`

	// File operations section
	fileTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Fichiers:")
	fileContent := `n: nouvelle note | N: dossier | D: supprimer | r: renommer
e: éditer | c: copier | p: coller | y: path | Y: contenu`

	// Organization section
	orgTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Organisation:")
	orgContent := `b: bookmark | B: voir bookmarks | Ctrl+R: récents | L: liens wiki`

	// Filters section
	filterTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Filtres:")
	filterContent := `m: .md only | .: hidden | s: tri | /: nom | F: note`

	// Interface section
	uiTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("213")).
		Bold(true).
		Render("Interface:")
	uiContent := `u/d: scroll | t: thème | ?: aide | q: quitter`

	helpText := helpStyle.Render("Esc ou ? pour fermer")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		navTitle,
		navContent,
		"",
		fileTitle,
		fileContent,
		"",
		orgTitle,
		orgContent,
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
		Padding(0, 1).
		Width(68)

	return modalStyle.Render(content)
}

// ========== Links Modal ==========

type linkItem struct {
	name   string
	path   string
	exists bool
}

func (l linkItem) Title() string {
	if l.exists {
		return "🔗 " + l.name
	}
	return "❓ " + l.name + " (n'existe pas)"
}

func (l linkItem) Description() string {
	if l.exists {
		return l.path
	}
	return "Appuyez sur Enter pour créer cette note"
}

func (l linkItem) FilterValue() string {
	return l.name
}

type linksModal struct {
	list blist.Model
}

func newLinksModal(links []string, rootDir string, width, height int) linksModal {
	items := make([]blist.Item, 0, len(links))
	for _, linkName := range links {
		notePath := findNoteByName(linkName, rootDir)
		items = append(items, linkItem{
			name:   linkName,
			path:   notePath,
			exists: notePath != "",
		})
	}

	l := blist.New(items, blist.NewDefaultDelegate(), width-10, height-10)
	l.Title = "Liens dans cette note"
	l.SetShowHelp(false)

	return linksModal{
		list: l,
	}
}

func (m linksModal) Update(msg tea.Msg) (linksModal, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m linksModal) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("81")).
		Bold(true).
		Render("🔗 Liens Wiki")

	listView := m.list.View()

	helpText := helpStyle.Render("Enter: ouvrir/créer • Esc: fermer")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
		"",
		helpText,
	)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Padding(1, 2).
		Width(70)

	return modalStyle.Render(content)
}
