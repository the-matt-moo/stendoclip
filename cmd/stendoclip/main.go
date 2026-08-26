package main

import (
	"github.com/mooreceipts/stendoclip/internal/app"
	"github.com/mooreceipts/stendoclip/internal/winapi"
)

var version = "dev"

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

	application, err := app.New(version)
	if err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
		return
	}
	defer application.Close()
	if err := application.Run(); err != nil {
		winapi.MessageBox("Stendoclip", err.Error())
	}
}
