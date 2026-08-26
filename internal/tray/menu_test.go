package tray

import (
	"strings"
	"testing"
)

func TestPointFromWParamPreservesNegativeCoordinates(t *testing.T) {
	x, y := int16(-120), int16(450)
	value := uintptr(uint16(x)) | uintptr(uint16(y))<<16
	if got := pointFromWParam(value); got.X != -120 || got.Y != 450 {
		t.Fatalf("point = %#v", got)
	}
}

func TestClipLabel(t *testing.T) {
	if got := clipLabel("first\n  second & third"); got != "first second && third" {
		t.Fatalf("label = %q", got)
	}
	if got := clipLabel(strings.Repeat("界", 61)); len([]rune(got)) != 60 || !strings.HasSuffix(got, "…") {
		t.Fatalf("truncated label = %q", got)
	}
	if got := clipLabel(" \n\t "); got != "(Empty clip)" {
		t.Fatalf("empty label = %q", got)
	}
}

func TestStartupValueQuotesPath(t *testing.T) {
	if got := startupValue(`C:\Program Files\Stendoclip\stendoclip.exe`); got != `"C:\Program Files\Stendoclip\stendoclip.exe"` {
		t.Fatalf("startup value = %q", got)
	}
}
