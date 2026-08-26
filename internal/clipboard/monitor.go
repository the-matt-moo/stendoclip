package clipboard

import (
	"errors"
	"fmt"

	"github.com/mooreceipts/stendoclip/internal/store"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

type Monitor struct {
	hwnd        winapi.HWND
	stack       *store.ClippingStack
	historyPath string
	maxBytes    int
	formats     sensitiveFormats
	started     bool
}

func NewMonitor(hwnd winapi.HWND, stack *store.ClippingStack, historyPath string, maxBytes int) (*Monitor, error) {
	formats, err := registerSensitiveFormats()
	if err != nil {
		return nil, err
	}
	return &Monitor{hwnd: hwnd, stack: stack, historyPath: historyPath, maxBytes: maxBytes, formats: formats}, nil
}

func (m *Monitor) Start() error {
	if m.started {
		return nil
	}
	if err := winapi.AddClipboardFormatListener(m.hwnd); err != nil {
		return err
	}
	m.started = true
	return nil
}

func (m *Monitor) Capture() error {
	text, err := ReadText(m.hwnd, m.maxBytes, m.formats)
	if errors.Is(err, ErrNoText) || errors.Is(err, ErrSensitive) || errors.Is(err, ErrTooLarge) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read clipboard: %w", err)
	}
	if !m.stack.Push(text) {
		return nil
	}
	if err := m.stack.Save(m.historyPath); err != nil {
		return fmt.Errorf("save clip stack: %w", err)
	}
	return nil
}

func (m *Monitor) Close() error {
	if !m.started {
		return nil
	}
	m.started = false
	return winapi.RemoveClipboardFormatListener(m.hwnd)
}
