package llm

import (
	"context"
	"fmt"
)

// QueryCustom calls a custom OpenAI-compatible endpoint (DISABLED)
func QueryCustom(ctx context.Context, query string, screenshotBase64 string, apiKey string, endpoint string) (string, error) {
	return "", fmt.Errorf("Custom API is disabled. Please use GPT-OSS-120B instead")
}
