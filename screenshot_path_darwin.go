//go:build darwin

package main

import "github.com/Kelen/Korner/internal/platform/darwin"

func getLastScreenshotPath() (string, error) {
	return darwin.GetLastScreenshotPath()
}
