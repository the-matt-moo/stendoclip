package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	MaxHistory    int    `json:"max_history"`
	MaxEntryBytes int    `json:"max_entry_bytes"`
	PasteDelayMs  int    `json:"paste_delay_ms"`
	TimeoutSecs   int    `json:"timeout_secs"`
	Wraparound    bool   `json:"wraparound"`
	DebugLog      bool   `json:"debug_log"`
	HotkeyOpen    string `json:"hotkey_open"`
	HotkeyPin     string `json:"hotkey_pin"`
	HistoryPath   string `json:"history_path"`
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
	}
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

func (c Config) Validate() error {
	var errs []error
	if c.MaxHistory < 1 {
		errs = append(errs, fmt.Errorf("max_history must be at least 1"))
	}
	if c.MaxEntryBytes < 1 {
		errs = append(errs, fmt.Errorf("max_entry_bytes must be at least 1"))
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
