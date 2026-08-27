# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.2] - 27-08-2026

### Changed
- Bezel body font range is now 12–48px
- Bezel height scales with body font while remaining inside 90% of monitor work area
- Counter footer uses a fixed 12px font, dimmer color, and added body separation
- Font-size tray actions redraw the visible bezel and reset its timeout
- Font-size tray menu remains available for repeated adjustments

### Fixed
- Font-size changes now work when `config.json` does not yet exist
- Font-size changes persist and apply directly without depending on config hot-reload

## [1.0.1] - 27-08-2026

### Added
- Font size controls in system tray menu: "Increase font size" and "Decrease font size" items
- Dynamic font size persistence to `config.json` (range: 8–96px, ±2px per click)
- Live redraw of bezel overlay when font size changes while popup is visible

### Changed
- Tray menu now includes font size adjustment options with visual feedback

## [1.0.0] - 2026-08-26

### Added
- Initial stable release
- Clipboard history with pinning and deletion
- Bezel overlay with keyboard navigation
- System tray controls (pause, clear, export, startup toggle)
- Hot-reload configuration via `config.json`
- Full key binding customization
- DPI-aware rendering (Win10/11 compatible)
- Markdown export of clipboard history
- Per-monitor DPI detection with manifest
