//go:build windows
// +build windows

package main

import (
	"log"
	"unsafe"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	registerHotKey   = user32Common.NewProc("RegisterHotKey")
	unregisterHotKey = user32Common.NewProc("UnregisterHotKey")
	getMessage       = user32Common.NewProc("GetMessageW")
)

const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008
	WM_HOTKEY   = 0x0312
)

// RegisterGlobalHotkey registers Ctrl+Alt+Q as global hotkey
func (a *App) RegisterGlobalHotkey() error {
	// Register Ctrl+Alt+Q (Q = 0x51)
	hotkeyID := 1
	ret, _, err := registerHotKey.Call(
		0,                            // NULL window handle (global)
		uintptr(hotkeyID),            // hotkey ID
		uintptr(MOD_CONTROL|MOD_ALT), // modifiers
		uintptr(0x51),                // Q key
	)

	if ret == 0 {
		log.Printf("Failed to register hotkey: %v", err)
		return err
	}

	log.Println("Global hotkey Ctrl+Alt+Q registered successfully")

	// Start listening for hotkey in a goroutine
	go a.listenForHotkey()

	return nil
}

// listenForHotkey listens for the global hotkey press
func (a *App) listenForHotkey() {
	var msg struct {
		hwnd    uintptr
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}

	for {
		ret, _, _ := getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)

		if ret == 0 {
			break
		}

		if msg.message == WM_HOTKEY {
			log.Println("Hotkey Ctrl+Alt+Q pressed!")
			// Trigger pie menu
			if a.ctx != nil {
				wailsruntime.WindowShow(a.ctx)
				wailsruntime.WindowSetAlwaysOnTop(a.ctx, true)
				wailsruntime.EventsEmit(a.ctx, "hotkey-triggered")
			}
		}
	}
}

// UnregisterGlobalHotkey unregisters the global hotkey
func (a *App) UnregisterGlobalHotkey() {
	hotkeyID := 1
	unregisterHotKey.Call(0, uintptr(hotkeyID))
	log.Println("Global hotkey unregistered")
}
