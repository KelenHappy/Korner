package main

import (
	"context"
	"log"
	"runtime"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after the frontend DOM is ready
func (a *App) domReady(ctx context.Context) {
	// Window will start hidden (StartHidden: true in main.go)
	// Log DPI scale for diagnostics
	logDPIInfo()
	// Initialize system tray after Wails is ready
	go a.InitSystemTray()
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

// GetDPIScale returns the current display DPI scaling factor
// Platform implementations:
// - Windows: Returns actual DPI scale (1.0, 1.25, 1.5, 2.0, etc.)
// - Other platforms: Returns 1.0 (not yet implemented)
func (a *App) GetDPIScale() float64 {
	return getDPIScale()
}

// GetScreenSize returns the actual physical screen dimensions
func (a *App) GetScreenSize() (int, int) {
	return getScreenSize()
}

// CaptureScreenshot captures a screenshot of the specified region
// and returns it as a data URL string (data:image/png;base64,...)
//
// Platform implementations:
// - macOS (darwin): Uses screencapture command (see app_screenshot_darwin.go)
// - Windows: Planned (see app_screenshot_others.go stub)
//
// Parameters:
//
//	x, y, width, height: Screenshot region in viewport coordinates (relative to window)
//	If width or height are 0, platform may open interactive selection UI
func (a *App) CaptureScreenshot(x, y, width, height int) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	// Get window position to convert viewport coords to screen coords
	winX, winY := wailsruntime.WindowGetPosition(a.ctx)

	// Get actual screen size and DPI for scaling
	screenWidth, screenHeight := getScreenSize()
	dpiScale := getDPIScale()

	// Calculate the scaling factor: actual screen size / viewport size
	// viewport size = screen size / DPI scale
	expectedViewportWidth := int(float64(screenWidth) / dpiScale)
	expectedViewportHeight := int(float64(screenHeight) / dpiScale)

	log.Printf("DEBUG: Window position: (%d, %d), Viewport coords: (%d, %d, %d, %d)\n", winX, winY, x, y, width, height)
	log.Printf("DEBUG: Screen size: %dx%d, DPI: %.2f, Expected viewport: %dx%d\n",
		screenWidth, screenHeight, dpiScale, expectedViewportWidth, expectedViewportHeight)

	// Add window offset to viewport coordinates
	screenX := x + winX
	screenY := y + winY
	log.Printf("DEBUG: Screen coords after window offset: (%d, %d, %d, %d)\n", screenX, screenY, width, height)

	return captureScreenshot(ctx, screenX, screenY, width, height)
}

// QueryLLM sends a query with screenshot to AMD GPT OSS 120B model.
// It calls the OpenAI-compatible endpoint specified in environment variables:
//   - AMD_LLM_ENDPOINT: Full URL to the chat completions endpoint
//   - AMD_API_KEY: Bearer token for authentication
//   - MODEL_NAME: Optional model identifier (defaults to "gpt-oss-120b")
//
// The screenshotBase64 can be:
//   - Empty string (text-only query)
//   - Raw base64 string (will be prefixed with data URL)
//   - Full data URL (data:image/png;base64,...)
//
// Returns the assistant's response text, or an error if the API call fails.
func (a *App) QueryLLM(query string, screenshotBase64 string) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return AMDQueryLLM(ctx, query, screenshotBase64)
}

// OpenDevTools opens the developer tools window
func (a *App) OpenDevTools() {
	wailsruntime.WindowShow(a.ctx)
	// Note: In production builds, DevTools must be enabled via build flags
	// This is a placeholder for triggering DevTools programmatically
}

// RegisterGlobalHotkey registers a global hotkey
func (a *App) RegisterGlobalHotkey() error {
	// TODO: Implement platform-specific global hotkey registration
	// Windows: Use RegisterHotKey Win32 API
	// macOS: Use Carbon Event Manager via CGO
	return nil
}
