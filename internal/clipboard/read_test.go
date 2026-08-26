package clipboard

import (
	"errors"
	"testing"
	"unicode/utf16"
)

func TestDecodeText(t *testing.T) {
	units := append(utf16.Encode([]rune("hello 世界")), 0, 'x')
	text, err := decodeText(units, 64)
	if err != nil {
		t.Fatal(err)
	}
	if text != "hello 世界" {
		t.Fatalf("got %q", text)
	}
	if _, err := decodeText(units, 4); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("expected ErrTooLarge, got %v", err)
	}
	if _, err := decodeText([]uint16{0}, 64); !errors.Is(err, ErrNoText) {
		t.Fatalf("expected ErrNoText, got %v", err)
	}
}
