# Stendoclip Implementation Plan

## Context

Lightweight Windows clipboard manager in Go, inspired by macOS Jumpcut. Captures plain-text clipboard history, lets user cycle through clips via hotkey + overlay bezel, paste into any app. System tray for management. All design decisions settled via grilling session.

## Architecture

Single-threaded Win32 message pump on main OS thread (`runtime.LockOSThread()`). Business logic in goroutines communicating via channels. Hidden window receives `WM_CLIPBOARDUPDATE` + `WM_HOTKEY`. Overlay is separate `WS_POPUP` window created on demand.

```
Main Thread (locked to OS thread)
  Hidden Window (HWND)
    WM_CLIPBOARDUPDATE -> read clipboard, push to store
    WM_HOTKEY          -> show overlay
    WM_APP+N           -> execute goroutine commands
    tray callbacks     -> menu handling

  Overlay Window (HWND, on demand)
    WM_PAINT   -> GDI draw (dark bezel, text preview, index)
    WM_KEYDOWN -> Up/Down/Enter/Esc/Del/Ctrl+P
    WM_TIMER   -> 5s auto-dismiss

Goroutines:
  Config watcher (fsnotify)
  Paste executor (delay + SendInput)
```

## Dependencies

| Module | Purpose |
|--------|---------|
| `golang.org/x/sys` | Win32 syscall wrappers |
| `github.com/fsnotify/fsnotify` | Config file hot-reload |

No systray library — direct `Shell_NotifyIconW` + `CreatePopupMenu` via `x/sys/windows`.

## Progress

- [x] Phase 1 — Scaffold + Core Data Structures (completed 25-08-2026)
- [x] Phase 2 — Win32 Foundation (completed 25-08-2026)
- [x] Phase 3 — Clipboard Monitoring (completed 26-08-2026)
- [x] Phase 4 — Hotkeys + Overlay + Paste (completed 26-08-2026)
- [x] Phase 5 — System Tray (completed 26-08-2026)
- [x] Phase 6 — Config Hot-Reload + Key Bindings (completed 26-08-2026)
- [x] Phase 7 — Polish + Release (completed 26-08-2026)

## Phases

### Phase 1: Scaffold + Core Data Structures
No Win32. Store, config, logger, tests.

Files: `go.mod`, `cmd/stendoclip/main.go`, `internal/store/stack.go`, `internal/store/stack_test.go`, `internal/store/persist.go`, `internal/store/persist_test.go`, `internal/config/config.go`, `internal/config/config_test.go`, `internal/logger/logger.go`, `Makefile`, `LICENSE`, `README.md`, `.gitignore`

Verify: `go test ./internal/store/... ./internal/config/...` pass, `go build` produces .exe.

### Phase 2: Win32 Foundation — Hidden Window + Message Loop
Single-instance mutex, hidden window, message pump with `MsgWaitForMultipleObjects`.

Files: `internal/winapi/dll.go`, `internal/winapi/types.go`, `internal/winapi/constants.go`, `internal/app/app.go`, `internal/app/instance.go`

Verify: .exe starts silently, visible in Task Manager, second instance shows balloon and exits.

Completed: raw Win32 bindings, hidden top-level owner window, event-backed command queue, `MsgWaitForMultipleObjects` pump, named mutex, and duplicate-instance notification balloon. Verified amd64 tests/vet/release build, 386 ABI tests, primary-process persistence, duplicate exit, and mutex release/relaunch.

### Phase 3: Clipboard Monitoring
`AddClipboardFormatListener`, read text with retry, dedup, persist, sensitive data skip.

Files: `internal/clipboard/monitor.go`, `internal/clipboard/read.go`, `internal/clipboard/write.go`, `internal/clipboard/sensitive.go`

Verify: Copy text in Notepad, `history.json` updates. Dedup works. Password manager clips skipped. >64KB skipped.

Completed: clipboard listener registration, retrying Unicode reads, privacy-marker filtering, size limits, stack deduplication, persistence, and reusable clipboard writes.

### Phase 4: Hotkeys + Overlay + Paste
`Ctrl+Shift+V` opens overlay, arrows cycle, Enter pastes via `SendInput`, Esc cancels, Del deletes, Ctrl+P pins.

Files: `internal/hotkey/hotkey.go`, `internal/overlay/overlay.go`, `internal/overlay/render.go`, `internal/overlay/position.go`, `internal/paste/paste.go`

Verify: Full cycling UX works. Paste lands in previously focused window. Multi-monitor overlay positioning correct.

Completed: configurable global hotkey registration, keyboard cycling with optional wraparound, native GDI bezel positioned on the paste target's monitor, pin/delete persistence, timeout dismissal, clipboard replacement, focus restoration, and synthesized `Ctrl+V`.

### Phase 5: System Tray
Custom .ico, context menu with recent clips, pins submenu, pause/clear/startup/about/quit.

Files: `internal/tray/tray.go`, `internal/tray/menu.go`, `internal/tray/icon.go`, `internal/tray/registry.go`, `assets/watergun_icon.ico`

Verify: Right-click tray shows menu. Click clip pastes. Pause toggle works. Start with Windows writes registry.

Completed: embedded custom tray and executable icons, dynamic recent clips and pins menus, click-to-paste, capture pause, unpinned clip clearing, per-user startup registration, About, and graceful Quit.

### Phase 6: Config Hot-Reload
fsnotify watcher, debounced re-parse, apply changes (max_history trim, hotkey re-register, etc.).

Files: `internal/config/watcher.go`, `internal/config/config.go`

Verify: Edit config.json while running, changes apply without restart.

Completed: fsnotify watcher with 500ms debounce, hot-reload of all settings (max_history, max_entry_bytes, paste_delay, timeout, wraparound, debug_log, hotkey re-register), and fully customizable `keys` config block for all bezel actions (open, previous, next, paste, cancel, delete, pin) with multi-binding support and legacy field fallback.

### Phase 7: Polish + Release
Graceful shutdown, DPI awareness, clipboard self-ownership skip, debounce rapid changes, `-H windowsgui` build flag.

Verify: Full integration test across all features. <15MB memory, zero CPU idle.

Completed: DPI per-monitor v2 manifest, clipboard self-ownership skip (ignores own paste writes), 100ms rapid-change debounce, Windows 10/11 compatibility GUIDs in manifest. All phases complete; version 1.0.0.

## Estimated Size
~2,700 lines Go code + ~530 lines tests. ~30 files total.

## Build
```
make release
```
