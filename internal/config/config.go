package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Keys struct {
	Open     []string `json:"open"`
	Previous []string `json:"previous"`
	Next     []string `json:"next"`
	Paste    []string `json:"paste"`
	Cancel   []string `json:"cancel"`
	Delete   []string `json:"delete"`
	Pin      []string `json:"pin"`
}

type Config struct {
	MaxHistory         int    `json:"max_history"`
	MaxEntryBytes      int    `json:"max_entry_bytes"`
	PasteDelayMs       int    `json:"paste_delay_ms"`
	TimeoutSecs        int    `json:"timeout_secs"`
	Wraparound         bool   `json:"wraparound"`
	TrimWhitespace     bool   `json:"trim_whitespace"`
	DebugLog           bool   `json:"debug_log"`
	HotkeyOpen         string `json:"hotkey_open"`
	HotkeyPin          string `json:"hotkey_pin"`
	HistoryPath        string `json:"history_path"`
	MarkdownExportPath string `json:"markdown_export_path"`
	BezelFontSize      int    `json:"bezel_font_size"`
	Keys               *Keys  `json:"keys,omitempty"`
}

func DefaultKeys() Keys {
	return Keys{
		Open:     []string{"Ctrl+Shift+V"},
		Previous: []string{"Up", "Left"},
		Next:     []string{"Down", "Right"},
		Paste:    []string{"Enter"},
		Cancel:   []string{"Escape"},
		Delete:   []string{"Delete"},
		Pin:      []string{"Ctrl+P"},
	}
}

func Defaults() Config {
	return Config{
		MaxHistory:    50,
		MaxEntryBytes: 65536,
		PasteDelayMs:  200,
		TimeoutSecs:   5,
		Wraparound:    true,
		HotkeyOpen:    "Ctrl+Shift+V",
		HotkeyPin:     "Ctrl+P",
		BezelFontSize: 18,
	}
}

func (c Config) ResolvedKeys() Keys {
	keys := DefaultKeys()
	if c.Keys == nil {
		// Honour legacy top-level fields when keys block is absent.
		if c.HotkeyOpen != Defaults().HotkeyOpen {
			keys.Open = []string{c.HotkeyOpen}
		}
		if c.HotkeyPin != Defaults().HotkeyPin {
			keys.Pin = []string{c.HotkeyPin}
		}
		return keys
	}
	if len(c.Keys.Open) > 0 {
		keys.Open = c.Keys.Open
	}
	if len(c.Keys.Previous) > 0 {
		keys.Previous = c.Keys.Previous
	}
	if len(c.Keys.Next) > 0 {
		keys.Next = c.Keys.Next
	}
	if len(c.Keys.Paste) > 0 {
		keys.Paste = c.Keys.Paste
	}
	if len(c.Keys.Cancel) > 0 {
		keys.Cancel = c.Keys.Cancel
	}
	if len(c.Keys.Delete) > 0 {
		keys.Delete = c.Keys.Delete
	}
	if len(c.Keys.Pin) > 0 {
		keys.Pin = c.Keys.Pin
	}
	return keys
}

func LoadConfig(path string) (Config, error) {
	cfg := Defaults()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func SaveConfig(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func UpdateFontSize(path string, newSize int) error {
	cfg, err := LoadConfig(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg = Defaults()
	} else if err != nil {
		return err
	}
	if newSize < 12 {
		newSize = 12
	}
	if newSize > 48 {
		newSize = 48
	}
	cfg.BezelFontSize = newSize
	return SaveConfig(path, cfg)
}

func UpdateTrimWhitespace(path string, enabled bool) error {
	cfg, err := LoadConfig(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg = Defaults()
	} else if err != nil {
		return err
	}
	cfg.TrimWhitespace = enabled
	return SaveConfig(path, cfg)
}

func (c Config) Validate() error {
	var errs []error
	if c.MaxHistory < 1 {
		errs = append(errs, fmt.Errorf("max_history must be at least 1"))
	}
	if c.MaxEntryBytes < 1 {
		errs = append(errs, fmt.Errorf("max_entry_bytes must be at least 1"))
	}
	if c.MaxEntryBytes > 10<<20 {
		errs = append(errs, fmt.Errorf("max_entry_bytes cannot exceed 10MB"))
	}
	if c.PasteDelayMs < 0 {
		errs = append(errs, fmt.Errorf("paste_delay_ms cannot be negative"))
	}
	if c.TimeoutSecs < 1 {
		errs = append(errs, fmt.Errorf("timeout_secs must be at least 1"))
	}
	if c.HotkeyOpen == "" {
		errs = append(errs, fmt.Errorf("hotkey_open cannot be empty"))
	}
	if c.HotkeyPin == "" {
		errs = append(errs, fmt.Errorf("hotkey_pin cannot be empty"))
	}
	return errors.Join(errs...)
}
