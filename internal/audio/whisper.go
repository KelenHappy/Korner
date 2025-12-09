package audio

import (
	"fmt"
	"log"
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
// Supports: wav, mp3, m4a, flac, ogg, opus, and other formats supported by ffmpeg
func (w *WhisperTranscriber) Transcribe(audioPath string, options TranscribeOptions) (string, error) {
	log.Printf("[Whisper] Starting transcription for: %s", audioPath)
	
	// Verify audio file exists
	if _, err := os.Stat(audioPath); err != nil {
		log.Printf("[Whisper] Error: Audio file not found: %s, error: %v", audioPath, err)
		return "", fmt.Errorf("音訊檔案不存在")
	}

	// Verify file format is supported
	ext := strings.ToLower(filepath.Ext(audioPath))
	supportedFormats := []string{".wav", ".mp3", ".m4a", ".flac", ".ogg", ".opus", ".aac", ".wma"}
	isSupported := false
	for _, format := range supportedFormats {
		if ext == format {
			isSupported = true
			break
		}
	}
	if !isSupported {
		log.Printf("[Whisper] Error: Unsupported format: %s, supported: %v", ext, supportedFormats)
		return "", fmt.Errorf("不支援的格式")
	}

	// Find Python - try multiple locations
	pythonCmd := findPython()
	if pythonCmd == "" {
		log.Printf("[Whisper] Error: Python not found in PATH or common locations")
		return "", fmt.Errorf("找不到 Python")
	}
	log.Printf("[Whisper] Using Python: %s", pythonCmd)

	// Get executable directory and create recordtext output directory
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("[Whisper] Error: Failed to get executable path: %v", err)
		return "", fmt.Errorf("系統錯誤")
	}
	exeDir := filepath.Dir(exePath)
	outputDir := filepath.Join(exeDir, "recordtext")
	
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("[Whisper] Error: Failed to create output directory %s: %v", outputDir, err)
		return "", fmt.Errorf("無法建立輸出目錄")
	}
	
	log.Printf("[Whisper] Output directory: %s", outputDir)
	
	// 執行 python -m whisper 命令，輸出 txt 格式，指定輸出目錄
	// 對於 mp3 等壓縮格式，Whisper 會自動使用 ffmpeg 解碼
	cmd := exec.Command(pythonCmd, "-m", "whisper", audioPath, 
		"--model", "tiny", 
		"--output_format", "txt",
		"--output_dir", outputDir)
	
	// 設定環境變數以支援 UTF-8 輸出（解決中文編碼問題）
	cmd.Env = append(os.Environ(),
		"PYTHONIOENCODING=utf-8",
		"PYTHONUTF8=1",
	)
	
	log.Printf("[Whisper] Executing Whisper command...")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		errorMsg := string(output)
		log.Printf("[Whisper] Error: Command failed: %v", err)
		log.Printf("[Whisper] Full output: %s", errorMsg)
		
		// 檢查是否是 ffmpeg 相關錯誤
		if strings.Contains(errorMsg, "ffmpeg") || strings.Contains(errorMsg, "RuntimeError") {
			log.Printf("[Whisper] Detected ffmpeg-related error")
			return "", fmt.Errorf("轉錄失敗")
		}
		return "", fmt.Errorf("轉錄失敗")
	}

	log.Printf("[Whisper] Command completed successfully")

	// 讀取生成的 .txt 檔案
	// Whisper 會生成 <filename>.txt (不含原副檔名)
	baseFilename := filepath.Base(audioPath)
	filenameWithoutExt := strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename))
	txtPath := filepath.Join(outputDir, filenameWithoutExt+".txt")
	
	log.Printf("[Whisper] Reading transcription from: %s", txtPath)
	data, err := os.ReadFile(txtPath)
	if err != nil {
		log.Printf("[Whisper] Error: Failed to read transcription file %s: %v", txtPath, err)
		return "", fmt.Errorf("無法讀取轉錄結果")
	}

	transcription := strings.TrimSpace(string(data))
	if transcription == "" {
		log.Printf("[Whisper] Warning: Transcription is empty for file: %s", audioPath)
		return "", fmt.Errorf("轉錄結果為空")
	}

	log.Printf("[Whisper] Transcription completed successfully, length: %d chars", len(transcription))
	return transcription, nil
}

// findPython tries to find Python executable
// Priority: 1. Virtual env in project, 2. System Python
func findPython() string {
	// First, try to find virtual environment in project directory
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		
		// Check for venv in project directory
		venvPaths := []string{
			filepath.Join(exeDir, "venv", "Scripts", "python.exe"),
			filepath.Join(exeDir, ".venv", "Scripts", "python.exe"),
			filepath.Join(exeDir, "python-env", "Scripts", "python.exe"),
		}
		
		for _, venvPath := range venvPaths {
			if _, err := os.Stat(venvPath); err == nil {
				return venvPath
			}
		}
	}
	
	// Try common Python commands in PATH
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
		log.Printf("[Whisper] Error: Python not found during initialization")
		log.Printf("[Whisper] Searched in PATH and common installation directories")
		return nil, fmt.Errorf("找不到 Python")
	}
	
	log.Printf("[Whisper] Transcriber initialized with Python: %s", pythonCmd)
	return &WhisperTranscriber{}, nil
}
