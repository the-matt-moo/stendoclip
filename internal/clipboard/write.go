package clipboard

import (
	"encoding/binary"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func WriteText(hwnd winapi.HWND, text string) (err error) {
	if strings.TrimSpace(text) == "" {
		return ErrNoText
	}
	encoded, err := windows.UTF16FromString(text)
	if err != nil {
		return err
	}
	if err = openWithRetry(hwnd); err != nil {
		return err
	}
	defer func() {
		if closeErr := winapi.CloseClipboard(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	if err = winapi.EmptyClipboard(); err != nil {
		return err
	}

	handle, err := winapi.GlobalAlloc(winapi.GMEMMoveable, uintptr(len(encoded)*2))
	if err != nil {
		return err
	}
	ownedByClipboard := false
	defer func() {
		if !ownedByClipboard {
			winapi.GlobalFree(handle)
		}
	}()
	data := make([]byte, len(encoded)*2)
	for i, unit := range encoded {
		binary.LittleEndian.PutUint16(data[i*2:], unit)
	}
	if err = winapi.GlobalWrite(handle, data); err != nil {
		return err
	}
	if err = winapi.SetClipboardData(winapi.CFUnicodeText, handle); err != nil {
		return err
	}
	ownedByClipboard = true
	return nil
}
