package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	name    string
	path    string
	isDir   bool
	size    int64
	modTime int64 // Unix timestamp
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

	// Format size
	var sizeStr string
	if f.size < 1024 {
		sizeStr = fmt.Sprintf("%d B", f.size)
	} else if f.size < 1024*1024 {
		sizeStr = fmt.Sprintf("%.1f KB", float64(f.size)/1024)
	} else {
		sizeStr = fmt.Sprintf("%.1f MB", float64(f.size)/(1024*1024))
	}

	// Format date (relative)
	modTime := time.Unix(f.modTime, 0)
	now := time.Now()
	diff := now.Sub(modTime)

	var timeStr string
	if diff < 24*time.Hour {
		timeStr = "Aujourd'hui"
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		timeStr = fmt.Sprintf("Il y a %d jour(s)", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		timeStr = fmt.Sprintf("Il y a %d semaine(s)", weeks)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		timeStr = fmt.Sprintf("Il y a %d mois", months)
	} else {
		timeStr = modTime.Format("02/01/2006")
	}

	return fmt.Sprintf("%s • %s", sizeStr, timeStr)
}

func (f fileItem) FilterValue() string {
	return f.name
}

// model represents the application state
type model struct {
	mode       viewMode
	rootDir    string
	currentDir string

	list              blist.Model
	baseItems         []blist.Item
	allFiles          []fileItem
	viewport          bviewport.Model
	showPreview       bool
	autoPreview       bool // Auto-preview on selection change
	lastSelectedIndex int  // Track selection changes

	searchActive bool
	searchQuery  string

	// filters
	mdOnly     bool // Show only .md files
	showHidden bool // Show hidden files (starting with .)
	sortMode   int  // 0=name, 1=date, 2=size

	// window size
	width  int
	height int

	// theme
	themeIndex int

	// modals
	showNoteModal      bool
	noteModal          noteModal
	showConfirmModal   bool
	confirmModal       confirmModal
	showRenameModal    bool
	renameModal        renameModal
	showHelpModal      bool
	helpModal          helpModal
	showCreateDirModal bool
	createDirModal     createDirModal
	showRecentModal    bool
	recentModal        recentFilesModal
	showBookmarksModal bool
	bookmarksModal     bookmarksModal
	showLinksModal     bool
	linksModal         linksModal

	// configuration & persistence
	config      *Config
	recentFiles []string
	bookmarks   []string

	// file operations
	clipboard     *FileClipboard
	clipboardMode string

	// status bar
	statusBar StatusBar

	// in-note search
	searchInNoteActive bool
	noteSearchQuery    string
	currentNoteRaw     string // Raw markdown content of current note
	currentNotePath    string

	// vim-style navigation
	lastKey       string
	pendingDelete bool     // for 'dd' double-tap
	navHistory    []string // Navigation history stack
	navIndex      int      // Current position in history
}

// setDir changes the current directory and updates the list
func (m *model) setDir(dir string) {
	// Add to navigation history
	if len(m.navHistory) == 0 {
		// First directory in history
		m.navHistory = append(m.navHistory, dir)
		m.navIndex = 0
	} else if m.navHistory[m.navIndex] != dir {
		// Truncate forward history if we're in the middle
		m.navHistory = m.navHistory[:m.navIndex+1]
		m.navHistory = append(m.navHistory, dir)
		m.navIndex = len(m.navHistory) - 1
	}

	m.currentDir = dir
	m.baseItems = readDir(dir)
	m.applyFilters()
}

// navBack navigates to previous directory in history
func (m *model) navBack() bool {
	if m.navIndex > 0 {
		m.navIndex--
		m.currentDir = m.navHistory[m.navIndex]
		m.baseItems = readDir(m.currentDir)
		m.applyFilters()
		return true
	}
	return false
}

// navForward navigates to next directory in history
func (m *model) navForward() bool {
	if m.navIndex < len(m.navHistory)-1 {
		m.navIndex++
		m.currentDir = m.navHistory[m.navIndex]
		m.baseItems = readDir(m.currentDir)
		m.applyFilters()
		return true
	}
	return false
}

// applyFilters applies active filters to the file list
func (m *model) applyFilters() {
	var filtered []blist.Item

	// Apply filters
	for _, item := range m.baseItems {
		if fi, ok := item.(fileItem); ok {
			// Filter hidden files
			if !m.showHidden && len(fi.name) > 0 && fi.name[0] == '.' {
				continue
			}

			// Filter .md only
			if m.mdOnly && !fi.isDir && filepath.Ext(fi.name) != ".md" {
				continue
			}

			filtered = append(filtered, item)
		}
	}

	// Apply sorting
	m.sortItems(filtered)

	m.list.SetItems(filtered)
}

// sortItems sorts items based on sortMode
func (m *model) sortItems(items []blist.Item) {
	switch m.sortMode {
	case 0: // Name (alphabetical)
		// Already sorted by name from readDir
	case 1: // Date (newest first)
		for i := 0; i < len(items); i++ {
			for j := i + 1; j < len(items); j++ {
				fi := items[i].(fileItem)
				fj := items[j].(fileItem)
				if fj.modTime > fi.modTime {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
	case 2: // Size (largest first)
		for i := 0; i < len(items); i++ {
			for j := i + 1; j < len(items); j++ {
				fi := items[i].(fileItem)
				fj := items[j].(fileItem)
				if fj.size > fi.size {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
	}
}

// ensureAllFilesScanned scans all files AND folders from rootDir if not already done
func (m *model) ensureAllFilesScanned() {
	if len(m.allFiles) > 0 {
		return
	}

	var files []fileItem
	filepath.WalkDir(m.rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// Skip root directory itself
		if path == m.rootDir {
			return nil
		}

		// Include both files and directories
		info, err := d.Info()
		if err != nil {
			return nil
		}

		files = append(files, fileItem{
			name:    d.Name(),
			path:    path,
			isDir:   d.IsDir(),
			size:    info.Size(),
			modTime: info.ModTime().Unix(),
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

// toggleBookmark adds or removes a bookmark
func (m *model) toggleBookmark(path string) bool {
	for i, bookmark := range m.bookmarks {
		if bookmark == path {
			m.bookmarks = append(m.bookmarks[:i], m.bookmarks[i+1:]...)
			return false
		}
	}
	m.bookmarks = append(m.bookmarks, path)
	return true
}

// isBookmarked checks if a path is bookmarked
func (m *model) isBookmarked(path string) bool {
	for _, bookmark := range m.bookmarks {
		if bookmark == path {
			return true
		}
	}
	return false
}

// trackRecentFile adds a file to recent files list
func (m *model) trackRecentFile(path string) {
	for i, recent := range m.recentFiles {
		if recent == path {
			m.recentFiles = append(m.recentFiles[:i], m.recentFiles[i+1:]...)
			break
		}
	}

	m.recentFiles = append([]string{path}, m.recentFiles...)

	maxRecent := 10
	if m.config != nil && m.config.Search.MaxRecentFiles > 0 {
		maxRecent = m.config.Search.MaxRecentFiles
	}
	if len(m.recentFiles) > maxRecent {
		m.recentFiles = m.recentFiles[:maxRecent]
	}
}

// updateStatusBar updates the status bar with current state
func (m *model) updateStatusBar() {
	m.statusBar.SetPath(m.currentDir)

	fileCount := 0
	dirCount := 0
	for _, item := range m.baseItems {
		if fi, ok := item.(fileItem); ok {
			if fi.isDir {
				dirCount++
			} else {
				fileCount++
			}
		}
	}
	m.statusBar.SetCounts(fileCount, dirCount)

	mode := "Browser"
	if m.searchActive {
		mode = "Search"
	} else if m.searchInNoteActive {
		mode = "Search in Note"
	}
	m.statusBar.SetMode(mode)

	var filters []string
	if m.mdOnly {
		filters = append(filters, "[.md only]")
	}
	if m.showHidden {
		filters = append(filters, "[hidden]")
	}
	if m.sortMode == 1 {
		filters = append(filters, "[↓ date]")
	} else if m.sortMode == 2 {
		filters = append(filters, "[↓ size]")
	}
	m.statusBar.SetFilters(filters)
}
