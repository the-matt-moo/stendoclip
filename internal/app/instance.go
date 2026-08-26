package app

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

const instanceMutexName = `Local\Stendoclip.SingleInstance.5C04078D-ABF7-43D3-8CC7-E64D9479AA19`

type Instance struct {
	handle windows.Handle
}

func AcquireInstance() (*Instance, bool, error) {
	return acquireInstance(instanceMutexName)
}

func acquireInstance(name string) (*Instance, bool, error) {
	namePtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, false, err
	}
	handle, err := windows.CreateMutex(nil, false, namePtr)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		if handle != 0 {
			_ = windows.CloseHandle(handle)
		}
		return nil, true, nil
	}
	if err != nil {
		if handle != 0 {
			_ = windows.CloseHandle(handle)
		}
		return nil, false, fmt.Errorf("create instance mutex: %w", err)
	}
	return &Instance{handle: handle}, false, nil
}

func (i *Instance) Close() error {
	if i == nil || i.handle == 0 {
		return nil
	}
	err := windows.CloseHandle(i.handle)
	i.handle = 0
	return err
}

func ShowAlreadyRunningNotification() error {
	instance, err := winapi.GetModuleHandle()
	if err != nil {
		return err
	}
	className, _ := windows.UTF16PtrFromString("STATIC")
	windowName, _ := windows.UTF16PtrFromString("Stendoclip")
	hwnd, err := winapi.CreateWindowEx(0, className, windowName, 0, 0, 0, 0, 0, 0, 0, instance, nil)
	if err != nil {
		return err
	}
	defer winapi.DestroyWindow(hwnd)

	icon, err := winapi.LoadIcon(0, winapi.IDIInformation)
	if err != nil {
		return err
	}
	data := winapi.NotifyIconData{
		HWnd:      hwnd,
		UID:       1,
		Flags:     winapi.NIFIcon | winapi.NIFTip | winapi.NIFInfo,
		Icon:      icon,
		InfoFlags: winapi.NIIFInfo,
	}
	data.CbSize = uint32(unsafe.Sizeof(data))
	copyUTF16(data.Tip[:], "Stendoclip")
	copyUTF16(data.InfoTitle[:], "Stendoclip")
	copyUTF16(data.Info[:], "Stendoclip is already running.")
	if err := winapi.ShellNotifyIcon(winapi.NIMAdd, &data); err != nil {
		return err
	}
	defer winapi.ShellNotifyIcon(winapi.NIMDelete, &data)
	data.TimeoutVersion = winapi.NotifyIconVersion4
	_ = winapi.ShellNotifyIcon(winapi.NIMSetVersion, &data)
	time.Sleep(3 * time.Second)
	return nil
}

func copyUTF16(destination []uint16, value string) {
	encoded, _ := windows.UTF16FromString(value)
	if len(encoded) > len(destination) {
		encoded = encoded[:len(destination)]
		encoded[len(encoded)-1] = 0
	}
	copy(destination, encoded)
}
