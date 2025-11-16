package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Editor     string       `json:"editor"`
	Theme      int          `json:"theme"`
	DefaultDir string       `json:"default_dir"`
	Filters    FilterConfig `json:"filters"`
	Search     SearchConfig `json:"search"`
}

type FilterConfig struct {
	MdOnly     bool `json:"md_only"`
	ShowHidden bool `json:"show_hidden"`
	SortMode   int  `json:"sort_mode"`
}

type SearchConfig struct {
	IncludeContent    bool `json:"include_content"`
	RespectGitignore  bool `json:"respect_gitignore"`
	MaxRecentFiles    int  `json:"max_recent_files"`
}

type SessionState struct {
	LastDirectory string   `json:"last_directory"`
	LastTheme     int      `json:"last_theme"`
	RecentFiles   []string `json:"recent_files"`
	Bookmarks     []string `json:"bookmarks"`
}

func getConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "notesmd")
}

func getConfigPath() string {
	return filepath.Join(getConfigDir(), "config.json")
}

func getStatePath() string {
	return filepath.Join(getConfigDir(), "state.json")
}

func ensureConfigDir() error {
	return os.MkdirAll(getConfigDir(), 0755)
}

func LoadConfig() (*Config, error) {
	if err := ensureConfigDir(); err != nil {
		return nil, err
	}

	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func DefaultConfig() *Config {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "notes")

	return &Config{
		Editor:     editor,
		Theme:      0,
		DefaultDir: defaultDir,
		Filters: FilterConfig{
			MdOnly:     false,
			ShowHidden: false,
			SortMode:   0,
		},
		Search: SearchConfig{
			IncludeContent:   true,
			RespectGitignore: true,
			MaxRecentFiles:   10,
		},
	}
}

func (c *Config) Save() error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getConfigPath(), data, 0644)
}

func LoadState() (*SessionState, error) {
	if err := ensureConfigDir(); err != nil {
		return nil, err
	}

	statePath := getStatePath()
	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &SessionState{
				RecentFiles: []string{},
				Bookmarks:   []string{},
			}, nil
		}
		return nil, err
	}

	var state SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func SaveState(s *SessionState) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getStatePath(), data, 0644)
}
