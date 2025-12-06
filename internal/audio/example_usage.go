package audio

import (
	"fmt"
	"time"
)

// ExampleBasicTranscription 展示基本轉錄用法
func ExampleBasicTranscription() {
	// 建立轉錄器（自動偵測 Whisper 和模型）
	transcriber, err := NewWhisperTranscriberAuto("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// 設定轉錄選項
	options := DefaultTranscribeOptions()
	options.Language = "zh" // 中文
	
	// 轉錄音訊檔案
	text, err := transcriber.Transcribe("recording.wav", options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("轉錄結果: %s\n", text)
}

// ExampleRecordAndTranscribe 展示錄音並轉錄
func ExampleRecordAndTranscribe() {
	// 建立整合錄音器
	recorder, err := NewRecorderWithTranscription("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer recorder.Close()
	
	// 開始錄音
	fmt.Println("開始錄音...")
	if err := recorder.StartRecording(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// 錄音 10 秒
	time.Sleep(10 * time.Second)
	
	// 停止並轉錄
	fmt.Println("停止錄音，開始轉錄...")
	options := DefaultTranscribeOptions()
	options.Language = "zh"
	
	audioPath, transcription, err := recorder.StopRecordingAndTranscribe(options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("音訊檔案: %s\n", audioPath)
	fmt.Printf("轉錄結果: %s\n", transcription)
}

// ExampleRecordWithProgress 展示帶進度的錄音轉錄
func ExampleRecordWithProgress() {
	recorder, err := NewRecorderWithTranscription("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer recorder.Close()
	
	// 開始錄音
	recorder.StartRecording()
	
	// 顯示錄音進度
	go func() {
		for recorder.IsRecording() {
			duration := recorder.GetDuration()
			fmt.Printf("\r錄音中... %.1f 秒", duration)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	// 錄音 5 秒
	time.Sleep(5 * time.Second)
	
	// 停止並轉錄（帶進度）
	fmt.Println("\n轉錄中...")
	options := DefaultTranscribeOptions()
	options.Language = "zh"
	
	audioPath, text, err := recorder.StopRecordingAndTranscribeWithProgress(
		options,
		func(progress string) {
			// 顯示轉錄進度
			fmt.Print(progress)
		},
	)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("\n完成！\n")
	fmt.Printf("音訊: %s\n", audioPath)
	fmt.Printf("文字: %s\n", text)
}

// ExampleQuickTranscription 展示快速轉錄
func ExampleQuickTranscription() {
	// 一行完成轉錄
	text, err := QuickTranscribe("audio.wav", "models/ggml-tiny.bin")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("轉錄結果: %s\n", text)
}

// ExampleMultiLanguage 展示多語言轉錄
func ExampleMultiLanguage() {
	transcriber, err := NewWhisperTranscriberAuto("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// 中文轉錄
	optionsChinese := DefaultTranscribeOptions()
	optionsChinese.Language = "zh"
	textChinese, _ := transcriber.Transcribe("chinese.wav", optionsChinese)
	fmt.Printf("中文: %s\n", textChinese)
	
	// 英文轉錄
	optionsEnglish := DefaultTranscribeOptions()
	optionsEnglish.Language = "en"
	textEnglish, _ := transcriber.Transcribe("english.wav", optionsEnglish)
	fmt.Printf("English: %s\n", textEnglish)
	
	// 自動偵測語言
	optionsAuto := DefaultTranscribeOptions()
	optionsAuto.Language = "auto"
	textAuto, _ := transcriber.Transcribe("unknown.wav", optionsAuto)
	fmt.Printf("Auto: %s\n", textAuto)
}

// ExampleSubtitleGeneration 展示字幕生成
func ExampleSubtitleGeneration() {
	transcriber, err := NewWhisperTranscriberAuto("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// 生成 SRT 字幕
	options := DefaultTranscribeOptions()
	options.Language = "zh"
	options.OutputSrt = true
	options.OutputVtt = true
	
	err = transcriber.TranscribeToFile("video_audio.wav", options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("字幕檔案已生成:")
	fmt.Println("- video_audio.srt")
	fmt.Println("- video_audio.vtt")
}

// ExampleTranslation 展示翻譯功能
func ExampleTranslation() {
	transcriber, err := NewWhisperTranscriberAuto("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// 轉錄並翻譯成英文
	options := DefaultTranscribeOptions()
	options.Language = "zh"
	options.Translate = true // 翻譯成英文
	
	text, err := transcriber.Transcribe("chinese.wav", options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("翻譯結果（英文）: %s\n", text)
}

// ExampleBatchTranscription 展示批次轉錄
func ExampleBatchTranscription() {
	transcriber, err := NewWhisperTranscriberAuto("tiny")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// 要轉錄的檔案列表
	audioFiles := []string{
		"recording1.wav",
		"recording2.wav",
		"recording3.wav",
	}
	
	options := DefaultTranscribeOptions()
	options.Language = "zh"
	
	// 批次轉錄
	for i, file := range audioFiles {
		fmt.Printf("轉錄 %d/%d: %s\n", i+1, len(audioFiles), file)
		
		text, err := transcriber.Transcribe(file, options)
		if err != nil {
			fmt.Printf("  錯誤: %v\n", err)
			continue
		}
		
		fmt.Printf("  結果: %s\n", text)
	}
}
