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
		got.HotkeyOpen != "Ctrl+Shift+V" || got.HotkeyPin != "Ctrl+P" || got.HistoryPath != "" ||
		got.BezelFontSize != 18 {
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

func TestResolvedKeysUsesDefaultsWhenNoKeysBlock(t *testing.T) {
	cfg := Defaults()
	keys := cfg.ResolvedKeys()
	if keys.Open[0] != "Ctrl+Shift+V" || keys.Pin[0] != "Ctrl+P" || len(keys.Previous) != 2 {
		t.Fatalf("resolved keys = %#v", keys)
	}
}

func TestResolvedKeysHonoursLegacyOverride(t *testing.T) {
	cfg := Defaults()
	cfg.HotkeyOpen = "Alt+V"
	keys := cfg.ResolvedKeys()
	if keys.Open[0] != "Alt+V" {
		t.Fatalf("open = %q", keys.Open[0])
	}
}

func TestResolvedKeysBlockOverridesDefaults(t *testing.T) {
	cfg := Defaults()
	cfg.Keys = &Keys{Cancel: []string{"Q"}}
	keys := cfg.ResolvedKeys()
	if keys.Cancel[0] != "Q" || keys.Open[0] != "Ctrl+Shift+V" {
		t.Fatalf("resolved = %#v", keys)
	}
}

func TestKeysBlockFromJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"keys":{"cancel":["Q"],"previous":["K"]}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	keys := cfg.ResolvedKeys()
	if keys.Cancel[0] != "Q" || keys.Previous[0] != "K" || keys.Open[0] != "Ctrl+Shift+V" {
		t.Fatalf("resolved = %#v", keys)
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
