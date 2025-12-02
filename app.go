package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

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
}

// AppSettings stores user configuration
type AppSettings struct {
	APIProvider  string `json:"apiProvider"` // "gptoss", "openai", "anthropic", "gemini", "custom"
	APIKey       string `json:"apiKey"`
	APIEndpoint  string `json:"apiEndpoint"`
	FloatingIcon string `json:"floatingIcon"`
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
func (a *App) QueryLLM(query string, screenshotBase64 string) (string, error) {
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

	switch a.settings.APIProvider {
	case "gptoss":
		// GPT-OSS-120B (默認優先)
		endpoint := a.settings.APIEndpoint
		if endpoint == "" {
			endpoint = "http://210.61.209.139:45014/v1/"
		}
		model = "gpt-oss-120b"
		result, err = llm.QueryGPTOSS(ctx, query, screenshotBase64, a.settings.APIKey, endpoint)
	case "openai":
		model = "gpt-4-vision-preview"
		result, err = llm.QueryOpenAI(ctx, query, screenshotBase64, a.settings.APIKey, model)
	case "anthropic":
		model = "claude-3-5-sonnet"
		result, err = llm.QueryAnthropic(ctx, query, screenshotBase64, a.settings.APIKey)
	case "gemini":
		model = "gemini-2.0-flash-lite"
		result, err = llm.QueryGemini(ctx, query, screenshotBase64, a.settings.APIKey)
	case "custom":
		if a.settings.APIEndpoint == "" {
			return "", fmt.Errorf("custom API endpoint not configured")
		}
		model = "custom"
		result, err = llm.QueryCustom(ctx, query, screenshotBase64, a.settings.APIKey, a.settings.APIEndpoint)
	default:
		return "", fmt.Errorf("unsupported API provider: %s", a.settings.APIProvider)
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
