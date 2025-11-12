package main

import (
	"context"
	"runtime"
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

// GetPlatform returns the current platform
func (a *App) GetPlatform() string {
	return runtime.GOOS
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
//	x, y, width, height: Screenshot region in screen coordinates
//	If width or height are 0, platform may open interactive selection UI
func (a *App) CaptureScreenshot(x, y, width, height int) (string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return captureScreenshot(ctx, x, y, width, height)
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

// RegisterGlobalHotkey registers a global hotkey
func (a *App) RegisterGlobalHotkey() error {
	// TODO: Implement platform-specific global hotkey registration
	// Windows: Use RegisterHotKey Win32 API
	// macOS: Use Carbon Event Manager via CGO
	return nil
}
