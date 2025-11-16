package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type FileClipboard struct {
	path string
	mode string
}

type pasteCompletedMsg struct {
	success bool
	message string
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

func copyDir(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, sourceInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func moveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		if err := copyFile(src, dst); err != nil {
			return err
		}
		return os.Remove(src)
	}
	return nil
}

func moveDir(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		if err := copyDir(src, dst); err != nil {
			return err
		}
		return os.RemoveAll(src)
	}
	return nil
}

func pasteFile(clipboard *FileClipboard, destDir string) tea.Cmd {
	return func() tea.Msg {
		if clipboard == nil {
			return pasteCompletedMsg{
				success: false,
				message: "Clipboard is empty",
			}
		}

		srcInfo, err := os.Stat(clipboard.path)
		if err != nil {
			return pasteCompletedMsg{
				success: false,
				message: fmt.Sprintf("Error: %v", err),
			}
		}

		baseName := filepath.Base(clipboard.path)
		dstPath := filepath.Join(destDir, baseName)

		if _, err := os.Stat(dstPath); err == nil {
			return pasteCompletedMsg{
				success: false,
				message: fmt.Sprintf("File already exists: %s", baseName),
			}
		}

		var opErr error
		if clipboard.mode == "copy" {
			if srcInfo.IsDir() {
				opErr = copyDir(clipboard.path, dstPath)
			} else {
				opErr = copyFile(clipboard.path, dstPath)
			}
		} else if clipboard.mode == "cut" {
			if srcInfo.IsDir() {
				opErr = moveDir(clipboard.path, dstPath)
			} else {
				opErr = moveFile(clipboard.path, dstPath)
			}
		}

		if opErr != nil {
			return pasteCompletedMsg{
				success: false,
				message: fmt.Sprintf("Error: %v", opErr),
			}
		}

		action := "Copied"
		if clipboard.mode == "cut" {
			action = "Moved"
		}

		return pasteCompletedMsg{
			success: true,
			message: fmt.Sprintf("%s: %s", action, baseName),
		}
	}
}
