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
| `bezel_font_size` | `18` | Bezel text size in pixels |

Configuration lives at `%AppData%\\Stendoclip\\config.json`. History defaults to `%AppData%\\Stendoclip\\history.json`; relative `history_path` values resolve from that directory.

All settings are hot-reloaded: save the file while Stendoclip is running and changes take effect within one second. Invalid config is logged and ignored.

Example:

```json
{
  "max_history": 100,
  "history_path": "C:\\Users\\me\\Documents\\stendoclip-history.json"
}
```

Stendoclip captures `CF_UNICODETEXT` only. It ignores clips larger than `max_entry_bytes` and clips marked with Windows or KeePass clipboard privacy formats.

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

## License

MIT © 2026 mooreceipts
