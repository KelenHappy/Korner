package llm

import (
	"context"
	"fmt"
)

// QueryOpenAI calls OpenAI's GPT-4 Vision API (DISABLED)
func QueryOpenAI(ctx context.Context, query string, screenshotBase64 string, apiKey string, model string) (string, error) {
	return "", fmt.Errorf("OpenAI API is disabled. Please use GPT-OSS-120B instead")
}
