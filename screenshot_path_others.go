//go:build !windows && !darwin

package main

import "fmt"

func getLastScreenshotPath() (string, error) {
	return "", fmt.Errorf("screenshot not supported on this platform")
}
