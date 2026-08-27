package main

import (
	"errors"
	"os"
	"path/filepath"
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

var version = "dev"

const openHotkeyID = 1

func resolveKeys(settings config.Config) (config.Keys, overlay.KeyBindings, error) {
	keys := settings.ResolvedKeys()
	bindings, err := overlay.ParseBindings(
		keys.Previous, keys.Next, keys.Paste, keys.Cancel, keys.Delete, keys.Pin,
	)
	return keys, bindings, err
}

func main() {
	unlock := app.LockMainThread()
	defer unlock()

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

	application, err := app.New(version)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer application.Close()

	capture, err := clipboard.NewMonitor(application.Window(), clips, historyPath, settings.MaxEntryBytes)
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

	paster := paste.New(application.Window(), time.Duration(settings.PasteDelayMs)*time.Millisecond, application.Post, func(err error) {
		log.Error("paste failed: %v", err)
	})
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
		application.Window(), assets.WatergunIcon, assets.CowImage, clips, historyPath, executable, version, paster.Execute,
		func(value bool) {
			paused = value
			log.Info("capture paused: %v", value)
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

			// Logger.
			log.SetDebug(cfg.DebugLog)

			// Paste delay.
			paster.SetDelay(time.Duration(cfg.PasteDelayMs) * time.Millisecond)

			// Bezel settings.
			bezel.SetTimeout(time.Duration(cfg.TimeoutSecs) * time.Second)
			bezel.SetWraparound(cfg.Wraparound)
			bezel.SetFontSize(int32(cfg.BezelFontSize))

			// Key bindings.
			_, newBindings, err := resolveKeys(cfg)
			if err != nil {
				log.Error("config reload key bindings: %v", err)
			} else {
				bezel.SetBindings(newBindings)
			}

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
