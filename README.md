# stendoclip

stendoclip is a lightweight Windows 11 clipboard manager inspired by Jumpcut. Phase 1 provides the clip stack, pinning, JSON persistence, configuration, and rotating file logging.

## Build

Requires Go 1.23 or newer.

```sh
make test
make build
make run
make release
```

`make release` produces a GUI-subsystem `stendoclip.exe`. Version metadata comes from `git describe` and is embedded with `-ldflags`.

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
| `hotkey_open` | `Ctrl+Shift+V` | Open the bezel or cycle clips |
| `hotkey_pin` | `Ctrl+P` | Pin selected entry |
| `history_path` | `""` | Custom clipboard-history file path |

Example:

```json
{
  "max_history": 100,
  "hotkey_open": "Ctrl+Alt+V",
  "history_path": "C:\\Users\\me\\Documents\\stendoclip-history.json"
}
```

## License

MIT © 2026 mooreceipts
