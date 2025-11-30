package llm

import (
	"context"
	"fmt"
)

// QueryAnthropic calls Anthropic's Claude API (DISABLED)
func QueryAnthropic(ctx context.Context, query string, screenshotBase64 string, apiKey string) (string, error) {
	return "", fmt.Errorf("Anthropic API is disabled. Please use GPT-OSS-120B instead")
}
