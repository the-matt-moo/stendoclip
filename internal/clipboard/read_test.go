package clipboard

import (
	"errors"
	"testing"
	"unicode/utf16"
)

func TestDecodeText(t *testing.T) {
	units := append(utf16.Encode([]rune("hello 世界")), 0, 'x')
	text, err := decodeText(units, 64, false)
	if err != nil {
		t.Fatal(err)
	}
	if text != "hello 世界" {
		t.Fatalf("got %q", text)
	}
	if _, err := decodeText(units, 4, false); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("expected ErrTooLarge, got %v", err)
	}
	if _, err := decodeText([]uint16{0}, 64, false); !errors.Is(err, ErrNoText) {
		t.Fatalf("expected ErrNoText, got %v", err)
	}
	if _, err := decodeText(utf16.Encode([]rune(" \n\t ")), 64, false); !errors.Is(err, ErrNoText) {
		t.Fatalf("expected blank text to be ignored, got %v", err)
	}
	if text, err := decodeText(utf16.Encode([]rune("  a  ")), 3, true); err != nil || text != "a" {
		t.Fatalf("trimmed text = %q, err %v", text, err)
	}
	if _, err := decodeText(utf16.Encode([]rune("  a  ")), 3, false); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("expected untrimmed text to exceed limit, got %v", err)
	}
}
