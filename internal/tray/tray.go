package tray

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mooreceipts/stendoclip/internal/store"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

const (
	CallbackMessage = winapi.WMApp + 1
	iconID          = 1
)

type Tray struct {
	hwnd           winapi.HWND
	icon           winapi.HICON
	stack          *store.ClippingStack
	historyPath    string
	executable     string
	version        string
	paste          func(winapi.HWND, string) error
	onPause        func(bool)
	onQuit         func() error
	onError        func(error)
	paused         bool
	added          bool
	taskbarMessage uint32
}

func New(hwnd winapi.HWND, iconData []byte, stack *store.ClippingStack, historyPath, executable, version string, paste func(winapi.HWND, string) error, onPause func(bool), onQuit func() error, onError func(error)) (*Tray, error) {
	icon, err := loadIcon(iconData)
	if err != nil {
		return nil, err
	}
	taskbarMessage, err := winapi.RegisterWindowMessage("TaskbarCreated")
	if err != nil {
		winapi.DestroyIcon(icon)
		return nil, err
	}
	tray := &Tray{
		hwnd: hwnd, icon: icon, stack: stack, historyPath: historyPath, executable: executable, version: version,
		paste: paste, onPause: onPause, onQuit: onQuit, onError: onError, taskbarMessage: taskbarMessage,
	}
	if err := tray.addIcon(); err != nil {
		winapi.DestroyIcon(icon)
		return nil, err
	}
	return tray, nil
}

func (t *Tray) HandleMessage(message uint32, wParam, lParam uintptr) bool {
	if message == t.taskbarMessage {
		t.added = false
		t.report(t.addIcon())
		return true
	}
	if message != CallbackMessage {
		return false
	}
	event := uint32(lParam & 0xffff)
	switch event {
	case winapi.WMContextMenu:
		point := pointFromWParam(wParam)
		if point.X == -1 && point.Y == -1 {
			t.showMenu(nil)
		} else {
			t.showMenu(&point)
		}
	case winapi.WMRButtonUp:
		t.showMenu(nil)
	}
	return true
}

func (t *Tray) addIcon() error {
	data := winapi.NotifyIconData{
		CbSize:      uint32(unsafe.Sizeof(winapi.NotifyIconData{})),
		HWnd:        t.hwnd,
		UID:         iconID,
		Flags:       winapi.NIFMessage | winapi.NIFIcon | winapi.NIFTip | winapi.NIFShowTip,
		CallbackMsg: CallbackMessage,
		Icon:        t.icon,
	}
	copyUTF16(data.Tip[:], "Stendoclip "+t.version)
	if err := winapi.ShellNotifyIcon(winapi.NIMAdd, &data); err != nil {
		return fmt.Errorf("add tray icon: %w", err)
	}
	t.added = true
	data.TimeoutVersion = winapi.NotifyIconVersion4
	if err := winapi.ShellNotifyIcon(winapi.NIMSetVersion, &data); err != nil {
		_ = winapi.ShellNotifyIcon(winapi.NIMDelete, &data)
		t.added = false
		return fmt.Errorf("set tray icon version: %w", err)
	}
	return nil
}

func (t *Tray) Close() error {
	var result error
	if t.added {
		data := winapi.NotifyIconData{CbSize: uint32(unsafe.Sizeof(winapi.NotifyIconData{})), HWnd: t.hwnd, UID: iconID}
		result = errors.Join(result, winapi.ShellNotifyIcon(winapi.NIMDelete, &data))
		t.added = false
	}
	if t.icon != 0 {
		winapi.DestroyIcon(t.icon)
		t.icon = 0
	}
	return result
}

func (t *Tray) showMenu(point *winapi.Point) {
	target := winapi.GetForegroundWindow()
	state, err := t.buildMenu()
	if err != nil {
		t.report(err)
		return
	}
	defer winapi.DestroyMenu(state.menu)
	if point == nil {
		cursor, err := winapi.GetCursorPos()
		if err != nil {
			t.report(err)
			return
		}
		point = &cursor
	}
	t.report(winapi.SetForegroundWindow(t.hwnd))
	command := winapi.TrackPopupMenu(
		state.menu, winapi.TPMRightButton|winapi.TPMNonotify|winapi.TPMReturnCmd, point.X, point.Y, t.hwnd,
	)
	winapi.PostMessage(t.hwnd, winapi.WMNull, 0, 0)
	if target != 0 {
		t.report(winapi.SetForegroundWindow(target))
	}
	if command != 0 {
		t.runCommand(command, state, target)
	}
}

func pointFromWParam(value uintptr) winapi.Point {
	return winapi.Point{X: int32(int16(value)), Y: int32(int16(value >> 16))}
}

func (t *Tray) report(err error) {
	if err != nil && t.onError != nil {
		t.onError(err)
	}
}

func copyUTF16(destination []uint16, value string) {
	encoded, _ := windows.UTF16FromString(value)
	copy(destination, encoded)
}
