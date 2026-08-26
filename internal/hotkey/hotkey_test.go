package hotkey

import (
	"testing"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func TestParse(t *testing.T) {
	got, err := Parse("Ctrl+Shift+V")
	if err != nil {
		t.Fatal(err)
	}
	if got.Modifiers != winapi.MODControl|winapi.MODShift || got.Key != 'V' {
		t.Fatalf("binding = %#v", got)
	}

	got, err = Parse("Alt+F12")
	if err != nil || got.Modifiers != winapi.MODAlt || got.Key != winapi.VKF12 {
		t.Fatalf("F12 binding = %#v, %v", got, err)
	}
}

func TestParseRejectsInvalid(t *testing.T) {
	for _, spec := range []string{"Ctrl", "Ctrl+V+X", "Ctrl+Ctrl+V", "Ctrl+F13", "Ctrl+"} {
		if _, err := Parse(spec); err == nil {
			t.Errorf("Parse(%q) succeeded", spec)
		}
	}
}
