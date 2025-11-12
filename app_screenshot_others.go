//go:build windows

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"syscall"

	"github.com/kbinani/screenshot"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	shcore               = syscall.NewLazyDLL("shcore.dll")
	procGetDC            = user32.NewProc("GetDC")
	procReleaseDC        = user32.NewProc("ReleaseDC")
	procGetDeviceCaps    = syscall.NewLazyDLL("gdi32.dll").NewProc("GetDeviceCaps")
	procGetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")
	procMonitorFromPoint = user32.NewProc("MonitorFromPoint")
)

const (
	LOGPIXELSX = 88
	LOGPIXELSY = 90
)

func captureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return captureScreenshotWindows(ctx, x, y, width, height)
}

// getDPIScale returns the current display DPI scaling factor
func getDPIScale() float64 {
	// Try to get DPI using GetDpiForMonitor (Windows 8.1+)
	hdc, _, _ := procGetDC.Call(0)
	if hdc != 0 {
		defer procReleaseDC.Call(0, hdc)

		dpiX, _, _ := procGetDeviceCaps.Call(hdc, LOGPIXELSX)
		if dpiX != 0 {
			scale := float64(dpiX) / 96.0
			log.Printf("DEBUG: Detected DPI scale: %.2f (DPI: %d)\n", scale, dpiX)
			return scale
		}
	}

	log.Printf("DEBUG: Could not detect DPI scale, using 1.0\n")
	return 1.0
}

// getScreenSize returns the actual physical screen dimensions
func getScreenSize() (int, int) {
	// Get the primary display bounds
	n := screenshot.NumActiveDisplays()
	if n == 0 {
		log.Printf("DEBUG: No displays found, using default 1920x1080\n")
		return 1920, 1080
	}

	bounds := screenshot.GetDisplayBounds(0)
	width := bounds.Dx()
	height := bounds.Dy()
	log.Printf("DEBUG: Screen size detected: %dx%d\n", width, height)
	return width, height
}

// logDPIInfo logs the current DPI scaling for debugging
func logDPIInfo() {
	scale := getDPIScale()
	log.Printf("DEBUG: System DPI Scale: %.2f\n", scale)
}

func captureScreenshotWindows(ctx context.Context, x, y, width, height int) (string, error) {
	log.Printf("DEBUG: captureScreenshotWindows called with physical coords: x=%d, y=%d, width=%d, height=%d\n", x, y, width, height)

	// Get the number of displays
	n := screenshot.NumActiveDisplays()
	if n == 0 {
		return "", fmt.Errorf("no active displays found")
	}

	// Get the primary display bounds
	bounds := screenshot.GetDisplayBounds(0)
	displayWidth := bounds.Dx()
	displayHeight := bounds.Dy()
	log.Printf("DEBUG: Primary display bounds: %+v (width=%d, height=%d)\n", bounds, displayWidth, displayHeight)

	var img *image.RGBA
	var err error

	// Check if this is a full-window capture
	// Frontend passes physical coordinates already scaled by devicePixelRatio
	isFullWindow := (x == 0 && y == 0 && width > int(float64(displayWidth)*0.9) && height > int(float64(displayHeight)*0.9))

	if isFullWindow {
		// Capture entire primary display at native resolution
		log.Printf("DEBUG: Detected full-window capture, capturing entire display at native resolution\n")
		img, err = screenshot.CaptureDisplay(0)
		if err != nil {
			log.Printf("DEBUG: screenshot.CaptureDisplay failed: %v\n", err)
			return "", fmt.Errorf("capture display failed: %w", err)
		}
	} else if width > 0 && height > 0 {
		// Capture specific region using physical coordinates from frontend
		// Frontend already scaled by devicePixelRatio, so use directly

		// Calculate end coordinates
		x2 := x + width
		y2 := y + height

		log.Printf("DEBUG: Region capture - before clamp: x=%d, y=%d, x2=%d, y2=%d (width=%d, height=%d)\n",
			x, y, x2, y2, width, height)

		// Clamp to screen bounds to prevent out-of-bounds capture
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		if x2 > displayWidth {
			x2 = displayWidth
		}
		if y2 > displayHeight {
			y2 = displayHeight
		}

		log.Printf("DEBUG: Region capture - after clamp: x=%d, y=%d, x2=%d, y2=%d\n",
			x, y, x2, y2)

		captureRect := image.Rect(x, y, x2, y2)
		log.Printf("DEBUG: Capturing region: %+v (actual size will be %dx%d)\n", captureRect, x2-x, y2-y)

		img, err = screenshot.CaptureRect(captureRect)
		if err != nil {
			log.Printf("DEBUG: screenshot.CaptureRect failed: %v\n", err)
			return "", fmt.Errorf("capture region failed: %w", err)
		}
	} else {
		// Capture entire primary display
		log.Printf("DEBUG: Capturing entire display (fallback)\n")
		img, err = screenshot.CaptureDisplay(0)
		if err != nil {
			log.Printf("DEBUG: screenshot.CaptureDisplay failed: %v\n", err)
			return "", fmt.Errorf("capture display failed: %w", err)
		}
	}

	if img == nil {
		return "", fmt.Errorf("captured image is nil")
	}

	log.Printf("DEBUG: Captured image size: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())

	// Encode to PNG in memory
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Printf("DEBUG: PNG encoding failed: %v\n", err)
		return "", fmt.Errorf("encode png: %w", err)
	}

	pngBytes := buf.Bytes()
	log.Printf("DEBUG: PNG encoded, size: %d bytes\n", len(pngBytes))

	// Convert to base64 data URL
	b64 := base64.StdEncoding.EncodeToString(pngBytes)
	dataURL := "data:image/png;base64," + b64

	log.Printf("DEBUG: Returning data URL with %d chars\n", len(dataURL))
	return dataURL, nil
}
