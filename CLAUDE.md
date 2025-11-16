# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`notesmd` is a terminal-based Markdown notes browser built with Go and the Bubble Tea TUI framework. It provides an interactive file explorer with live Markdown preview and fuzzy search capabilities.

## Building and Running

**Build the application:**
```bash
go build -o nmd ./cmd/notesmd
```

**Run directly:**
```bash
go run ./cmd/notesmd [directory]
```

**Run the compiled binary:**
```bash
./nmd [directory]
```

If no directory is provided, it defaults to the current directory.

## Development Commands

**Install dependencies:**
```bash
go mod download
```

**Update dependencies:**
```bash
go mod tidy
```

**Run with specific directory:**
```bash
go run ./cmd/notesmd ~/Documents/notes
```

## Architecture

### Modular Architecture

This application uses a clean modular architecture with code organized by responsibility. All source files are in `cmd/notesmd/`.

### Core Components

**Bubble Tea Model (`model` struct)**
- Manages application state including view mode, directory navigation, search state, and theme
- Contains two main view modes: `modeHome` (welcome screen) and `modeBrowser` (file explorer)
- Integrates three Bubble Tea components:
  - `list.Model` for file/directory listing
  - `viewport.Model` for Markdown preview with scrolling
  - Custom search and note creation states

**View System**
- `viewHome()`: Welcome screen with ASCII art title and keyboard shortcuts
- `viewBrowser()`: Split-pane layout (30% file list, 70% preview) with responsive sizing

**State Management**
- `modeHome` / `modeBrowser`: Top-level view states
- `searchActive`: Global fuzzy search mode
- `creatingNote`: Quick note creation workflow
- `showPreview`: Toggle for Markdown rendering
- `themeIndex`: Cycling through color palettes

### Key Architectural Patterns

**Responsive Layout Calculation**
- Dynamic width/height allocation based on `tea.WindowSizeMsg`
- 30/70 split between list and preview with border compensation
- Fallback dimensions (120x40) when size unavailable

**Fuzzy Search Implementation**
- On-demand file scanning with `ensureAllFilesScanned()`
- Uses `github.com/sahilm/fuzzy` for matching
- Scans entire directory tree from `rootDir`, filters by relative paths

**External Editor Integration**
- Spawns `$EDITOR` (defaults to `nvim`) for file editing
- Uses `editorDoneMsg` / `editorErrorMsg` for async editor completion
- Refreshes preview after editor closes

**Markdown Rendering**
- `github.com/charmbracelet/glamour` for terminal-based Markdown rendering
- Uses "dark" theme constant
- Falls back to raw content for non-Markdown files

### Navigation Flow

```
Start → Home Screen (modeHome)
  ↓ [ENTER]
Browser Mode (modeBrowser)
  ├─ [→/l/Enter] → Enter directory
  ├─ [←/h] → Parent directory
  ├─ [o] → Toggle preview
  ├─ [e] → Open in $EDITOR
  ├─ [n] → Create new note
  └─ [/] → Global fuzzy search
```

### File Structure

**Current directory structure:**
```
notesmd/
├── cmd/notesmd/
│   ├── main.go          # Entry point + Init/View (69 lines)
│   ├── model.go         # Model struct + types (149 lines)
│   ├── update.go        # Update logic (213 lines)
│   ├── notes.go         # Note creation handler (48 lines)
│   ├── modal.go         # Note creation modal (172 lines)
│   ├── view_home.go     # Home screen (59 lines)
│   ├── view_browser.go  # Browser view with modal overlay (119 lines)
│   ├── theme.go         # Styles + ASCII (59 lines)
│   ├── fs.go           # File operations (56 lines)
│   └── editor.go       # External editor (33 lines)
├── internal/            # Reserved for future use
│   ├── fs/
│   ├── preview/
│   └── ui/
├── go.mod
├── go.sum
└── nmd                  # Compiled binary (977 lines total)
```

**Module Organization:**
- `main.go`: Program entry point, initialization, and view delegation
- `model.go`: Application state, data structures, and model methods
- `update.go`: All state update logic (Update, updateHome, updateBrowser, handleSearchKey)
- `notes.go`: Note creation logic (handleModalKey)
- `modal.go`: Reusable note creation modal component with textinput/textarea
- `view_home.go`: Welcome screen rendering
- `view_browser.go`: File browser rendering with split layout and modal overlay
- `theme.go`: Lipgloss styles, color palette, and ASCII art
- `fs.go`: File system operations (readDir, loadMarkdown)
- `editor.go`: External editor integration

## Key Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - Reusable TUI components (list, viewport)
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/charmbracelet/glamour` - Markdown rendering
- `github.com/sahilm/fuzzy` - Fuzzy search matching

## Development Notes

**State Mutation Patterns**
- Most model updates return `(tea.Model, tea.Cmd)` following Bubble Tea conventions
- Helper methods like `setDir()` directly mutate model state (not pure functions)
- Search and note creation use specialized key handlers (`handleSearchKey`, `handleCreateNoteKey`)

**Theming System**
- `titlePalette` array defines 5 color schemes: ["208", "196", "46", "51", "201"]
- Theme cycling via `toggleTheme()` method
- Dynamic style application using `lipgloss.Color()` at render time

**UI Component Lifecycle**
- `Init()` enters alternate screen mode
- `Update()` handles resize, editor completion, and key events
- `View()` delegates to mode-specific view functions
- Components sized via `SetWidth/SetHeight` on window resize

**File Operations**
- Directory reading via `os.ReadDir()` wrapped in `readDir()` helper
- Markdown loading with `glamour.Render()` in `loadMarkdown()`
- Note creation writes to `currentDir` with auto-generated title header

**Note Creation Modal (Professional Implementation)**
1. User presses `n` → Modal opens centered on screen
2. Modal contains:
   - `textinput` component for note name
   - `textarea` component for multi-line Markdown content
   - Focused field switching with `TAB`
3. Keybindings:
   - `TAB` → Switch between name and content fields
   - `Enter` in textarea → New line (not submit)
   - `Ctrl+S` or `Ctrl+Enter` → Save note and close modal
   - `Esc` → Cancel and close modal
4. On save:
   - Create file with `.md` extension
   - Add `# {title}` header automatically
   - Refresh directory list and auto-select new note

## Roadmap

### ✅ Completed: Code Refactoring
Monolith `main.go` (770 lines) has been successfully restructured into 9 modular files (840 lines total):
- Improved readability and maintainability
- Professional structure ready for GitHub publication
- Easier collaboration and code review
- Better showcase for portfolio/recruiters

### ✅ Completed: Full Note Creation Modal
Professional centered modal implementation (172 lines in `modal.go`):
- Centered modal with border and padding
- Two interactive components:
  - `textinput` for note name (with placeholder and char limit)
  - `textarea` for multi-line Markdown content (5000 char limit)
- Smart field switching with `TAB`
- `Ctrl+S`/`Ctrl+Enter` to save
- `Esc` to cancel
- Centered overlay rendering with lipgloss.Place
- Auto-refresh and selection after creation

**Project Status:** Core features complete and ready for open-source publication! 🎉
