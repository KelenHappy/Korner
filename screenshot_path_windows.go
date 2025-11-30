//go:build windows

package main

import "github.com/Kelen/Korner/internal/platform/windows"

func getLastScreenshotPath() (string, error) {
	return windows.GetLastScreenshotPath()
}
