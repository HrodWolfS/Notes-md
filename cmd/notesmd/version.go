package main

import "fmt"

// Version information (set via ldflags during build)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "manual"
)

// Version returns the version string
func Version() string {
	return version
}

// FullVersion returns detailed version information
func FullVersion() string {
	return fmt.Sprintf("NotesMD %s\nCommit: %s\nBuilt: %s\nBy: %s",
		version, commit, date, builtBy)
}
