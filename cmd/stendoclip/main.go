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
	settings, err := config.LoadConfig(filepath.Join(dataDir, "config.json"))
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

	paster := paste.New(application.Window(), time.Duration(settings.PasteDelayMs)*time.Millisecond, application.Post, func(err error) {
		log.Error("paste failed: %v", err)
	})
	bezel, err := overlay.New(
		clips, historyPath, settings.HotkeyPin, time.Duration(settings.TimeoutSecs)*time.Second, settings.Wraparound,
		paster.Execute, func(err error) { log.Error("bezel failed: %v", err) },
	)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer bezel.Close()
	if _, err := hotkey.Register(application.Window(), openHotkeyID, settings.HotkeyOpen); err != nil {
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
	systemTray, err := tray.New(
		application.Window(), assets.WatergunIcon, clips, historyPath, executable, version, paster.Execute,
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

	log.Info("started %s", version)
	if err := application.Run(); err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
	}
}
