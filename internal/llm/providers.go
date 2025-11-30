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

// QueryOpenAI calls OpenAI's GPT-4 Vision API
func QueryOpenAI(ctx context.Context, query string, screenshotBase64 string, apiKey string, model string) (string, error) {
	if model == "" {
		model = "gpt-4-vision-preview"
	}

	endpoint := "https://api.openai.com/v1/chat/completions"

	messages := []openAIMessage{
		{
			Role: "user",
			Content: []openAIContent{
				{
					Type: "text",
					Text: query,
				},
			},
		},
	}

	// Add image if provided
	if screenshotBase64 != "" {
		imageURL := normalizeDataURL(screenshotBase64)
		messages[0].Content = append(messages[0].Content, openAIContent{
			Type: "image_url",
			ImageURL: &openAIImageURL{
				URL: imageURL,
			},
		})
	}

	reqPayload := openAIRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpClient := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no response from API")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// QueryAnthropic calls Anthropic's Claude API
func QueryAnthropic(ctx context.Context, query string, screenshotBase64 string, apiKey string) (string, error) {
	endpoint := "https://api.anthropic.com/v1/messages"

	content := []anthropicContent{
		{
			Type: "text",
			Text: query,
		},
	}

	// Add image if provided
	if screenshotBase64 != "" {
		// Remove data URL prefix if present
		base64Data := screenshotBase64
		if strings.Contains(base64Data, ",") {
			parts := strings.Split(base64Data, ",")
			if len(parts) > 1 {
				base64Data = parts[1]
			}
		}

		content = append(content, anthropicContent{
			Type: "image",
			Source: &anthropicImageSource{
				Type:      "base64",
				MediaType: "image/png",
				Data:      base64Data,
			},
		})
	}

	reqPayload := anthropicRequest{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpClient := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", errors.New("no response from API")
	}

	return strings.TrimSpace(result.Content[0].Text), nil
}

// QueryGemini calls Google's Gemini API
func QueryGemini(ctx context.Context, query string, screenshotBase64 string, apiKey string) (string, error) {
	model := "gemini-2.0-flash-lite"
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	log.Printf("[Gemini] Using model: %s", model)
	log.Printf("[Gemini] Has screenshot: %v", screenshotBase64 != "")

	parts := []geminiPart{
		{
			Text: query,
		},
	}

	// Add image if provided
	if screenshotBase64 != "" {
		// Remove data URL prefix if present
		base64Data := screenshotBase64
		if strings.Contains(base64Data, ",") {
			splitParts := strings.Split(base64Data, ",")
			if len(splitParts) > 1 {
				base64Data = splitParts[1]
			}
		}

		log.Printf("[Gemini] Adding image data, length: %d", len(base64Data))
		parts = append(parts, geminiPart{
			InlineData: &geminiInlineData{
				MimeType: "image/png",
				Data:     base64Data,
			},
		})
	}

	reqPayload := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: parts,
			},
		},
		GenerationConfig: &geminiGenerationConfig{
			Temperature:     0.7,
			MaxOutputTokens: 1024,
		},
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		log.Printf("[Gemini] ERROR marshaling request: %v", err)
		return "", fmt.Errorf("marshal request: %w", err)
	}

	log.Printf("[Gemini] Request body length: %d", len(body))

	httpClient := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		log.Printf("[Gemini] ERROR creating request: %v", err)
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	log.Printf("[Gemini] Sending request to API...")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("[Gemini] ERROR sending request: %v", err)
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[Gemini] Response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[Gemini] ERROR response body: %s", string(bodyBytes))
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Gemini] ERROR reading response body: %v", err)
		return "", fmt.Errorf("read response: %w", err)
	}

	log.Printf("[Gemini] Response body: %s", string(bodyBytes))

	var result geminiResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Printf("[Gemini] ERROR decoding response: %v", err)
		return "", fmt.Errorf("decode response: %w", err)
	}

	log.Printf("[Gemini] Candidates count: %d", len(result.Candidates))
	if len(result.Candidates) > 0 {
		log.Printf("[Gemini] Parts count: %d", len(result.Candidates[0].Content.Parts))
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		log.Printf("[Gemini] ERROR: no response content in API response")
		return "", errors.New("no response from API")
	}

	responseText := strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text)
	log.Printf("[Gemini] Response text length: %d", len(responseText))
	return responseText, nil
}

// QueryCustom calls a custom OpenAI-compatible endpoint
func QueryCustom(ctx context.Context, query string, screenshotBase64 string, apiKey string, endpoint string) (string, error) {
	messages := []openAIMessage{
		{
			Role: "user",
			Content: []openAIContent{
				{
					Type: "text",
					Text: query,
				},
			},
		},
	}

	// Add image if provided
	if screenshotBase64 != "" {
		imageURL := normalizeDataURL(screenshotBase64)
		messages[0].Content = append(messages[0].Content, openAIContent{
			Type: "image_url",
			ImageURL: &openAIImageURL{
				URL: imageURL,
			},
		})
	}

	reqPayload := openAIRequest{
		Model:       "gpt-4-vision-preview",
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpClient := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no response from API")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// normalizeDataURL ensures the screenshot is a valid data URL
func normalizeDataURL(b64 string) string {
	b64 = strings.TrimSpace(b64)
	if b64 == "" {
		return ""
	}
	if strings.HasPrefix(b64, "data:") {
		return b64
	}
	return "data:image/png;base64," + b64
}

// OpenAI API structures
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openAIMessage struct {
	Role    string          `json:"role"`
	Content []openAIContent `json:"content"`
}

type openAIContent struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *openAIImageURL `json:"image_url,omitempty"`
}

type openAIImageURL struct {
	URL string `json:"url"`
}

type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
}

type openAIChoice struct {
	Message openAIMessageResponse `json:"message"`
}

type openAIMessageResponse struct {
	Content string `json:"content"`
}

// Anthropic API structures
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string             `json:"role"`
	Content []anthropicContent `json:"content"`
}

type anthropicContent struct {
	Type   string                `json:"type"`
	Text   string                `json:"text,omitempty"`
	Source *anthropicImageSource `json:"source,omitempty"`
}

type anthropicImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type anthropicResponse struct {
	Content []anthropicContentResponse `json:"content"`
}

type anthropicContentResponse struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Gemini API structures
type geminiRequest struct {
	Contents         []geminiContent         `json:"contents"`
	GenerationConfig *geminiGenerationConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text       string            `json:"text,omitempty"`
	InlineData *geminiInlineData `json:"inline_data,omitempty"`
}

type geminiInlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type geminiGenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}
