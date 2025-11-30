//go:build windows
// +build windows

package windows

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	procGetWindowLongW = user32.NewProc("GetWindowLongW")
	procSetWindowLongW = user32.NewProc("SetWindowLongW")
	procFindWindowW    = user32.NewProc("FindWindowW")
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
		log.Printf("Window not found: %s", windowTitle)
		return nil
	}

	// Get current window style
	style, _, _ := procGetWindowLongW.Call(hwnd, GWL_STYLE)

	// Remove WS_THICKFRAME and WS_MAXIMIZEBOX to disable snap
	newStyle := style &^ WS_THICKFRAME &^ WS_MAXIMIZEBOX

	procSetWindowLongW.Call(hwnd, GWL_STYLE, newStyle)
	log.Printf("Disabled Windows Snap for: %s", windowTitle)

	return nil
}
