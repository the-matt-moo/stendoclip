package overlay

import (
	"fmt"
	"unsafe"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func position(target winapi.HWND) (x, y, width, height int32, err error) {
	monitor := winapi.MonitorFromWindow(target, winapi.MonitorDefaultToNearest)
	if monitor == 0 {
		return 0, 0, 0, 0, fmt.Errorf("find paste target monitor")
	}
	info := winapi.MonitorInfo{CbSize: uint32(unsafe.Sizeof(winapi.MonitorInfo{}))}
	if err := winapi.GetMonitorInfo(monitor, &info); err != nil {
		return 0, 0, 0, 0, err
	}
	x, y, width, height = centered(info.Work)
	return x, y, width, height, nil
}

func centered(work winapi.Rect) (x, y, width, height int32) {
	workWidth := work.Right - work.Left
	workHeight := work.Bottom - work.Top
	width, height = 640, 190
	if width > workWidth-40 {
		width = workWidth - 40
	}
	if height > workHeight-40 {
		height = workHeight - 40
	}
	x = work.Left + (workWidth-width)/2
	y = work.Top + (workHeight-height)/2
	return
}
