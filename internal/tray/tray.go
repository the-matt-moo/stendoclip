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
	hwnd               winapi.HWND
	instance           winapi.HINSTANCE
	icon               winapi.HICON
	aboutImage         []byte
	stack              *store.ClippingStack
	historyPath        string
	markdownExportPath string
	configPath         string
	currentFontSize    int
	trimWhitespace     bool
	onFontSizeChange   func(int)
	onTrimChange       func(bool)
	executable         string
	version            string
	paste              func(winapi.HWND, string) error
	onPause            func(bool)
	onQuit             func() error
	onError            func(error)
	paused             bool
	added              bool
	taskbarMessage     uint32
}

func New(hwnd winapi.HWND, iconData, aboutImage []byte, stack *store.ClippingStack, historyPath, markdownExportPath, configPath, executable, version string, currentFontSize int, trimWhitespace bool, paste func(winapi.HWND, string) error, onPause func(bool), onFontSizeChange func(int), onTrimChange func(bool), onQuit func() error, onError func(error)) (*Tray, error) {
	instance, err := winapi.GetModuleHandle()
	if err != nil {
		return nil, err
	}
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
		hwnd: hwnd, instance: instance, icon: icon, aboutImage: aboutImage,
		stack: stack, historyPath: historyPath, markdownExportPath: markdownExportPath, configPath: configPath, currentFontSize: currentFontSize, trimWhitespace: trimWhitespace, onFontSizeChange: onFontSizeChange, onTrimChange: onTrimChange, executable: executable, version: version,
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

func (t *Tray) SetMarkdownExportPath(path string) { t.markdownExportPath = path }

func (t *Tray) SetFontSize(size int) { t.currentFontSize = size }

func (t *Tray) SetTrimWhitespace(enabled bool) { t.trimWhitespace = enabled }

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
	if point == nil {
		cursor, err := winapi.GetCursorPos()
		if err != nil {
			t.report(err)
			return
		}
		point = &cursor
	}
	state, err := t.buildMenu()
	if err != nil {
		t.report(err)
		return
	}
	defer winapi.DestroyMenu(state.menu)
	for {
		t.report(winapi.SetForegroundWindow(t.hwnd))
		command := winapi.TrackPopupMenu(
			state.menu, winapi.TPMRightButton|winapi.TPMNonotify|winapi.TPMReturnCmd, point.X, point.Y, t.hwnd,
		)
		if command == 0 {
			break
		}
		if command == incFontID || command == decFontID {
			t.runCommand(command, state, target)
			continue
		}
		winapi.PostMessage(t.hwnd, winapi.WMNull, 0, 0)
		if target != 0 {
			t.report(winapi.SetForegroundWindow(target))
		}
		t.runCommand(command, state, target)
		return
	}
	winapi.PostMessage(t.hwnd, winapi.WMNull, 0, 0)
	if target != 0 {
		t.report(winapi.SetForegroundWindow(target))
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
