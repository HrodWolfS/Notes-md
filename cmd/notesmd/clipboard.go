package main

import (
	"os"
)

type clipboardCopiedMsg struct {
	message string
}

func copyFilePath(path string) clipboardCopiedMsg {
	return clipboardCopiedMsg{message: "Path: " + path}
}

func copyFileContent(path string) clipboardCopiedMsg {
	content, err := readFileContent(path)
	if err != nil {
		return clipboardCopiedMsg{message: "Error: " + err.Error()}
	}
	return clipboardCopiedMsg{message: "Content copied (" + string(rune(len(content))) + " bytes)"}
}

func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
