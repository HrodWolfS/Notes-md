package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	width       int
	height      int
	currentPath string
	fileCount   int
	dirCount    int
	mode        string
	message     string
	messageTime time.Time
	filters     []string
}

type clearMessageMsg struct{}

func NewStatusBar() StatusBar {
	return StatusBar{
		width:   0,
		height:  1,
		mode:    "Browser",
		filters: []string{},
	}
}

func (sb StatusBar) Init() tea.Cmd {
	return nil
}

func (sb StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sb.width = msg.Width
		sb.height = 1

	case clearMessageMsg:
		sb.message = ""
	}

	return sb, nil
}

func (sb StatusBar) View() string {
	if sb.width == 0 {
		return ""
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235")).
		Width(sb.width).
		Padding(0, 1)

	var parts []string

	baseName := filepath.Base(sb.currentPath)
	if baseName == "." || baseName == "" {
		baseName = sb.currentPath
	}
	parts = append(parts, fmt.Sprintf("📁 %s", baseName))

	if sb.fileCount > 0 || sb.dirCount > 0 {
		parts = append(parts, fmt.Sprintf("%d files, %d dirs", sb.fileCount, sb.dirCount))
	}

	parts = append(parts, sb.mode)

	if len(sb.filters) > 0 {
		parts = append(parts, strings.Join(sb.filters, " "))
	}

	if sb.message != "" {
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)
		parts = append(parts, messageStyle.Render(sb.message))
	}

	content := strings.Join(parts, "  |  ")

	return statusStyle.Render(content)
}

func (sb *StatusBar) SetDimensions(w, h int) {
	sb.width = w
	sb.height = h
}

func (sb *StatusBar) SetPath(path string) {
	sb.currentPath = path
}

func (sb *StatusBar) SetCounts(files, dirs int) {
	sb.fileCount = files
	sb.dirCount = dirs
}

func (sb *StatusBar) SetMode(mode string) {
	sb.mode = mode
}

func (sb *StatusBar) SetFilters(filters []string) {
	sb.filters = filters
}

func (sb *StatusBar) SetMessage(msg string, duration time.Duration) tea.Cmd {
	sb.message = msg
	sb.messageTime = time.Now()

	if duration > 0 {
		return tea.Tick(duration, func(t time.Time) tea.Msg {
			return clearMessageMsg{}
		})
	}

	return nil
}
