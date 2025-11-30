package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// QueryGPTOSS calls AMD GPT-OSS-120B via vLLM endpoint
func QueryGPTOSS(ctx context.Context, query string, screenshotBase64 string, apiKey string, endpoint string) (string, error) {
	if endpoint == "" {
		endpoint = "http://210.61.209.139:45014/v1/chat/completions"
	} else if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint = strings.TrimSuffix(endpoint, "/") + "/chat/completions"
	}

	log.Printf("[GPT-OSS] Using endpoint: %s", endpoint)

	messages := []OpenAIMessage{
		{
			Role: "system",
			Content: []OpenAIContent{
				{
					Type: "text",
					Text: "You are a helpful AI assistant. Analyze images and answer questions accurately. When an image is provided, describe what you see in detail.",
				},
			},
		},
		{
			Role: "user",
			Content: []OpenAIContent{
				{
					Type: "text",
					Text: query,
				},
			},
		},
	}

	if screenshotBase64 != "" {
		imageURL := normalizeDataURL(screenshotBase64)
		messages[1].Content = append(messages[1].Content, OpenAIContent{
			Type: "image_url",
			ImageURL: &OpenAIImageURL{
				URL: imageURL,
			},
		})
		log.Printf("[GPT-OSS] Added image to request (base64 length: %d)", len(screenshotBase64))
	} else {
		log.Printf("[GPT-OSS] No image provided")
	}

	reqPayload := OpenAIRequest{
		Model:       "gpt-oss-120b",
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpClient := &http.Client{Timeout: 120 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	if apiKey == "" {
		apiKey = "dummy-key"
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[GPT-OSS] Sending request...")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[GPT-OSS] ERROR response: %s", string(bodyBytes))
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no response from API")
	}

	responseText := strings.TrimSpace(result.Choices[0].Message.Content)
	log.Printf("[GPT-OSS] Success! Response length: %d", len(responseText))
	return responseText, nil
}
