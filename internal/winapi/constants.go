package winapi

const (
	WMNull            = 0x0000
	WMDestroy         = 0x0002
	WMPaint           = 0x000F
	WMClose           = 0x0010
	WMQuit            = 0x0012
	WMEraseBackground = 0x0014
	WMContextMenu     = 0x007B
	WMKeyDown         = 0x0100
	WMSysKeyDown      = 0x0104
	WMTimer           = 0x0113
	WMRButtonUp       = 0x0205
	WMHotkey          = 0x0312
	WMClipboardUpdate = 0x031D
	WMApp             = 0x8000

	CFUnicodeText = 13
	GMEMMoveable  = 0x0002

	PMRemove   = 0x0001
	QSAllInput = 0x04FF

	WaitObject0 = 0x00000000
	WaitFailed  = 0xFFFFFFFF
	Infinite    = 0xFFFFFFFF

	WSCaption         = 0x00C00000
	WSSysMenu         = 0x00080000
	WSPopup           = 0x80000000
	WSExTopmost       = 0x00000008
	WSExToolWindow    = 0x00000080
	WSExLayered       = 0x00080000
	WSExDlgModalFrame = 0x00000001

	SWHide = 0
	SWShow = 5

	SWPNoActivate = 0x0010
	SWPShowWindow = 0x0040

	MODAlt      = 0x0001
	MODControl  = 0x0002
	MODShift    = 0x0004
	MODWin      = 0x0008
	MODNoRepeat = 0x4000

	VKBack    = 0x08
	VKTab     = 0x09
	VKReturn  = 0x0D
	VKShift   = 0x10
	VKControl = 0x11
	VKMenu    = 0x12
	VKEscape  = 0x1B
	VKSpace   = 0x20
	VKEnd     = 0x23
	VKHome    = 0x24
	VKLeft    = 0x25
	VKUp      = 0x26
	VKRight   = 0x27
	VKDown    = 0x28
	VKDelete  = 0x2E
	VKLWin    = 0x5B
	VKRWin    = 0x5C
	VK0       = 0x30
	VKA       = 0x41
	VKP       = 0x50
	VKF1      = 0x70
	VKF12     = 0x7B

	MonitorDefaultToNearest = 0x00000002

	LayeredWindowAlpha = 0x00000002

	InputKeyboard = 1
	KeyEventKeyUp = 0x0002

	Transparent    = 1
	DefaultGUIFont = 17
	DTWordBreak    = 0x0010
	DTRight        = 0x0002
	DTBottom       = 0x0008
	DTSingleLine   = 0x0020
	DTNoPrefix     = 0x0800
	DTEndEllipsis  = 0x8000
	DTCenter       = 0x0001
	DTCalcRect     = 0x0400

	FWNormal          = 400
	DefaultCharset    = 1
	OutDefaultPrecis  = 0
	ClipDefaultPrecis = 0
	ClearTypeQuality  = 5
	DefaultPitch      = 0

	DIFlagNormal = 0x0003

	WMLButtonDown = 0x0201

	NIMAdd        = 0x00000000
	NIMDelete     = 0x00000002
	NIMSetVersion = 0x00000004

	NIFMessage = 0x00000001
	NIFIcon    = 0x00000002
	NIFTip     = 0x00000004
	NIFInfo    = 0x00000010
	NIFShowTip = 0x00000080

	NIIFInfo = 0x00000001

	MFString    = 0x00000000
	MFGray      = 0x00000001
	MFChecked   = 0x00000008
	MFPopup     = 0x00000010
	MFSeparator = 0x00000800

	TPMRightButton = 0x0002
	TPMNonotify    = 0x0080
	TPMReturnCmd   = 0x0100

	MBUserIcon = 0x00000080

	ImageIcon      = 1
	LRLoadFromFile = 0x0010
	LRDefaultSize  = 0x0040

	NotifyIconVersion4 = 4
	IDIInformation     = 32516
)

var HWNDTopmost = HWND(^uintptr(0)) // (HWND)-1
