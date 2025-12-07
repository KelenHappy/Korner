package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OllamaRequest represents the request to Ollama API
type OllamaRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Images []string `json:"images,omitempty"`
	Stream bool     `json:"stream"`
}

// OllamaResponse represents the response from Ollama API
type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// ExtractTextFromImage uses Ollama's vision model to extract text from an image
func ExtractTextFromImage(ctx context.Context, imageBase64 string, endpoint string) (string, error) {
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}
	
	// Remove data URL prefix if present
	imageData := imageBase64
	if strings.HasPrefix(imageData, "data:image/") {
		parts := strings.SplitN(imageData, ",", 2)
		if len(parts) == 2 {
			imageData = parts[1]
		}
	}
	
	apiURL := strings.TrimSuffix(endpoint, "/") + "/api/generate"
	log.Printf("[Ollama OCR] Endpoint: %s", endpoint)
	log.Printf("[Ollama OCR] API URL: %s", apiURL)
	
	// Test connectivity first
	testURL := strings.TrimSuffix(endpoint, "/") + "/api/tags"
	log.Printf("[Ollama OCR] Testing connectivity to: %s", testURL)
	testResp, testErr := http.Get(testURL)
	if testErr != nil {
		log.Printf("[Ollama OCR] Connectivity test failed: %v", testErr)
	} else {
		log.Printf("[Ollama OCR] Connectivity test successful: %d", testResp.StatusCode)
		testResp.Body.Close()
	}
	
	// Prompt for OCR
	prompt := `請仔細觀察這張圖片，並提取圖片中的所有文字內容。
要求：
1. 按照原文順序提取所有可見的文字
2. 保持原有的格式和結構
3. 如果有表格，請保持表格結構
4. 只輸出提取的文字，不要添加任何解釋或說明
5. 如果圖片中沒有文字，請描述圖片內容`
	
	reqPayload := OllamaRequest{
		Model:  "qwen3-vl:4b",
		Prompt: prompt,
		Images: []string{imageData},
		Stream: false,
	}
	
	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}
	
	log.Printf("[Ollama OCR] Sending request to: %s", apiURL)
	log.Printf("[Ollama OCR] Image data length: %d bytes", len(imageData))
	
	// Create HTTP client that bypasses proxy for localhost
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		Proxy: func(req *http.Request) (*url.URL, error) {
			// Don't use proxy for localhost/127.0.0.1
			host := req.URL.Host
			if host == "127.0.0.1:11434" || host == "127.0.0.1" || 
			   host == "localhost:11434" || host == "localhost" {
				log.Printf("[Ollama OCR] Bypassing proxy for: %s", host)
				return nil, nil
			}
			// Use system proxy for other requests
			return http.ProxyFromEnvironment(req)
		},
	}
	
	httpClient := &http.Client{
		Timeout:   120 * time.Second,
		Transport: transport,
	}
	
	// Retry logic for connection issues
	var resp *http.Response
	var reqErr error
	maxRetries := 5
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Create new reader for each attempt (important!)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("create request: %w", err)
		}
		
		req.Header.Set("Content-Type", "application/json")
		
		log.Printf("[Ollama OCR] Attempt %d/%d...", attempt, maxRetries)
		resp, reqErr = httpClient.Do(req)
		
		if reqErr == nil {
			log.Printf("[Ollama OCR] Connection successful on attempt %d", attempt)
			break
		}
		
		log.Printf("[Ollama OCR] Attempt %d failed: %v", attempt, reqErr)
		if attempt < maxRetries {
			waitTime := time.Duration(attempt*2) * time.Second
			log.Printf("[Ollama OCR] Waiting %v before retry...", waitTime)
			time.Sleep(waitTime)
		}
	}
	
	if reqErr != nil {
		return "", fmt.Errorf("http request failed after %d attempts: %w", maxRetries, reqErr)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[Ollama OCR] ERROR response: %s", string(bodyBytes))
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	
	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	
	extractedText := strings.TrimSpace(result.Response)
	log.Printf("[Ollama OCR] Extracted text length: %d", len(extractedText))
	
	if extractedText == "" || strings.Contains(extractedText, "沒有文字") {
		return "", nil // No text found in image
	}
	
	return extractedText, nil
}


// QueryOllamaWithWebSearch queries Ollama directly using generate API
func QueryOllamaWithWebSearch(ctx context.Context, query string, endpoint string, language string) (string, error) {
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}

	apiURL := strings.TrimSuffix(endpoint, "/") + "/api/generate"
	log.Printf("[Ollama] Endpoint: %s", endpoint)
	log.Printf("[Ollama] Query: %s", query)

	// Build prompt based on language
	prompt := query
	if language == "zh-TW" || language == "zh" {
		prompt = `規則：
1. 純文字，不用 Markdown
2. 用數字列表（1. 2. 3.）
3. 空行分段
4. 請用繁體中文直接回答以下問題：

` + query
	} else {
		prompt = `Rules:
1. Pure text, no Markdown
2. Use numbered lists (1. 2. 3.)
3. Empty line breaks
4. Please answer the following question directly:

` + query
	}

	reqPayload := OllamaRequest{
		Model:  "qwen3-vl:4b",
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	log.Printf("[Ollama] Sending request to: %s", apiURL)

	// Create HTTP client that bypasses proxy for localhost
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		Proxy: func(req *http.Request) (*url.URL, error) {
			host := req.URL.Host
			if host == "127.0.0.1:11434" || host == "127.0.0.1" ||
				host == "localhost:11434" || host == "localhost" {
				return nil, nil
			}
			return http.ProxyFromEnvironment(req)
		},
	}

	httpClient := &http.Client{
		Timeout:   180 * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[Ollama] ERROR response: %s", string(bodyBytes))
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	responseText := strings.TrimSpace(result.Response)
	log.Printf("[Ollama] Response length: %d", len(responseText))

	return responseText, nil
}
