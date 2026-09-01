package clipboard

import (
	"errors"
	"fmt"
	"time"

	"github.com/mooreceipts/stendoclip/internal/store"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

const debounceInterval = 100 * time.Millisecond

type Monitor struct {
	hwnd           winapi.HWND
	stack          *store.ClippingStack
	historyPath    string
	maxBytes       int
	trimWhitespace bool
	formats        sensitiveFormats
	started        bool
	lastCapture    time.Time
}

func NewMonitor(hwnd winapi.HWND, stack *store.ClippingStack, historyPath string, maxBytes int, trimWhitespace bool) (*Monitor, error) {
	formats, err := registerSensitiveFormats()
	if err != nil {
		return nil, err
	}
	return &Monitor{hwnd: hwnd, stack: stack, historyPath: historyPath, maxBytes: maxBytes, trimWhitespace: trimWhitespace, formats: formats}, nil
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

func (m *Monitor) SetMaxBytes(n int) { m.maxBytes = n }

func (m *Monitor) SetTrimWhitespace(enabled bool) { m.trimWhitespace = enabled }

func (m *Monitor) Capture() error {
	// Skip our own clipboard writes.
	if winapi.GetClipboardOwner() == m.hwnd {
		return nil
	}
	// Debounce rapid clipboard changes.
	now := time.Now()
	if now.Sub(m.lastCapture) < debounceInterval {
		return nil
	}
	m.lastCapture = now

	text, err := ReadText(m.hwnd, m.maxBytes, m.formats, m.trimWhitespace)
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
