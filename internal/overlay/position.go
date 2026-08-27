package overlay

import (
	"fmt"
	"unsafe"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func position(target winapi.HWND, fontSize int32) (x, y, width, height int32, err error) {
	monitor := winapi.MonitorFromWindow(target, winapi.MonitorDefaultToNearest)
	if monitor == 0 {
		return 0, 0, 0, 0, fmt.Errorf("find paste target monitor")
	}
	info := winapi.MonitorInfo{CbSize: uint32(unsafe.Sizeof(winapi.MonitorInfo{}))}
	if err := winapi.GetMonitorInfo(monitor, &info); err != nil {
		return 0, 0, 0, 0, err
	}
	x, y, width, height = centered(info.Work, fontSize)
	return x, y, width, height, nil
}

func centered(work winapi.Rect, fontSize int32) (x, y, width, height int32) {
	workWidth := work.Right - work.Left
	workHeight := work.Bottom - work.Top

	// Scale height with font size: base 190 at font 18, grows proportionally.
	width = 640
	height = 190 * fontSize / 18
	if height < 190 {
		height = 190
	}

	// Cap at 90% of screen.
	maxWidth := workWidth * 9 / 10
	maxHeight := workHeight * 9 / 10
	if width > maxWidth {
		width = maxWidth
	}
	if height > maxHeight {
		height = maxHeight
	}

	x = work.Left + (workWidth-width)/2
	y = work.Top + (workHeight-height)/2
	return
}
