package main

import (
	_ "embed"
	"log"

	"github.com/getlantern/systray"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/appicon.png
var iconData []byte

// InitSystemTray initializes the system tray icon and menu
func (a *App) InitSystemTray() {
	go systray.Run(func() {
		onReady(a)
	}, onExit)
}

func onReady(app *App) {
	// Set the icon (Windows will show this in notification area)
	systray.SetIcon(iconData)
	systray.SetTitle("Korner")
	systray.SetTooltip("Korner - AI Screenshot Assistant")

	// Create menu items
	mScreenshot := systray.AddMenuItem("📸 New Screenshot", "Take a new screenshot")
	systray.AddSeparator()
	mShow := systray.AddMenuItem("Show Window", "Show the main window")
	mHide := systray.AddMenuItem("Hide Window", "Hide the main window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mScreenshot.ClickedCh:
				// Trigger screenshot
				if app.ctx != nil {
					app.TriggerScreenshot()
				}

			case <-mShow.ClickedCh:
				// Show window
				if app.ctx != nil {
					wailsruntime.WindowShow(app.ctx)
					wailsruntime.WindowSetAlwaysOnTop(app.ctx, true)
				}

			case <-mHide.ClickedCh:
				// Hide window
				if app.ctx != nil {
					wailsruntime.WindowHide(app.ctx)
				}

			case <-mQuit.ClickedCh:
				// Quit application
				if app.ctx != nil {
					wailsruntime.Quit(app.ctx)
				}
				systray.Quit()
				return
			}
		}
	}()

	log.Println("System tray initialized")
}

func onExit() {
	// Cleanup when systray exits
	log.Println("System tray exiting")
}
