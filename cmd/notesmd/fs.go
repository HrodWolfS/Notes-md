package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/glamour"
)

// readDir reads a directory and returns list items
func readDir(dir string) []blist.Item {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []blist.Item{fileItem{
			name:  fmt.Sprintf("[Erreur: %v]", err),
			path:  dir,
			isDir: false,
		}}
	}

	var res []blist.Item
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		res = append(res, fileItem{
			name:    info.Name(),
			path:    filepath.Join(dir, info.Name()),
			isDir:   info.IsDir(),
			size:    info.Size(),
			modTime: info.ModTime().Unix(),
		})
	}

	return res
}

// loadMarkdown loads and renders a Markdown file
func loadMarkdown(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Erreur de lecture du fichier:\n%s\n\n%v", path, err)
	}

	if filepath.Ext(path) != ".md" {
		return string(data)
	}

	out, err := glamour.Render(string(data), markdownTheme)
	if err != nil {
		return fmt.Sprintf("Erreur de rendu Markdown pour %s:\n%v", path, err)
	}

	return out
}

// loadMarkdownRaw loads raw markdown content without rendering
func loadMarkdownRaw(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// loadMarkdownWithHighlight loads and renders markdown with search highlights
func loadMarkdownWithHighlight(path string, query string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Erreur de lecture du fichier:\n%s\n\n%v", path, err)
	}

	content := string(data)

	if filepath.Ext(path) != ".md" {
		return highlightMatches(content, query)
	}

	// Highlight matches in raw markdown before rendering
	highlighted := highlightMatches(content, query)

	out, err := glamour.Render(highlighted, markdownTheme)
	if err != nil {
		return fmt.Sprintf("Erreur de rendu Markdown pour %s:\n%v", path, err)
	}

	return out
}

// highlightMatches highlights matches with markdown bold and emoji markers
func highlightMatches(content string, query string) string {
	if query == "" {
		return content
	}

	// Case-insensitive search
	queryLower := strings.ToLower(query)
	contentLower := strings.ToLower(content)

	var result strings.Builder
	lastIndex := 0

	for {
		index := strings.Index(contentLower[lastIndex:], queryLower)
		if index == -1 {
			result.WriteString(content[lastIndex:])
			break
		}

		actualIndex := lastIndex + index
		result.WriteString(content[lastIndex:actualIndex])

		// Wrap match with markdown bold and lightning emoji for visibility
		result.WriteString("**⚡ ")
		result.WriteString(content[actualIndex : actualIndex+len(query)])
		result.WriteString(" ⚡**")

		lastIndex = actualIndex + len(query)
	}

	return result.String()
}
