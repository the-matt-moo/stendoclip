package app

import (
	"runtime"
	"testing"
	"time"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func TestCommandWakesMessagePump(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	a, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	if winapi.IsWindowVisible(a.Window()) {
		t.Fatal("owner window is visible")
	}

	posted := make(chan error, 1)
	go func() {
		time.Sleep(20 * time.Millisecond)
		posted <- a.Quit()
	}()
	if err := a.Run(); err != nil {
		t.Fatal(err)
	}
	if err := <-posted; err != nil {
		t.Fatal(err)
	}
	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	if err := a.Post(func() {}); err == nil {
		t.Fatal("post accepted after close")
	}
}
