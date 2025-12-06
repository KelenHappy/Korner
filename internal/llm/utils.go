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
	// This handles formats like "analysis...assistantfinal<actual response>"
	markers := []string{"assistantfinal", "assistant_final", "final_response"}
	for _, marker := range markers {
		lowerText := strings.ToLower(cleaned)
		idx := strings.Index(lowerText, marker)
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
	
	// Strategy 2: Remove <think>...</think> blocks (DeepSeek-style reasoning)
	for {
		startIdx := strings.Index(strings.ToLower(cleaned), "<think>")
		if startIdx == -1 {
			break
		}
		endIdx := strings.Index(strings.ToLower(cleaned[startIdx:]), "</think>")
		if endIdx == -1 {
			// No closing tag, remove from <think> to end of line
			lineEnd := strings.Index(cleaned[startIdx:], "\n")
			if lineEnd != -1 {
				cleaned = cleaned[:startIdx] + cleaned[startIdx+lineEnd+1:]
			} else {
				cleaned = cleaned[:startIdx]
			}
		} else {
			// Remove the entire <think>...</think> block
			cleaned = cleaned[:startIdx] + cleaned[startIdx+endIdx+len("</think>"):]
		}
	}
	
	// Strategy 3: Remove <reasoning>...</reasoning> blocks
	for {
		startIdx := strings.Index(strings.ToLower(cleaned), "<reasoning>")
		if startIdx == -1 {
			break
		}
		endIdx := strings.Index(strings.ToLower(cleaned[startIdx:]), "</reasoning>")
		if endIdx == -1 {
			lineEnd := strings.Index(cleaned[startIdx:], "\n")
			if lineEnd != -1 {
				cleaned = cleaned[:startIdx] + cleaned[startIdx+lineEnd+1:]
			} else {
				cleaned = cleaned[:startIdx]
			}
		} else {
			cleaned = cleaned[:startIdx] + cleaned[startIdx+endIdx+len("</reasoning>"):]
		}
	}
	
	// Strategy 4: Remove common reasoning prefixes at the start
	// (like "analysis", "thinking", etc. that appear before the actual response)
	prefixes := []string{"analysis", "thinking", "reasoning"}
	for _, prefix := range prefixes {
		lowerCleaned := strings.ToLower(cleaned)
		if strings.HasPrefix(lowerCleaned, prefix) {
			// Find where the actual numbered list starts (e.g., "1. ")
			if idx := strings.Index(cleaned, "1."); idx != -1 && idx < 500 {
				cleaned = cleaned[idx:]
				break
			}
		}
	}
	
	// Trim whitespace but preserve internal formatting
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
