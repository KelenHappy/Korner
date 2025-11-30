//go:build !windows && !darwin
// +build !windows,!darwin

package common

import (
	"context"
	"fmt"
	"log"
)

// CaptureScreenshot is not implemented on unsupported platforms
func CaptureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return "", fmt.Errorf("screenshot not implemented on this platform")
}

// GetDPIScale returns default DPI scale
func GetDPIScale() float64 {
	return 1.0
}

// GetScreenSize returns default screen size
func GetScreenSize() (int, int) {
	return 1920, 1080
}

// LogDPIInfo logs DPI information
func LogDPIInfo() {
	log.Printf("DEBUG: Platform not supported for DPI detection")
}
