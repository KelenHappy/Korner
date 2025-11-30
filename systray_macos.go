//go:build darwin

package main

// For macOS, we stub out InitSystemTray due to linker conflicts with Wails
// The systray library and Wails v2.11.0 both define AppDelegate on macOS,
// causing duplicate symbol errors during linking (https://github.com/wailsapp/wails/issues/XXXX)
func (a *App) InitSystemTray() {
	// Disabled on macOS due to Wails linker issue
	// TODO: Re-enable after Wails fixes the duplicate AppDelegate symbol issue
}
