package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// This file provides an implementation for calling an OpenAI-compatible AMD endpoint
// to perform Vision-enabled chat completions. It does NOT override the existing
// App.QueryLLM method in app.go to avoid duplicate symbol errors. Instead, it exposes
// helper functions that app.go can call.
//
// Typical usage (inside your existing App.QueryLLM in app.go):
//
//   func (a *App) QueryLLM(query string, screenshotBase64 string) (string, error) {
//       // Prefer the context captured at startup if available; otherwise, use background.
//       ctx := a.ctx
//       if ctx == nil {
//           ctx = context.Background()
//       }
//       return AMDQueryLLM(ctx, query, screenshotBase64)
//   }
//
// Environment variables used:
// - AMD_LLM_ENDPOINT (required): OpenAI-compatible endpoint URL, e.g. https://your-amd-endpoint/v1/chat/completions
// - AMD_API_KEY     (required): API key for the AMD endpoint
// - MODEL_NAME      (optional): Model name; defaults to "gpt-oss-120b" if not set
//
// The request conforms to OpenAI Chat Completions schema with Vision content:
// messages: [
//   { "role": "user", "content": [
//       { "type": "text", "text": "<your text>" },
//       { "type": "image_url", "image_url": { "url": "data:image/png;base64,..." } }
//   ]}
// ]

// AMDQueryLLM invokes the AMD OpenAI-compatible endpoint using the provided query
// and an optional screenshot as a base64-encoded PNG (without or with data URL prefix).
// Returns the assistant message content from the first choice.
func AMDQueryLLM(ctx context.Context, query string, screenshotBase64 string) (string, error) {
	endpoint := strings.TrimSpace(os.Getenv("AMD_LLM_ENDPOINT"))
	apiKey := strings.TrimSpace(os.Getenv("AMD_API_KEY"))
	model := strings.TrimSpace(os.Getenv("MODEL_NAME"))
	if model == "" {
		model = "gpt-oss-120b"
	}
	if endpoint == "" {
		return "", errors.New("missing AMD_LLM_ENDPOINT environment variable")
	}
	if apiKey == "" {
		return "", errors.New("missing AMD_API_KEY environment variable")
	}

	reqPayload := openAIChatRequest{
		Model:       model,
		Messages:    buildVisionMessages(query, screenshotBase64),
		Temperature: 0.2,
		// You may tune these fields based on latency/quality tradeoffs:
		// MaxTokens:   1024,
		// TopP:        1.0,
		// Stream:      false,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	// Use per-request timeout; if ctx already has a deadline, rely on it.
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	var parsed openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return "", fmt.Errorf("llm error (%d): %s", resp.StatusCode, parsed.Error.Message)
		}
		return "", fmt.Errorf("llm error status: %d", resp.StatusCode)
	}

	if len(parsed.Choices) == 0 {
		return "", errors.New("no choices returned from LLM")
	}

	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

// buildVisionMessages constructs a single user message with text and optional image content.
func buildVisionMessages(query string, screenshotBase64 string) []oaMessage {
	parts := []oaContentPart{
		{
			Type: "text",
			Text: strings.TrimSpace(query),
		},
	}

	imgDataURL := normalizeDataURL(screenshotBase64)
	if imgDataURL != "" {
		parts = append(parts, oaContentPart{
			Type: "image_url",
			ImageURL: &oaImageURL{
				URL:    imgDataURL,
				Detail: "auto",
			},
		})
	}

	return []oaMessage{
		{
			Role:    "user",
			Content: parts,
		},
	}
}

// normalizeDataURL ensures the screenshot is a valid data URL. If the input is empty,
// returns empty string. If already a "data:" URL, returns as-is. Otherwise, prefixes
// "data:image/png;base64,".
func normalizeDataURL(b64 string) string {
	b64 = strings.TrimSpace(b64)
	if b64 == "" {
		return ""
	}
	if strings.HasPrefix(b64, "data:") {
		return b64
	}
	// Assume PNG; adjust if you support other formats.
	return "data:image/png;base64," + b64
}

/***************
 OpenAI-compatible request/response schema
***************/

type openAIChatRequest struct {
	Model       string      `json:"model"`
	Messages    []oaMessage `json:"messages"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Temperature float64     `json:"temperature,omitempty"`
	TopP        float64     `json:"top_p,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
	Stop        []string    `json:"stop,omitempty"`
	Frequency   float64     `json:"frequency_penalty,omitempty"`
	Presence    float64     `json:"presence_penalty,omitempty"`
	Tools       interface{} `json:"tools,omitempty"`
	ToolChoice  interface{} `json:"tool_choice,omitempty"`
	LogitBias   interface{} `json:"logit_bias,omitempty"`
	User        string      `json:"user,omitempty"`
	ResponseFmt interface{} `json:"response_format,omitempty"`
}

type oaMessage struct {
	Role    string          `json:"role"` // "system" | "user" | "assistant"
	Content []oaContentPart `json:"content"`
}

type oaContentPart struct {
	Type     string      `json:"type"` // "text" | "image_url"
	Text     string      `json:"text,omitempty"`
	ImageURL *oaImageURL `json:"image_url,omitempty"`
}

type oaImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"` // "low" | "high" | "auto"
}

type openAIChatResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []openAIChatChoice `json:"choices"`
	Usage   *openAIUsage       `json:"usage,omitempty"`
	Error   *openAIError       `json:"error,omitempty"`
}

type openAIChatChoice struct {
	Index        int         `json:"index"`
	Message      oaChoiceMsg `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type oaChoiceMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code,omitempty"`
}
