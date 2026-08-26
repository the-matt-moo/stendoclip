package tray

import (
	"fmt"
	"os"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func loadIcon(data []byte) (winapi.HICON, error) {
	file, err := os.CreateTemp("", "stendoclip-*.ico")
	if err != nil {
		return 0, err
	}
	path := file.Name()
	defer os.Remove(path)
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return 0, err
	}
	if err := file.Close(); err != nil {
		return 0, err
	}
	icon, err := winapi.LoadIconFile(path)
	if err != nil {
		return 0, fmt.Errorf("load tray icon: %w", err)
	}
	return icon, nil
}
