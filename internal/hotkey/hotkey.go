package hotkey

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

type Binding struct {
	Modifiers uint32
	Key       uint32
}

func Parse(spec string) (Binding, error) {
	var binding Binding
	parts := strings.Split(spec, "+")
	for _, raw := range parts {
		part := strings.ToUpper(strings.TrimSpace(raw))
		switch part {
		case "CTRL", "CONTROL":
			if has(binding.Modifiers, winapi.MODControl) {
				return Binding{}, fmt.Errorf("hotkey %q repeats Ctrl", spec)
			}
			binding.Modifiers |= winapi.MODControl
		case "SHIFT":
			if has(binding.Modifiers, winapi.MODShift) {
				return Binding{}, fmt.Errorf("hotkey %q repeats Shift", spec)
			}
			binding.Modifiers |= winapi.MODShift
		case "ALT":
			if has(binding.Modifiers, winapi.MODAlt) {
				return Binding{}, fmt.Errorf("hotkey %q repeats Alt", spec)
			}
			binding.Modifiers |= winapi.MODAlt
		case "WIN", "WINDOWS":
			if has(binding.Modifiers, winapi.MODWin) {
				return Binding{}, fmt.Errorf("hotkey %q repeats Win", spec)
			}
			binding.Modifiers |= winapi.MODWin
		default:
			if binding.Key != 0 {
				return Binding{}, fmt.Errorf("hotkey %q has multiple keys", spec)
			}
			key, ok := parseKey(part)
			if !ok {
				return Binding{}, fmt.Errorf("unsupported hotkey key %q", raw)
			}
			binding.Key = key
		}
	}
	if binding.Key == 0 {
		return Binding{}, fmt.Errorf("hotkey %q has no key", spec)
	}
	return binding, nil
}

func Register(hwnd winapi.HWND, id int32, spec string) (Binding, error) {
	binding, err := Parse(spec)
	if err != nil {
		return Binding{}, err
	}
	if err := winapi.RegisterHotKey(hwnd, id, binding.Modifiers|winapi.MODNoRepeat, binding.Key); err != nil {
		return Binding{}, fmt.Errorf("register hotkey %q: %w", spec, err)
	}
	return binding, nil
}

func Matches(binding Binding, key uint32) bool {
	if binding.Key != key {
		return false
	}
	var modifiers uint32
	if winapi.GetKeyState(winapi.VKControl) {
		modifiers |= winapi.MODControl
	}
	if winapi.GetKeyState(winapi.VKShift) {
		modifiers |= winapi.MODShift
	}
	if winapi.GetKeyState(winapi.VKMenu) {
		modifiers |= winapi.MODAlt
	}
	if winapi.GetKeyState(winapi.VKLWin) || winapi.GetKeyState(winapi.VKRWin) {
		modifiers |= winapi.MODWin
	}
	return modifiers == binding.Modifiers
}

func parseKey(key string) (uint32, bool) {
	if len(key) == 1 && ((key[0] >= 'A' && key[0] <= 'Z') || (key[0] >= '0' && key[0] <= '9')) {
		return uint32(key[0]), true
	}
	if strings.HasPrefix(key, "F") {
		n, err := strconv.Atoi(strings.TrimPrefix(key, "F"))
		if err == nil && n >= 1 && n <= 12 {
			return winapi.VKF1 + uint32(n-1), true
		}
	}
	keys := map[string]uint32{
		"BACKSPACE": winapi.VKBack,
		"DELETE":    winapi.VKDelete,
		"ENTER":     winapi.VKReturn,
		"ESC":       winapi.VKEscape,
		"ESCAPE":    winapi.VKEscape,
		"SPACE":     winapi.VKSpace,
		"TAB":       winapi.VKTab,
		"UP":        winapi.VKUp,
		"DOWN":      winapi.VKDown,
		"LEFT":      winapi.VKLeft,
		"RIGHT":     winapi.VKRight,
		"HOME":      winapi.VKHome,
		"END":       winapi.VKEnd,
	}
	value, ok := keys[key]
	return value, ok
}

func has(value, flag uint32) bool { return value&flag != 0 }
