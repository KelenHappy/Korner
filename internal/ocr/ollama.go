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

// QueryOllama queries Ollama with optional image support
func QueryOllama(ctx context.Context, query string, imageBase64 string, endpoint string, language string) (string, error) {
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}

	apiURL := strings.TrimSuffix(endpoint, "/") + "/api/generate"
	log.Printf("[Ollama] Endpoint: %s", endpoint)
	log.Printf("[Ollama] Query: %s", query)
	log.Printf("[Ollama] Has image: %v", imageBase64 != "")

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

	// Prepare request with optional image
	reqPayload := OllamaRequest{
		Model:  "qwen3-vl:4b",
		Prompt: prompt,
		Stream: false,
	}

	// Add image if provided
	if imageBase64 != "" {
		// Remove data URL prefix if present
		imageData := imageBase64
		if strings.HasPrefix(imageData, "data:image/") {
			parts := strings.SplitN(imageData, ",", 2)
			if len(parts) == 2 {
				imageData = parts[1]
			}
		}
		reqPayload.Images = []string{imageData}
		log.Printf("[Ollama] Image data length: %d bytes", len(imageData))
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


// QueryOllamaWithWebSearch queries Ollama with web search results (when web search is enabled)
func QueryOllamaWithWebSearch(ctx context.Context, query string, endpoint string, language string) (string, error) {
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}

	log.Printf("[Ollama+WebSearch] Query: %s", query)

	// Step 1: Perform web search
	searchResults, err := WebSearch(ctx, query)
	if err != nil {
		log.Printf("[Ollama+WebSearch] Web search failed: %v", err)
		searchResults = nil
	}

	// Step 2: Format search results
	searchContext := FormatSearchResultsForLLM(searchResults, language)
	currentTime := time.Now().Format("2006年1月2日 15:04")

	// Step 3: Build prompt with search context
	var prompt string
	if language == "zh-TW" || language == "zh" {
		prompt = fmt.Sprintf(`你是台灣的 AI 助手。現在時間：%s（UTC+8）

網路搜尋結果：
%s

規則：純文字、數字列表、繁體中文、不提其他國家
如果搜尋結果不足，說「搜尋結果有限」

問題：%s`, currentTime, searchContext, query)
	} else {
		prompt = fmt.Sprintf(`You are an AI assistant. Current time: %s (UTC+8)

Search results:
%s

Rules: plain text, numbered lists, don't mention unrelated regions
If results insufficient, say so

Question: %s`, currentTime, searchContext, query)
	}

	// Step 4: Send to Ollama
	apiURL := strings.TrimSuffix(endpoint, "/") + "/api/generate"
	reqPayload := OllamaRequest{
		Model:  "qwen3-vl:4b",
		Prompt: prompt,
		Stream: false,
	}

	body, _ := json.Marshal(reqPayload)

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			host := req.URL.Host
			if strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "localhost") {
				return nil, nil
			}
			return http.ProxyFromEnvironment(req)
		},
	}
	httpClient := &http.Client{Timeout: 180 * time.Second, Transport: transport}

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
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return strings.TrimSpace(result.Response), nil
}

// WebSearch performs a web search using DuckDuckGo HTML search for better results
func WebSearch(ctx context.Context, query string) (*WebSearchResponse, error) {
	log.Printf("[WebSearch] Searching for: %s", query)

	// Use DuckDuckGo HTML lite version for better search results
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "zh-TW,zh;q=0.9,en;q=0.8")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	htmlContent := string(body)
	results := parseHTMLSearchResults(htmlContent)

	log.Printf("[WebSearch] Found %d results", len(results))

	return &WebSearchResponse{
		Query:   query,
		Results: results,
	}, nil
}

// parseHTMLSearchResults extracts search results from DuckDuckGo HTML response
func parseHTMLSearchResults(html string) []SearchResult {
	results := make([]SearchResult, 0)

	// Find all result blocks - look for result__a class (title links)
	// Pattern: <a rel="nofollow" class="result__a" href="...">Title</a>
	// And snippet: <a class="result__snippet" href="...">Snippet text</a>

	remaining := html
	for i := 0; i < 10; i++ { // Limit to 10 results
		// Find result title
		titleStart := strings.Index(remaining, `class="result__a"`)
		if titleStart == -1 {
			break
		}

		remaining = remaining[titleStart:]

		// Extract URL from href
		hrefStart := strings.Index(remaining, `href="`)
		if hrefStart == -1 {
			break
		}
		hrefStart += 6
		hrefEnd := strings.Index(remaining[hrefStart:], `"`)
		if hrefEnd == -1 {
			break
		}
		resultURL := remaining[hrefStart : hrefStart+hrefEnd]

		// Decode DuckDuckGo redirect URL
		if strings.Contains(resultURL, "uddg=") {
			if decoded, err := url.QueryUnescape(resultURL); err == nil {
				if idx := strings.Index(decoded, "uddg="); idx != -1 {
					resultURL = decoded[idx+5:]
					if ampIdx := strings.Index(resultURL, "&"); ampIdx != -1 {
						resultURL = resultURL[:ampIdx]
					}
				}
			}
		}

		// Extract title text
		titleTagStart := strings.Index(remaining, ">")
		if titleTagStart == -1 {
			break
		}
		titleTagEnd := strings.Index(remaining[titleTagStart:], "</a>")
		if titleTagEnd == -1 {
			break
		}
		title := remaining[titleTagStart+1 : titleTagStart+titleTagEnd]
		title = stripHTMLTags(title)
		title = strings.TrimSpace(title)

		// Move past this result
		remaining = remaining[titleTagStart+titleTagEnd:]

		// Try to find snippet
		snippet := ""
		snippetStart := strings.Index(remaining, `class="result__snippet"`)
		nextResultStart := strings.Index(remaining, `class="result__a"`)

		if snippetStart != -1 && (nextResultStart == -1 || snippetStart < nextResultStart) {
			snippetRemaining := remaining[snippetStart:]
			snippetTagStart := strings.Index(snippetRemaining, ">")
			if snippetTagStart != -1 {
				snippetTagEnd := strings.Index(snippetRemaining[snippetTagStart:], "</a>")
				if snippetTagEnd != -1 {
					snippet = snippetRemaining[snippetTagStart+1 : snippetTagStart+snippetTagEnd]
					snippet = stripHTMLTags(snippet)
					snippet = strings.TrimSpace(snippet)
				}
			}
		}

		if title != "" {
			results = append(results, SearchResult{
				Title:   title,
				URL:     resultURL,
				Snippet: snippet,
			})
		}
	}

	return results
}

// stripHTMLTags removes HTML tags from a string
func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(r)
		}
	}
	// Also decode common HTML entities
	output := result.String()
	output = strings.ReplaceAll(output, "&amp;", "&")
	output = strings.ReplaceAll(output, "&lt;", "<")
	output = strings.ReplaceAll(output, "&gt;", ">")
	output = strings.ReplaceAll(output, "&quot;", "\"")
	output = strings.ReplaceAll(output, "&#39;", "'")
	output = strings.ReplaceAll(output, "&nbsp;", " ")
	return output
}

// WebSearchResponse represents the aggregated search response
type WebSearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// FormatSearchResultsForLLM formats search results into context for LLM
func FormatSearchResultsForLLM(response *WebSearchResponse, language string) string {
	if response == nil || len(response.Results) == 0 {
		if language == "zh-TW" || language == "zh" {
			return "未找到相關搜尋結果。"
		}
		return "No search results found."
	}

	var sb strings.Builder

	for i, result := range response.Results {
		if i >= 5 {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
		if result.Snippet != "" && result.Snippet != result.Title {
			sb.WriteString(fmt.Sprintf("   %s\n", result.Snippet))
		}
		if result.URL != "" {
			sb.WriteString(fmt.Sprintf("   來源: %s\n", result.URL))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
