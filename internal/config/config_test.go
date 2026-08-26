package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	got := Defaults()
	if got.MaxHistory != 50 || got.MaxEntryBytes != 65536 || got.PasteDelayMs != 200 ||
		got.TimeoutSecs != 5 || !got.Wraparound || got.DebugLog ||
		got.HotkeyOpen != "Ctrl+Shift+V" || got.HotkeyPin != "Ctrl+P" || got.HistoryPath != "" {
		t.Fatalf("unexpected defaults: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("defaults invalid: %v", err)
	}
}

func TestPartialJSONMergesDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"max_history":75,"debug_log":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.MaxHistory != 75 || !got.DebugLog || got.MaxEntryBytes != 65536 || got.HotkeyOpen != "Ctrl+Shift+V" || !got.Wraparound {
		t.Fatalf("partial merge = %#v", got)
	}
}

func TestInvalidValuesRejected(t *testing.T) {
	tests := []Config{
		func() Config { c := Defaults(); c.MaxHistory = 0; return c }(),
		func() Config { c := Defaults(); c.MaxEntryBytes = 0; return c }(),
		func() Config { c := Defaults(); c.PasteDelayMs = -1; return c }(),
		func() Config { c := Defaults(); c.TimeoutSecs = 0; return c }(),
		func() Config { c := Defaults(); c.HotkeyOpen = ""; return c }(),
		func() Config { c := Defaults(); c.HotkeyPin = ""; return c }(),
	}
	for i, cfg := range tests {
		if err := cfg.Validate(); err == nil {
			t.Errorf("case %d accepted invalid config: %#v", i, cfg)
		}
	}
}
