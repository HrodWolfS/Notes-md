package main

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// Editor messages
type editorDoneMsg struct{}
type editorErrorMsg string

// openInEditor opens a file in the external editor
func openInEditor(path string) tea.Cmd {
	return func() tea.Msg {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nvim"
		}

		cmd := exec.Command(editor, path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return editorErrorMsg(err.Error())
		}

		return editorDoneMsg{}
	}
}
