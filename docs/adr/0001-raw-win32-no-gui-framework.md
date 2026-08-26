# Raw Win32 API instead of a GUI framework

The bezel and system tray use direct Win32 API calls (`CreateWindowExW`, `Shell_NotifyIconW`, GDI `DrawText`) rather than a Go GUI framework like Fyne or lxn/walk.

The bezel is a single borderless popup showing 3 lines of text and an index counter — no layout engine, widget tree, or GPU rendering needed. Fyne adds ~15MB to the binary and pulls in OpenGL; walk adds CGo complexity. Raw Win32 keeps the binary small, avoids dependencies, and gives full control over the overlay's layered-window transparency and keyboard capture behavior. The project is Windows-only, so cross-platform abstraction has no value.

## Considered Options

- **Fyne**: Cross-platform, but adds OpenGL dependency, inflates binary from ~5MB to ~20MB, and its window model doesn't map cleanly to a `WS_EX_TOOLWINDOW` overlay.
- **lxn/walk**: Windows-native, but wraps COM/OLE and adds CGo build complexity for a UI that's simpler than walk's simplest widget.
- **Wails/webview**: Embeds a browser for the overlay — absurd overhead for 3 lines of text.
