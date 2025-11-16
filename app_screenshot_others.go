//go:build !windows && !darwin
// +build !windows,!darwin

package main

import (
	"context"
	"fmt"
)

func captureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return "", fmt.Errorf("screenshot not implemented on this platform")
}

func getDPIScale() float64 {
	return 1.0
}

func getScreenSize() (int, int) {
	return 1920, 1080
}

func logDPIInfo() {
	// No-op
}
