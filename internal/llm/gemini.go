package llm

import (
	"context"
	"fmt"
)

// QueryGemini calls Google's Gemini API (DISABLED)
func QueryGemini(ctx context.Context, query string, screenshotBase64 string, apiKey string, language string) (string, error) {
	return "", fmt.Errorf("Gemini API is disabled. Please use GPT-OSS-120B instead")
}
