# Korner — Cross-Platform AI Desktop Assistant

Korner is a lightweight Agentic AI desktop assistant that allows users to capture screenshots and interact with AI instantly without switching between multiple applications or browser tabs. It provides seamless integration with an **AMD Instinct MI300X GPU running GPT OSS 120B**, enabling advanced text and visual analysis directly on your desktop.

## Features

* **Region Screenshot Capture**
  Select a specific area of your screen for AI analysis using a transparent overlay and drag-to-select interface.

* **AI Questioning with Images**
  Combine screenshots and user text queries and send them to the GPT OSS 120B Vision API for instant answers.

* **Floating Response Window**
  View AI responses in a movable, resizable window with options to copy, pin, or close the response.

* **Multi-Platform Consistency**
  Korner provides a consistent UI/UX experience across Windows and macOS, following native platform conventions.

* **Versatile Use Cases**
  * Academic: Capture code snippets, charts, or lecture slides and receive explanations or summaries.
  * Work: Summarize meeting notes, generate agendas, or assist with project management.
  * Personal: Plan travel itineraries, organize health routines, or track learning progress.
  * Entertainment: Receive recommendations for events, exhibitions, and leisure activities.

* **Flexible Interaction Modes**
  Supports text, voice, and visual input. Users can type queries, record audio, or submit screenshots for analysis.

## Platforms

* **Windows 10/11**

## Build
* **Windows 10/11**  
[Install Go](https://go.dev/dl/)    
[Install Python](https://www.python.org/)   
[Install Ollama](https://ollama.com/download)  
[Install FFmpeg](https://ffmpeg.org/)  

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails build --clean
```

## Install(IMPORTANT)
* **Windows 10/11** 
```sh
python -m pip install -U openai-whisper
ollama run qwen3-vl:4b
```

## Technical Overview

* **Framework:** Wails v2 (Go 1.22 + Vue 3)
* **Frontend:** Vue 3 + Tailwind CSS
* **Backend:** Go (handles screenshot capture, hotkeys, and LLM communication)
* **LLM Integration:** Connects to AMD GPT OSS 120B via OpenAI-compatible API
* **Packaging:**

  * Windows: `.exe` with system tray icon
  * macOS: Universal `.app` (Intel & Apple Silicon)

## Privacy and Security

* Screenshots are kept in memory only and are never written to disk.
* No user data is stored or shared without explicit consent.

## Future Enhancements

Planned upgrades include:

* Personalized recommendations based on user activity
* Proactive alerts and reminders
* Advanced multi-modal reasoning for more intelligent assistant capabilities

## Usage

1. Launch Korner and place the floating icon anywhere on your screen.
2. Use the global hotkey or click the icon to start a screenshot or AI interaction.
3. Submit your query (text, voice, or image).
4. View responses in the floating window and interact as needed.

## Goal

Korner aims to provide a truly integrated AI assistant for desktop users. By reducing the need to switch between applications or browser-based LLMs, it enhances productivity, learning, and daily life efficiency across personal and professional tasks.

## Demo

https://hackmd.io/@carina2992/rJU3JG5eZl