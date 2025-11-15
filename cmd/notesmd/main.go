package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	blist "github.com/charmbracelet/bubbles/list"
	bviewport "github.com/charmbracelet/bubbles/viewport"
)

func main() {
	startDir := "."
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	}

	absDir, err := filepath.Abs(startDir)
	if err != nil {
		fmt.Println("Erreur de chemin :", err)
		os.Exit(1)
	}

	m := initialModel(absDir)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}

// initialModel creates and returns the initial application model
func initialModel(absDir string) model {
	items := []blist.Item{}

	l := blist.New(items, blist.NewDefaultDelegate(), 0, 0)
	l.Title = "Fichiers"
	l.SetShowHelp(false)

	vp := bviewport.New(0, 0)
	vp.SetContent("Appuie sur 'o' pour prévisualiser un fichier Markdown.")

	return model{
		mode:       modeHome,
		rootDir:    absDir,
		currentDir: absDir,
		list:       l,
		baseItems:  items,
		viewport:   vp,
	}
}

// Init initializes the Bubble Tea program
func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// View delegates to the appropriate view function based on mode
func (m model) View() string {
	switch m.mode {
	case modeHome:
		return m.viewHome()
	case modeBrowser:
		return m.viewBrowser()
	default:
		return ""
	}
}
