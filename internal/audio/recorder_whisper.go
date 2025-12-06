package audio

import (
	"fmt"
	"time"
)

// RecorderWithTranscription combines audio recording with Whisper transcription
type RecorderWithTranscription struct {
	recorder    *Recorder
	transcriber *WhisperTranscriber
}

// NewRecorderWithTranscription creates a new recorder with transcription capability
func NewRecorderWithTranscription(modelName string) (*RecorderWithTranscription, error) {
	// Create recorder
	recorder, err := NewRecorder()
	if err != nil {
		return nil, fmt.Errorf("failed to create recorder: %w", err)
	}
	
	// Create transcriber with auto-detection
	transcriber, err := NewWhisperTranscriberAuto(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to create transcriber: %w", err)
	}
	
	return &RecorderWithTranscription{
		recorder:    recorder,
		transcriber: transcriber,
	}, nil
}

// StartRecording starts recording audio
func (r *RecorderWithTranscription) StartRecording() error {
	return r.recorder.StartRecording()
}

// StopRecordingAndTranscribe stops recording and transcribes the audio
func (r *RecorderWithTranscription) StopRecordingAndTranscribe(options TranscribeOptions) (audioPath string, transcription string, err error) {
	// Stop recording
	audioPath, err = r.recorder.StopRecording()
	if err != nil {
		return "", "", fmt.Errorf("failed to stop recording: %w", err)
	}
	
	// Wait a bit for file to be fully written
	time.Sleep(100 * time.Millisecond)
	
	// Transcribe
	transcription, err = r.transcriber.Transcribe(audioPath, options)
	if err != nil {
		return audioPath, "", fmt.Errorf("failed to transcribe: %w", err)
	}
	
	return audioPath, transcription, nil
}

// StopRecordingAndTranscribeWithProgress stops recording and transcribes with progress callback
func (r *RecorderWithTranscription) StopRecordingAndTranscribeWithProgress(options TranscribeOptions, progressCallback func(string)) (audioPath string, transcription string, err error) {
	// Stop recording
	audioPath, err = r.recorder.StopRecording()
	if err != nil {
		return "", "", fmt.Errorf("failed to stop recording: %w", err)
	}
	
	// Wait a bit for file to be fully written
	time.Sleep(100 * time.Millisecond)
	
	// Transcribe with progress
	transcription, err = r.transcriber.TranscribeWithProgress(audioPath, options, progressCallback)
	if err != nil {
		return audioPath, "", fmt.Errorf("failed to transcribe: %w", err)
	}
	
	return audioPath, transcription, nil
}

// IsRecording returns whether currently recording
func (r *RecorderWithTranscription) IsRecording() bool {
	return r.recorder.IsRecording()
}

// GetDuration returns recording duration in seconds
func (r *RecorderWithTranscription) GetDuration() float64 {
	return r.recorder.GetDuration()
}

// Close cleans up resources
func (r *RecorderWithTranscription) Close() error {
	return r.recorder.Close()
}

// RecordAndTranscribe is a convenience function that records for a duration and transcribes
func RecordAndTranscribe(duration time.Duration, modelName string, language string) (audioPath string, transcription string, err error) {
	// Create recorder with transcription
	recorder, err := NewRecorderWithTranscription(modelName)
	if err != nil {
		return "", "", err
	}
	defer recorder.Close()
	
	// Start recording
	if err := recorder.StartRecording(); err != nil {
		return "", "", err
	}
	
	fmt.Printf("Recording for %.0f seconds...\n", duration.Seconds())
	
	// Wait for duration
	time.Sleep(duration)
	
	// Stop and transcribe
	options := DefaultTranscribeOptions()
	options.Language = language
	
	audioPath, transcription, err = recorder.StopRecordingAndTranscribe(options)
	if err != nil {
		return audioPath, "", err
	}
	
	return audioPath, transcription, nil
}
