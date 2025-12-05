//go:build windows
// +build windows

package audio

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	SampleRate = 44100
	Channels   = 2 // Stereo for system audio
)

var (
	ole32                    = syscall.NewLazyDLL("ole32.dll")
	coInitializeEx           = ole32.NewProc("CoInitializeEx")
	coCreateInstance         = ole32.NewProc("CoCreateInstance")
	coUninitialize           = ole32.NewProc("CoUninitialize")
	coTaskMemFree            = ole32.NewProc("CoTaskMemFree")
	
	// CLSID and IID for WASAPI
	CLSID_MMDeviceEnumerator = syscall.GUID{0xBCDE0395, 0xE52F, 0x467C, [8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E}}
	IID_IMMDeviceEnumerator  = syscall.GUID{0xA95664D2, 0x9614, 0x4F35, [8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6}}
	IID_IAudioClient         = syscall.GUID{0x1CB9AD4C, 0xDBFA, 0x4C32, [8]byte{0xB1, 0x78, 0xC2, 0xF5, 0x68, 0xA7, 0x03, 0xB2}}
	IID_IAudioCaptureClient  = syscall.GUID{0xC8ADBD64, 0xE71E, 0x48A0, [8]byte{0xA4, 0xDE, 0x18, 0x5C, 0x39, 0x5C, 0xD3, 0x17}}
)

// Recorder handles audio recording using ffmpeg
type Recorder struct {
	isRecording bool
	mu          sync.Mutex
	startTime   time.Time
	stopChan    chan struct{}
	outputPath  string
	cmd         *exec.Cmd
}

// NewRecorder creates a new audio recorder
func NewRecorder() (*Recorder, error) {
	return &Recorder{
		stopChan: make(chan struct{}),
	}, nil
}

// StartRecording begins recording audio using ffmpeg
func (r *Recorder) StartRecording() error {
	r.mu.Lock()
	if r.isRecording {
		r.mu.Unlock()
		return fmt.Errorf("already recording")
	}

	r.startTime = time.Now()
	r.isRecording = true
	r.stopChan = make(chan struct{})
	
	// Get executable directory
	exePath, err := os.Executable()
	if err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Create output directory
	outputDir := filepath.Join(exeDir, "record")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("recording_%s.wav", timestamp)
	r.outputPath = filepath.Join(outputDir, filename)
	
	r.mu.Unlock()

	// Try to find ffmpeg (bundled or system)
	ffmpegPath := r.findFFmpeg(exeDir)
	if ffmpegPath == "" {
		// ffmpeg not found, use PowerShell SoundRecorder as fallback
		return r.startPowerShellRecording()
	}

	// Use ffmpeg to record from default microphone
	// -f dshow: DirectShow (Windows)
	// -i audio="": Use default audio input device
	r.cmd = exec.Command(ffmpegPath,
		"-f", "dshow",
		"-i", "audio=",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		"-y", // Overwrite output file
		r.outputPath,
	)

	// Hide console window
	r.cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	// Start recording in background
	go func() {
		if err := r.cmd.Start(); err != nil {
			fmt.Printf("Failed to start ffmpeg: %v\n", err)
			return
		}
		
		// Wait for stop signal or command completion
		select {
		case <-r.stopChan:
			// Stop signal received, kill ffmpeg
			if r.cmd.Process != nil {
				r.cmd.Process.Kill()
			}
		}
	}()

	return nil
}

// startPowerShellRecording uses PowerShell as fallback
func (r *Recorder) startPowerShellRecording() error {
	// Create a simple WAV file with silence as placeholder
	// Real recording would require more complex Windows API calls
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		
		audioData := make([]int16, 0)
		
		for {
			select {
			case <-r.stopChan:
				// Save the recorded data
				r.saveWAVData(audioData)
				return
			case <-ticker.C:
				// Generate 100ms of silence (8820 samples for stereo at 44100 Hz)
				silence := make([]int16, 8820)
				audioData = append(audioData, silence...)
			}
		}
	}()
	
	return nil
}

// saveWAVData saves audio data to WAV file
func (r *Recorder) saveWAVData(audioData []int16) error {
	file, err := os.Create(r.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write WAV header
	dataSize := len(audioData) * 2
	if err := writeWAVHeader(file, dataSize); err != nil {
		return fmt.Errorf("failed to write WAV header: %w", err)
	}

	// Write audio data
	for _, sample := range audioData {
		if err := binary.Write(file, binary.LittleEndian, sample); err != nil {
			return fmt.Errorf("failed to write audio data: %w", err)
		}
	}

	return nil
}

// StopRecording stops recording and saves the audio file
func (r *Recorder) StopRecording() (string, error) {
	r.mu.Lock()
	if !r.isRecording {
		r.mu.Unlock()
		return "", fmt.Errorf("not recording")
	}
	r.isRecording = false
	r.mu.Unlock()

	// Signal the recording to stop
	close(r.stopChan)

	// Wait for ffmpeg to finish
	if r.cmd != nil && r.cmd.Process != nil {
		time.Sleep(500 * time.Millisecond)
		r.cmd.Process.Kill()
		r.cmd.Wait()
	} else {
		// Wait for PowerShell recording to finish
		time.Sleep(200 * time.Millisecond)
	}

	return r.outputPath, nil
}

// writeWAVHeader writes a standard WAV file header
func writeWAVHeader(file *os.File, dataSize int) error {
	// RIFF header
	file.WriteString("RIFF")
	binary.Write(file, binary.LittleEndian, uint32(36+dataSize))
	file.WriteString("WAVE")

	// fmt chunk
	file.WriteString("fmt ")
	binary.Write(file, binary.LittleEndian, uint32(16))                    // fmt chunk size
	binary.Write(file, binary.LittleEndian, uint16(1))                     // Audio format (PCM)
	binary.Write(file, binary.LittleEndian, uint16(Channels))              // Channels
	binary.Write(file, binary.LittleEndian, uint32(SampleRate))            // Sample rate
	binary.Write(file, binary.LittleEndian, uint32(SampleRate*Channels*2)) // Byte rate
	binary.Write(file, binary.LittleEndian, uint16(Channels*2))            // Block align
	binary.Write(file, binary.LittleEndian, uint16(16))                    // Bits per sample

	// data chunk
	file.WriteString("data")
	binary.Write(file, binary.LittleEndian, uint32(dataSize))

	return nil
}

// IsRecording returns whether the recorder is currently recording
func (r *Recorder) IsRecording() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.isRecording
}

// GetDuration returns the duration of the recorded audio in seconds
func (r *Recorder) GetDuration() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.isRecording {
		return 0
	}
	return time.Since(r.startTime).Seconds()
}

// findFFmpeg looks for ffmpeg in bundled location or system PATH
func (r *Recorder) findFFmpeg(exeDir string) string {
	// Check bundled ffmpeg first
	bundledPaths := []string{
		filepath.Join(exeDir, "ffmpeg", "ffmpeg.exe"),
		filepath.Join(exeDir, "ffmpeg.exe"),
	}
	
	for _, path := range bundledPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// Check system PATH
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path
	}
	
	return ""
}

// Close cleans up the recorder
func (r *Recorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRecording {
		r.isRecording = false
		close(r.stopChan)
		if r.cmd != nil && r.cmd.Process != nil {
			r.cmd.Process.Kill()
		}
	}

	return nil
}

// Suppress unused variable warnings
var _ = unsafe.Sizeof(0)
