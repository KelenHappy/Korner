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
	"io"
)

// CaptureScreenshot captures a screenshot on macOS and saves it to build/bin directory
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

	// Copy file to screenshots directory
	if err := copyFile(tmpPath, filePath); err != nil {
		log.Printf("WARNING: Could not save screenshot file: %v", err)
	} else {
		log.Printf("Screenshot saved to: %s", filePath)
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return "data:image/png;base64," + b64, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
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
