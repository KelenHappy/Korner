package meeting

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kelen/Korner/internal/audio"
)

// Summary represents a meeting summary result
type Summary struct {
	Content       string
	Transcription string
	AudioPath     string
	Duration      time.Duration
}

// Generator handles meeting summary generation
type Generator struct {
	transcriber *audio.WhisperTranscriber
}

// NewGenerator creates a new meeting summary generator
func NewGenerator() (*Generator, error) {
	transcriber, err := audio.NewWhisperTranscriberAuto("tiny")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Whisper: %w", err)
	}

	return &Generator{
		transcriber: transcriber,
	}, nil
}

// Generate transcribes audio and returns the transcription
func (g *Generator) Generate(ctx context.Context, audioPath string, language string) (*Summary, error) {
	log.Printf("[Meeting] Starting transcription for: %s", audioPath)

	// 檢查檔案
	fileInfo, err := os.Stat(audioPath)
	if err != nil {
		return nil, fmt.Errorf("無法讀取音訊檔案: %w", err)
	}
	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("音訊檔案是空的，請重新錄音")
	}

	log.Printf("[Meeting] File size: %d bytes", fileInfo.Size())

	// 轉錄音訊
	options := audio.DefaultTranscribeOptions()
	if language != "" {
		options.Language = language
	} else {
		options.Language = "zh"
	}

	log.Printf("[Meeting] Starting Whisper transcription...")
	transcription, err := g.transcriber.Transcribe(audioPath, options)
	if err != nil {
		log.Printf("[Meeting] Transcription error: %v", err)
		return nil, fmt.Errorf("轉錄失敗: %w", err)
	}

	if transcription == "" {
		log.Printf("[Meeting] Transcription is empty")
		return nil, fmt.Errorf("轉錄結果是空的，請檢查音訊檔案")
	}

	log.Printf("[Meeting] Transcription completed, length: %d", len(transcription))
	if len(transcription) > 100 {
		log.Printf("[Meeting] Transcription preview: %s...", transcription[:100])
	}

	return &Summary{
		Transcription: transcription,
		AudioPath:     audioPath,
	}, nil
}

// GenerateSummaryPrompt generates the prompt for meeting summary based on language
func GenerateSummaryPrompt(language string, transcription string) string {
	currentTime := time.Now().Format("2006-01-02 15:04")

	if language == "zh-TW" || language == "zh" {
		return fmt.Sprintf(`你是一位專業的會議記錄助理。請根據以下會議錄音的轉錄內容，生成一份完整且實用的會議智慧摘要。

請仔細分析會議內容，並按照以下格式輸出：

會議智慧摘要
=============

會議基本資訊
-----------
會議主題：[從對話中推斷出的會議主題]

主要討論內容
-----------
1. [討論主題]
   討論內容：[簡要說明討論的內容]
   關鍵觀點：[列出重要的觀點或意見]
2...

重要決議與共識
-------------
- [決議 1 - 具體說明決定了什麼]
- [決議 2 ...
行動項目與追蹤事項
-----------------
序號 | 行動項目 | 負責人 | 預計完成時間 | 優先級 | 狀態
-----|----------|--------|--------------|--------|------
1    | [項目]   | [人員] | [期限]       | [高/中/低] | 待處理

如果沒有明確的行動項目，請寫「本次會議未產生明確的行動項目」
待解決問題與風險
---------------
- [問題或風險 1]
...

關鍵結論與下一步
---------------
會議核心結論：[總結會議最重要的結論]
下一步行動：[說明接下來需要做什麼]
下次會議建議：[如果有提到，說明下次會議的時間或議題]

注意事項：
1. 請用繁體中文回覆
`, currentTime, transcription)
	}

	// English prompt
	return fmt.Sprintf(`You are a professional meeting assistant. Please generate a comprehensive and practical meeting summary based on the following transcription.

Please analyze the meeting content carefully and output in the following format:

Meeting Summary
===============

Basic Information
----------------
Meeting Topic: [Infer the meeting topic from the conversation]
Meeting Time: %s
Participants: [Identify participants from the conversation, or write "Not explicitly mentioned"]

Main Discussion Points
---------------------
1. [Discussion Topic]
   Content: [Brief description of the discussion]
   Key Points: [List important viewpoints or opinions]
...

Important Decisions & Consensus
-------------------------------
- [Decision 1 - Specify what was decided]
...

If no clear decisions were made, write "No clear decisions were reached in this meeting"

Action Items & Tracking
-----------------------
No. | Action Item | Owner | Due Date | Priority | Status
----|-------------|-------|----------|----------|--------
1   | [Item]      | [Person] | [Date] | [High/Med/Low] | Pending

If no clear action items, write "No clear action items were generated in this meeting"

Pending Issues & Risks
---------------------
- [Issue or Risk 1]
...

Key Conclusions & Next Steps
----------------------------
Core Conclusion: [Summarize the most important conclusion]
Next Actions: [Describe what needs to be done next]
Next Meeting Suggestion: [If mentioned, specify the time or agenda for the next meeting]

Notes:
1. Please respond in English
"`, currentTime, transcription)
}
