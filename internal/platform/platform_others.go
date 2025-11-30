//go:build !windows && !darwin

package platform

import (
	"context"

	"github.com/Kelen/Korner/internal/platform/common"
)

type commonPlatform struct{}

func New() Platform {
	return &commonPlatform{}
}

func (p *commonPlatform) CaptureScreenshot(ctx context.Context, x, y, width, height int) (string, error) {
	return common.CaptureScreenshot(ctx, x, y, width, height)
}

func (p *commonPlatform) GetDPIScale() float64 {
	return common.GetDPIScale()
}

func (p *commonPlatform) GetScreenSize() (int, int) {
	return common.GetScreenSize()
}

func (p *commonPlatform) LogDPIInfo() {
	common.LogDPIInfo()
}
