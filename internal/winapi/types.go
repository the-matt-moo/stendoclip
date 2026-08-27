package winapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type (
	HWND      = windows.Handle
	HINSTANCE = windows.Handle
	HICON     = windows.Handle
	HCURSOR   = windows.Handle
	HBRUSH    = windows.Handle
	HMENU     = windows.Handle
	HDC       = windows.Handle
	HGDIOBJ   = windows.Handle
	HMONITOR  = windows.Handle
)

type Point struct {
	X int32
	Y int32
}

type Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
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

type OpenFileName struct {
	StructSize      uint32
	Owner           HWND
	Instance        HINSTANCE
	Filter          *uint16
	CustomFilter    *uint16
	MaxCustomFilter uint32
	FilterIndex     uint32
	File            *uint16
	MaxFile         uint32
	FileTitle       *uint16
	MaxFileTitle    uint32
	InitialDir      *uint16
	Title           *uint16
	Flags           uint32
	FileOffset      uint16
	FileExtension   uint16
	DefaultExt      *uint16
	CustData        uintptr
	Hook            uintptr
	TemplateName    *uint16
	Reserved        uintptr
	ReservedFlags   uint32
	FlagsEx         uint32
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

type PaintStruct struct {
	DC        HDC
	Erase     int32
	Paint     Rect
	Restore   int32
	IncUpdate int32
	Reserved  [32]byte
}

type MonitorInfo struct {
	CbSize  uint32
	Monitor Rect
	Work    Rect
	Flags   uint32
}

type KeyboardInput struct {
	VK        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type MouseInput struct {
	X         int32
	Y         int32
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type Input struct {
	Type uint32
	Data MouseInput
}

type MsgBoxParams struct {
	CbSize         uint32
	Owner          HWND
	Instance       HINSTANCE
	Text           *uint16
	Caption        *uint16
	Style          uint32
	Icon           uintptr // MAKEINTRESOURCE(id)
	ContextHelpId  uintptr
	MsgBoxCallback uintptr
	LanguageId     uint32
}

func NewKeyboardInput(key uint16, flags uint32) Input {
	input := Input{Type: InputKeyboard}
	*(*KeyboardInput)(unsafe.Pointer(&input.Data)) = KeyboardInput{VK: key, Flags: flags}
	return input
}
