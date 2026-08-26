# Single-threaded Win32 message pump with goroutine channels

The main thread is locked to the OS thread (`runtime.LockOSThread()`) and runs a single Win32 message loop handling both `WM_CLIPBOARDUPDATE` and `WM_HOTKEY`. Business logic (config reload, paste execution) runs in goroutines that communicate back to the main thread via channels + a Windows Event object that wakes `MsgWaitForMultipleObjects`.

Win32 requires clipboard listener and hotkey registration on the same thread as the message pump. Running two separate message loops (one for clipboard, one for hotkeys) would require cross-thread synchronization and duplicated pump logic. A single hidden window receiving both message types is simpler and matches how native Windows apps work. Goroutines handle the async work (fsnotify, paste delays) without blocking the pump.
