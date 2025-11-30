package platform

import "context"

// Platform defines the interface for platform-specific operations
type Platform interface {
	// Screenshot operations
	CaptureScreenshot(ctx context.Context, x, y, width, height int) (string, error)
	GetDPIScale() float64
	GetScreenSize() (int, int)
	LogDPIInfo()
}
