//go:build darwin

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// getDPIScale returns the DPI scale for macOS
// macOS handles DPI scaling automatically, so we return 1.0
func getDPIScale() float64 {
	return 1.0
}

// getScreenSize returns the actual physical screen dimensions for macOS
// For now, returns a default value - should be implemented using Cocoa APIs
func getScreenSize() (int, int) {
	// TODO: Implement using Cocoa APIs to get actual screen dimensions
	return 1920, 1080
}

func captureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
// captureScreenshot is the macOS (darwin) implementation that captures either a specified rectangle
// or, if width/height are not positive, opens the interactive system UI for selection.
//
// It returns a data URL string: "data:image/png;base64,<...>"
func captureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return captureScreenshotDarwin(ctx, x, y, width, height)
}

// captureScreenshotDarwin performs the actual screencapture call on macOS.
//
// Behavior:
// - Region mode: if width > 0 && height > 0, captures the rectangle at (x,y,width,height) without UI.
// - Interactive mode: otherwise, opens the native selection UI (-i).
// Notes:
// - Requires the "screencapture" utility available on macOS.
// - Uses a temporary PNG file, which is removed after reading.
// - Adds "-x" to suppress sounds and "-t png" to enforce PNG format.
// - Returns a data URL string for convenient frontend usage.
func captureScreenshotDarwin(ctx context.Context, x, y, width, height int) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// Apply a sensible timeout for the operation (especially for interactive mode).
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("snapask-%d.png", time.Now().UnixNano()))
	// Ensure cleanup even if we fail later.
	defer func() { _ = os.Remove(tmpPath) }()

	args := []string{
		"-x", // no sounds
		"-t", "png",
	}

	if width > 0 && height > 0 {
		// Region capture without UI
		args = append(args, "-R", fmt.Sprintf("%d,%d,%d,%d", x, y, width, height))
	} else {
		// Interactive selection UI
		args = append(args, "-i")
	}

	args = append(args, tmpPath)

	cmd := exec.CommandContext(ctx, "screencapture", args...)
	// Do not inherit stdin/stdout/stderr to keep it quiet
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("screencapture failed: %w", err)
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("read screenshot file: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return "data:image/png;base64," + b64, nil
}
