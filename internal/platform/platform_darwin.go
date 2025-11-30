//go:build darwin

package platform

import (
	"context"

	"github.com/Kelen/Korner/internal/platform/darwin"
)

type darwinPlatform struct{}

func New() Platform {
	return &darwinPlatform{}
}

func (p *darwinPlatform) CaptureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return darwin.CaptureScreenshot(ctx, x, y, width, height)
}

func (p *darwinPlatform) GetDPIScale() float64 {
	return darwin.GetDPIScale()
}

func (p *darwinPlatform) GetScreenSize() (int, int) {
	return darwin.GetScreenSize()
}

func (p *darwinPlatform) LogDPIInfo() {
	darwin.LogDPIInfo()
}
