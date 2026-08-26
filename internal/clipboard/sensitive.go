package clipboard

import (
	"encoding/binary"
	"fmt"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

type sensitiveFormats struct {
	exclude        uint32
	viewerIgnore   uint32
	includeHistory uint32
	uploadToCloud  uint32
}

func registerSensitiveFormats() (sensitiveFormats, error) {
	var formats sensitiveFormats
	var err error
	if formats.exclude, err = winapi.RegisterClipboardFormat("ExcludeClipboardContentFromMonitorProcessing"); err != nil {
		return formats, fmt.Errorf("register exclusion format: %w", err)
	}
	if formats.viewerIgnore, err = winapi.RegisterClipboardFormat("Clipboard Viewer Ignore"); err != nil {
		return formats, fmt.Errorf("register viewer-ignore format: %w", err)
	}
	if formats.includeHistory, err = winapi.RegisterClipboardFormat("CanIncludeInClipboardHistory"); err != nil {
		return formats, fmt.Errorf("register history format: %w", err)
	}
	if formats.uploadToCloud, err = winapi.RegisterClipboardFormat("CanUploadToCloudClipboard"); err != nil {
		return formats, fmt.Errorf("register cloud format: %w", err)
	}
	return formats, nil
}

func isSensitive(formats sensitiveFormats) bool {
	if winapi.IsClipboardFormatAvailable(formats.exclude) || winapi.IsClipboardFormatAvailable(formats.viewerIgnore) {
		return true
	}
	return clipboardDWORDIsZero(formats.includeHistory) || clipboardDWORDIsZero(formats.uploadToCloud)
}

func clipboardDWORDIsZero(format uint32) bool {
	if !winapi.IsClipboardFormatAvailable(format) {
		return false
	}
	handle, err := winapi.GetClipboardData(format)
	if err != nil || winapi.GlobalSize(handle) < 4 {
		return false
	}
	data, err := winapi.GlobalRead(handle, 4)
	return err == nil && binary.LittleEndian.Uint32(data) == 0
}
