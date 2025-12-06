//go:build windows
// +build windows

package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// FindWhisperExecutable tries to find whisper.cpp executable on Windows
func FindWhisperExecutable() (string, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	
	// Check common Windows locations
	candidates := []string{
		// In current directory
		filepath.Join(cwd, "whisper.exe"),
		filepath.Join(cwd, "main.exe"),
		
		// In whisper.cpp subdirectory
		filepath.Join(cwd, "whisper.cpp", "main.exe"),
		filepath.Join(cwd, "whisper.cpp", "build", "bin", "Release", "main.exe"),
		filepath.Join(cwd, "whisper.cpp", "build", "bin", "main.exe"),
		
		// In build directory
		filepath.Join(cwd, "build", "whisper.exe"),
		filepath.Join(cwd, "build", "main.exe"),
		filepath.Join(cwd, "build", "bin", "main.exe"),
		
		// In models directory (sometimes bundled together)
		filepath.Join(cwd, "models", "main.exe"),
	}
	
	// Check each candidate
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	
	// Try system PATH
	if path, err := exec.LookPath("main.exe"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("whisper.exe"); err == nil {
		return path, nil
	}
	
	return "", fmt.Errorf("whisper.cpp executable not found")
}

// FindWhisperModel tries to find whisper model file on Windows
func FindWhisperModel(modelName string) (string, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	
	// If modelName doesn't have extension, add .bin
	if filepath.Ext(modelName) == "" {
		modelName = modelName + ".bin"
	}
	
	// Check common model locations
	candidates := []string{
		// In models directory
		filepath.Join(cwd, "models", modelName),
		filepath.Join(cwd, "models", "ggml-"+modelName),
		
		// In whisper.cpp/models directory
		filepath.Join(cwd, "whisper.cpp", "models", modelName),
		filepath.Join(cwd, "whisper.cpp", "models", "ggml-"+modelName),
		
		// In current directory
		filepath.Join(cwd, modelName),
		filepath.Join(cwd, "ggml-"+modelName),
	}
	
	// Check each candidate
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	
	return "", fmt.Errorf("whisper model not found: %s", modelName)
}

// NewWhisperTranscriberAuto creates a Whisper transcriber with auto-detection on Windows
func NewWhisperTranscriberAuto(modelName string) (*WhisperTranscriber, error) {
	// Find whisper executable
	whisperPath, err := FindWhisperExecutable()
	if err != nil {
		return nil, fmt.Errorf("failed to find whisper executable: %w", err)
	}
	
	// Find model file
	modelPath, err := FindWhisperModel(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to find model: %w", err)
	}
	
	fmt.Printf("Using whisper: %s\n", whisperPath)
	fmt.Printf("Using model: %s\n", modelPath)
	
	return NewWhisperTranscriber(modelPath, whisperPath)
}

// TranscribeWithProgress transcribes audio with progress callback (Windows specific)
func (w *WhisperTranscriber) TranscribeWithProgress(audioPath string, options TranscribeOptions, progressCallback func(string)) (string, error) {
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
	
	// Create command
	cmd := exec.Command(w.whisperPath, args...)
	
	// Hide console window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	
	// Capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stderr pipe: %w", err)
	}
	
	// Start command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start whisper: %w", err)
	}
	
	// Read output in background
	outputChan := make(chan string, 100)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				line := string(buf[:n])
				outputChan <- line
				if progressCallback != nil {
					progressCallback(line)
				}
			}
			if err != nil {
				break
			}
		}
	}()
	
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				line := string(buf[:n])
				outputChan <- line
				if progressCallback != nil {
					progressCallback(line)
				}
			}
			if err != nil {
				break
			}
		}
	}()
	
	// Wait for command to finish
	err = cmd.Wait()
	close(outputChan)
	
	if err != nil {
		return "", fmt.Errorf("whisper failed: %w", err)
	}
	
	// Collect all output
	var fullOutput string
	for line := range outputChan {
		fullOutput += line
	}
	
	// Parse output to extract transcription
	text := w.parseWhisperOutput(fullOutput)
	
	return text, nil
}
