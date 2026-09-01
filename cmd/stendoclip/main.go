package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mooreceipts/stendoclip/assets"
	"github.com/mooreceipts/stendoclip/internal/app"
	"github.com/mooreceipts/stendoclip/internal/clipboard"
	"github.com/mooreceipts/stendoclip/internal/config"
	"github.com/mooreceipts/stendoclip/internal/hotkey"
	"github.com/mooreceipts/stendoclip/internal/logger"
	"github.com/mooreceipts/stendoclip/internal/overlay"
	"github.com/mooreceipts/stendoclip/internal/paste"
	"github.com/mooreceipts/stendoclip/internal/store"
	"github.com/mooreceipts/stendoclip/internal/tray"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

// version is overridden from VERSION by the Makefile for release builds.
var version = "1.0.8"

const openHotkeyID = 1

func resolveKeys(settings config.Config) (config.Keys, overlay.KeyBindings, error) {
	keys := settings.ResolvedKeys()
	bindings, err := overlay.ParseBindings(
		keys.Previous, keys.Next, keys.Paste, keys.Cancel, keys.Delete, keys.Pin,
	)
	return keys, bindings, err
}

func buildAboutText(version string, keys config.Keys) string {
	return strings.Join([]string{
		fmt.Sprintf("Stendoclip %s", version),
		"",
		"Keyboard-first clipboard manager for Windows.",
		"",
		"Keybindings:",
		"Open bezel: " + keys.Open[0],
		"Previous clip: " + strings.Join(keys.Previous, " / "),
		"Next clip: " + strings.Join(keys.Next, " / "),
		"Paste clip: " + strings.Join(keys.Paste, " / "),
		"Cancel: " + strings.Join(keys.Cancel, " / "),
		"Delete clip: " + strings.Join(keys.Delete, " / "),
		"Pin clip: " + strings.Join(keys.Pin, " / "),
		"",
		"Created by Matt Moo",
		"License: MIT",
		"https://github.com/the-matt-moo/stendoclip",
	}, "\n")
}

func main() {
	unlock := app.LockMainThread()
	defer unlock()
	if err := winapi.EnablePerMonitorDPI(); err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}

	instance, alreadyRunning, err := app.AcquireInstance()
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	if alreadyRunning {
		if err := app.ShowAlreadyRunningNotification(); err != nil {
			winapi.MessageBox("Stendoclip", err.Error())
		}
		return
	}
	defer instance.Close()

	dataDir, err := os.UserConfigDir()
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	dataDir = filepath.Join(dataDir, "Stendoclip")
	configPath := filepath.Join(dataDir, "config.json")
	settings, err := config.LoadConfig(configPath)
	if errors.Is(err, os.ErrNotExist) {
		settings = config.Defaults()
	} else if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}

	log, err := logger.New(filepath.Join(dataDir, "stendoclip.log"))
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer log.Close()
	log.SetDebug(settings.DebugLog)

	historyPath := settings.HistoryPath
	if historyPath == "" {
		historyPath = filepath.Join(dataDir, "history.json")
	} else if !filepath.IsAbs(historyPath) {
		historyPath = filepath.Join(dataDir, historyPath)
	}
	clips, err := store.Load(historyPath)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	clips.SetMax(settings.MaxHistory)
	clips.SetMaxBytes(settings.MaxEntryBytes)
	markdownExportPath := settings.MarkdownExportPath
	if markdownExportPath != "" && !filepath.IsAbs(markdownExportPath) {
		markdownExportPath = filepath.Join(dataDir, markdownExportPath)
	}

	application, err := app.New(version)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer application.Close()

	capture, err := clipboard.NewMonitor(application.Window(), clips, historyPath, settings.MaxEntryBytes, settings.TrimWhitespace)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	paused := false
	application.SetClipboardHandler(func() {
		if paused {
			return
		}
		if err := capture.Capture(); err != nil {
			log.Error("capture failed: %v", err)
		}
	})
	if err := capture.Start(); err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer capture.Close()

	keys, bindings, err := resolveKeys(settings)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	aboutText := buildAboutText(version, keys)

	paster := paste.New(application.Window(), time.Duration(settings.PasteDelayMs)*time.Millisecond, application.Post, func(err error) {
		log.Error("paste failed: %v", err)
	}, settings.TrimWhitespace)
	bezel, err := overlay.New(
		clips, historyPath, bindings, time.Duration(settings.TimeoutSecs)*time.Second, settings.Wraparound,
		int32(settings.BezelFontSize),
		paster.Execute, func(err error) { log.Error("bezel failed: %v", err) },
	)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer bezel.Close()

	openSpec := keys.Open[0]
	if _, err := hotkey.Register(application.Window(), openHotkeyID, openSpec); err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer winapi.UnregisterHotKey(application.Window(), openHotkeyID)

	application.SetHotkeyHandler(func(id uintptr) {
		if id == openHotkeyID {
			if err := bezel.Open(); err != nil {
				log.Error("open bezel failed: %v", err)
			}
		}
	})

	executable, err := os.Executable()
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	// GDI+ for About window PNG rendering.
	gdipToken, err := winapi.GdiplusStartup()
	if err != nil {
		log.Error("gdiplus init: %v", err)
	} else {
		defer winapi.GdiplusShutdown(gdipToken)
	}

	systemTray, err := tray.New(
		application.Window(), assets.WatergunIcon, assets.CowImage, clips, historyPath, markdownExportPath, configPath, executable, version, aboutText, settings.BezelFontSize, settings.TrimWhitespace, paster.Execute,
		func(value bool) {
			paused = value
			log.Info("capture paused: %v", value)
		},
		func(size int) {
			bezel.SetFontSize(int32(size))
			log.Info("font size changed to %d", size)
		},
		func(enabled bool) {
			capture.SetTrimWhitespace(enabled)
			paster.SetTrimWhitespace(enabled)
			log.Info("trim whitespace: %v", enabled)
		},
		application.Quit,
		func(err error) { log.Error("tray failed: %v", err) },
	)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer systemTray.Close()
	application.SetMessageHandler(systemTray.HandleMessage)

	// Config hot-reload watcher.
	watcher, err := config.Watch(configPath, 500*time.Millisecond, func(cfg config.Config) {
		application.Post(func() {
			// Stack limits.
			clips.SetMax(cfg.MaxHistory)
			clips.SetMaxBytes(cfg.MaxEntryBytes)
			capture.SetMaxBytes(cfg.MaxEntryBytes)

			// Logger and export location.
			log.SetDebug(cfg.DebugLog)
			exportPath := cfg.MarkdownExportPath
			if exportPath != "" && !filepath.IsAbs(exportPath) {
				exportPath = filepath.Join(dataDir, exportPath)
			}
			systemTray.SetMarkdownExportPath(exportPath)

			// Paste delay.
			paster.SetDelay(time.Duration(cfg.PasteDelayMs) * time.Millisecond)
			capture.SetTrimWhitespace(cfg.TrimWhitespace)
			paster.SetTrimWhitespace(cfg.TrimWhitespace)
			systemTray.SetTrimWhitespace(cfg.TrimWhitespace)

			// Bezel settings.
			bezel.SetTimeout(time.Duration(cfg.TimeoutSecs) * time.Second)
			bezel.SetWraparound(cfg.Wraparound)
			bezel.SetFontSize(int32(cfg.BezelFontSize))
			systemTray.SetFontSize(cfg.BezelFontSize)

			// Key bindings.
			_, newBindings, err := resolveKeys(cfg)
			if err != nil {
				log.Error("config reload key bindings: %v", err)
			} else {
				bezel.SetBindings(newBindings)
			}
			systemTray.SetAboutText(buildAboutText(version, cfg.ResolvedKeys()))

			// Global hotkey re-register.
			newKeys := cfg.ResolvedKeys()
			newOpen := newKeys.Open[0]
			if newOpen != openSpec {
				winapi.UnregisterHotKey(application.Window(), openHotkeyID)
				if _, err := hotkey.Register(application.Window(), openHotkeyID, newOpen); err != nil {
					log.Error("config reload hotkey: %v", err)
					// Re-register old one as fallback.
					hotkey.Register(application.Window(), openHotkeyID, openSpec)
				} else {
					log.Info("hotkey changed to %s", newOpen)
				}
			}

			log.Info("config reloaded")
		})
	}, func(err error) {
		log.Error("config watcher: %v", err)
	})
	if err != nil {
		// Non-fatal: watcher failure doesn't prevent running.
		log.Error("config watcher failed to start: %v", err)
	} else {
		defer watcher.Close()
	}

	log.Info("started %s", version)
	if err := application.Run(); err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
	}
}
