package winapi

import "golang.org/x/sys/windows"

type (
	HWND      = windows.Handle
	HINSTANCE = windows.Handle
	HICON     = windows.Handle
	HCURSOR   = windows.Handle
	HBRUSH    = windows.Handle
	HMENU     = windows.Handle
)

type Point struct {
	X int32
	Y int32
}

type Msg struct {
	HWnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      Point
}

type WndClassEx struct {
	CbSize     uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   HINSTANCE
	Icon       HICON
	Cursor     HCURSOR
	Background HBRUSH
	MenuName   *uint16
	ClassName  *uint16
	IconSm     HICON
}

type NotifyIconData struct {
	CbSize         uint32
	HWnd           HWND
	UID            uint32
	Flags          uint32
	CallbackMsg    uint32
	Icon           HICON
	Tip            [128]uint16
	State          uint32
	StateMask      uint32
	Info           [256]uint16
	TimeoutVersion uint32
	InfoTitle      [64]uint16
	InfoFlags      uint32
	GUIDItem       windows.GUID
	BalloonIcon    HICON
}
