package winapi

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestStructSizes(t *testing.T) {
	want := map[string][3]uintptr{
		"386":   {48, 28, 956},
		"amd64": {80, 48, 976},
	}[runtime.GOARCH]
	if want == [3]uintptr{} {
		t.Skip("unsupported architecture")
	}
	got := [3]uintptr{
		unsafe.Sizeof(WndClassEx{}),
		unsafe.Sizeof(Msg{}),
		unsafe.Sizeof(NotifyIconData{}),
	}
	if got != want {
		t.Fatalf("struct sizes = %v, want %v", got, want)
	}
}
