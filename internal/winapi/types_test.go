package winapi

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestStructSizes(t *testing.T) {
	want := map[string][8]uintptr{
		"386":   {48, 28, 956, 64, 40, 16, 24, 28},
		"amd64": {80, 48, 976, 72, 40, 24, 32, 40},
	}[runtime.GOARCH]
	if want == [8]uintptr{} {
		t.Skip("unsupported architecture")
	}
	got := [8]uintptr{
		unsafe.Sizeof(WndClassEx{}),
		unsafe.Sizeof(Msg{}),
		unsafe.Sizeof(NotifyIconData{}),
		unsafe.Sizeof(PaintStruct{}),
		unsafe.Sizeof(MonitorInfo{}),
		unsafe.Sizeof(KeyboardInput{}),
		unsafe.Sizeof(MouseInput{}),
		unsafe.Sizeof(Input{}),
	}
	if got != want {
		t.Fatalf("struct sizes = %v, want %v", got, want)
	}
}
