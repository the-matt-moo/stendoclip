# Implementation Plan — stendoclip review fixes (from 9b684f0...HEAD review)

Goal: fix the code-review findings without changing user-visible behavior beyond bug fixes.
No new dependencies. After each step: `go build ./... && go vet ./... && go test ./...` must pass.
Do NOT commit at the end — report the diff summary instead.

## 1. Trim must apply before the size check (Spec finding)

Currently `monitor.go` trims AFTER `ReadText` already rejected over-limit text, so a clip that
trimmed under `maxBytes` but measures over untrimmed (e.g. large whitespace padding) is dropped as
`ErrTooLarge`.

- `internal/clipboard/read.go`
  - Change `ReadText(hwnd winapi.HWND, maxBytes int, formats sensitiveFormats)` to add a
    trailing `trim bool` parameter.
  - Change `decodeText(units []uint16, maxBytes int)` → `decodeText(units []uint16, maxBytes int, trim bool)`.
  - In `decodeText`, after `text := windows.UTF16ToString(units[:end])`, insert
    `if trim { text = strings.TrimSpace(text) }` BEFORE the blank check and size check.
- `internal/clipboard/monitor.go` `Capture()`:
  - Call `ReadText(m.hwnd, m.maxBytes, m.formats, m.trimWhitespace)`.
  - Delete the now-redundant `if m.trimWhitespace { text = strings.TrimSpace(text) }` block.
- `internal/clipboard/read_test.go`: update the three `decodeText(units, N)` calls to pass
  `false` as the new trim arg. Add a trim-before-size case:
  `decodeText(utf16.Encode([]rune("  a  ")), 3, true)` must return `"a"` (raw 5 bytes > 3, trimmed 1 byte <= 3),
  and the same call with `trim=false` must return `ErrTooLarge`.

## 2. Deduplicate the blank-text check (Standards: Duplicated Code)

The same `strings.TrimSpace(text) == ""` shape exists in 4 places. Export the one already in
`store` and reuse it. Use the CONTEXT.md vocabulary (Clip, not Entry) for the name.

- `internal/store/stack.go`: rename `isBlankText(text string) bool` →
  `func IsBlankClip(text string) bool` (exported). Update its use in `Push`.
- `internal/store/persist.go`: rename `filterEmptyEntries` → `filterBlankClips`
  (2 call sites in `Load`); it calls `IsBlankClip` instead of `isBlankText`.
- `internal/clipboard/read.go`: replace the blank check with `store.IsBlankClip(text)`
  (package `clipboard` already imports `store` via monitor.go — same package, no new import).
- `internal/clipboard/write.go`: `if store.IsBlankClip(text) { return ErrNoText }`;
  add the `store` import.
- `internal/paste/paste.go`: in `Execute`, after `text = strings.TrimSpace(text)`,
  use `store.IsBlankClip(text)` for the `return nil` guard; add the `store` import.
  No import cycle: `store` imports no internal packages.
- Keep the blank-clip REJECTION semantics identical everywhere (read → ErrNoText,
  write → ErrNoText, push → false, paste → nil).
- Rename the test helpers/names added in this diff that use "Entry" vocabulary for
  blank-clip behavior: `TestLoadDropsBlankEntries` → `TestLoadDropsBlankClips`,
  `TestRejectsBlankOrOversizedEntries` → `TestRejectsBlankOrOversizedClips`
  (check which files these live in; update only what exists).

## 3. About window: class registration must retry after failure (Spec finding)

`getAboutClass` fail-locks on `sync.Once` — a transient first failure (RegisterClassEx /
CreateSolidBrush) breaks the About window for the whole session.

- `internal/tray/about.go`:
  - Remove `aboutClassOnce sync.Once` and `aboutClassErr error` package vars.
  - Add `aboutClassMu sync.Mutex` (keep `aboutClassName *uint16`).
  - Rewrite `getAboutClass(instance)` to lock the mutex, return early if
    `aboutClassName != nil`, otherwise register the class and set `aboutClassName`
    ONLY on success. On any failure cleanup (delete the brush on RegisterClassEx
    failure) and return the error — the next call retries because the latch is unset.
  - Keep the comment about never unregistering the class (still true).
- `internal/tray/about_test.go`: existing `TestAboutCanBeReopened` must keep passing
  unchanged. No new test needed (the failure path can't be forced without mocking winapi).

## 4. No code change: load-time blank cleanup (Spec: scope creep note)

`Load()` filtering blank entries out of an existing history file is undocumented data
mutation. Keep the behavior (it enforces 1.0.6's "empty clips are neither stored nor
pasted" for pre-1.0.6 files — removing it would resurrect blanks in the stack) but
document it in the CHANGELOG entry for this release.

## 5. No code change: trim checkbox on hot-reload

The tray menu is rebuilt on every right-click (`buildMenu`), and the checkbox reads the
current `trimWhitespace` at build time. Nothing to fix.

## 6. Release chores (do in this order, do NOT commit)

- `VERSION`: `1.0.6` → `1.0.7`.
- `CHANGELOG.md`: add on top:
  `## [1.0.7] - 31-08-2026` with a `### Fixed` section:
  - Whitespace trimming now applies before the size limit check, so clips that trim under
    the limit are no longer dropped (Capture).
  - About window class registration retries on failure instead of failing for the session.
  - Existing history files are cleaned of blank entries when loaded (documents prior behavior).
- `README.md`: verify nothing there describes behavior this change alters; if accurate, leave
  untouched.
- Final report: `git status` + `git diff --stat`, tests/vet/build results, and a ready-to-paste
  commit message (e.g. `fix: trim before size check, shared blank-clip check, retry about class (1.0.7)`).