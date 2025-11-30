package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"fmt"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	settings *AppSettings
}

// AppSettings stores user configuration
type AppSettings struct {
	APIProvider  string `json:"apiProvider"` // "openai", "anthropic", "gemini", "custom"
	APIKey       string `json:"apiKey"`
	APIEndpoint  string `json:"apiEndpoint"`
	FloatingIcon string `json:"floatingIcon"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		settings: &AppSettings{
			APIProvider:  "openai",
			APIKey:       "",
			APIEndpoint:  "",
			FloatingIcon: "🌸",
		},
	}
	// Load settings from file
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
			APIProvider:  "openai",
			APIKey:       "",
			APIEndpoint:  "",
			FloatingIcon: "🌸",
		}
	}
	return *a.settings
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after the frontend DOM is ready
func (a *App) domReady(ctx context.Context) {
	// TODO: Debug why app crashes on startup
	// Temporarily disabling all operations to test basic startup
	log.Println("[domReady] Starting...")

	// Log DPI scale for diagnostics
	logDPIInfo()

	// Disable Windows Snap feature
	// disableSnap()

	// Initialize system tray
	// go a.InitSystemTray()

	// Register global hotkey (Ctrl+Alt+Q)
	// go func() {
	// 	err := a.RegisterGlobalHotkey()
	// 	if err != nil {
	// 		log.Printf("Failed to register global hotkey: %v", err)
	// 	}
	// }()

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
	return getDPIScale()
}

// GetScreenSize returns the actual physical screen dimensions
func (a *App) GetScreenSize() (int, int) {
	return getScreenSize()
}

// CaptureScreenshot captures a screenshot of the specified region
func (a *App) CaptureScreenshot(x, y, width, height int) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	winX, winY := wailsruntime.WindowGetPosition(a.ctx)
	screenWidth, screenHeight := getScreenSize()
	dpiScale := getDPIScale()

	expectedViewportWidth := int(float64(screenWidth) / dpiScale)
	expectedViewportHeight := int(float64(screenHeight) / dpiScale)

	log.Printf("DEBUG: Window position: (%d, %d), Viewport coords: (%d, %d, %d, %d)\n", winX, winY, x, y, width, height)
	log.Printf("DEBUG: Screen size: %dx%d, DPI: %.2f, Expected viewport: %dx%d\n",
		screenWidth, screenHeight, dpiScale, expectedViewportWidth, expectedViewportHeight)

	screenX := x + winX
	screenY := y + winY
	log.Printf("DEBUG: Screen coords after window offset: (%d, %d, %d, %d)\n", screenX, screenY, width, height)

	return captureScreenshot(ctx, screenX, screenY, width, height)
}

// QueryLLM sends a query with screenshot to the configured LLM provider
func (a *App) QueryLLM(query string, screenshotBase64 string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	log.Printf("[QueryLLM] Starting query with provider: %s", a.settings.APIProvider)
	log.Printf("[QueryLLM] Query length: %d, Screenshot length: %d", len(query), len(screenshotBase64))

	if a.settings == nil {
		log.Printf("[QueryLLM] ERROR: settings is nil")
		return "", fmt.Errorf("Settings not initialized. Please configure your API settings.")
	}

	if a.settings.APIKey == "" {
		log.Printf("[QueryLLM] ERROR: API key is empty for provider: %s", a.settings.APIProvider)
		return "", fmt.Errorf("API key not configured. Please set your API key in Settings.")
	}

	log.Printf("[QueryLLM] API Key present: %d characters", len(a.settings.APIKey))

	// Route to appropriate provider
	var result string
	var err error

	switch a.settings.APIProvider {
	case "openai":
		log.Printf("[QueryLLM] Calling OpenAI API...")
		result, err = QueryOpenAI(ctx, query, screenshotBase64, a.settings.APIKey, "gpt-4-vision-preview")
	case "anthropic":
		log.Printf("[QueryLLM] Calling Anthropic API...")
		result, err = QueryAnthropic(ctx, query, screenshotBase64, a.settings.APIKey)
	case "gemini":
		log.Printf("[QueryLLM] Calling Gemini API...")
		result, err = QueryGemini(ctx, query, screenshotBase64, a.settings.APIKey)
	case "custom":
		if a.settings.APIEndpoint == "" {
			log.Printf("[QueryLLM] ERROR: Custom endpoint is empty")
			return "", fmt.Errorf("custom API endpoint not configured")
		}
		log.Printf("[QueryLLM] Calling custom endpoint: %s", a.settings.APIEndpoint)
		result, err = QueryCustom(ctx, query, screenshotBase64, a.settings.APIKey, a.settings.APIEndpoint)
	default:
		log.Printf("[QueryLLM] ERROR: Unsupported provider: %s", a.settings.APIProvider)
		return "", fmt.Errorf("unsupported API provider: %s", a.settings.APIProvider)
	}

	if err != nil {
		log.Printf("[QueryLLM] ERROR from provider: %v", err)
		return "", err
	}

	log.Printf("[QueryLLM] Success! Response length: %d", len(result))
	return result, nil
}

// OpenDevTools opens the developer tools window
func (a *App) OpenDevTools() {
	wailsruntime.WindowShow(a.ctx)
}
