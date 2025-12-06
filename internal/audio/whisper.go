package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WhisperTranscriber handles audio transcription using Python Whisper
type WhisperTranscriber struct{}

// TranscribeOptions contains options for transcription
type TranscribeOptions struct {
	Language string // Language code (e.g., "en", "zh", "auto" for auto-detect)
}

// DefaultTranscribeOptions returns default transcription options
func DefaultTranscribeOptions() TranscribeOptions {
	return TranscribeOptions{
		Language: "zh",
	}
}

// Transcribe transcribes an audio file and returns the text
func (w *WhisperTranscriber) Transcribe(audioPath string, options TranscribeOptions) (string, error) {
	// Verify audio file exists
	if _, err := os.Stat(audioPath); err != nil {
		return "", fmt.Errorf("audio file not found: %w", err)
	}

	// Find Python - try multiple locations
	pythonCmd := findPython()
	if pythonCmd == "" {
		return "", fmt.Errorf("Python not found. Please install Python and add to PATH")
	}

	// Get executable directory and create recordtext output directory
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	outputDir := filepath.Join(exeDir, "recordtext")
	
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// 執行 python -m whisper 命令，輸出 txt 格式，指定輸出目錄
	cmd := exec.Command(pythonCmd, "-m", "whisper", audioPath, 
		"--model", "tiny", 
		"--output_format", "txt",
		"--output_dir", outputDir)
	
	// 設定環境變數以支援 UTF-8 輸出（解決中文編碼問題）
	cmd.Env = append(os.Environ(),
		"PYTHONIOENCODING=utf-8",
		"PYTHONUTF8=1",
	)
	
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return "", fmt.Errorf("Whisper 轉錄失敗: %w\n輸出: %s", err, string(output))
	}

	// 讀取生成的 .txt 檔案
	// Whisper 會生成 <filename>.txt (不含原副檔名)
	baseFilename := filepath.Base(audioPath)
	filenameWithoutExt := strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename))
	txtPath := filepath.Join(outputDir, filenameWithoutExt+".txt")
	
	data, err := os.ReadFile(txtPath)
	if err != nil {
		return "", fmt.Errorf("無法讀取轉錄結果 %s: %w", txtPath, err)
	}

	return strings.TrimSpace(string(data)), nil
}

// findPython tries to find Python executable
func findPython() string {
	// Try common Python commands
	candidates := []string{"python", "python3", "py"}
	
	for _, cmd := range candidates {
		if path, err := exec.LookPath(cmd); err == nil {
			return path
		}
	}
	
	// Try common installation paths on Windows
	commonPaths := []string{
		`C:\Python313\python.exe`,
		`C:\Python312\python.exe`,
		`C:\Python311\python.exe`,
		`C:\Python310\python.exe`,
		`C:\Users\` + os.Getenv("USERNAME") + `\AppData\Local\Programs\Python\Python313\python.exe`,
		`C:\Users\` + os.Getenv("USERNAME") + `\AppData\Local\Programs\Python\Python312\python.exe`,
		`C:\Users\` + os.Getenv("USERNAME") + `\AppData\Local\Programs\Python\Python311\python.exe`,
	}
	
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return ""
}



// NewWhisperTranscriberAuto creates a Whisper transcriber
func NewWhisperTranscriberAuto(modelName string) (*WhisperTranscriber, error) {
	pythonCmd := findPython()
	if pythonCmd == "" {
		return nil, fmt.Errorf("Python not found. Please install Python and Whisper: pip install openai-whisper")
	}
	
	return &WhisperTranscriber{}, nil
}
