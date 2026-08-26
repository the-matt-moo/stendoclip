package paste

import (
	"testing"
	"unsafe"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func TestInputs(t *testing.T) {
	got := inputs()
	want := []winapi.KeyboardInput{
		{VK: winapi.VKControl},
		{VK: 'V'},
		{VK: 'V', Flags: winapi.KeyEventKeyUp},
		{VK: winapi.VKControl, Flags: winapi.KeyEventKeyUp},
	}
	if len(got) != len(want) {
		t.Fatalf("input count = %d", len(got))
	}
	for i := range got {
		keyboard := *(*winapi.KeyboardInput)(unsafe.Pointer(&got[i].Data))
		if got[i].Type != winapi.InputKeyboard || keyboard.VK != want[i].VK || keyboard.Flags != want[i].Flags {
			t.Errorf("input %d = type %d, %#v", i, got[i].Type, keyboard)
		}
	}
}
