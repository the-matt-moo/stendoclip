# stendoclip

![Stendoclip social card](assets/social_card.jpg)

stendoclip is a lightweight Windows 11 clipboard manager inspired by Jumpcut. It captures plain-text clipboard changes and provides a keyboard-driven bezel plus system-tray controls for cycling, pinning, deleting, and pasting clips. Configuration is hot-reloaded — edit `config.json` while running and changes apply instantly.

## Controls

| Key | Action |
|---|---|
| `Ctrl+Shift+V` | Open the bezel; press again to cycle forward |
| Arrow keys | Cycle clips |
| `Enter` | Paste the selected clip |
| `Esc` | Cancel |
| `Delete` | Delete the selected clip |
| `Ctrl+P` | Pin or unpin the selected clip |

All keys are configurable (see [Key bindings](#key-bindings) below).

## System tray

Right-click the watergun tray icon to paste recent clips or pins, pause capture, clear unpinned clips, toggle Start with Windows, view version information, or quit. The icon is embedded in both the tray and compiled executable from `assets/watergun_icon.ico`, so no external asset file is required at runtime.

### Export history to markdown

Right-click the tray icon, select **Export history to Markdown...**, then choose any output path. Or set `markdown_export_path` in `config.json` to export directly to that path without prompting. Relative paths resolve from `%AppData%\\Stendoclip\\`. The export contains full clip text, pinned status, and timestamps. Use a file inside an Obsidian vault for backups, cloud sync, or knowledge-base import.

## Build

Requires Go 1.23 or newer.

```sh
make test
make build
make run
make release
```

`make release` produces a GUI-subsystem `stendoclip.exe`. Version metadata comes from `VERSION` and is embedded with `-ldflags`. Run `make resources` only after replacing the `.ico`; it regenerates committed 386 and amd64 Windows resources.

## Configuration

Configuration is JSON. Unspecified fields retain these defaults:

| Field | Default | Purpose |
|---|---:|---|
| `max_history` | `50` | Maximum unpinned clips in the stack |
| `max_entry_bytes` | `65536` | Maximum UTF-8 bytes per clipping |
| `paste_delay_ms` | `200` | Delay before simulated paste |
| `timeout_secs` | `5` | Selection overlay timeout |
| `wraparound` | `true` | Wrap when cycling entries |
| `debug_log` | `false` | Enable debug messages |
| `hotkey_open` | `Ctrl+Shift+V` | Open the bezel or cycle clips (legacy; prefer `keys.open`) |
| `hotkey_pin` | `Ctrl+P` | Pin selected entry (legacy; prefer `keys.pin`) |
| `history_path` | `""` | Custom clipboard-history file path |
| `markdown_export_path` | `""` | Default Markdown export path; blank prompts on export |
| `bezel_font_size` | `18` | Bezel text size in pixels |

Configuration lives at `%AppData%\\Stendoclip\\config.json`.

### User customizations

| Option | Hot-reload | Notes |
|---|---|---|
| `max_history` | yes | Max unpinned clips kept in memory/history |
| `max_entry_bytes` | yes | Clip size cap, 1..10MB guardrail |
| `paste_delay_ms` | yes | Delay before simulated Ctrl+V |
| `timeout_secs` | yes | Bezel auto-close timeout |
| `wraparound` | yes | Wrap clip navigation |
| `debug_log` | yes | Extra logging |
| `bezel_font_size` | yes | Bezel text size in pixels, min 8 |
| `keys` | yes | Full key-binding override for bezel actions |
| `hotkey_open` / `hotkey_pin` | yes | Legacy single-hotkey fallback only |
| `history_path` | restart | Move persisted history/pins file |
| `markdown_export_path` | yes | Default export target; blank opens Save dialog |

Clips and pins are persisted to `%AppData%\\Stendoclip\\history.json` (customisable via `history_path`) and survive restarts. Pinned clips are exempt from the `max_history` cap. Relative `history_path` values resolve from the config directory.

All settings except `history_path` are hot-reloaded: save the file while Stendoclip is running and changes take effect within one second. Invalid config is logged and ignored.

Example:

```json
{
  "max_history": 100,
  "max_entry_bytes": 262144,
  "paste_delay_ms": 150,
  "timeout_secs": 8,
  "wraparound": false,
  "debug_log": true,
  "bezel_font_size": 22,
  "history_path": "C:\\Users\\me\\Documents\\stendoclip-history.json",
  "markdown_export_path": "C:\\Users\\me\\Obsidian\\Clipboard History.md",
  "keys": {
    "open": ["Ctrl+Alt+V"],
    "previous": ["K", "Up"],
    "next": ["J", "Down"],
    "cancel": ["Q", "Escape"]
  }
}
```

### What stendoclip captures

Stendoclip captures `CF_UNICODETEXT` only:
- **Images, files, formatted content**: silently ignored (no-op, not stored)
- **Text larger than `max_entry_bytes`**: rejected, not stored
- **Clips marked sensitive**: Windows/KeePass privacy formats are detected and ignored

All rejection cases are silent — nothing is logged unless `debug_log: true`.


## Key bindings

All bezel keys can be overridden via an optional `keys` object. Each action accepts an array of key specs. Omitted actions keep their defaults. Key specs use `Modifier+Key` format (e.g. `Ctrl+Shift+V`, `Alt+F12`, `Enter`, `Up`).

Defaults:

| Action | Default keys |
|---|---|
| `open` | `["Ctrl+Shift+V"]` |
| `previous` | `["Up", "Left"]` |
| `next` | `["Down", "Right"]` |
| `paste` | `["Enter"]` |
| `cancel` | `["Escape"]` |
| `delete` | `["Delete"]` |
| `pin` | `["Ctrl+P"]` |

Example override in `config.json`:

```json
{
  "keys": {
    "open": ["Ctrl+Alt+V"],
    "previous": ["K", "Up"],
    "next": ["J", "Down"],
    "cancel": ["Q", "Escape"]
  }
}
```

When `keys` is present, the legacy `hotkey_open` and `hotkey_pin` fields are ignored. Supported key names: `A`–`Z`, `0`–`9`, `F1`–`F12`, `Up`, `Down`, `Left`, `Right`, `Enter`, `Escape`, `Delete`, `Backspace`, `Space`, `Tab`, `Home`, `End`. Modifiers: `Ctrl`, `Shift`, `Alt`, `Win`.

## FAQ

**What happens when I copy an image or file?**

- Stendoclip ignores non-text content silently. Only plain Unicode text (`CF_UNICODETEXT`) is captured. Images, files, formatted text, and other formats are discarded without logging (unless `debug_log: true`).

**Where are my clips stored?**

- In `%AppData%\Stendoclip\history.json` by default. Pinned clips and unpinned clips coexist in the same file and survive restarts. Use `history_path` in config to move it elsewhere.

**Can I sync clips between computers?**

- Not built-in. You can manually copy `history.json` between machines, or use `history_path` to point to a cloud-synced folder (OneDrive, Dropbox, etc.), but stendoclip doesn't auto-sync.

**Is my clipboard data private?**

- Stendoclip stores clips as plaintext JSON locally. Only your machine can read it. Sensitive data (Windows privacy-marked content, KeePass) is detected and not stored. If you don't want a clip captured, manually clear it or pause capture via the tray menu.

**How do I uninstall?**

- Delete `stendoclip.exe`. Optionally remove the config and history from `%AppData%\Stendoclip\`. To disable autostart, toggle "Start with Windows" in the tray menu or manually delete the registry entry at `HKCU\Software\Microsoft\Windows\CurrentVersion\Run\Stendoclip`.

**Can I backup my clips?**

- Yes. `history.json` is a valid JSON file. Copy it anywhere or add it to version control. To restore, replace the file in `%AppData%\Stendoclip\` and restart stendoclip.

**What if my config is invalid?**

- Stendoclip logs the error and ignores the bad config on startup/reload. It falls back to defaults and uses the last known good config in memory. Check the log at `%AppData%\Stendoclip\stendoclip.log` for details.

**Can I customize every hotkey?**

- Yes. Use the `keys` object in `config.json` to override any bezel action (open, previous, next, paste, cancel, delete, pin). Each accepts an array of key specs for multi-binding.

**Does stendoclip consume a lot of memory?**

- No. Runtime is ~15MB for a 50-clip history. Memory grows linearly with `max_history` and `max_entry_bytes`; at defaults, a 500-clip stack is still <50MB.

## License

MIT © 2026 mooreceipts
