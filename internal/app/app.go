package app

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

var (
	windowProcCallback = windows.NewCallback(windowProc)
	activeApp          *App
	activeAppMu        sync.Mutex
)

type App struct {
	mu           sync.Mutex
	hwnd         winapi.HWND
	instance     winapi.HINSTANCE
	className    *uint16
	commandEvent windows.Handle
	commands     chan func()
	closed       bool
}

func New(version string) (*App, error) {
	instance, err := winapi.GetModuleHandle()
	if err != nil {
		return nil, err
	}
	event, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("create command event: %w", err)
	}

	className, err := windows.UTF16PtrFromString(fmt.Sprintf("Stendoclip.Hidden.%d", os.Getpid()))
	if err != nil {
		windows.CloseHandle(event)
		return nil, err
	}
	title, err := windows.UTF16PtrFromString("Stendoclip " + version)
	if err != nil {
		windows.CloseHandle(event)
		return nil, err
	}

	a := &App{
		instance:     instance,
		className:    className,
		commandEvent: event,
		commands:     make(chan func(), 64),
	}
	class := winapi.WndClassEx{
		WndProc:   windowProcCallback,
		Instance:  instance,
		ClassName: className,
	}
	class.CbSize = uint32(unsafe.Sizeof(class))
	if _, err := winapi.RegisterClassEx(&class); err != nil {
		windows.CloseHandle(event)
		return nil, err
	}

	activeAppMu.Lock()
	activeApp = a
	activeAppMu.Unlock()
	hwnd, err := winapi.CreateWindowEx(0, className, title, 0, 0, 0, 0, 0, winapi.HWNDMessage, 0, instance, nil)
	if err != nil {
		activeAppMu.Lock()
		activeApp = nil
		activeAppMu.Unlock()
		_ = winapi.UnregisterClass(className, instance)
		windows.CloseHandle(event)
		return nil, err
	}
	a.hwnd = hwnd
	return a, nil
}

func (a *App) Run() error {
	for {
		result, err := winapi.MsgWaitForMultipleObjects(
			[]windows.Handle{a.commandEvent}, false, winapi.Infinite, winapi.QSAllInput,
		)
		if err != nil {
			return err
		}
		switch result {
		case winapi.WaitObject0:
			a.drainCommands()
		case winapi.WaitObject0 + 1:
			var message winapi.Msg
			for winapi.PeekMessage(&message, 0, 0, 0, winapi.PMRemove) {
				if message.Message == winapi.WMQuit {
					return nil
				}
				winapi.TranslateMessage(&message)
				winapi.DispatchMessage(&message)
			}
		default:
			return fmt.Errorf("unexpected message wait result: 0x%x", result)
		}
	}
}

func (a *App) Post(command func()) error {
	if command == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.closed {
		return errors.New("app closed")
	}
	select {
	case a.commands <- command:
		if err := windows.SetEvent(a.commandEvent); err != nil {
			return fmt.Errorf("signal command event: %w", err)
		}
		return nil
	default:
		return errors.New("command queue full")
	}
}

func (a *App) Quit() error {
	return a.Post(func() { winapi.PostQuitMessage(0) })
}

func (a *App) Window() winapi.HWND {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.hwnd
}

func (a *App) Close() error {
	a.mu.Lock()
	if a.closed {
		a.mu.Unlock()
		return nil
	}
	a.closed = true
	hwnd := a.hwnd
	a.hwnd = 0
	event := a.commandEvent
	a.commandEvent = 0
	a.mu.Unlock()

	var result error
	if hwnd != 0 {
		if err := winapi.DestroyWindow(hwnd); err != nil {
			result = errors.Join(result, err)
		}
	}
	if err := winapi.UnregisterClass(a.className, a.instance); err != nil {
		result = errors.Join(result, err)
	}
	if err := windows.CloseHandle(event); err != nil {
		result = errors.Join(result, err)
	}
	activeAppMu.Lock()
	if activeApp == a {
		activeApp = nil
	}
	activeAppMu.Unlock()
	return result
}

func (a *App) drainCommands() {
	for {
		select {
		case command := <-a.commands:
			command()
		default:
			return
		}
	}
}

func windowProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	switch message {
	case winapi.WMClose:
		_ = winapi.DestroyWindow(winapi.HWND(hwnd))
		return 0
	case winapi.WMDestroy:
		activeAppMu.Lock()
		if activeApp != nil {
			activeApp.mu.Lock()
			if activeApp.hwnd == winapi.HWND(hwnd) {
				activeApp.hwnd = 0
			}
			activeApp.mu.Unlock()
		}
		activeAppMu.Unlock()
		winapi.PostQuitMessage(0)
		return 0
	default:
		return winapi.DefWindowProc(winapi.HWND(hwnd), message, wParam, lParam)
	}
}

func LockMainThread() func() {
	runtime.LockOSThread()
	return runtime.UnlockOSThread
}
