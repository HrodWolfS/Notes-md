package main

import (
	"fmt"
	"os"
	"path/filepath"

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
			name:  info.Name(),
			path:  filepath.Join(dir, info.Name()),
			isDir: info.IsDir(),
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
