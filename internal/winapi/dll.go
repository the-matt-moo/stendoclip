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
	gdi32    = windows.NewLazySystemDLL("gdi32.dll")

	procRegisterClassExW              = user32.NewProc("RegisterClassExW")
	procUnregisterClassW              = user32.NewProc("UnregisterClassW")
	procCreateWindowExW               = user32.NewProc("CreateWindowExW")
	procDestroyWindow                 = user32.NewProc("DestroyWindow")
	procDefWindowProcW                = user32.NewProc("DefWindowProcW")
	procPeekMessageW                  = user32.NewProc("PeekMessageW")
	procTranslateMessage              = user32.NewProc("TranslateMessage")
	procDispatchMessageW              = user32.NewProc("DispatchMessageW")
	procPostQuitMessage               = user32.NewProc("PostQuitMessage")
	procMsgWaitForMultipleObjects     = user32.NewProc("MsgWaitForMultipleObjects")
	procLoadIconW                     = user32.NewProc("LoadIconW")
	procLoadImageW                    = user32.NewProc("LoadImageW")
	procDestroyIcon                   = user32.NewProc("DestroyIcon")
	procCreatePopupMenu               = user32.NewProc("CreatePopupMenu")
	procDestroyMenu                   = user32.NewProc("DestroyMenu")
	procAppendMenuW                   = user32.NewProc("AppendMenuW")
	procTrackPopupMenu                = user32.NewProc("TrackPopupMenu")
	procGetCursorPos                  = user32.NewProc("GetCursorPos")
	procPostMessageW                  = user32.NewProc("PostMessageW")
	procMessageBoxW                   = user32.NewProc("MessageBoxW")
	procAddClipboardFormatListener    = user32.NewProc("AddClipboardFormatListener")
	procRemoveClipboardFormatListener = user32.NewProc("RemoveClipboardFormatListener")
	procOpenClipboard                 = user32.NewProc("OpenClipboard")
	procCloseClipboard                = user32.NewProc("CloseClipboard")
	procEmptyClipboard                = user32.NewProc("EmptyClipboard")
	procGetClipboardData              = user32.NewProc("GetClipboardData")
	procSetClipboardData              = user32.NewProc("SetClipboardData")
	procIsClipboardFormatAvailable    = user32.NewProc("IsClipboardFormatAvailable")
	procRegisterClipboardFormatW      = user32.NewProc("RegisterClipboardFormatW")
	procRegisterHotKey                = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey              = user32.NewProc("UnregisterHotKey")
	procRegisterWindowMessageW        = user32.NewProc("RegisterWindowMessageW")
	procGetForegroundWindow           = user32.NewProc("GetForegroundWindow")
	procSetForegroundWindow           = user32.NewProc("SetForegroundWindow")
	procSetFocus                      = user32.NewProc("SetFocus")
	procShowWindow                    = user32.NewProc("ShowWindow")
	procIsWindowVisible               = user32.NewProc("IsWindowVisible")
	procSetWindowPos                  = user32.NewProc("SetWindowPos")
	procInvalidateRect                = user32.NewProc("InvalidateRect")
	procGetClientRect                 = user32.NewProc("GetClientRect")
	procSetTimer                      = user32.NewProc("SetTimer")
	procKillTimer                     = user32.NewProc("KillTimer")
	procMonitorFromWindow             = user32.NewProc("MonitorFromWindow")
	procGetMonitorInfoW               = user32.NewProc("GetMonitorInfoW")
	procBeginPaint                    = user32.NewProc("BeginPaint")
	procEndPaint                      = user32.NewProc("EndPaint")
	procFillRect                      = user32.NewProc("FillRect")
	procDrawTextW                     = user32.NewProc("DrawTextW")
	procSetLayeredWindowAttributes    = user32.NewProc("SetLayeredWindowAttributes")
	procGetKeyState                   = user32.NewProc("GetKeyState")
	procSendInput                     = user32.NewProc("SendInput")
	procGetModuleHandleW              = kernel32.NewProc("GetModuleHandleW")
	procGlobalAlloc                   = kernel32.NewProc("GlobalAlloc")
	procGlobalFree                    = kernel32.NewProc("GlobalFree")
	procGlobalLock                    = kernel32.NewProc("GlobalLock")
	procGlobalUnlock                  = kernel32.NewProc("GlobalUnlock")
	procGlobalSize                    = kernel32.NewProc("GlobalSize")
	procRtlMoveMemory                 = kernel32.NewProc("RtlMoveMemory")
	procCreateSolidBrush              = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject                  = gdi32.NewProc("DeleteObject")
	procGetStockObject                = gdi32.NewProc("GetStockObject")
	procSelectObject                  = gdi32.NewProc("SelectObject")
	procSetBkMode                     = gdi32.NewProc("SetBkMode")
	procSetTextColor                  = gdi32.NewProc("SetTextColor")
	procShellNotifyIconW              = shell32.NewProc("Shell_NotifyIconW")
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

func RegisterWindowMessage(name string) (uint32, error) {
	pointer, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	r, _, callErr := procRegisterWindowMessageW.Call(uintptr(unsafe.Pointer(pointer)))
	if r == 0 {
		return 0, syscallError("RegisterWindowMessageW", callErr)
	}
	return uint32(r), nil
}

func RegisterHotKey(hwnd HWND, id int32, modifiers, key uint32) error {
	r, _, callErr := procRegisterHotKey.Call(uintptr(hwnd), uintptr(id), uintptr(modifiers), uintptr(key))
	if r == 0 {
		return syscallError("RegisterHotKey", callErr)
	}
	return nil
}

func UnregisterHotKey(hwnd HWND, id int32) error {
	r, _, callErr := procUnregisterHotKey.Call(uintptr(hwnd), uintptr(id))
	if r == 0 {
		return syscallError("UnregisterHotKey", callErr)
	}
	return nil
}

func GetForegroundWindow() HWND {
	r, _, _ := procGetForegroundWindow.Call()
	return HWND(r)
}

func SetForegroundWindow(hwnd HWND) error {
	r, _, callErr := procSetForegroundWindow.Call(uintptr(hwnd))
	if r == 0 {
		return syscallError("SetForegroundWindow", callErr)
	}
	return nil
}

func SetFocus(hwnd HWND) {
	procSetFocus.Call(uintptr(hwnd))
}

func ShowWindow(hwnd HWND, command int32) {
	procShowWindow.Call(uintptr(hwnd), uintptr(command))
}

func IsWindowVisible(hwnd HWND) bool {
	r, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
	return r != 0
}

func SetWindowPos(hwnd, insertAfter HWND, x, y, width, height int32, flags uint32) error {
	r, _, callErr := procSetWindowPos.Call(
		uintptr(hwnd), uintptr(insertAfter), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(flags),
	)
	if r == 0 {
		return syscallError("SetWindowPos", callErr)
	}
	return nil
}

func InvalidateRect(hwnd HWND) {
	procInvalidateRect.Call(uintptr(hwnd), 0, 0)
}

func GetClientRect(hwnd HWND) (Rect, error) {
	var rect Rect
	r, _, callErr := procGetClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
	if r == 0 {
		return Rect{}, syscallError("GetClientRect", callErr)
	}
	return rect, nil
}

func SetTimer(hwnd HWND, id uintptr, milliseconds uint32) error {
	r, _, callErr := procSetTimer.Call(uintptr(hwnd), id, uintptr(milliseconds), 0)
	if r == 0 {
		return syscallError("SetTimer", callErr)
	}
	return nil
}

func KillTimer(hwnd HWND, id uintptr) {
	procKillTimer.Call(uintptr(hwnd), id)
}

func MonitorFromWindow(hwnd HWND, flags uint32) HMONITOR {
	r, _, _ := procMonitorFromWindow.Call(uintptr(hwnd), uintptr(flags))
	return HMONITOR(r)
}

func GetMonitorInfo(monitor HMONITOR, info *MonitorInfo) error {
	r, _, callErr := procGetMonitorInfoW.Call(uintptr(monitor), uintptr(unsafe.Pointer(info)))
	if r == 0 {
		return syscallError("GetMonitorInfoW", callErr)
	}
	return nil
}

func BeginPaint(hwnd HWND, paint *PaintStruct) HDC {
	r, _, _ := procBeginPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(paint)))
	return HDC(r)
}

func EndPaint(hwnd HWND, paint *PaintStruct) {
	procEndPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(paint)))
}

func FillRect(dc HDC, rect *Rect, brush HBRUSH) {
	procFillRect.Call(uintptr(dc), uintptr(unsafe.Pointer(rect)), uintptr(brush))
}

func DrawText(dc HDC, text string, rect *Rect, format uint32) error {
	if text == "" {
		return nil
	}
	encoded, err := windows.UTF16FromString(text)
	if err != nil {
		return err
	}
	r, _, callErr := procDrawTextW.Call(
		uintptr(dc), uintptr(unsafe.Pointer(&encoded[0])), uintptr(len(encoded)-1), uintptr(unsafe.Pointer(rect)), uintptr(format),
	)
	if r == 0 {
		return syscallError("DrawTextW", callErr)
	}
	return nil
}

func SetLayeredWindowAlpha(hwnd HWND, alpha byte) error {
	r, _, callErr := procSetLayeredWindowAttributes.Call(uintptr(hwnd), 0, uintptr(alpha), LayeredWindowAlpha)
	if r == 0 {
		return syscallError("SetLayeredWindowAttributes", callErr)
	}
	return nil
}

func GetKeyState(key uint32) bool {
	r, _, _ := procGetKeyState.Call(uintptr(key))
	return uint16(r)&0x8000 != 0
}

func SendInput(inputs []Input) error {
	if len(inputs) == 0 {
		return nil
	}
	r, _, callErr := procSendInput.Call(
		uintptr(len(inputs)), uintptr(unsafe.Pointer(&inputs[0])), unsafe.Sizeof(inputs[0]),
	)
	if int(r) != len(inputs) {
		return syscallError("SendInput", callErr)
	}
	return nil
}

func CreateSolidBrush(color uint32) (HBRUSH, error) {
	r, _, callErr := procCreateSolidBrush.Call(uintptr(color))
	if r == 0 {
		return 0, syscallError("CreateSolidBrush", callErr)
	}
	return HBRUSH(r), nil
}

func DeleteObject(object HGDIOBJ) {
	procDeleteObject.Call(uintptr(object))
}

func GetStockObject(id int32) HGDIOBJ {
	r, _, _ := procGetStockObject.Call(uintptr(id))
	return HGDIOBJ(r)
}

func SelectObject(dc HDC, object HGDIOBJ) HGDIOBJ {
	r, _, _ := procSelectObject.Call(uintptr(dc), uintptr(object))
	return HGDIOBJ(r)
}

func SetBkMode(dc HDC, mode int32) {
	procSetBkMode.Call(uintptr(dc), uintptr(mode))
}

func SetTextColor(dc HDC, color uint32) {
	procSetTextColor.Call(uintptr(dc), uintptr(color))
}

func RGB(red, green, blue byte) uint32 {
	return uint32(red) | uint32(green)<<8 | uint32(blue)<<16
}

func AddClipboardFormatListener(hwnd HWND) error {
	r, _, callErr := procAddClipboardFormatListener.Call(uintptr(hwnd))
	if r == 0 {
		return syscallError("AddClipboardFormatListener", callErr)
	}
	return nil
}

func RemoveClipboardFormatListener(hwnd HWND) error {
	r, _, callErr := procRemoveClipboardFormatListener.Call(uintptr(hwnd))
	if r == 0 {
		return syscallError("RemoveClipboardFormatListener", callErr)
	}
	return nil
}

func OpenClipboard(hwnd HWND) error {
	r, _, callErr := procOpenClipboard.Call(uintptr(hwnd))
	if r == 0 {
		return syscallError("OpenClipboard", callErr)
	}
	return nil
}

func CloseClipboard() error {
	r, _, callErr := procCloseClipboard.Call()
	if r == 0 {
		return syscallError("CloseClipboard", callErr)
	}
	return nil
}

func EmptyClipboard() error {
	r, _, callErr := procEmptyClipboard.Call()
	if r == 0 {
		return syscallError("EmptyClipboard", callErr)
	}
	return nil
}

func GetClipboardData(format uint32) (windows.Handle, error) {
	r, _, callErr := procGetClipboardData.Call(uintptr(format))
	if r == 0 {
		return 0, syscallError("GetClipboardData", callErr)
	}
	return windows.Handle(r), nil
}

func SetClipboardData(format uint32, handle windows.Handle) error {
	r, _, callErr := procSetClipboardData.Call(uintptr(format), uintptr(handle))
	if r == 0 {
		return syscallError("SetClipboardData", callErr)
	}
	return nil
}

func IsClipboardFormatAvailable(format uint32) bool {
	r, _, _ := procIsClipboardFormatAvailable.Call(uintptr(format))
	return r != 0
}

func RegisterClipboardFormat(name string) (uint32, error) {
	pointer, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	r, _, callErr := procRegisterClipboardFormatW.Call(uintptr(unsafe.Pointer(pointer)))
	if r == 0 {
		return 0, syscallError("RegisterClipboardFormatW", callErr)
	}
	return uint32(r), nil
}

func GlobalAlloc(flags uint32, bytes uintptr) (windows.Handle, error) {
	r, _, callErr := procGlobalAlloc.Call(uintptr(flags), bytes)
	if r == 0 {
		return 0, syscallError("GlobalAlloc", callErr)
	}
	return windows.Handle(r), nil
}

func GlobalFree(handle windows.Handle) {
	procGlobalFree.Call(uintptr(handle))
}

func GlobalRead(handle windows.Handle, size uintptr) ([]byte, error) {
	if size == 0 {
		return nil, nil
	}
	if size > uintptr(^uint(0)>>1) {
		return nil, fmt.Errorf("global memory block too large: %d bytes", size)
	}
	pointer, _, callErr := procGlobalLock.Call(uintptr(handle))
	if pointer == 0 {
		return nil, syscallError("GlobalLock", callErr)
	}
	defer procGlobalUnlock.Call(uintptr(handle))
	data := make([]byte, int(size))
	procRtlMoveMemory.Call(uintptr(unsafe.Pointer(&data[0])), pointer, size)
	return data, nil
}

func GlobalWrite(handle windows.Handle, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	pointer, _, callErr := procGlobalLock.Call(uintptr(handle))
	if pointer == 0 {
		return syscallError("GlobalLock", callErr)
	}
	defer procGlobalUnlock.Call(uintptr(handle))
	procRtlMoveMemory.Call(pointer, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	return nil
}

func GlobalSize(handle windows.Handle) uintptr {
	r, _, _ := procGlobalSize.Call(uintptr(handle))
	return r
}

func LoadIcon(instance HINSTANCE, resourceID uint16) (HICON, error) {
	r, _, callErr := procLoadIconW.Call(uintptr(instance), uintptr(resourceID))
	if r == 0 {
		return 0, syscallError("LoadIconW", callErr)
	}
	return HICON(r), nil
}

func LoadIconFile(path string) (HICON, error) {
	pointer, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	r, _, callErr := procLoadImageW.Call(
		0, uintptr(unsafe.Pointer(pointer)), ImageIcon, 0, 0, LRLoadFromFile|LRDefaultSize,
	)
	if r == 0 {
		return 0, syscallError("LoadImageW", callErr)
	}
	return HICON(r), nil
}

func DestroyIcon(icon HICON) {
	procDestroyIcon.Call(uintptr(icon))
}

func CreatePopupMenu() (HMENU, error) {
	r, _, callErr := procCreatePopupMenu.Call()
	if r == 0 {
		return 0, syscallError("CreatePopupMenu", callErr)
	}
	return HMENU(r), nil
}

func DestroyMenu(menu HMENU) {
	procDestroyMenu.Call(uintptr(menu))
}

func AppendMenu(menu HMENU, flags uint32, id uintptr, text string) error {
	var pointer *uint16
	var err error
	if text != "" {
		pointer, err = windows.UTF16PtrFromString(text)
		if err != nil {
			return err
		}
	}
	r, _, callErr := procAppendMenuW.Call(uintptr(menu), uintptr(flags), id, uintptr(unsafe.Pointer(pointer)))
	if r == 0 {
		return syscallError("AppendMenuW", callErr)
	}
	return nil
}

func TrackPopupMenu(menu HMENU, flags uint32, x, y int32, owner HWND) uint32 {
	r, _, _ := procTrackPopupMenu.Call(
		uintptr(menu), uintptr(flags), uintptr(x), uintptr(y), 0, uintptr(owner), 0,
	)
	return uint32(r)
}

func GetCursorPos() (Point, error) {
	var point Point
	r, _, callErr := procGetCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	if r == 0 {
		return Point{}, syscallError("GetCursorPos", callErr)
	}
	return point, nil
}

func PostMessage(hwnd HWND, message uint32, wParam, lParam uintptr) {
	procPostMessageW.Call(uintptr(hwnd), uintptr(message), wParam, lParam)
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
	messageBox(title, message, 0x10)
}

func MessageBoxInfo(title, message string) {
	messageBox(title, message, 0x40)
}

func messageBox(title, message string, flags uint32) {
	titlePtr, _ := windows.UTF16PtrFromString(title)
	messagePtr, _ := windows.UTF16PtrFromString(message)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), uintptr(flags))
}

func syscallError(name string, err error) error {
	if err == nil || err == windows.ERROR_SUCCESS {
		return fmt.Errorf("%s failed", name)
	}
	return fmt.Errorf("%s: %w", name, err)
}
