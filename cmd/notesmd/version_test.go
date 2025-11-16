package main

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()
	if v == "" {
		t.Error("Version() should not return empty string")
	}
}

func TestFullVersion(t *testing.T) {
	fv := FullVersion()
	if fv == "" {
		t.Error("FullVersion() should not return empty string")
	}

	// Should contain "NotesMD"
	if !strings.Contains(fv, "NotesMD") {
		t.Error("FullVersion() should contain 'NotesMD'")
	}
}
