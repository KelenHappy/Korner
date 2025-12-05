package llm

import (
	"testing"
)

func TestCleanResponseText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple assistantfinal marker",
			input:    "analysisUser says hi.assistantfinal你好！有什麼我可以幫忙的嗎？",
			expected: "你好！有什麼我可以幫忙的嗎？",
		},
		{
			name:     "Complex reasoning with assistantfinal",
			input:    "User says \"hi\". Need to respond in Traditional Chinese, friendly greeting.assistantfinal嗨！你好，有什麼我可以幫忙的嗎？",
			expected: "嗨！你好，有什麼我可以幫忙的嗎？",
		},
		{
			name:     "Analysis prefix only",
			input:    "analysisThis is a simple greeting.assistant你好",
			expected: "你好",
		},
		{
			name:     "No markers - clean text",
			input:    "Hello! How can I help you?",
			expected: "Hello! How can I help you?",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only markers",
			input:    "analysisassistantfinal",
			expected: "analysisassistantfinal", // Should return original if nothing left
		},
		{
			name:     "Reasoning with period",
			input:    "User asks a question. Need to respond clearly. Hello, I'm here to help!",
			expected: "Hello, I'm here to help!",
		},
		{
			name:     "Chinese response with English reasoning",
			input:    "Need to respond in Chinese.你好！很高興見到你。",
			expected: "你好！很高興見到你。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanResponseText(tt.input)
			if result != tt.expected {
				t.Errorf("cleanResponseText() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestContainsChinese(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Pure Chinese",
			input:    "你好世界",
			expected: true,
		},
		{
			name:     "Mixed Chinese and English",
			input:    "Hello 你好",
			expected: true,
		},
		{
			name:     "Pure English",
			input:    "Hello World",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Numbers and symbols",
			input:    "123!@#",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsChinese(tt.input)
			if result != tt.expected {
				t.Errorf("containsChinese(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
