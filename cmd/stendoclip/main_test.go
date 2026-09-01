package main

import (
	"strings"
	"testing"

	"github.com/mooreceipts/stendoclip/internal/config"
)

func TestBuildAboutTextIncludesKeybindings(t *testing.T) {
	text := buildAboutText("1.2.3", config.Keys{
		Open:     []string{"Ctrl+Shift+V"},
		Previous: []string{"Up", "Left"},
		Next:     []string{"Down", "Right"},
		Paste:    []string{"Enter"},
		Cancel:   []string{"Escape"},
		Delete:   []string{"Delete"},
		Pin:      []string{"Ctrl+P"},
	})

	for _, want := range []string{
		"Stendoclip 1.2.3",
		"Open bezel: Ctrl+Shift+V",
		"Previous clip: Up / Left",
		"Pin clip: Ctrl+P",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %q", want, text)
		}
	}
}
