package winapi

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shell32  = windows.NewLazySystemDLL("shell32.dll")

	procRegisterClassExW          = user32.NewProc("RegisterClassExW")
	procUnregisterClassW          = user32.NewProc("UnregisterClassW")
	procCreateWindowExW           = user32.NewProc("CreateWindowExW")
	procDestroyWindow             = user32.NewProc("DestroyWindow")
	procDefWindowProcW            = user32.NewProc("DefWindowProcW")
	procPeekMessageW              = user32.NewProc("PeekMessageW")
	procTranslateMessage          = user32.NewProc("TranslateMessage")
	procDispatchMessageW          = user32.NewProc("DispatchMessageW")
	procPostQuitMessage           = user32.NewProc("PostQuitMessage")
	procMsgWaitForMultipleObjects = user32.NewProc("MsgWaitForMultipleObjects")
	procLoadIconW                 = user32.NewProc("LoadIconW")
	procMessageBoxW               = user32.NewProc("MessageBoxW")
	procGetModuleHandleW          = kernel32.NewProc("GetModuleHandleW")
	procShellNotifyIconW          = shell32.NewProc("Shell_NotifyIconW")
)

func RegisterClassEx(class *WndClassEx) (uint16, error) {
	r, _, callErr := procRegisterClassExW.Call(uintptr(unsafe.Pointer(class)))
	if r == 0 {
		return 0, syscallError("RegisterClassExW", callErr)
	}
	return uint16(r), nil
}

func UnregisterClass(className *uint16, instance HINSTANCE) error {
	r, _, callErr := procUnregisterClassW.Call(uintptr(unsafe.Pointer(className)), uintptr(instance))
	if r == 0 {
		return syscallError("UnregisterClassW", callErr)
	}
	return nil
}

func CreateWindowEx(exStyle uint32, className, windowName *uint16, style uint32, x, y, width, height int32, parent HWND, menu HMENU, instance HINSTANCE, param unsafe.Pointer) (HWND, error) {
	r, _, callErr := procCreateWindowExW.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr(style),
		uintptr(x), uintptr(y), uintptr(width), uintptr(height),
		uintptr(parent), uintptr(menu), uintptr(instance), uintptr(param),
	)
	if r == 0 {
		return 0, syscallError("CreateWindowExW", callErr)
	}
	return HWND(r), nil
}

func DestroyWindow(hwnd HWND) error {
	r, _, callErr := procDestroyWindow.Call(uintptr(hwnd))
	if r == 0 {
		return syscallError("DestroyWindow", callErr)
	}
	return nil
}

func DefWindowProc(hwnd HWND, message uint32, wParam, lParam uintptr) uintptr {
	r, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(message), wParam, lParam)
	return r
}

func PeekMessage(message *Msg, hwnd HWND, min, max, remove uint32) bool {
	r, _, _ := procPeekMessageW.Call(
		uintptr(unsafe.Pointer(message)), uintptr(hwnd), uintptr(min), uintptr(max), uintptr(remove),
	)
	return r != 0
}

func TranslateMessage(message *Msg) {
	procTranslateMessage.Call(uintptr(unsafe.Pointer(message)))
}

func DispatchMessage(message *Msg) uintptr {
	r, _, _ := procDispatchMessageW.Call(uintptr(unsafe.Pointer(message)))
	return r
}

func PostQuitMessage(exitCode int32) {
	procPostQuitMessage.Call(uintptr(exitCode))
}

func MsgWaitForMultipleObjects(handles []windows.Handle, waitAll bool, timeout, wakeMask uint32) (uint32, error) {
	var pointer uintptr
	if len(handles) > 0 {
		pointer = uintptr(unsafe.Pointer(&handles[0]))
	}
	var wait uintptr
	if waitAll {
		wait = 1
	}
	r, _, callErr := procMsgWaitForMultipleObjects.Call(
		uintptr(len(handles)), pointer, wait, uintptr(timeout), uintptr(wakeMask),
	)
	if uint32(r) == WaitFailed {
		return WaitFailed, syscallError("MsgWaitForMultipleObjects", callErr)
	}
	return uint32(r), nil
}

func LoadIcon(instance HINSTANCE, resourceID uint16) (HICON, error) {
	r, _, callErr := procLoadIconW.Call(uintptr(instance), uintptr(resourceID))
	if r == 0 {
		return 0, syscallError("LoadIconW", callErr)
	}
	return HICON(r), nil
}

func GetModuleHandle() (HINSTANCE, error) {
	r, _, callErr := procGetModuleHandleW.Call(0)
	if r == 0 {
		return 0, syscallError("GetModuleHandleW", callErr)
	}
	return HINSTANCE(r), nil
}

func ShellNotifyIcon(message uint32, data *NotifyIconData) error {
	r, _, callErr := procShellNotifyIconW.Call(uintptr(message), uintptr(unsafe.Pointer(data)))
	if r == 0 {
		return syscallError("Shell_NotifyIconW", callErr)
	}
	return nil
}

func MessageBox(title, message string) {
	titlePtr, _ := windows.UTF16PtrFromString(title)
	messagePtr, _ := windows.UTF16PtrFromString(message)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}

func syscallError(name string, err error) error {
	if err == nil || err == windows.ERROR_SUCCESS {
		return fmt.Errorf("%s failed", name)
	}
	return fmt.Errorf("%s: %w", name, err)
}
