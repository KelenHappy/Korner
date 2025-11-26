//go:build !windows
// +build !windows

package main

import "log"

// RegisterGlobalHotkey is not implemented on non-Windows platforms yet
func (a *App) RegisterGlobalHotkey() error {
	log.Println("Global hotkey not implemented on this platform")
	return nil
}

// UnregisterGlobalHotkey is not implemented on non-Windows platforms yet
func (a *App) UnregisterGlobalHotkey() {
	log.Println("Global hotkey not implemented on this platform")
}

// disableSnap is a no-op on non-Windows platforms
func disableSnap() {
	// Not needed on non-Windows platforms
}
