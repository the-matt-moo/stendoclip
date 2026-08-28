package tray

import (
	"runtime"
	"testing"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func TestAboutCanBeReopened(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	instance, err := winapi.GetModuleHandle()
	if err != nil {
		t.Fatal(err)
	}
	for attempt := 1; attempt <= 2; attempt++ {
		showAbout(instance, nil, "test")
		activeAboutMu.Lock()
		about := activeAbout
		activeAboutMu.Unlock()
		if about == nil {
			t.Fatalf("attempt %d did not create the About window", attempt)
		}
		if err := winapi.DestroyWindow(about.hwnd); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLayoutAboutFitsClient(t *testing.T) {
	for _, client := range []winapi.Rect{
		{Right: 424, Bottom: 401},
		{Right: 228, Bottom: 210},
	} {
		image, text := layoutAbout(client, 100, true)
		if image.Left < 0 || image.Top < 0 || image.Right > client.Right || image.Bottom > client.Bottom {
			t.Fatalf("image %v outside client %v", image, client)
		}
		if text.Left < 0 || text.Top < 0 || text.Right > client.Right || text.Bottom > client.Bottom {
			t.Fatalf("text %v outside client %v", text, client)
		}
		if text.Top < image.Bottom {
			t.Fatalf("text starts before image ends: image=%v text=%v", image, text)
		}
	}
}
