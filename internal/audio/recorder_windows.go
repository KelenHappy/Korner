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

// RecordMode defines what audio sources to record
type RecordMode int

const (
	RecordMicrophone RecordMode = iota // Only microphone
	RecordSystem                       // Only system audio
	RecordBoth                         // Both microphone and system audio (mixed)
)

// Recorder handles audio recording using ffmpeg and WASAPI
type Recorder struct {
	isRecording    bool
	mu             sync.Mutex
	startTime      time.Time
	stopChan       chan struct{}
	outputPath     string
	cmd            *exec.Cmd
	mode           RecordMode
	wasapiRecorder *WASAPILoopbackRecorder // For system audio
	micTempPath    string                   // Temp path for mic recording
	sysTempPath    string                   // Temp path for system audio
}

// NewRecorder creates a new audio recorder with default mode (both mic + system)
func NewRecorder() (*Recorder, error) {
	return NewRecorderWithMode(RecordBoth)
}

// NewRecorderWithMode creates a new audio recorder with specified mode
func NewRecorderWithMode(mode RecordMode) (*Recorder, error) {
	return &Recorder{
		stopChan: make(chan struct{}),
		mode:     mode,
	}, nil
}

// StartRecording begins recording audio using ffmpeg (system audio + microphone mixed)
func (r *Recorder) StartRecording() error {
	r.mu.Lock()
	if r.isRecording {
		r.mu.Unlock()
		return fmt.Errorf("already recording")
	}

	r.startTime = time.Now()
	r.isRecording = true
	r.stopChan = make(chan struct{})
	
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create output directory in current working directory
	outputDir := filepath.Join(cwd, "record")
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
	ffmpegPath := r.findFFmpeg(cwd)
	if ffmpegPath == "" {
		return fmt.Errorf("ffmpeg not found. Please install ffmpeg")
	}

	// Use the ATR2100x USB Microphone for mic input
	micDeviceID := "@device_cm_{33D9A762-90C8-11D0-BD43-00A0C911CE86}\\wave_{984DAF9D-7658-4468-8B48-A8A188347AC3}"
	
	// Build ffmpeg command based on recording mode
	var args []string
	
	switch r.mode {
	case RecordMicrophone:
		// Only microphone using dshow
		args = []string{
			"-f", "dshow",
			"-rtbufsize", "100M",
			"-i", fmt.Sprintf("audio=%s", micDeviceID),
			"-acodec", "pcm_s16le",
			"-ar", "44100",
			"-ac", "1", // Mono
			"-y",
			r.outputPath,
		}
		fmt.Printf("Recording microphone (ATR2100x)...\n")
		
	case RecordSystem:
		// System audio using WASAPI loopback (like OBS)
		r.wasapiRecorder = NewWASAPILoopbackRecorder()
		err := r.wasapiRecorder.StartRecording(r.outputPath)
		if err != nil {
			return fmt.Errorf("failed to start WASAPI loopback recording: %w", err)
		}
		fmt.Println("Recording system audio using WASAPI loopback (like OBS)...")
		return nil
		
	case RecordBoth:
		// Record both system audio and microphone simultaneously
		// Create temp paths
		timestamp := time.Now().Format("20060102_150405")
		r.micTempPath = filepath.Join(filepath.Dir(r.outputPath), fmt.Sprintf("temp_mic_%s.wav", timestamp))
		r.sysTempPath = filepath.Join(filepath.Dir(r.outputPath), fmt.Sprintf("temp_sys_%s.wav", timestamp))
		
		// Start WASAPI loopback for system audio
		r.wasapiRecorder = NewWASAPILoopbackRecorder()
		err := r.wasapiRecorder.StartRecording(r.sysTempPath)
		if err != nil {
			return fmt.Errorf("failed to start system audio recording: %w", err)
		}
		
		// Start ffmpeg for microphone
		args = []string{
			"-f", "dshow",
			"-rtbufsize", "100M",
			"-i", fmt.Sprintf("audio=%s", micDeviceID),
			"-acodec", "pcm_s16le",
			"-ar", "44100",
			"-ac", "1", // Mono
			"-y",
			r.micTempPath,
		}
		fmt.Printf("Recording system audio (WASAPI) + microphone (ffmpeg) simultaneously...\n")
	}
	
	fmt.Printf("ffmpeg command: %s %v\n", ffmpegPath, args)
	r.cmd = exec.Command(ffmpegPath, args...)

	// Hide console window
	r.cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	// Get stdin pipe to send 'q' command for graceful shutdown
	stdin, err := r.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	// Capture stderr for debugging
	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start ffmpeg
	if err := r.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Log ffmpeg output for debugging in background
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				fmt.Printf("ffmpeg: %s", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// Wait for stop signal in background
	go func() {
		<-r.stopChan
		// Send 'q' to ffmpeg for graceful shutdown
		stdin.Write([]byte("q"))
		stdin.Close()
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
	mode := r.mode
	r.mu.Unlock()

	// If recording both (system + mic)
	if mode == RecordBoth && r.wasapiRecorder != nil && r.cmd != nil {
		// Stop WASAPI recorder
		sysPath, err := r.wasapiRecorder.StopRecording()
		if err != nil {
			return "", fmt.Errorf("failed to stop system audio: %w", err)
		}
		
		// Stop ffmpeg (microphone)
		close(r.stopChan)
		done := make(chan error, 1)
		go func() {
			done <- r.cmd.Wait()
		}()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			r.cmd.Process.Kill()
			<-done
		}
		
		// Mix the two files using ffmpeg
		fmt.Println("Mixing system audio and microphone...")
		err = r.mixAudioFiles(sysPath, r.micTempPath, r.outputPath)
		
		// Clean up temp files
		os.Remove(sysPath)
		os.Remove(r.micTempPath)
		
		if err != nil {
			return "", fmt.Errorf("failed to mix audio: %w", err)
		}
		
		r.wasapiRecorder = nil
		return r.outputPath, nil
	}

	// If using WASAPI loopback recorder only
	if r.wasapiRecorder != nil {
		outputPath, err := r.wasapiRecorder.StopRecording()
		r.wasapiRecorder = nil
		return outputPath, err
	}

	// Signal the recording to stop (this will send 'q' to ffmpeg)
	close(r.stopChan)

	// Wait for ffmpeg to finish gracefully
	if r.cmd != nil && r.cmd.Process != nil {
		// Wait for process to exit (with timeout)
		done := make(chan error, 1)
		go func() {
			done <- r.cmd.Wait()
		}()
		
		select {
		case err := <-done:
			// Process exited gracefully
			if err != nil {
				fmt.Printf("ffmpeg exited with error: %v\n", err)
			}
		case <-time.After(3 * time.Second):
			// Timeout, force kill
			fmt.Println("ffmpeg timeout, forcing kill...")
			r.cmd.Process.Kill()
			<-done
		}
	} else {
		// Wait for PowerShell recording to finish
		time.Sleep(200 * time.Millisecond)
	}

	return r.outputPath, nil
}

// mixAudioFiles mixes two audio files using ffmpeg
func (r *Recorder) mixAudioFiles(systemAudioPath, micPath, outputPath string) error {
	// Find ffmpeg
	cwd, _ := os.Getwd()
	ffmpegPath := r.findFFmpeg(cwd)
	if ffmpegPath == "" {
		return fmt.Errorf("ffmpeg not found")
	}

	// Mix using amerge filter
	cmd := exec.Command(ffmpegPath,
		"-i", systemAudioPath,
		"-i", micPath,
		"-filter_complex", "[0:a][1:a]amerge=inputs=2[aout]",
		"-map", "[aout]",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg mix failed: %w\nOutput: %s", err, string(output))
	}

	return nil
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

// getDefaultAudioDevice tries to get the default audio input device name
func (r *Recorder) getDefaultAudioDevice() (string, error) {
	devices, err := ListAudioDevices()
	if err != nil || len(devices) == 0 {
		return "", fmt.Errorf("no audio devices found")
	}
	
	// Return the first available device
	return devices[0], nil
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
