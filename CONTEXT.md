# Stendoclip

A lightweight Windows clipboard manager that captures plain-text clipboard history and lets users cycle through clips via a hotkey-activated overlay bezel.

## Language

**Clip**:
A single plain-text clipboard entry captured from the system clipboard.
_Avoid_: Clipping, snippet, item, entry

**Stack**:
The ordered collection of clips, newest first, with a configurable maximum size.
_Avoid_: History, buffer, ring

**Bezel**:
The semi-transparent overlay popup that appears when the user invokes the hotkey, showing the current clip preview and position index.
_Avoid_: Overlay, popup, tooltip, window

**Pin**:
A clip marked for permanent retention, stored separately from the stack and exempt from the max-size cap.
_Avoid_: Bookmark, favorite, saved

**Cycling**:
Navigating forward or backward through the stack while the bezel is visible, with wraparound at both ends.
_Avoid_: Scrolling, browsing, navigating

**Paste Target**:
The window that had focus before the bezel was invoked — the destination for the selected clip.
_Avoid_: Target window, foreground window, destination

**Capture**:
The act of reading the system clipboard after a change event and pushing the text onto the stack.
_Avoid_: Monitor, watch, listen, record

**Pause**:
A user-toggled state where clipboard change events are received but ignored — no capture occurs.
_Avoid_: Suspend, disable, mute
