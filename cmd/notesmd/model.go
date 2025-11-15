package main

import (
	"os"
	"path/filepath"

	blist "github.com/charmbracelet/bubbles/list"
	bviewport "github.com/charmbracelet/bubbles/viewport"
	fuzzy "github.com/sahilm/fuzzy"
)

// viewMode represents the current view state
type viewMode int

const (
	modeHome viewMode = iota
	modeBrowser
)

// fileItem represents a file or directory in the list
type fileItem struct {
	name  string
	path  string
	isDir bool
}

func (f fileItem) Title() string {
	if f.isDir {
		return "📁 " + f.name + "/"
	}

	ext := filepath.Ext(f.name)
	if ext == ".md" {
		return "📝 " + f.name
	}

	return "📄 " + f.name
}

func (f fileItem) Description() string {
	if f.isDir {
		return "Dossier"
	}
	return f.path
}

func (f fileItem) FilterValue() string {
	return f.name
}

// model represents the application state
type model struct {
	mode       viewMode
	rootDir    string
	currentDir string

	list        blist.Model
	baseItems   []blist.Item
	allFiles    []fileItem
	viewport    bviewport.Model
	showPreview bool

	searchActive bool
	searchQuery  string

	// filters
	mdOnly bool // Show only .md files

	// window size
	width  int
	height int

	// theme
	themeIndex int

	// modals
	showNoteModal    bool
	noteModal        noteModal
	showConfirmModal bool
	confirmModal     confirmModal
	showRenameModal  bool
	renameModal      renameModal
	showHelpModal    bool
	helpModal        helpModal

	// vim-style navigation
	lastKey       string
	pendingDelete bool // for 'dd' double-tap
}

// setDir changes the current directory and updates the list
func (m *model) setDir(dir string) {
	m.currentDir = dir
	m.baseItems = readDir(dir)
	m.applyFilters()
}

// applyFilters applies active filters to the file list
func (m *model) applyFilters() {
	if !m.mdOnly {
		m.list.SetItems(m.baseItems)
		return
	}

	// Filter to show only .md files and directories
	var filtered []blist.Item
	for _, item := range m.baseItems {
		if fi, ok := item.(fileItem); ok {
			if fi.isDir || filepath.Ext(fi.name) == ".md" {
				filtered = append(filtered, item)
			}
		}
	}
	m.list.SetItems(filtered)
}

// ensureAllFilesScanned scans all files from rootDir if not already done
func (m *model) ensureAllFilesScanned() {
	if len(m.allFiles) > 0 {
		return
	}

	var files []fileItem
	filepath.WalkDir(m.rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, fileItem{
			name:  d.Name(),
			path:  path,
			isDir: false,
		})
		return nil
	})

	m.allFiles = files
}

// buildSearchResults filters files based on search query
func (m *model) buildSearchResults() {
	if m.searchQuery == "" {
		items := make([]blist.Item, 0, len(m.allFiles))
		for _, f := range m.allFiles {
			items = append(items, f)
		}
		m.list.SetItems(items)
		return
	}

	var candidates []string
	for _, f := range m.allFiles {
		rel, err := filepath.Rel(m.rootDir, f.path)
		if err != nil {
			rel = f.name
		}
		candidates = append(candidates, rel)
	}

	matches := fuzzy.Find(m.searchQuery, candidates)
	if len(matches) == 0 {
		m.list.SetItems([]blist.Item{})
		return
	}

	var filtered []blist.Item
	for _, match := range matches {
		filtered = append(filtered, m.allFiles[match.Index])
	}
	m.list.SetItems(filtered)
}

// toggleTheme cycles through available theme colors
func (m *model) toggleTheme() {
	if len(titlePalette) == 0 {
		return
	}
	m.themeIndex = (m.themeIndex + 1) % len(titlePalette)
}
