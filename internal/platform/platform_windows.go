//go:build windows

package platform

import (
	"context"

	"github.com/Kelen/Korner/internal/platform/windows"
)

type windowsPlatform struct{}

func New() Platform {
	return &windowsPlatform{}
}

func (p *windowsPlatform) CaptureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return windows.CaptureScreenshot(ctx, x, y, width, height)
}

func (p *windowsPlatform) GetDPIScale() float64 {
	return windows.GetDPIScale()
}

func (p *windowsPlatform) GetScreenSize() (int, int) {
	return windows.GetScreenSize()
}

func (p *windowsPlatform) LogDPIInfo() {
	windows.LogDPIInfo()
}

// DisableWindowSnap disables Windows Snap feature
func DisableWindowSnap(windowTitle string) error {
	return windows.DisableWindowSnap(windowTitle)
}
