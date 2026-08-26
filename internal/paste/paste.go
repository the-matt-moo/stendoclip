package paste

import (
	"fmt"
	"time"

	"github.com/mooreceipts/stendoclip/internal/clipboard"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

type Executor struct {
	hwnd    winapi.HWND
	delay   time.Duration
	post    func(func()) error
	onError func(error)
}

func New(hwnd winapi.HWND, delay time.Duration, post func(func()) error, onError func(error)) *Executor {
	return &Executor{hwnd: hwnd, delay: delay, post: post, onError: onError}
}

func (e *Executor) Execute(target winapi.HWND, text string) error {
	if target == 0 {
		return fmt.Errorf("paste target is unavailable")
	}
	if err := winapi.SetForegroundWindow(target); err != nil {
		return fmt.Errorf("restore paste target: %w", err)
	}
	if err := clipboard.WriteText(e.hwnd, text); err != nil {
		return fmt.Errorf("write selected clip: %w", err)
	}
	time.AfterFunc(e.delay, func() {
		if err := e.post(func() {
			if err := send(target); err != nil {
				e.report(err)
			}
		}); err != nil {
			e.report(err)
		}
	})
	return nil
}

func send(target winapi.HWND) error {
	if winapi.GetForegroundWindow() != target {
		return fmt.Errorf("paste target lost focus")
	}
	if err := winapi.SendInput(inputs()); err != nil {
		return fmt.Errorf("send Ctrl+V: %w", err)
	}
	return nil
}

func inputs() []winapi.Input {
	return []winapi.Input{
		winapi.NewKeyboardInput(winapi.VKControl, 0),
		winapi.NewKeyboardInput('V', 0),
		winapi.NewKeyboardInput('V', winapi.KeyEventKeyUp),
		winapi.NewKeyboardInput(winapi.VKControl, winapi.KeyEventKeyUp),
	}
}

func (e *Executor) report(err error) {
	if e.onError != nil {
		e.onError(err)
	}
}
