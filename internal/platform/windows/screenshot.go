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
	"os"
	"path/filepath"
	"time"

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

// CaptureScreenshot captures a screenshot and saves it to build/bin directory
// Returns both the file path and base64 data URL
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

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("encode png: %w", err)
	}

	// Save to build/bin directory
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("WARNING: Could not get executable path: %v", err)
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	
	// Create screenshots directory if it doesn't exist
	screenshotsDir := filepath.Join(exeDir, "screenshots")
	if err := os.MkdirAll(screenshotsDir, 0755); err != nil {
		log.Printf("WARNING: Could not create screenshots directory: %v", err)
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("screenshot_%d.png", time.Now().UnixNano())
	filePath := filepath.Join(screenshotsDir, filename)

	// Save file
	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		log.Printf("WARNING: Could not save screenshot file: %v", err)
	} else {
		log.Printf("Screenshot saved to: %s", filePath)
	}

	// Return base64 data URL for frontend display
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + b64, nil
}

// GetLastScreenshotPath returns the path to the most recent screenshot
func GetLastScreenshotPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	screenshotsDir := filepath.Join(exeDir, "screenshots")

	entries, err := os.ReadDir(screenshotsDir)
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("no screenshots found")
	}

	// Get the most recent file
	var latestFile string
	var latestTime time.Time

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
			latestFile = filepath.Join(screenshotsDir, entry.Name())
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("no valid screenshots found")
	}

	return latestFile, nil
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
