package llm

import "strings"

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

// containsChinese checks if a string contains Chinese characters
func containsChinese(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

// cleanResponseText removes internal reasoning markers and artifacts from model output
func cleanResponseText(text string) string {
	if text == "" {
		return text
	}
	
	cleaned := text
	
	// Strategy 1: Look for "assistantfinal" marker and extract everything after it
	// This is the most reliable way to get the actual response
	markers := []string{"assistantfinal", "assistant"}
	for _, marker := range markers {
		// Case-insensitive search
		lowerText := strings.ToLower(cleaned)
		idx := strings.Index(lowerText, strings.ToLower(marker))
		if idx != -1 {
			// Extract everything after the marker
			afterMarker := cleaned[idx+len(marker):]
			afterMarker = strings.TrimSpace(afterMarker)
			if afterMarker != "" {
				cleaned = afterMarker
				break
			}
		}
	}
	
	// Strategy 2: Remove common prefixes like "analysis", "User says", etc.
	prefixPatterns := []string{
		"analysis",
		"user says",
		"user asks",
		"we need to",
		"need to respond",
		"<think>",
		"</think>",
		"<reasoning>",
		"</reasoning>",
	}
	
	maxIterations := 10
	for i := 0; i < maxIterations; i++ {
		found := false
		lowerCleaned := strings.ToLower(cleaned)
		
		for _, pattern := range prefixPatterns {
			if strings.HasPrefix(lowerCleaned, pattern) {
				cleaned = cleaned[len(pattern):]
				cleaned = strings.TrimSpace(cleaned)
				found = true
				break
			}
		}
		
		if !found {
			break
		}
	}
	
	// Strategy 3: If text starts with a sentence describing what to do,
	// try to find where the actual response starts (usually after a period or newline)
	if strings.Contains(strings.ToLower(cleaned[:min(100, len(cleaned))]), "respond") ||
		strings.Contains(strings.ToLower(cleaned[:min(100, len(cleaned))]), "answer") {
		
		// Look for the first sentence break followed by actual content
		for i := 0; i < len(cleaned)-1; i++ {
			if cleaned[i] == '.' || cleaned[i] == '!' || cleaned[i] == '?' {
				// Check if next part looks like actual response (starts with capital or Chinese)
				remaining := strings.TrimSpace(cleaned[i+1:])
				if len(remaining) > 0 {
					firstChar := []rune(remaining)[0]
					// If starts with capital letter or Chinese character, this is likely the response
					if (firstChar >= 'A' && firstChar <= 'Z') || (firstChar >= 0x4E00 && firstChar <= 0x9FFF) {
						cleaned = remaining
						break
					}
				}
			}
		}
	}
	
	// Remove any leading special characters or artifacts
	cleaned = strings.TrimLeft(cleaned, ".:;,!?-_=+")
	cleaned = strings.TrimSpace(cleaned)
	
	// If we accidentally removed everything, return original
	if cleaned == "" && text != "" {
		return text
	}
	
	return cleaned
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
