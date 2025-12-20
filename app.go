package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Kelen/Korner/internal/audio"
	"github.com/Kelen/Korner/internal/history"
	"github.com/Kelen/Korner/internal/llm"
	"github.com/Kelen/Korner/internal/ocr"
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
	APIProvider    string `json:"apiProvider"` // "ollama", "gptoss", "openai", "anthropic", "gemini"
	APIKey         string `json:"apiKey"`
	APIEndpoint    string `json:"apiEndpoint"`
	FloatingIcon   string `json:"floatingIcon"`
	Language       string `json:"language"`       // "en" or "zh-TW"
	OllamaEndpoint string `json:"ollamaEndpoint"` // Ollama server endpoint
}

// NewApp creates a new App application struct
func NewApp() *App {
	historyMgr, err := history.NewManager()
	if err != nil {
		log.Printf("Warning: failed to initialize history manager: %v", err)
	}

	app := &App{
		settings: &AppSettings{
			APIProvider:    "ollama", // 默認使用 Ollama
			APIKey:         "dummy-key",
			APIEndpoint:    "http://210.61.209.139:45014/v1/", // GPT-OSS 端點
			FloatingIcon:   "🌸",
			Language:       "zh-TW",                     // 默認語言
			OllamaEndpoint: "http://127.0.0.1:11434", // Ollama 端點
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
			APIProvider:    "ollama",
			APIKey:         "dummy-key",
			APIEndpoint:    "http://210.61.209.139:45014/v1/",
			FloatingIcon:   "🌸",
			Language:       "zh-TW",
			OllamaEndpoint: "http://127.0.0.1:11434",
		}
	}
	return *a.settings
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("[startup] App context initialized")
	
	// Check Ollama connectivity
	go func() {
		time.Sleep(2 * time.Second) // Wait for app to fully initialize
		endpoint := a.settings.OllamaEndpoint
		if endpoint == "" {
			endpoint = "http://127.0.0.1:11434"
		}
		
		testURL := strings.TrimSuffix(endpoint, "/") + "/api/tags"
		log.Printf("[startup] Testing Ollama connectivity at: %s", testURL)
		
		for i := 0; i < 3; i++ {
			resp, err := http.Get(testURL)
			if err == nil {
				log.Printf("[startup] Ollama is ready! Status: %d", resp.StatusCode)
				resp.Body.Close()
				return
			}
			log.Printf("[startup] Ollama not ready (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
		}
		log.Printf("[startup] Warning: Could not connect to Ollama after 3 attempts")
	}()
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

// ReadScreenshotAsBase64 reads a screenshot file and returns it as base64
func (a *App) ReadScreenshotAsBase64(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read screenshot file: %w", err)
	}
	
	// Convert to base64
	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str, nil
}

// ExtractTextFromScreenshot uses Ollama to extract text from screenshot
func (a *App) ExtractTextFromScreenshot(screenshotBase64 string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	
	endpoint := a.settings.OllamaEndpoint
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}
	
	log.Printf("[OCR] Extracting text from screenshot using Ollama")
	log.Printf("[OCR] Endpoint: %s", endpoint)
	log.Printf("[OCR] Screenshot base64 length: %d", len(screenshotBase64))
	
	extractedText, err := ocr.ExtractTextFromImage(ctx, screenshotBase64, endpoint)
	if err != nil {
		log.Printf("[OCR] Failed to extract text: %v", err)
		return "", err
	}
	
	log.Printf("[OCR] Successfully extracted text (length: %d)", len(extractedText))
	return extractedText, nil
}

// QueryLLM sends a query with screenshot to the configured LLM provider
func (a *App) QueryLLM(query string, screenshotBase64 string, language string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	log.Printf("[QueryLLM] Starting query with provider: %s", a.settings.APIProvider)
	log.Printf("[QueryLLM] Query length: %d, Screenshot base64 length: %d", len(query), len(screenshotBase64))
	
	// If screenshot is provided and provider is NOT Ollama, extract text using Ollama first
	var extractedText string
	var shouldSendImage bool = false
	
	if screenshotBase64 != "" && a.settings.APIProvider != "ollama" {
		log.Printf("[QueryLLM] Screenshot provided (length: %d), extracting text with Ollama...", len(screenshotBase64))
		text, err := a.ExtractTextFromScreenshot(screenshotBase64)
		if err != nil {
			log.Printf("[QueryLLM] Warning: OCR failed, continuing without extracted text: %v", err)
		} else if text != "" {
			extractedText = text
			log.Printf("[QueryLLM] OCR extracted text: %s", extractedText[:min(100, len(extractedText))])
			
			// Append extracted text to query
			if query != "" {
				query = query + "\n\n[圖片中的文字內容]\n" + extractedText
			} else {
				query = "[圖片中的文字內容]\n" + extractedText
			}
		}
		
		// Check if provider supports multimodal
		// GPT-OSS doesn't support images, only Gemini does
		if a.settings.APIProvider == "gemini" {
			shouldSendImage = true
		}
	}

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
	case "ollama":
		// Ollama 本地模型
		endpoint := a.settings.OllamaEndpoint
		if endpoint == "" {
			endpoint = "http://127.0.0.1:11434"
		}
		model = "qwen3-vl:4b"
		// Ollama 支持圖片，如果有截圖就發送
		imageToSend := ""
		if screenshotBase64 != "" {
			imageToSend = screenshotBase64
		}
		result, err = ocr.QueryOllama(ctx, query, imageToSend, endpoint, language)
	case "gptoss":
		// GPT-OSS-120B (不支持圖片，只發送文字)
		endpoint := a.settings.APIEndpoint
		if endpoint == "" {
			endpoint = "http://210.61.209.139:45014/v1/"
		}
		model = "gpt-oss-120b"
		// Don't send image to GPT-OSS, only send extracted text
		result, err = llm.QueryGPTOSS(ctx, query, "", a.settings.APIKey, endpoint, language)
	case "gemini":
		model = "gemini-2.0-flash-lite"
		// Gemini supports multimodal, send image if available
		imageToSend := ""
		if shouldSendImage {
			imageToSend = screenshotBase64
		}
		result, err = llm.QueryGemini(ctx, query, imageToSend, a.settings.APIKey, language)
	default:
		// Default to Ollama
		endpoint := a.settings.OllamaEndpoint
		if endpoint == "" {
			endpoint = "http://127.0.0.1:11434"
		}
		model = "qwen3-vl:4b"
		imageToSend := ""
		if screenshotBase64 != "" {
			imageToSend = screenshotBase64
		}
		result, err = ocr.QueryOllama(ctx, query, imageToSend, endpoint, language)
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

// QueryLLMWithWebSearch sends a query with web search enabled using Ollama
func (a *App) QueryLLMWithWebSearch(query string, screenshotBase64 string, language string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	log.Printf("[QueryLLMWithWebSearch] Starting query with web search enabled")
	log.Printf("[QueryLLMWithWebSearch] Query: %s", query)

	// Get Ollama endpoint
	endpoint := a.settings.OllamaEndpoint
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}

	// If screenshot is provided, extract text first
	if screenshotBase64 != "" {
		log.Printf("[QueryLLMWithWebSearch] Screenshot provided, extracting text...")
		text, err := a.ExtractTextFromScreenshot(screenshotBase64)
		if err != nil {
			log.Printf("[QueryLLMWithWebSearch] Warning: OCR failed: %v", err)
		} else if text != "" {
			query = query + "\n\n[圖片中的文字內容]\n" + text
		}
	}

	// Use Ollama with web search
	result, err := ocr.QueryOllamaWithWebSearch(ctx, query, endpoint, language)
	if err != nil {
		log.Printf("[QueryLLMWithWebSearch] ERROR: %v", err)
		return "", err
	}

	log.Printf("[QueryLLMWithWebSearch] Success! Response length: %d", len(result))

	// Save to history
	if a.history != nil {
		screenshotPath, _ := getLastScreenshotPath()
		conv := history.Conversation{
			Timestamp:      time.Now(),
			Question:       query + " [聯網搜尋]",
			Answer:         result,
			ScreenshotPath: screenshotPath,
			Provider:       "ollama",
			Model:          "qwen3-vl:4b",
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

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

