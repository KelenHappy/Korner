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
	"sync"
	"time"
)

var (
	modelNameCache      = make(map[string]string)
	modelNameCacheMutex sync.RWMutex
)

// ModelsResponse represents the vLLM /v1/models response
type ModelsResponse struct {
	Data []ModelData `json:"data"`
}

type ModelData struct {
	ID string `json:"id"`
}

// getModelName fetches the actual model name from vLLM server
func getModelName(baseURL string, apiKey string) (string, error) {
	// Check cache first
	modelNameCacheMutex.RLock()
	if cached, ok := modelNameCache[baseURL]; ok {
		modelNameCacheMutex.RUnlock()
		return cached, nil
	}
	modelNameCacheMutex.RUnlock()

	// Fetch from server
	modelsURL := strings.TrimSuffix(baseURL, "/") + "/models"
	log.Printf("[GPT-OSS] Fetching available models from: %s", modelsURL)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, modelsURL, nil)
	if err != nil {
		return "", fmt.Errorf("create models request: %w", err)
	}

	if apiKey == "" {
		apiKey = "dummy-key"
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch models failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("models API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return "", fmt.Errorf("decode models response: %w", err)
	}

	if len(modelsResp.Data) == 0 {
		return "", errors.New("no models available on server")
	}

	modelName := modelsResp.Data[0].ID
	log.Printf("[GPT-OSS] Found model: %s", modelName)

	// Cache the result
	modelNameCacheMutex.Lock()
	modelNameCache[baseURL] = modelName
	modelNameCacheMutex.Unlock()

	return modelName, nil
}

// QueryGPTOSS calls AMD GPT-OSS-120B via vLLM endpoint
func QueryGPTOSS(ctx context.Context, query string, screenshotBase64 string, apiKey string, endpoint string, language string) (string, error) {
	if endpoint == "" {
		endpoint = "http://210.61.209.139:45014/v1/chat/completions"
	} else if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint = strings.TrimSuffix(endpoint, "/") + "/chat/completions"
	}

	// Extract base URL for model fetching
	baseURL := strings.TrimSuffix(endpoint, "/chat/completions")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL = strings.TrimSuffix(baseURL, "/") + "/v1"
	}

	log.Printf("[GPT-OSS] Using endpoint: %s", endpoint)

	// Fetch the actual model name from the server
	modelName, err := getModelName(baseURL, apiKey)
	if err != nil {
		log.Printf("[GPT-OSS] Warning: Could not fetch model name, using default: %v", err)
		modelName = "gpt-oss-120b" // Fallback to default
	}

	// Build system prompt based on language
	// Use plain text format only (no markdown)
	systemPrompt := `You are a helpful AI assistant. Follow these rules:
1. Use plain text format and \n only, no markdown syntax
2. Use numbered lists (1. 2. 3.) for steps or items
3. Use blank lines to separate sections
4. Provide clear, direct answers
5. When analyzing images, describe what you see concisely`

if language == "zh-TW" || language == "zh" || containsChinese(query) {
    systemPrompt = `你是 AI 助手。規則：
1. 純文字和\n，不用 Markdown
2. 用數字列表（1. 2. 3.）
3. 空行分段
4. 直接回答
5. 簡潔描述圖片
6. 用繁體中文回答`
}

	messages := []OpenAIMessage{
		{
			Role: "system",
			Content: []OpenAIContent{
				{
					Type: "text",
					Text: systemPrompt,
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

	// Build request payload
	// Note: Some vLLM servers may not support all OpenAI parameters
	// Start with basic parameters and add more if needed
	reqPayload := OpenAIRequest{
		Model:       modelName,
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
	}
	
	// Optionally add advanced parameters (may not be supported by all vLLM servers)
	// Uncomment if your vLLM server supports these:
	// reqPayload.TopP = 0.9
	// reqPayload.FrequencyPenalty = 0.3
	// reqPayload.PresencePenalty = 0.3

	log.Printf("[GPT-OSS] Using model: %s", modelName)

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
	log.Printf("[GPT-OSS] Raw response length: %d", len(responseText))
	
	if responseText == "" {
		log.Printf("[GPT-OSS] WARNING: Empty response from API")
		return "", errors.New("empty response from API")
	}
	
	// Log first 200 chars to see the format
	preview := responseText
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	log.Printf("[GPT-OSS] Raw response preview: %q", preview)
	
	// Clean up internal reasoning markers from the response
	originalText := responseText
	responseText = cleanResponseText(responseText)
	
	if responseText != originalText {
		log.Printf("[GPT-OSS] Cleaned response (removed %d chars)", len(originalText)-len(responseText))
		cleanedPreview := responseText
		if len(cleanedPreview) > 200 {
			cleanedPreview = cleanedPreview[:200] + "..."
		}
		log.Printf("[GPT-OSS] Cleaned response preview: %q", cleanedPreview)
	}
	
	log.Printf("[GPT-OSS] Success! Final response length: %d", len(responseText))
	return responseText, nil
}
