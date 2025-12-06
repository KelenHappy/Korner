package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Kelen/Korner/internal/audio"
)

func main() {
	fmt.Println("=== Audio Recording Test ===")
	fmt.Println("This will test recording in 3 modes:")
	fmt.Println("1. Microphone only")
	fmt.Println("2. System audio only")
	fmt.Println("3. Both mixed (microphone + system audio)")
	fmt.Println()
	
	// List available devices
	fmt.Println("Listing audio devices...")
	devices, err := audio.ListAudioDevices()
	if err != nil {
		fmt.Printf("Error listing devices: %v\n", err)
	} else {
		fmt.Printf("Found %d audio devices:\n", len(devices))
		for i, device := range devices {
			fmt.Printf("  %d. %s\n", i+1, device)
		}
	}
	fmt.Println()
	
	// Test all three modes
	modes := []struct {
		name string
		mode audio.RecordMode
		duration int
	}{
		{"Microphone Only", audio.RecordMicrophone, 3},
		{"System Audio Only", audio.RecordSystem, 3},
		{"Both Mixed", audio.RecordBoth, 5},
	}
	
	for i, test := range modes {
		fmt.Printf("\n=== Test %d/%d: %s ===\n", i+1, len(modes), test.name)
		
		// Create recorder with specific mode
		recorder, err := audio.NewRecorderWithMode(test.mode)
		if err != nil {
			fmt.Printf("Error creating recorder: %v\n", err)
			continue
		}
		
		// Start recording
		fmt.Printf("Starting recording (%d seconds)...\n", test.duration)
		err = recorder.StartRecording()
		if err != nil {
			fmt.Printf("Error starting recording: %v\n", err)
			recorder.Close()
			continue
		}
		
		if test.mode == audio.RecordMicrophone || test.mode == audio.RecordBoth {
			fmt.Println("Please speak into your microphone...")
		}
		if test.mode == audio.RecordSystem || test.mode == audio.RecordBoth {
			fmt.Println("Please play some audio on your computer...")
		}
		
		// Record for specified duration
		for j := 1; j <= test.duration; j++ {
			time.Sleep(1 * time.Second)
			duration := recorder.GetDuration()
			fmt.Printf("  %d/%d seconds (%.1fs recorded)\n", j, test.duration, duration)
		}
		
		// Stop recording
		fmt.Println("Stopping recording...")
		outputPath, err := recorder.StopRecording()
		if err != nil {
			fmt.Printf("Error stopping recording: %v\n", err)
			recorder.Close()
			continue
		}
		
		// Check file size
		if info, err := os.Stat(outputPath); err == nil {
			fmt.Printf("✓ Saved to: %s\n", outputPath)
			fmt.Printf("✓ File size: %d bytes (%.2f KB)\n", info.Size(), float64(info.Size())/1024)
			if info.Size() > 10000 {
				fmt.Println("✓ Recording appears successful!")
			} else {
				fmt.Println("⚠ Warning: File size is very small, recording may have failed")
			}
		}
		
		recorder.Close()
		time.Sleep(500 * time.Millisecond)
	}
	
	fmt.Println("\n=== All tests complete ===")
	fmt.Println("Check the 'record' folder for output files")
}
