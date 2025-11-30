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
