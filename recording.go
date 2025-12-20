package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kelen/Korner/internal/audio"
	"github.com/Kelen/Korner/internal/document"
	"github.com/Kelen/Korner/internal/history"
	"github.com/Kelen/Korner/internal/meeting"
	"github.com/Kelen/Korner/internal/ocr"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// StartRecording starts audio recording
func (a *App) StartRecording() error {
	if a.recorder == nil {
		recorder, err := audio.NewRecorder()
		if err != nil {
			return fmt.Errorf("failed to create recorder: %w", err)
		}
		a.recorder = recorder
	}

	return a.recorder.StartRecording()
}

// StopRecording stops audio recording and returns the file path
func (a *App) StopRecording() (string, error) {
	if a.recorder == nil {
		return "", fmt.Errorf("recorder not initialized")
	}

	return a.recorder.StopRecording()
}

// IsRecording returns whether audio is currently being recorded
func (a *App) IsRecording() bool {
	if a.recorder == nil {
		return false
	}
	return a.recorder.IsRecording()
}

// GetRecordingDuration returns the current recording duration in seconds
func (a *App) GetRecordingDuration() float64 {
	if a.recorder == nil {
		return 0
	}
	return a.recorder.GetDuration()
}

// OpenRecordingFolder opens the folder containing recordings
func (a *App) OpenRecordingFolder() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	recordDir := filepath.Join(filepath.Dir(exePath), "record")
	
	if err := os.MkdirAll(recordDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return exec.Command("explorer", recordDir).Start()
}

// SelectAudioFile opens a file dialog to select an audio file
func (a *App) SelectAudioFile() (string, error) {
	// Open file dialog
	filePath, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "選擇音訊檔案",
		Filters: []wailsruntime.FileFilter{
			{
				DisplayName: "音訊檔案 (*.wav, *.mp3, *.m4a, *.flac, *.ogg, *.opus, *.aac)",
				Pattern:     "*.wav;*.mp3;*.m4a;*.flac;*.ogg;*.opus;*.aac;*.wma",
			},
			{
				DisplayName: "所有檔案 (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to open file dialog: %w", err)
	}

	if filePath == "" {
		log.Println("[Audio] User cancelled file selection")
		return "", nil // Not an error, user just cancelled
	}

	log.Printf("[Audio] Selected file: %s", filePath)
	return filePath, nil
}

// GenerateMeetingSummary transcribes audio and generates a meeting summary
func (a *App) GenerateMeetingSummary(audioPath string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	// 1. 轉錄音訊
	generator, err := meeting.NewGenerator()
	if err != nil {
		return "", fmt.Errorf("failed to initialize meeting generator: %w\n\n請確保已安裝 Python 和 Whisper:\npip install openai-whisper", err)
	}

	language := a.settings.Language
	if language == "" {
		language = "zh-TW"
	}

	result, err := generator.Generate(ctx, audioPath, language)
	if err != nil {
		return "", err
	}

	// 2. 使用 Ollama 生成會議摘要（不需要聯網）
	summaryPrompt := meeting.GenerateSummaryPrompt(language, result.Transcription)

	ollamaEndpoint := a.settings.OllamaEndpoint
	if ollamaEndpoint == "" {
		ollamaEndpoint = "http://127.0.0.1:11434"
	}

	// 使用 QueryOllama 而非 QueryOllamaWithWebSearch，因為摘要不需要聯網
	summary, err := ocr.QueryOllama(ctx, summaryPrompt, "", ollamaEndpoint, language)
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	log.Printf("[MeetingSummary] Summary generated successfully")

	// 3. 保存到歷史記錄
	if a.history != nil {
		conv := history.Conversation{
			Timestamp:      time.Now(),
			Question:       "會議摘要 - " + filepath.Base(audioPath),
			Answer:         summary,
			ScreenshotPath: audioPath,
			Provider:       a.settings.APIProvider,
			Model:          "whisper-tiny + ollama",
		}
		if err := a.history.Save(conv); err != nil {
			log.Printf("Warning: failed to save meeting summary to history: %v", err)
		}
	}

	return summary, nil
}

// SelectDocumentFiles opens a file dialog to select document files
func (a *App) SelectDocumentFiles() ([]string, error) {
	filePaths, err := wailsruntime.OpenMultipleFilesDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "選擇文件",
		Filters: []wailsruntime.FileFilter{
			{
				DisplayName: "文件檔案 (*.txt, *.md, *.pdf, *.json, *.csv)",
				Pattern:     "*.txt;*.md;*.pdf;*.json;*.csv;*.log",
			},
			{
				DisplayName: "所有檔案 (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open file dialog: %w", err)
	}

	if len(filePaths) == 0 {
		log.Println("[Document] No files selected")
		return []string{}, nil
	}

	log.Printf("[Document] Selected %d files", len(filePaths))
	return filePaths, nil
}

// ReadDocumentFile reads a document file and returns its text content
func (a *App) ReadDocumentFile(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	log.Printf("[Document] Reading file: %s (type: %s)", filePath, ext)

	switch ext {
	case ".pdf":
		return document.ExtractPDFText(filePath)
	case ".txt", ".md", ".json", ".csv", ".log":
		return document.ReadTextFile(filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}
