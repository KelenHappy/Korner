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
	systray.SetTooltip("Korner - AI Screenshot Assistant (Ctrl+Alt+Q)")

	// Create menu items
	mShowMenu := systray.AddMenuItem("Show Menu (Ctrl+Alt+Q)", "Show Pie Menu")
	mScreenshot := systray.AddMenuItem("📸 Screenshot", "Take a screenshot")
	mAskQuestion := systray.AddMenuItem("💬 Ask Question", "Ask a question")
	systray.AddSeparator()
	mSettings := systray.AddMenuItem("⚙️ Settings", "Open settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mShowMenu.ClickedCh:
				// Show Pie Menu at mouse position
				if app.ctx != nil {
					showPieMenuAtMouse(app)
				}

			case <-mScreenshot.ClickedCh:
				// Show menu and trigger screenshot
				if app.ctx != nil {
					showPieMenuAtMouse(app)
					wailsruntime.EventsEmit(app.ctx, "trigger-screenshot")
				}

			case <-mAskQuestion.ClickedCh:
				// Show menu and trigger ask question
				if app.ctx != nil {
					showPieMenuAtMouse(app)
					wailsruntime.EventsEmit(app.ctx, "trigger-question")
				}

			case <-mSettings.ClickedCh:
				// Show menu and trigger settings
				if app.ctx != nil {
					showPieMenuAtMouse(app)
					wailsruntime.EventsEmit(app.ctx, "trigger-settings")
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

func showPieMenuAtMouse(app *App) {
	if app.ctx != nil {
		wailsruntime.WindowShow(app.ctx)
		wailsruntime.WindowSetAlwaysOnTop(app.ctx, true)
		// Emit event to show menu - position will be handled by frontend
		wailsruntime.EventsEmit(app.ctx, "trigger-pie-menu", 0, 0)
	}
}
