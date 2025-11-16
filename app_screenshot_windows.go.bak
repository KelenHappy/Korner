//go:build windows || linux

package main

import (
	"context"
	"fmt"
	"runtime"
)

// captureScreenshot is a platform stub for Windows and Linux.
// It returns a descriptive "not implemented" error until the
// platform-specific implementations are added.
//
// Windows (planned):
//   - Use Win32 GDI (BitBlt) or github.com/kbinani/screenshot to capture the screen.
//   - Crop to the requested rectangle (x, y, width, height).
//   - Encode to PNG and return as a data URL "data:image/png;base64,<...>".
//
// Linux (planned - GNOME/KDE via Flatpak):
//   - Use org.freedesktop.portal.Screenshot (xdg-desktop-portal) to request a screenshot.
//   - The portal provides the UI/permission flow and returns a screenshot file/pipe.
//   - Read the resulting image, encode to base64, and return as a data URL.
//
// NOTE: The macOS implementation exists in a separate file (app_screenshot_darwin.go).
func captureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return "", fmt.Errorf("captureScreenshot not implemented on %s (pending Windows/Linux backends)", runtime.GOOS)
}
