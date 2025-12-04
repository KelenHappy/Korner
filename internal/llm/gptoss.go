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
func QueryGPTOSS(ctx context.Context, query string, screenshotBase64 string, apiKey string, endpoint string) (string, error) {
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
		Model:       modelName,
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
	}

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
	log.Printf("[GPT-OSS] Success! Response length: %d", len(responseText))
	return responseText, nil
}
