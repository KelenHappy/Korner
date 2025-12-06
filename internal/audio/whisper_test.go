package audio

import (
	"fmt"
	"testing"
)

// TestWhisperTranscriber tests the Whisper transcriber
func TestWhisperTranscriber(t *testing.T) {
	// Skip if whisper not available
	whisperPath, err := FindWhisperExecutable()
	if err != nil {
		t.Skip("Whisper executable not found, skipping test")
	}
	
	modelPath, err := FindWhisperModel("tiny")
	if err != nil {
		t.Skip("Whisper tiny model not found, skipping test")
	}
	
	t.Logf("Using whisper: %s", whisperPath)
	t.Logf("Using model: %s", modelPath)
	
	// Create transcriber
	transcriber, err := NewWhisperTranscriber(modelPath, whisperPath)
	if err != nil {
		t.Fatalf("Failed to create transcriber: %v", err)
	}
	
	// Test with a sample audio file (you need to provide one)
	audioPath := "../../record/recording_20251206_202101.wav"
	
	options := DefaultTranscribeOptions()
	options.Language = "zh" // Chinese
	
	text, err := transcriber.Transcribe(audioPath, options)
	if err != nil {
		t.Fatalf("Failed to transcribe: %v", err)
	}
	
	t.Logf("Transcription: %s", text)
	
	if text == "" {
		t.Error("Transcription is empty")
	}
}

// TestWhisperAutoDetect tests auto-detection of whisper and model
func TestWhisperAutoDetect(t *testing.T) {
	transcriber, err := NewWhisperTranscriberAuto("tiny")
	if err != nil {
		t.Skip("Whisper or model not found, skipping test")
	}
	
	t.Logf("Transcriber created successfully")
	t.Logf("Model: %s", transcriber.modelPath)
	t.Logf("Whisper: %s", transcriber.whisperPath)
}

// TestQuickTranscribe tests the quick transcribe function
func TestQuickTranscribe(t *testing.T) {
	modelPath, err := FindWhisperModel("tiny")
	if err != nil {
		t.Skip("Whisper tiny model not found, skipping test")
	}
	
	audioPath := "../../record/recording_20251206_202101.wav"
	
	text, err := QuickTranscribe(audioPath, modelPath)
	if err != nil {
		t.Fatalf("Failed to quick transcribe: %v", err)
	}
	
	t.Logf("Quick transcription: %s", text)
}

// ExampleWhisperTranscriber demonstrates basic usage
func ExampleWhisperTranscriber() {
	// Create transcriber with auto-detection
	transcriber, err := NewWhisperTranscriberAuto("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Transcribe audio file
	options := DefaultTranscribeOptions()
	options.Language = "en"
	
	text, err := transcriber.Transcribe("audio.wav", options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Transcription: %s\n", text)
}

// ExampleQuickTranscribe demonstrates quick transcription
func ExampleQuickTranscribe() {
	text, err := QuickTranscribe("audio.wav", "models/ggml-tiny.bin")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Transcription: %s\n", text)
}
