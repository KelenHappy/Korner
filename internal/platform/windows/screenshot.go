//go:build windows

package windows

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/kbinani/screenshot"
)

var (
	procGetDC            = user32.NewProc("GetDC")
	procReleaseDC        = user32.NewProc("ReleaseDC")
	procGetDeviceCaps    = gdi32.NewProc("GetDeviceCaps")
	procGetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")
	procMonitorFromPoint = user32.NewProc("MonitorFromPoint")
)

const (
	LOGPIXELSX = 88
	LOGPIXELSY = 90
)

// CaptureScreenshot captures a screenshot of the specified region
func CaptureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	log.Printf("DEBUG: Capture called: x=%d, y=%d, width=%d, height=%d\n", x, y, width, height)

	n := screenshot.NumActiveDisplays()
	if n == 0 {
		return "", fmt.Errorf("no active displays found")
	}

	bounds := screenshot.GetDisplayBounds(0)
	displayWidth := bounds.Dx()
	displayHeight := bounds.Dy()

	var img *image.RGBA
	var err error

	if width > 0 && height > 0 {
		x2 := x + width
		y2 := y + height

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

		captureRect := image.Rect(x, y, x2, y2)
		img, err = screenshot.CaptureRect(captureRect)
		if err != nil {
			return "", fmt.Errorf("capture region failed: %w", err)
		}
	} else {
		img, err = screenshot.CaptureDisplay(0)
		if err != nil {
			return "", fmt.Errorf("capture display failed: %w", err)
		}
	}

	if img == nil {
		return "", fmt.Errorf("captured image is nil")
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("encode png: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + b64, nil
}

// GetDPIScale returns the DPI scaling factor
func GetDPIScale() float64 {
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

// GetScreenSize returns the screen dimensions
func GetScreenSize() (int, int) {
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

// LogDPIInfo logs DPI information
func LogDPIInfo() {
	scale := GetDPIScale()
	log.Printf("DEBUG: System DPI Scale: %.2f\n", scale)
}
