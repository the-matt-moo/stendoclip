package winapi

const (
	WMClose   = 0x0010
	WMDestroy = 0x0002
	WMQuit    = 0x0012

	PMRemove   = 0x0001
	QSAllInput = 0x04FF

	WaitObject0 = 0x00000000
	WaitFailed  = 0xFFFFFFFF
	Infinite    = 0xFFFFFFFF

	NIMAdd        = 0x00000000
	NIMDelete     = 0x00000002
	NIMSetVersion = 0x00000004

	NIFMessage = 0x00000001
	NIFIcon    = 0x00000002
	NIFTip     = 0x00000004
	NIFInfo    = 0x00000010

	NIIFInfo = 0x00000001

	NotifyIconVersion4 = 4
	IDIInformation     = 32516
)

var HWNDMessage = HWND(^uintptr(2)) // (HWND)-3
