package tray

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

const (
	aboutWidth       = 440
	aboutHeight      = 440
	aboutPadding     = 16
	aboutImageWidth  = 256
	aboutImageHeight = 145 // 920:521 aspect ratio
	aboutImageGap    = 12
)

var (
	aboutProcCallback = windows.NewCallback(aboutProc)
	activeAbout       *aboutWindow
	activeAboutMu     sync.Mutex
)

type aboutWindow struct {
	hwnd      winapi.HWND
	instance  winapi.HINSTANCE
	className *uint16
	image     winapi.GpImage
	text      string
}

func showAbout(ownerInstance winapi.HINSTANCE, imageData []byte, version string) {
	activeAboutMu.Lock()
	if activeAbout != nil {
		hwnd := activeAbout.hwnd
		activeAboutMu.Unlock()
		_ = winapi.SetForegroundWindow(hwnd)
		return
	}
	activeAboutMu.Unlock()

	className, err := windows.UTF16PtrFromString(fmt.Sprintf("Stendoclip.About.%d", os.Getpid()))
	if err != nil {
		return
	}
	title, _ := windows.UTF16PtrFromString("About Stendoclip")

	var image winapi.GpImage
	if len(imageData) > 0 {
		image, _ = winapi.GdipLoadImageFromBytes(imageData)
	}

	about := &aboutWindow{
		instance:  ownerInstance,
		className: className,
		image:     image,
		text: fmt.Sprintf(
			"Stendoclip %s\n\nKeyboard-first clipboard manager for Windows.\n\nCreated by Matt Moo\nLicense: MIT\nhttps://github.com/the-matt-moo/stendoclip",
			version,
		),
	}

	brush, err := winapi.CreateSolidBrush(winapi.RGB(255, 255, 255))
	if err != nil {
		if image != 0 {
			winapi.GdipDisposeImage(image)
		}
		return
	}
	class := winapi.WndClassEx{
		WndProc:    aboutProcCallback,
		Instance:   ownerInstance,
		ClassName:  className,
		Background: winapi.HBRUSH(brush),
	}
	class.CbSize = uint32(unsafe.Sizeof(class))
	if _, err := winapi.RegisterClassEx(&class); err != nil {
		winapi.DeleteObject(winapi.HGDIOBJ(brush))
		if image != 0 {
			winapi.GdipDisposeImage(image)
		}
		return
	}

	activeAboutMu.Lock()
	activeAbout = about
	activeAboutMu.Unlock()

	monitor := winapi.MonitorFromWindow(0, winapi.MonitorDefaultToNearest)
	info := winapi.MonitorInfo{CbSize: uint32(unsafe.Sizeof(winapi.MonitorInfo{}))}
	x, y := int32(100), int32(100)
	if err := winapi.GetMonitorInfo(monitor, &info); err == nil {
		workW := info.Work.Right - info.Work.Left
		workH := info.Work.Bottom - info.Work.Top
		x = info.Work.Left + (workW-aboutWidth)/2
		y = info.Work.Top + (workH-aboutHeight)/2
	}

	hwnd, err := winapi.CreateWindowEx(
		winapi.WSExDlgModalFrame,
		className, title,
		winapi.WSCaption|winapi.WSSysMenu|winapi.WSPopup,
		x, y, aboutWidth, aboutHeight,
		0, 0, ownerInstance, nil,
	)
	if err != nil {
		activeAboutMu.Lock()
		activeAbout = nil
		activeAboutMu.Unlock()
		_ = winapi.UnregisterClass(className, ownerInstance)
		if image != 0 {
			winapi.GdipDisposeImage(image)
		}
		return
	}
	about.hwnd = hwnd
	winapi.ShowWindow(hwnd, winapi.SWShow)
	_ = winapi.SetForegroundWindow(hwnd)
}

func aboutProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	activeAboutMu.Lock()
	about := activeAbout
	activeAboutMu.Unlock()
	if about == nil {
		return winapi.DefWindowProc(winapi.HWND(hwnd), message, wParam, lParam)
	}

	switch message {
	case winapi.WMPaint:
		about.paint()
		return 0
	case winapi.WMKeyDown:
		if wParam == winapi.VKEscape || wParam == winapi.VKReturn {
			_ = winapi.DestroyWindow(winapi.HWND(hwnd))
		}
		return 0
	case winapi.WMLButtonDown:
		_ = winapi.DestroyWindow(winapi.HWND(hwnd))
		return 0
	case winapi.WMClose:
		_ = winapi.DestroyWindow(winapi.HWND(hwnd))
		return 0
	case winapi.WMDestroy:
		activeAboutMu.Lock()
		if activeAbout != nil && activeAbout.hwnd == winapi.HWND(hwnd) {
			inst := activeAbout.instance
			cls := activeAbout.className
			img := activeAbout.image
			activeAbout = nil
			activeAboutMu.Unlock()
			if img != 0 {
				winapi.GdipDisposeImage(img)
			}
			_ = winapi.UnregisterClass(cls, inst)
		} else {
			activeAboutMu.Unlock()
		}
		return 0
	}
	return winapi.DefWindowProc(winapi.HWND(hwnd), message, wParam, lParam)
}

func (a *aboutWindow) paint() {
	var paint winapi.PaintStruct
	dc := winapi.BeginPaint(a.hwnd, &paint)
	if dc == 0 {
		return
	}
	defer winapi.EndPaint(a.hwnd, &paint)

	client, err := winapi.GetClientRect(a.hwnd)
	if err != nil {
		return
	}

	winapi.SetBkMode(dc, winapi.Transparent)
	winapi.SetTextColor(dc, winapi.RGB(30, 30, 30))
	previousFont := winapi.SelectObject(dc, winapi.GetStockObject(winapi.DefaultGUIFont))
	defer winapi.SelectObject(dc, previousFont)

	// Measure the DPI-scaled system font, then fit artwork into the remaining client area.
	// Windows can provide a smaller client rectangle on differently scaled monitors.
	measure := winapi.Rect{Left: aboutPadding, Right: client.Right - aboutPadding, Bottom: 10000}
	_ = winapi.DrawText(dc, a.text, &measure, winapi.DTCenter|winapi.DTWordBreak|winapi.DTNoPrefix|winapi.DTCalcRect)
	imageRect, textRect := layoutAbout(client, measure.Bottom, a.image != 0)

	if a.image != 0 && imageRect.Right > imageRect.Left && imageRect.Bottom > imageRect.Top {
		_ = winapi.GdipDrawImageRect(dc, a.image, imageRect.Left, imageRect.Top, imageRect.Right-imageRect.Left, imageRect.Bottom-imageRect.Top)
	}
	_ = winapi.DrawText(dc, a.text, &textRect, winapi.DTCenter|winapi.DTWordBreak|winapi.DTNoPrefix)
}

func layoutAbout(client winapi.Rect, textHeight int32, hasImage bool) (image, text winapi.Rect) {
	text.Left, text.Right = aboutPadding, client.Right-aboutPadding
	text.Top, text.Bottom = aboutPadding, client.Bottom-aboutPadding
	if !hasImage {
		return image, text
	}

	availableHeight := client.Bottom - 2*aboutPadding - aboutImageGap - textHeight
	imageHeight := min(aboutImageHeight, max(int32(0), availableHeight))
	imageWidth := imageHeight * aboutImageWidth / aboutImageHeight
	if maxWidth := client.Right - 2*aboutPadding; imageWidth > maxWidth {
		imageWidth = max(int32(0), maxWidth)
		imageHeight = imageWidth * aboutImageHeight / aboutImageWidth
	}
	image.Left = (client.Right - imageWidth) / 2
	image.Top = aboutPadding
	image.Right = image.Left + imageWidth
	image.Bottom = image.Top + imageHeight
	text.Top = image.Bottom + aboutImageGap
	return image, text
}
