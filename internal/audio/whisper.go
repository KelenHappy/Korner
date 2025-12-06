package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WhisperTranscriber handles audio transcription using whisper.cpp
type WhisperTranscriber struct {
	modelPath   string
	whisperPath string
}

// NewWhisperTranscriber creates a new Whisper transcriber
// modelPath: path to the ggml model file (e.g., "models/ggml-tiny.bin")
// whisperPath: path to whisper.cpp executable (optional, will search in PATH if empty)
func NewWhisperTranscriber(modelPath string, whisperPath string) (*WhisperTranscriber, error) {
	// Find whisper executable
	if whisperPath == "" {
		// Try to find in common locations
		candidates := []string{
			"whisper",
			"main",
			"./whisper",
			"./main",
			"./whisper.cpp/main",
			"./build/bin/main",
		}
		
		var found bool
		for _, candidate := range candidates {
			if path, err := exec.LookPath(candidate); err == nil {
				whisperPath = path
				found = true
				break
			}
		}
		
		if !found {
			return nil, fmt.Errorf("whisper.cpp executable not found. Please provide whisperPath or add to PATH")
		}
	}
	
	// Verify whisper executable exists
	if _, err := os.Stat(whisperPath); err != nil {
		return nil, fmt.Errorf("whisper executable not found at %s: %w", whisperPath, err)
	}
	
	// Verify model file exists
	if _, err := os.Stat(modelPath); err != nil {
		return nil, fmt.Errorf("model file not found at %s: %w", modelPath, err)
	}
	
	return &WhisperTranscriber{
		modelPath:   modelPath,
		whisperPath: whisperPath,
	}, nil
}

// TranscribeOptions contains options for transcription
type TranscribeOptions struct {
	Language    string  // Language code (e.g., "en", "zh", "auto" for auto-detect)
	Threads     int     // Number of threads to use (0 = auto)
	Translate   bool    // Translate to English
	OutputTxt   bool    // Output .txt file
	OutputSrt   bool    // Output .srt subtitle file
	OutputVtt   bool    // Output .vtt subtitle file
	OutputJson  bool    // Output .json file
	MaxLen      int     // Maximum segment length in characters (0 = no limit)
	SplitOnWord bool    // Split on word rather than token
	Temperature float64 // Temperature for sampling (0.0 = greedy, higher = more random)
}

// DefaultTranscribeOptions returns default transcription options
func DefaultTranscribeOptions() TranscribeOptions {
	return TranscribeOptions{
		Language:    "auto",
		Threads:     0,
		Translate:   false,
		OutputTxt:   true,
		OutputSrt:   false,
		OutputVtt:   false,
		OutputJson:  false,
		MaxLen:      0,
		SplitOnWord: true,
		Temperature: 0.0,
	}
}

// Transcribe transcribes an audio file and returns the text
func (w *WhisperTranscriber) Transcribe(audioPath string, options TranscribeOptions) (string, error) {
	// Verify audio file exists
	if _, err := os.Stat(audioPath); err != nil {
		return "", fmt.Errorf("audio file not found: %w", err)
	}
	
	// Build command arguments
	args := []string{
		"-m", w.modelPath,
		"-f", audioPath,
	}
	
	// Add language option
	if options.Language != "" && options.Language != "auto" {
		args = append(args, "-l", options.Language)
	}
	
	// Add threads option
	if options.Threads > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", options.Threads))
	}
	
	// Add translate option
	if options.Translate {
		args = append(args, "-tr")
	}
	
	// Add output format options
	if options.OutputTxt {
		args = append(args, "-otxt")
	}
	if options.OutputSrt {
		args = append(args, "-osrt")
	}
	if options.OutputVtt {
		args = append(args, "-ovtt")
	}
	if options.OutputJson {
		args = append(args, "-oj")
	}
	
	// Add max length option
	if options.MaxLen > 0 {
		args = append(args, "-ml", fmt.Sprintf("%d", options.MaxLen))
	}
	
	// Add split on word option
	if options.SplitOnWord {
		args = append(args, "-sow")
	}
	
	// Add temperature option
	if options.Temperature > 0 {
		args = append(args, "-t", fmt.Sprintf("%.2f", options.Temperature))
	}
	
	// Run whisper.cpp
	cmd := exec.Command(w.whisperPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper.cpp failed: %w\nOutput: %s", err, string(output))
	}
	
	// Parse output to extract transcription
	text := w.parseWhisperOutput(string(output))
	
	return text, nil
}

// TranscribeToFile transcribes an audio file and saves to output files
func (w *WhisperTranscriber) TranscribeToFile(audioPath string, options TranscribeOptions) error {
	_, err := w.Transcribe(audioPath, options)
	return err
}

// parseWhisperOutput extracts the transcription text from whisper.cpp output
func (w *WhisperTranscriber) parseWhisperOutput(output string) string {
	lines := strings.Split(output, "\n")
	var transcription strings.Builder
	
	// Look for lines that contain transcription (usually after timestamps)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and info lines
		if line == "" || strings.HasPrefix(line, "whisper_") || 
		   strings.HasPrefix(line, "system_info") || strings.HasPrefix(line, "main:") {
			continue
		}
		
		// Look for timestamp pattern [00:00:00.000 --> 00:00:00.000]
		if strings.Contains(line, "-->") {
			// Extract text after timestamp
			parts := strings.SplitN(line, "]", 2)
			if len(parts) == 2 {
				text := strings.TrimSpace(parts[1])
				if text != "" {
					transcription.WriteString(text)
					transcription.WriteString(" ")
				}
			}
		}
	}
	
	return strings.TrimSpace(transcription.String())
}

// ReadTranscriptionFile reads the transcription from a .txt file generated by whisper.cpp
func ReadTranscriptionFile(audioPath string) (string, error) {
	// whisper.cpp generates .txt file with same name as audio file
	txtPath := strings.TrimSuffix(audioPath, filepath.Ext(audioPath)) + ".txt"
	
	data, err := os.ReadFile(txtPath)
	if err != nil {
		return "", fmt.Errorf("failed to read transcription file: %w", err)
	}
	
	return strings.TrimSpace(string(data)), nil
}

// QuickTranscribe is a convenience function for quick transcription with default options
func QuickTranscribe(audioPath string, modelPath string) (string, error) {
	transcriber, err := NewWhisperTranscriber(modelPath, "")
	if err != nil {
		return "", err
	}
	
	return transcriber.Transcribe(audioPath, DefaultTranscribeOptions())
}
