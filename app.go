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
	// Log DPI scale for diagnostics
	logDPIInfo()
	
	// Initialize system tray
	go a.InitSystemTray()
	
	// Register global hotkey (Ctrl+Alt+Q)
	go func() {
		err := a.RegisterGlobalHotkey()
		if err != nil {
			log.Printf("Failed to register global hotkey: %v", err)
		}
	}()
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

// QueryLLM sends a query with screenshot to AMD GPT OSS 120B model
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
}
