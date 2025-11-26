//go:build windows
// +build windows

package main

import (
	"syscall"
	"unsafe"
)

// Shared Windows DLL instances
var (
	user32Common       = syscall.NewLazyDLL("user32.dll")
	shcore             = syscall.NewLazyDLL("shcore.dll")
	gdi32              = syscall.NewLazyDLL("gdi32.dll")
	procGetWindowLongW = user32Common.NewProc("GetWindowLongW")
	procSetWindowLongW = user32Common.NewProc("SetWindowLongW")
	procFindWindowW    = user32Common.NewProc("FindWindowW")
)

const (
	GWL_STYLE      = ^uintptr(15) // -16 as uintptr
	WS_THICKFRAME  = 0x00040000
	WS_MAXIMIZEBOX = 0x00010000
)

// DisableWindowSnap disables Windows Snap feature for the given window
func DisableWindowSnap(windowTitle string) error {
	titlePtr, _ := syscall.UTF16PtrFromString(windowTitle)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return nil
	}

	// Get current window style
	style, _, _ := procGetWindowLongW.Call(hwnd, GWL_STYLE)

	// Remove WS_THICKFRAME and WS_MAXIMIZEBOX to disable snap
	newStyle := style &^ WS_THICKFRAME &^ WS_MAXIMIZEBOX

	procSetWindowLongW.Call(hwnd, GWL_STYLE, newStyle)

	return nil
}


// disableSnap is called to disable Windows Snap for this app
func disableSnap() {
	DisableWindowSnap("Korner - AI Assistant")
}
