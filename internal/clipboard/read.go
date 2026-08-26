package clipboard

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sys/windows"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

var (
	ErrNoText    = errors.New("clipboard has no text")
	ErrSensitive = errors.New("clipboard content is marked sensitive")
	ErrTooLarge  = errors.New("clipboard text exceeds size limit")
)

func ReadText(hwnd winapi.HWND, maxBytes int, formats sensitiveFormats) (text string, err error) {
	if maxBytes < 1 {
		maxBytes = 65536
	}
	if err = openWithRetry(hwnd); err != nil {
		return "", err
	}
	defer func() {
		if closeErr := winapi.CloseClipboard(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if isSensitive(formats) {
		return "", ErrSensitive
	}
	if !winapi.IsClipboardFormatAvailable(winapi.CFUnicodeText) {
		return "", ErrNoText
	}
	handle, err := winapi.GetClipboardData(winapi.CFUnicodeText)
	if err != nil {
		return "", err
	}
	size := winapi.GlobalSize(handle)
	if size == 0 {
		return "", ErrNoText
	}
	unitCount := size / 2
	limit := uintptr(maxBytes)
	if limit < ^uintptr(0) {
		limit++
	}
	if unitCount > limit {
		unitCount = limit
	}
	data, err := winapi.GlobalRead(handle, unitCount*2)
	if err != nil {
		return "", err
	}
	units := make([]uint16, len(data)/2)
	for i := range units {
		units[i] = binary.LittleEndian.Uint16(data[i*2:])
	}
	return decodeText(units, maxBytes)
}

func decodeText(units []uint16, maxBytes int) (string, error) {
	end := 0
	for end < len(units) && units[end] != 0 {
		end++
	}
	text := windows.UTF16ToString(units[:end])
	if text == "" {
		return "", ErrNoText
	}
	if len([]byte(text)) > maxBytes {
		return "", ErrTooLarge
	}
	return text, nil
}

func openWithRetry(hwnd winapi.HWND) error {
	var err error
	for attempt := 0; attempt < 5; attempt++ {
		if err = winapi.OpenClipboard(hwnd); err == nil {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("open clipboard after retries: %w", err)
}
