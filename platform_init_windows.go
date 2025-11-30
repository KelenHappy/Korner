//go:build windows

package main

import (
	"log"

	"github.com/Kelen/Korner/internal/platform"
)

func initPlatform() {
	// Disable Windows Snap feature
	if err := platform.DisableWindowSnap("Korner - AI Assistant"); err != nil {
		log.Printf("Failed to disable Windows Snap: %v", err)
	}
}
