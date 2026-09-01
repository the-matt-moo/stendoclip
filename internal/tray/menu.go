package tray

import (
	"fmt"
	"strings"

	"github.com/mooreceipts/stendoclip/internal/config"
	"github.com/mooreceipts/stendoclip/internal/export"
	"github.com/mooreceipts/stendoclip/internal/store"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

const (
	recentBase = 1000
	pinsBase   = 2000
	pauseID    = 3001
	clearID    = 3002
	startupID  = 3003
	aboutID    = 3004
	exportID   = 3005
	trimID     = 3006
	incFontID  = 3007
	decFontID  = 3008
	quitID     = 3009
	menuLimit  = 10
)

type menuState struct {
	menu   winapi.HMENU
	recent []store.Entry
	pins   []store.Entry
}

func (t *Tray) buildMenu() (state menuState, err error) {
	state.menu, err = winapi.CreatePopupMenu()
	if err != nil {
		return state, err
	}
	failed := true
	defer func() {
		if failed {
			winapi.DestroyMenu(state.menu)
		}
	}()

	if err = winapi.AppendMenu(state.menu, winapi.MFString|winapi.MFGray, 0, "Recent clips"); err != nil {
		return state, err
	}
	state.recent = t.stack.HistoryEntries()
	if len(state.recent) > menuLimit {
		state.recent = state.recent[:menuLimit]
	}
	if len(state.recent) == 0 {
		if err = winapi.AppendMenu(state.menu, winapi.MFString|winapi.MFGray, 0, "(No recent clips)"); err != nil {
			return state, err
		}
	}
	for i, entry := range state.recent {
		if err = winapi.AppendMenu(state.menu, winapi.MFString, recentBase+uintptr(i), clipLabel(entry.Text)); err != nil {
			return state, err
		}
	}

	state.pins = t.stack.PinnedEntries()
	if len(state.pins) > menuLimit {
		state.pins = state.pins[:menuLimit]
	}
	pinsMenu, err := winapi.CreatePopupMenu()
	if err != nil {
		return state, err
	}
	if len(state.pins) == 0 {
		if err = winapi.AppendMenu(pinsMenu, winapi.MFString|winapi.MFGray, 0, "(No pins)"); err != nil {
			winapi.DestroyMenu(pinsMenu)
			return state, err
		}
	}
	for i, entry := range state.pins {
		if err = winapi.AppendMenu(pinsMenu, winapi.MFString, pinsBase+uintptr(i), clipLabel(entry.Text)); err != nil {
			winapi.DestroyMenu(pinsMenu)
			return state, err
		}
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFPopup, uintptr(pinsMenu), "Pins"); err != nil {
		winapi.DestroyMenu(pinsMenu)
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFSeparator, 0, ""); err != nil {
		return state, err
	}

	pauseFlags := uint32(winapi.MFString)
	if t.paused {
		pauseFlags |= winapi.MFChecked
	}
	if err = winapi.AppendMenu(state.menu, pauseFlags, pauseID, "Pause capture"); err != nil {
		return state, err
	}
	clearFlags := uint32(winapi.MFString)
	if len(state.recent) == 0 {
		clearFlags |= winapi.MFGray
	}
	if err = winapi.AppendMenu(state.menu, clearFlags, clearID, "Clear recent clips"); err != nil {
		return state, err
	}

	enabled, startupErr := startupEnabled(t.executable)
	if startupErr != nil {
		t.report(fmt.Errorf("read startup setting: %w", startupErr))
	}
	startupFlags := uint32(winapi.MFString)
	if enabled {
		startupFlags |= winapi.MFChecked
	}
	if err = winapi.AppendMenu(state.menu, startupFlags, startupID, "Start with Windows"); err != nil {
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFSeparator, 0, ""); err != nil {
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFString, incFontID, "Increase font size"); err != nil {
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFString, decFontID, "Decrease font size"); err != nil {
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFSeparator, 0, ""); err != nil {
		return state, err
	}
	trimFlags := uint32(winapi.MFString)
	if t.trimWhitespace {
		trimFlags |= winapi.MFChecked
	}
	if err = winapi.AppendMenu(state.menu, trimFlags, trimID, "Trim leading/trailing whitespace"); err != nil {
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFString, exportID, "Export history to Markdown..."); err != nil {
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFString, aboutID, "About Stendoclip"); err != nil {
		return state, err
	}
	if err = winapi.AppendMenu(state.menu, winapi.MFString, quitID, "Quit"); err != nil {
		return state, err
	}
	failed = false
	return state, nil
}

func (t *Tray) runCommand(command uint32, state menuState, target winapi.HWND) {
	switch {
	case command >= recentBase && int(command-recentBase) < len(state.recent):
		t.report(t.paste(target, state.recent[command-recentBase].Text))
	case command >= pinsBase && int(command-pinsBase) < len(state.pins):
		t.report(t.paste(target, state.pins[command-pinsBase].Text))
	case command == pauseID:
		t.paused = !t.paused
		if t.onPause != nil {
			t.onPause(t.paused)
		}
	case command == clearID:
		t.stack.ClearHistory()
		t.report(t.stack.Save(t.historyPath))
	case command == startupID:
		enabled, err := startupEnabled(t.executable)
		if err == nil {
			err = setStartup(t.executable, !enabled)
		}
		t.report(err)
	case command == trimID:
		t.trimWhitespace = !t.trimWhitespace
		if t.onTrimChange != nil {
			t.onTrimChange(t.trimWhitespace)
		}
		t.report(config.UpdateTrimWhitespace(t.configPath, t.trimWhitespace))
	case command == exportID:
		path := t.markdownExportPath
		if path == "" {
			selectedPath, ok, err := winapi.SaveMarkdownPath(t.hwnd)
			if err != nil {
				t.report(err)
				return
			}
			if !ok {
				return
			}
			path = selectedPath
		}
		t.report(export.ToMarkdown(t.stack, path))
	case command == incFontID:
		t.adjustFontSize(2)
	case command == decFontID:
		t.adjustFontSize(-2)
	case command == aboutID:
		showAbout(t.instance, t.aboutImage, t.aboutText)
	case command == quitID:
		if t.onQuit != nil {
			t.report(t.onQuit())
		}
	}
}

func (t *Tray) adjustFontSize(delta int) {
	newSize := t.currentFontSize + delta
	if newSize < 12 {
		newSize = 12
	}
	if newSize > 96 {
		newSize = 96
	}
	if newSize == t.currentFontSize {
		return
	}
	t.currentFontSize = newSize
	if t.onFontSizeChange != nil {
		t.onFontSizeChange(newSize)
	}
	t.report(config.UpdateFontSize(t.configPath, newSize))
}

func clipLabel(text string) string {
	label := strings.Join(strings.Fields(text), " ")
	if label == "" {
		label = "(Empty clip)"
	}
	runes := []rune(label)
	if len(runes) > 60 {
		label = string(runes[:59]) + "…"
	}
	return strings.ReplaceAll(label, "&", "&&")
}
