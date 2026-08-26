package overlay

import (
	"testing"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func TestCenteredOnNegativeMonitor(t *testing.T) {
	x, y, width, height := centered(winapi.Rect{Left: -1920, Top: -200, Right: 0, Bottom: 880})
	if x != -1280 || y != 245 || width != 640 || height != 190 {
		t.Fatalf("position = %d,%d %dx%d", x, y, width, height)
	}

	x, y, width, height = centered(winapi.Rect{Left: 0, Top: 0, Right: 500, Bottom: 180})
	if x != 20 || y != 20 || width != 460 || height != 140 {
		t.Fatalf("small work area position = %d,%d %dx%d", x, y, width, height)
	}
}

func TestNextIndex(t *testing.T) {
	tests := []struct {
		current, direction, length int
		wrap                       bool
		want                       int
	}{
		{0, 1, 3, true, 1},
		{2, 1, 3, true, 0},
		{0, -1, 3, true, 2},
		{2, 1, 3, false, 2},
		{0, -1, 3, false, 0},
		{0, 1, 0, true, -1},
	}
	for _, test := range tests {
		if got := nextIndex(test.current, test.direction, test.length, test.wrap); got != test.want {
			t.Errorf("nextIndex(%d, %d, %d, %v) = %d, want %d", test.current, test.direction, test.length, test.wrap, got, test.want)
		}
	}
}
