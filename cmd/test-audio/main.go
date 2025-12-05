package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kelen/Korner/internal/audio"
)

func main() {
	fmt.Println("=== Korner Audio Recorder Test ===")
	fmt.Println()

	// Create recorder
	recorder, err := audio.NewRecorder()
	if err != nil {
		log.Fatalf("Failed to create recorder: %v", err)
	}
	defer recorder.Close()

	fmt.Println("Press ENTER to start recording...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Start recording
	if err := recorder.StartRecording(); err != nil {
		log.Fatalf("Failed to start recording: %v", err)
	}

	fmt.Println("🎤 Recording... Press ENTER to stop")

	// Show duration while recording
	go func() {
		for recorder.IsRecording() {
			duration := recorder.GetDuration()
			fmt.Printf("\r⏱️  Duration: %.1f seconds", duration)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Stop recording
	fmt.Println("\n\n⏹️  Stopping recording...")
	filePath, err := recorder.StopRecording()
	if err != nil {
		log.Fatalf("Failed to stop recording: %v", err)
	}

	fmt.Printf("✅ Recording saved to: %s\n", filePath)
	fmt.Printf("📊 Duration: %.2f seconds\n", recorder.GetDuration())
}
