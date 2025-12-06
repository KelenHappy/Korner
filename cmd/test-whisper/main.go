package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"Korner/internal/audio"
)

func main() {
	// 命令行參數
	audioFile := flag.String("file", "", "Audio file to transcribe")
	modelName := flag.String("model", "tiny", "Model name (tiny, base, small, medium, large)")
	language := flag.String("lang", "auto", "Language code (zh, en, auto, etc.)")
	record := flag.Bool("record", false, "Record audio before transcribing")
	duration := flag.Int("duration", 10, "Recording duration in seconds")
	translate := flag.Bool("translate", false, "Translate to English")
	outputSrt := flag.Bool("srt", false, "Output SRT subtitle file")
	
	flag.Parse()
	
	// 如果需要錄音
	if *record {
		fmt.Printf("Recording for %d seconds...\n", *duration)
		audioPath, transcription, err := audio.RecordAndTranscribe(
			time.Duration(*duration)*time.Second,
			*modelName,
			*language,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("\n✅ Recording saved: %s\n", audioPath)
		fmt.Printf("📝 Transcription: %s\n", transcription)
		return
	}
	
	// 如果沒有指定檔案，顯示幫助
	if *audioFile == "" {
		fmt.Println("Whisper.cpp Transcription Tool")
		fmt.Println("\nUsage:")
		fmt.Println("  Transcribe audio file:")
		fmt.Println("    test-whisper -file audio.wav -lang zh")
		fmt.Println("\n  Record and transcribe:")
		fmt.Println("    test-whisper -record -duration 10 -lang zh")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(0)
	}
	
	// 檢查檔案是否存在
	if _, err := os.Stat(*audioFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Audio file not found: %s\n", *audioFile)
		os.Exit(1)
	}
	
	// 建立轉錄器
	fmt.Printf("Initializing Whisper with model: %s\n", *modelName)
	transcriber, err := audio.NewWhisperTranscriberAuto(*modelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nPlease ensure:\n")
		fmt.Fprintf(os.Stderr, "1. Whisper.cpp is installed (main.exe)\n")
		fmt.Fprintf(os.Stderr, "2. Model file exists (ggml-%s.bin)\n", *modelName)
		os.Exit(1)
	}
	
	// 設定轉錄選項
	options := audio.DefaultTranscribeOptions()
	options.Language = *language
	options.Translate = *translate
	options.OutputSrt = *outputSrt
	
	// 開始轉錄
	fmt.Printf("Transcribing: %s\n", *audioFile)
	fmt.Printf("Language: %s\n", options.Language)
	
	startTime := time.Now()
	
	// 使用帶進度的轉錄
	text, err := transcriber.TranscribeWithProgress(*audioFile, options, func(progress string) {
		fmt.Print(progress)
	})
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
		os.Exit(1)
	}
	
	elapsed := time.Since(startTime)
	
	// 顯示結果
	fmt.Printf("\n\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("✅ Transcription completed in %.2f seconds\n", elapsed.Seconds())
	fmt.Printf("%s\n\n", strings.Repeat("=", 60))
	fmt.Printf("📝 Result:\n%s\n\n", text)
	
	// 如果輸出了 SRT 檔案，顯示訊息
	if *outputSrt {
		srtFile := *audioFile[:len(*audioFile)-len(filepath.Ext(*audioFile))] + ".srt"
		fmt.Printf("📄 SRT file saved: %s\n", srtFile)
	}
}
