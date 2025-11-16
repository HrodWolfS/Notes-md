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
	config, err := LoadConfig()
	if err != nil {
		config = DefaultConfig()
	}

	state, err := LoadState()
	if err != nil {
		state = &SessionState{
			RecentFiles: []string{},
			Bookmarks:   []string{},
		}
	}

	startDir := "."
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	} else if state.LastDirectory != "" {
		startDir = state.LastDirectory
	} else if config.DefaultDir != "" {
		startDir = config.DefaultDir
	}

	absDir, err := filepath.Abs(startDir)
	if err != nil {
		fmt.Println("Erreur de chemin :", err)
		os.Exit(1)
	}

	m := initialModel(absDir, config, state)

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}

	if finalModel, ok := finalModel.(model); ok {
		saveState := &SessionState{
			LastDirectory: finalModel.currentDir,
			LastTheme:     finalModel.themeIndex,
			RecentFiles:   finalModel.recentFiles,
			Bookmarks:     finalModel.bookmarks,
		}
		SaveState(saveState)
	}
}

// initialModel creates and returns the initial application model
func initialModel(absDir string, config *Config, state *SessionState) model {
	items := []blist.Item{}

	l := blist.New(items, blist.NewDefaultDelegate(), 0, 0)
	l.Title = "Fichiers"
	l.SetShowHelp(false)

	vp := bviewport.New(0, 0)
	vp.SetContent("")

	themeIndex := 0
	if state.LastTheme >= 0 && state.LastTheme < len(titlePalette) {
		themeIndex = state.LastTheme
	} else if config.Theme >= 0 && config.Theme < len(titlePalette) {
		themeIndex = config.Theme
	}

	recentFiles := state.RecentFiles
	if recentFiles == nil {
		recentFiles = []string{}
	}

	bookmarks := state.Bookmarks
	if bookmarks == nil {
		bookmarks = []string{}
	}

	return model{
		mode:              modeHome,
		rootDir:           absDir,
		currentDir:        absDir,
		list:              l,
		baseItems:         items,
		viewport:          vp,
		autoPreview:       true,
		lastSelectedIndex: -1,
		config:            config,
		recentFiles:       recentFiles,
		bookmarks:         bookmarks,
		statusBar:         NewStatusBar(),
		themeIndex:        themeIndex,
		mdOnly:            config.Filters.MdOnly,
		showHidden:        config.Filters.ShowHidden,
		sortMode:          config.Filters.SortMode,
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
