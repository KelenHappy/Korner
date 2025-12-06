package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Kelen/Korner/internal/audio"
	"github.com/Kelen/Korner/internal/history"
	"github.com/Kelen/Korner/internal/llm"
	"github.com/Kelen/Korner/internal/platform"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	settings *AppSettings
	platform platform.Platform
	history  *history.Manager
	recorder *audio.Recorder
}

// AppSettings stores user configuration
type AppSettings struct {
	APIProvider  string `json:"apiProvider"` // "gptoss", "openai", "anthropic", "gemini", "custom"
	APIKey       string `json:"apiKey"`
	APIEndpoint  string `json:"apiEndpoint"`
	FloatingIcon string `json:"floatingIcon"`
	Language     string `json:"language"` // "en" or "zh-TW"
}

// NewApp creates a new App application struct
func NewApp() *App {
	historyMgr, err := history.NewManager()
	if err != nil {
		log.Printf("Warning: failed to initialize history manager: %v", err)
	}

	app := &App{
		settings: &AppSettings{
			APIProvider:  "gptoss", // 默認使用 GPT-OSS-120B
			APIKey:       "dummy-key",
			APIEndpoint:  "http://210.61.209.139:45014/v1/", // 默認端點
			FloatingIcon: "🌸",
			Language:     "zh-TW", // 默認語言
		},
		platform: platform.New(),
		history:  historyMgr,
	}
	app.loadSettings()
	return app
}

// getSettingsPath returns the path to the settings file
func (a *App) getSettingsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get home directory: %v", err)
		return "korner-settings.json"
	}
	return filepath.Join(homeDir, ".korner-settings.json")
}

// loadSettings loads settings from file
func (a *App) loadSettings() {
	settingsPath := a.getSettingsPath()
	data, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		log.Printf("Could not read settings file (will use defaults): %v", err)
		return
	}

	var settings AppSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		log.Printf("Could not parse settings file: %v", err)
		return
	}

	a.settings = &settings
	log.Printf("Loaded settings: provider=%s", a.settings.APIProvider)
}

// SaveSettings saves settings to file
func (a *App) SaveSettings(settings AppSettings) error {
	a.settings = &settings

	settingsPath := a.getSettingsPath()
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(settingsPath, data, 0600); err != nil {
		return err
	}

	log.Printf("Saved settings: provider=%s", a.settings.APIProvider)
	return nil
}

// GetSettings returns current settings
func (a *App) GetSettings() AppSettings {
	if a.settings == nil {
		return AppSettings{
			APIProvider:  "gptoss",
			APIKey:       "dummy-key",
			APIEndpoint:  "http://210.61.209.139:45014/v1/",
			FloatingIcon: "🌸",
			Language:     "zh-TW",
		}
	}
	return *a.settings
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("[startup] App context initialized")
}

// domReady is called after the frontend DOM is ready
func (a *App) domReady(ctx context.Context) {
	log.Println("[domReady] Starting...")

	// Log DPI scale for diagnostics
	a.platform.LogDPIInfo()

	// Platform-specific initialization
	initPlatform()

	// Show window
	wailsruntime.WindowShow(ctx)
	log.Println("[domReady] Complete")
}

// ShowWindow shows the application window
func (a *App) ShowWindow() {
	wailsruntime.WindowShow(a.ctx)
	wailsruntime.WindowSetAlwaysOnTop(a.ctx, true)
}

// HideWindow hides the application window
func (a *App) HideWindow() {
	wailsruntime.WindowHide(a.ctx)
}

// TriggerScreenshot triggers the screenshot overlay from system tray
func (a *App) TriggerScreenshot() {
	wailsruntime.WindowShow(a.ctx)
	wailsruntime.WindowSetAlwaysOnTop(a.ctx, true)
	wailsruntime.EventsEmit(a.ctx, "trigger-screenshot")
}

// GetPlatform returns the current platform
func (a *App) GetPlatform() string {
	return runtime.GOOS
}

// GetWindowPosition returns the current window's screen position (x, y)
func (a *App) GetWindowPosition() (int, int) {
	return wailsruntime.WindowGetPosition(a.ctx)
}

// SetWindowPosition sets the window's screen position to (x, y)
func (a *App) SetWindowPosition(x, y int) {
	wailsruntime.WindowSetPosition(a.ctx, x, y)
}

// PositionWindowAt positions the Pie Menu window at the given screen coordinates
func (a *App) PositionWindowAt(x, y int) {
	wailsruntime.WindowSetPosition(a.ctx, x, y)
}

// GetDPIScale returns the current display DPI scaling factor
func (a *App) GetDPIScale() float64 {
	return a.platform.GetDPIScale()
}

// GetScreenSize returns the actual physical screen dimensions
func (a *App) GetScreenSize() (int, int) {
	return a.platform.GetScreenSize()
}

// CaptureScreenshot captures a screenshot of the specified region
func (a *App) CaptureScreenshot(x, y, width, height int) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	winX, winY := wailsruntime.WindowGetPosition(a.ctx)
	screenWidth, screenHeight := a.platform.GetScreenSize()
	dpiScale := a.platform.GetDPIScale()

	expectedViewportWidth := int(float64(screenWidth) / dpiScale)
	expectedViewportHeight := int(float64(screenHeight) / dpiScale)

	log.Printf("DEBUG: Window position: (%d, %d), Viewport coords: (%d, %d, %d, %d)\n", winX, winY, x, y, width, height)
	log.Printf("DEBUG: Screen size: %dx%d, DPI: %.2f, Expected viewport: %dx%d\n",
		screenWidth, screenHeight, dpiScale, expectedViewportWidth, expectedViewportHeight)

	screenX := x + winX
	screenY := y + winY
	log.Printf("DEBUG: Screen coords after window offset: (%d, %d, %d, %d)\n", screenX, screenY, width, height)

	// Capture screenshot (saves to build/bin/screenshots/)
	return a.platform.CaptureScreenshot(ctx, screenX, screenY, width, height)
}

// GetLastScreenshotPath returns the file path of the most recent screenshot
func (a *App) GetLastScreenshotPath() (string, error) {
	return getLastScreenshotPath()
}

// QueryLLM sends a query with screenshot to the configured LLM provider
func (a *App) QueryLLM(query string, screenshotBase64 string, language string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	log.Printf("[QueryLLM] Starting query with provider: %s", a.settings.APIProvider)

	if a.settings == nil {
		return "", fmt.Errorf("Settings not initialized. Please configure your API settings.")
	}

	// Only check API key for providers that require it (not gptoss)
	if a.settings.APIProvider != "gptoss" && a.settings.APIKey == "" {
		return "", fmt.Errorf("API key not configured. Please set your API key in Settings.")
	}

	var result string
	var err error
	var model string

	// Use provided language or fall back to settings
	if language == "" {
		language = a.settings.Language
	}
	if language == "" {
		language = "zh-TW" // Default to Chinese
	}

	switch a.settings.APIProvider {
	case "gptoss":
		// GPT-OSS-120B (默認優先)
		endpoint := a.settings.APIEndpoint
		if endpoint == "" {
			endpoint = "http://210.61.209.139:45014/v1/"
		}
		model = "gpt-oss-120b"
		result, err = llm.QueryGPTOSS(ctx, query, screenshotBase64, a.settings.APIKey, endpoint, language)
	case "gemini":
		model = "gemini-2.0-flash-lite"
		result, err = llm.QueryGemini(ctx, query, screenshotBase64, a.settings.APIKey, language)
	default:
		// Default to GPT-OSS
		endpoint := a.settings.APIEndpoint
		if endpoint == "" {
			endpoint = "http://210.61.209.139:45014/v1/"
		}
		model = "gpt-oss-120b"
		result, err = llm.QueryGPTOSS(ctx, query, screenshotBase64, a.settings.APIKey, endpoint, language)
	}

	if err != nil {
		log.Printf("[QueryLLM] ERROR: %v", err)
		return "", err
	}

	log.Printf("[QueryLLM] Success! Response length: %d", len(result))

	// Save to history
	if a.history != nil {
		screenshotPath, _ := getLastScreenshotPath()
		conv := history.Conversation{
			Timestamp:      time.Now(),
			Question:       query,
			Answer:         result,
			ScreenshotPath: screenshotPath,
			Provider:       a.settings.APIProvider,
			Model:          model,
		}
		if err := a.history.Save(conv); err != nil {
			log.Printf("Warning: failed to save conversation to history: %v", err)
		}
	}

	return result, nil
}

// OpenDevTools opens the developer tools window
func (a *App) OpenDevTools() {
	wailsruntime.WindowShow(a.ctx)
}

// GetRecentHistory returns the most recent N conversations
func (a *App) GetRecentHistory(limit int) ([]history.Conversation, error) {
	if a.history == nil {
		return nil, fmt.Errorf("history manager not initialized")
	}
	return a.history.GetRecent(limit)
}

// GetTodayHistory returns all conversations from today
func (a *App) GetTodayHistory() ([]history.Conversation, error) {
	if a.history == nil {
		return nil, fmt.Errorf("history manager not initialized")
	}
	return a.history.GetToday()
}

// GetAllHistory returns all conversations
func (a *App) GetAllHistory() ([]history.Conversation, error) {
	if a.history == nil {
		return nil, fmt.Errorf("history manager not initialized")
	}
	return a.history.GetAll()
}

// DeleteHistoryItem deletes a conversation by ID
func (a *App) DeleteHistoryItem(id string) error {
	if a.history == nil {
		return fmt.Errorf("history manager not initialized")
	}
	return a.history.Delete(id)
}

// ClearHistory deletes all history
func (a *App) ClearHistory() error {
	if a.history == nil {
		return fmt.Errorf("history manager not initialized")
	}
	return a.history.Clear()
}

// ExportHistoryToText exports all conversations to a text file
func (a *App) ExportHistoryToText(outputPath string) error {
	if a.history == nil {
		return fmt.Errorf("history manager not initialized")
	}
	return a.history.ExportToText(outputPath)
}

// StartRecording starts audio recording
func (a *App) StartRecording() error {
	if a.recorder == nil {
		recorder, err := audio.NewRecorder()
		if err != nil {
			return fmt.Errorf("failed to create recorder: %w", err)
		}
		a.recorder = recorder
	}

	if err := a.recorder.StartRecording(); err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}

	log.Println("[Audio] Recording started")
	return nil
}

// StopRecording stops audio recording and returns the file path
func (a *App) StopRecording() (string, error) {
	if a.recorder == nil {
		return "", fmt.Errorf("recorder not initialized")
	}

	filePath, err := a.recorder.StopRecording()
	if err != nil {
		return "", fmt.Errorf("failed to stop recording: %w", err)
	}

	log.Printf("[Audio] Recording saved to: %s", filePath)
	return filePath, nil
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
	// Get the executable directory
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	
	// Record directory is relative to executable
	recordDir := filepath.Join(exeDir, "record")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(recordDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open folder in Windows Explorer
	cmd := exec.Command("explorer", recordDir)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open folder: %w", err)
	}

	log.Printf("[Audio] Opening folder: %s", recordDir)
	return nil
}

// GenerateMeetingSummary transcribes audio and generates a meeting summary
func (a *App) GenerateMeetingSummary(audioPath string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	log.Printf("[MeetingSummary] Starting transcription for: %s", audioPath)

	// 1. 使用 Whisper 轉錄音訊
	transcriber, err := audio.NewWhisperTranscriberAuto("tiny")
	if err != nil {
		return "", fmt.Errorf("failed to initialize Whisper: %w\n\n請確保:\n1. Whisper.cpp 已安裝 (main.exe)\n2. Tiny 模型已下載 (ggml-tiny.bin)", err)
	}

	options := audio.DefaultTranscribeOptions()
	options.Language = "zh" // 中文
	options.Threads = 4     // 使用 4 個執行緒加速

	transcription, err := transcriber.Transcribe(audioPath, options)
	if err != nil {
		return "", fmt.Errorf("failed to transcribe audio: %w", err)
	}

	if transcription == "" {
		return "", fmt.Errorf("transcription is empty, please check the audio file")
	}

	log.Printf("[MeetingSummary] Transcription completed, length: %d", len(transcription))

	// 2. 使用 LLM 生成會議摘要
	language := a.settings.Language
	if language == "" {
		language = "zh-TW"
	}

	summaryPrompt := fmt.Sprintf(`請根據以下會議錄音的轉錄內容，生成一份簡潔的會議摘要。

摘要應包含：
1. 會議主題
2. 主要討論點（3-5 點）
3. 重要決議或行動項目
4. 關鍵結論

轉錄內容：
%s

請用繁體中文回覆，格式清晰易讀。`, transcription)

	summary, err := llm.QueryGPTOSS(ctx, summaryPrompt, "", a.settings.APIKey, a.settings.APIEndpoint, language)
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	log.Printf("[MeetingSummary] Summary generated successfully")

	// 3. 保存到歷史記錄（不寫入 md 檔案）
	if a.history != nil {
		conv := history.Conversation{
			Timestamp:      time.Now(),
			Question:       "會議摘要 - " + filepath.Base(audioPath),
			Answer:         summary,
			ScreenshotPath: audioPath, // 使用音訊路徑作為參考
			Provider:       a.settings.APIProvider,
			Model:          "whisper-tiny + " + a.settings.APIProvider,
		}
		if err := a.history.Save(conv); err != nil {
			log.Printf("Warning: failed to save meeting summary to history: %v", err)
		}
	}

	return summary, nil
}
