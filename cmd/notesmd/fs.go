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

// loadMarkdown loads and renders a Markdown file with word wrap
func loadMarkdown(path string, width int) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Erreur de lecture du fichier:\n%s\n\n%v", path, err)
	}

	if filepath.Ext(path) != ".md" {
		return string(data)
	}

	// Create renderer with word wrap at specified width
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(markdownTheme),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return fmt.Sprintf("Erreur de création du renderer Markdown:\n%v", err)
	}

	out, err := renderer.Render(string(data))
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
func loadMarkdownWithHighlight(path string, query string, width int) string {
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

	// Create renderer with word wrap
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(markdownTheme),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return fmt.Sprintf("Erreur de création du renderer Markdown:\n%v", err)
	}

	out, err := renderer.Render(highlighted)
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

// parseWikiLinks extracts all wiki-style links [[...]] from content
func parseWikiLinks(content string) []string {
	var links []string
	seen := make(map[string]bool)

	// Find all [[...]] patterns
	start := 0
	for {
		openIdx := strings.Index(content[start:], "[[")
		if openIdx == -1 {
			break
		}
		openIdx += start

		closeIdx := strings.Index(content[openIdx:], "]]")
		if closeIdx == -1 {
			break
		}
		closeIdx += openIdx

		// Extract link text
		linkText := content[openIdx+2 : closeIdx]
		linkText = strings.TrimSpace(linkText)

		// Add to list if not already seen
		if linkText != "" && !seen[linkText] {
			links = append(links, linkText)
			seen[linkText] = true
		}

		start = closeIdx + 2
	}

	return links
}

// findNoteByName searches for a note file by name in rootDir and subdirectories
func findNoteByName(name string, rootDir string) string {
	// Add .md extension if not present
	if filepath.Ext(name) == "" {
		name += ".md"
	}

	var foundPath string

	// Walk through all files
	filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// Check if filename matches (case-insensitive)
		if strings.EqualFold(filepath.Base(path), name) {
			foundPath = path
			return filepath.SkipAll // Stop searching once found
		}

		return nil
	})

	return foundPath
}

// convertWikiLinks converts [[Note]] to markdown links [Note](path)
func convertWikiLinks(content string, rootDir string) string {
	result := content

	// Find all [[...]] patterns and replace them
	start := 0
	for {
		openIdx := strings.Index(result[start:], "[[")
		if openIdx == -1 {
			break
		}
		openIdx += start

		closeIdx := strings.Index(result[openIdx:], "]]")
		if closeIdx == -1 {
			break
		}
		closeIdx += openIdx

		// Extract link text
		linkText := result[openIdx+2 : closeIdx]
		linkText = strings.TrimSpace(linkText)

		if linkText != "" {
			// Find the actual file
			notePath := findNoteByName(linkText, rootDir)

			// Create markdown link with special marker for styling
			var replacement string
			if notePath != "" {
				// Use emoji to make it stand out
				replacement = fmt.Sprintf("🔗 [**%s**](%s)", linkText, notePath)
			} else {
				// Non-existent note - different styling
				replacement = fmt.Sprintf("🔗 [**%s**](#missing)", linkText)
			}

			// Replace [[...]] with markdown link
			result = result[:openIdx] + replacement + result[closeIdx+2:]
			start = openIdx + len(replacement)
		} else {
			start = closeIdx + 2
		}
	}

	return result
}

// loadMarkdownWithLinks loads markdown and converts wiki-style links
func loadMarkdownWithLinks(path string, rootDir string, width int) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Erreur de lecture du fichier:\n%s\n\n%v", path, err)
	}

	content := string(data)

	if filepath.Ext(path) != ".md" {
		return content
	}

	// Convert wiki links before rendering
	contentWithLinks := convertWikiLinks(content, rootDir)

	// Create renderer with word wrap
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(markdownTheme),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return fmt.Sprintf("Erreur de création du renderer Markdown:\n%v", err)
	}

	out, err := renderer.Render(contentWithLinks)
	if err != nil {
		return fmt.Sprintf("Erreur de rendu Markdown pour %s:\n%v", path, err)
	}

	return out
}
