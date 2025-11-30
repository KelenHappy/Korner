//go:build darwin

package darwin

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CaptureScreenshot captures a screenshot on macOS
// If width and height are positive, captures the specified rectangle
// Otherwise, opens the interactive system UI for selection
func CaptureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("korner-%d.png", time.Now().UnixNano()))
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

// GetDPIScale returns the DPI scale for macOS
// macOS handles DPI scaling automatically, so we return 1.0
func GetDPIScale() float64 {
	return 1.0
}

// GetScreenSize returns the screen dimensions for macOS
func GetScreenSize() (int, int) {
	// TODO: Implement using Cocoa APIs to get actual screen dimensions
	return 1920, 1080
}

// LogDPIInfo logs DPI information
func LogDPIInfo() {
	scale := GetDPIScale()
	log.Printf("DEBUG: System DPI Scale: %.2f\n", scale)
}
