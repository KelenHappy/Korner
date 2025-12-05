package audio

import (
	"testing"
	"time"
)

func TestRecorder(t *testing.T) {
	// Create recorder
	recorder, err := NewRecorder()
	if err != nil {
		t.Fatalf("Failed to create recorder: %v", err)
	}
	defer recorder.Close()

	// Start recording
	if err := recorder.StartRecording(); err != nil {
		t.Fatalf("Failed to start recording: %v", err)
	}

	// Check if recording
	if !recorder.IsRecording() {
		t.Error("Recorder should be recording")
	}

	// Record for 2 seconds
	time.Sleep(2 * time.Second)

	// Stop recording
	filePath, err := recorder.StopRecording()
	if err != nil {
		t.Fatalf("Failed to stop recording: %v", err)
	}

	t.Logf("Recording saved to: %s", filePath)

	// Check duration
	duration := recorder.GetDuration()
	if duration < 1.5 || duration > 2.5 {
		t.Errorf("Expected duration around 2 seconds, got %.2f", duration)
	}

	// Check if not recording anymore
	if recorder.IsRecording() {
		t.Error("Recorder should not be recording")
	}
}
