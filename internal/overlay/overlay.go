package overlay

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mooreceipts/stendoclip/internal/hotkey"
	"github.com/mooreceipts/stendoclip/internal/store"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

const timerID = 1

var (
	windowProcCallback = windows.NewCallback(windowProc)
	active             *Controller
	activeMu           sync.Mutex
)

type Controller struct {
	hwnd        winapi.HWND
	instance    winapi.HINSTANCE
	className   *uint16
	brush       winapi.HBRUSH
	stack       *store.ClippingStack
	historyPath string
	timeout     time.Duration
	wraparound  bool
	pinBinding  hotkey.Binding
	paste       func(winapi.HWND, string) error
	onError     func(error)
	target      winapi.HWND
	index       int
	visible     bool
}

func New(stack *store.ClippingStack, historyPath, pinSpec string, timeout time.Duration, wraparound bool, paste func(winapi.HWND, string) error, onError func(error)) (*Controller, error) {
	instance, err := winapi.GetModuleHandle()
	if err != nil {
		return nil, err
	}
	pinBinding, err := hotkey.Parse(pinSpec)
	if err != nil {
		return nil, fmt.Errorf("parse pin hotkey: %w", err)
	}
	className, err := windows.UTF16PtrFromString(fmt.Sprintf("Stendoclip.Bezel.%d", os.Getpid()))
	if err != nil {
		return nil, err
	}
	brush, err := winapi.CreateSolidBrush(winapi.RGB(28, 28, 30))
	if err != nil {
		return nil, err
	}
	controller := &Controller{
		instance: instance, className: className, brush: brush, stack: stack, historyPath: historyPath,
		timeout: timeout, wraparound: wraparound, pinBinding: pinBinding, paste: paste, onError: onError,
	}
	class := winapi.WndClassEx{WndProc: windowProcCallback, Instance: instance, ClassName: className}
	class.CbSize = uint32(unsafe.Sizeof(class))
	if _, err := winapi.RegisterClassEx(&class); err != nil {
		winapi.DeleteObject(winapi.HGDIOBJ(brush))
		return nil, err
	}

	activeMu.Lock()
	active = controller
	activeMu.Unlock()
	hwnd, err := winapi.CreateWindowEx(
		winapi.WSExTopmost|winapi.WSExToolWindow|winapi.WSExLayered,
		className, className, winapi.WSPopup, 0, 0, 0, 0, 0, 0, instance, nil,
	)
	if err != nil {
		activeMu.Lock()
		active = nil
		activeMu.Unlock()
		_ = winapi.UnregisterClass(className, instance)
		winapi.DeleteObject(winapi.HGDIOBJ(brush))
		return nil, err
	}
	controller.hwnd = hwnd
	if err := winapi.SetLayeredWindowAlpha(hwnd, 235); err != nil {
		_ = controller.Close()
		return nil, err
	}
	return controller, nil
}

func (c *Controller) Open() error {
	if c.stack.Len() == 0 {
		return nil
	}
	if c.visible {
		c.cycle(1)
		return nil
	}
	c.target = winapi.GetForegroundWindow()
	c.index = 0
	x, y, width, height, err := position(c.target)
	if err != nil {
		return err
	}
	if err := winapi.SetWindowPos(c.hwnd, winapi.HWNDTopmost, x, y, width, height, winapi.SWPNoActivate); err != nil {
		return err
	}
	c.visible = true
	winapi.ShowWindow(c.hwnd, winapi.SWShow)
	if err := winapi.SetForegroundWindow(c.hwnd); err != nil {
		c.hide()
		return err
	}
	winapi.SetFocus(c.hwnd)
	winapi.InvalidateRect(c.hwnd)
	return c.resetTimer()
}

func (c *Controller) Close() error {
	c.hide()
	var result error
	if c.hwnd != 0 {
		if err := winapi.DestroyWindow(c.hwnd); err != nil {
			result = errors.Join(result, err)
		}
		c.hwnd = 0
	}
	if err := winapi.UnregisterClass(c.className, c.instance); err != nil {
		result = errors.Join(result, err)
	}
	if c.brush != 0 {
		winapi.DeleteObject(winapi.HGDIOBJ(c.brush))
		c.brush = 0
	}
	activeMu.Lock()
	if active == c {
		active = nil
	}
	activeMu.Unlock()
	return result
}

func (c *Controller) handleKey(key uint32, repeated bool) {
	if hotkey.Matches(c.pinBinding, key) {
		if repeated {
			return
		}
		entry := c.stack.Get(c.index)
		c.stack.TogglePin(c.index)
		c.find(entry.Text)
		c.save()
		c.refresh()
		return
	}
	switch key {
	case winapi.VKUp, winapi.VKLeft:
		c.cycle(-1)
	case winapi.VKDown, winapi.VKRight:
		c.cycle(1)
	case winapi.VKReturn:
		entry := c.stack.Get(c.index)
		target := c.target
		c.hide()
		if c.paste != nil {
			c.report(c.paste(target, entry.Text))
		}
	case winapi.VKEscape:
		c.dismiss()
	case winapi.VKDelete:
		c.stack.Delete(c.index)
		c.save()
		if c.stack.Len() == 0 {
			c.dismiss()
			return
		}
		if c.index >= c.stack.Len() {
			c.index = c.stack.Len() - 1
		}
		c.refresh()
	}
}

func (c *Controller) cycle(direction int) {
	c.index = nextIndex(c.index, direction, c.stack.Len(), c.wraparound)
	c.refresh()
}

func (c *Controller) find(text string) {
	for i := 0; i < c.stack.Len(); i++ {
		if c.stack.Get(i).Text == text {
			c.index = i
			return
		}
	}
}

func (c *Controller) refresh() {
	winapi.InvalidateRect(c.hwnd)
	c.report(c.resetTimer())
}

func (c *Controller) resetTimer() error {
	winapi.KillTimer(c.hwnd, timerID)
	return winapi.SetTimer(c.hwnd, timerID, uint32(c.timeout/time.Millisecond))
}

func (c *Controller) hide() {
	if c.hwnd == 0 {
		return
	}
	winapi.KillTimer(c.hwnd, timerID)
	winapi.ShowWindow(c.hwnd, winapi.SWHide)
	c.visible = false
}

func (c *Controller) dismiss() {
	target := c.target
	c.hide()
	if target != 0 {
		c.report(winapi.SetForegroundWindow(target))
	}
}

func (c *Controller) save() { c.report(c.stack.Save(c.historyPath)) }

func (c *Controller) report(err error) {
	if err != nil && c.onError != nil {
		c.onError(err)
	}
}

func nextIndex(current, direction, length int, wrap bool) int {
	if length == 0 {
		return -1
	}
	next := current + direction
	if wrap {
		return ((next % length) + length) % length
	}
	if next < 0 {
		return 0
	}
	if next >= length {
		return length - 1
	}
	return next
}

func windowProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	activeMu.Lock()
	controller := active
	activeMu.Unlock()
	if controller == nil {
		return winapi.DefWindowProc(winapi.HWND(hwnd), message, wParam, lParam)
	}
	switch message {
	case winapi.WMPaint:
		controller.paint()
		return 0
	case winapi.WMEraseBackground:
		return 1
	case winapi.WMKeyDown, winapi.WMSysKeyDown:
		controller.handleKey(uint32(wParam), lParam&(1<<30) != 0)
		return 0
	case winapi.WMTimer:
		if wParam == timerID {
			controller.dismiss()
		}
		return 0
	case winapi.WMClose:
		controller.dismiss()
		return 0
	}
	return winapi.DefWindowProc(winapi.HWND(hwnd), message, wParam, lParam)
}
