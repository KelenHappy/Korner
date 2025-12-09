//go:build windows
// +build windows

package audio

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

// WASAPILoopbackRecorder records system audio using WASAPI loopback (like OBS)
type WASAPILoopbackRecorder struct {
	isRecording bool
	mu          sync.Mutex
	audioData   []byte
	outputPath  string
	startTime   time.Time
	stopChan    chan struct{}
	format      *wca.WAVEFORMATEX
}

// NewWASAPILoopbackRecorder creates a new WASAPI loopback recorder
func NewWASAPILoopbackRecorder() *WASAPILoopbackRecorder {
	return &WASAPILoopbackRecorder{
		audioData: make([]byte, 0),
		stopChan:  make(chan struct{}),
	}
}

// StartRecording starts recording system audio using WASAPI loopback
func (r *WASAPILoopbackRecorder) StartRecording(outputPath string) error {
	r.mu.Lock()
	if r.isRecording {
		r.mu.Unlock()
		return fmt.Errorf("already recording")
	}
	r.isRecording = true
	r.outputPath = outputPath
	r.audioData = make([]byte, 0)
	r.startTime = time.Now()
	r.stopChan = make(chan struct{})
	r.mu.Unlock()

	// Start recording in background
	go r.recordLoop()

	return nil
}

// recordLoop captures audio using WASAPI loopback
func (r *WASAPILoopbackRecorder) recordLoop() {
	// Initialize COM
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		fmt.Printf("CoInitializeEx failed: %v\n", err)
		return
	}
	defer ole.CoUninitialize()

	// Create device enumerator
	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		fmt.Printf("CoCreateInstance failed: %v\n", err)
		return
	}
	defer mmde.Release()

	// Get default audio endpoint (render device for loopback)
	var mmd *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		fmt.Printf("GetDefaultAudioEndpoint failed: %v\n", err)
		return
	}
	defer mmd.Release()

	// Activate audio client
	var ac *wca.IAudioClient
	if err := mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &ac); err != nil {
		fmt.Printf("Activate failed: %v\n", err)
		return
	}
	defer ac.Release()

	// Get mix format
	var wfx *wca.WAVEFORMATEX
	if err := ac.GetMixFormat(&wfx); err != nil {
		fmt.Printf("GetMixFormat failed: %v\n", err)
		return
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(wfx)))
	r.format = wfx

	fmt.Printf("Audio format: %d Hz, %d channels, %d bits\n", wfx.NSamplesPerSec, wfx.NChannels, wfx.WBitsPerSample)

	// Initialize audio client in loopback mode
	var defaultPeriod, minimumPeriod wca.REFERENCE_TIME
	if err := ac.GetDevicePeriod(&defaultPeriod, &minimumPeriod); err != nil {
		fmt.Printf("GetDevicePeriod failed: %v\n", err)
		return
	}

	if err := ac.Initialize(wca.AUDCLNT_SHAREMODE_SHARED, wca.AUDCLNT_STREAMFLAGS_LOOPBACK, defaultPeriod, 0, wfx, nil); err != nil {
		fmt.Printf("Initialize failed: %v\n", err)
		return
	}

	// Get buffer size
	var bufferFrameCount uint32
	if err := ac.GetBufferSize(&bufferFrameCount); err != nil {
		fmt.Printf("GetBufferSize failed: %v\n", err)
		return
	}
	fmt.Printf("Buffer size: %d frames\n", bufferFrameCount)

	// Get capture client
	var acc *wca.IAudioCaptureClient
	if err := ac.GetService(wca.IID_IAudioCaptureClient, &acc); err != nil {
		fmt.Printf("GetService failed: %v\n", err)
		return
	}
	defer acc.Release()

	// Start recording
	if err := ac.Start(); err != nil {
		fmt.Printf("Start failed: %v\n", err)
		return
	}

	// Capture loop - use a ticker for consistent timing
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	fmt.Println("WASAPI loopback recording started (like OBS)...")

	for {
		select {
		case <-r.stopChan:
			ac.Stop()
			return
		case <-ticker.C:
			// Get next packet size
			var packetLength uint32
			if err := acc.GetNextPacketSize(&packetLength); err != nil {
				continue
			}

			for packetLength > 0 {
				var data *byte
				var numFramesAvailable uint32
				var flags uint32

				if err := acc.GetBuffer(&data, &numFramesAvailable, &flags, nil, nil); err != nil {
					break
				}

				if numFramesAvailable > 0 {
					// Calculate data size
					bytesPerFrame := int(wfx.NBlockAlign)
					dataSize := int(numFramesAvailable) * bytesPerFrame

					r.mu.Lock()
					if data != nil && (flags&0x2) == 0 { // AUDCLNT_BUFFERFLAGS_SILENT = 0x2
						audioBytes := unsafe.Slice(data, dataSize)

						// Convert 32-bit float to 16-bit PCM if needed
						if wfx.WBitsPerSample == 32 {
							numSamples := dataSize / 4
							for i := 0; i < numSamples; i++ {
								bits := binary.LittleEndian.Uint32(audioBytes[i*4 : i*4+4])
								floatVal := *(*float32)(unsafe.Pointer(&bits))
								
								// Clamp and convert to 16-bit
								if floatVal > 1.0 {
									floatVal = 1.0
								} else if floatVal < -1.0 {
									floatVal = -1.0
								}
								sample := int16(floatVal * 32767)
								
								var buf [2]byte
								binary.LittleEndian.PutUint16(buf[:], uint16(sample))
								r.audioData = append(r.audioData, buf[:]...)
							}
						} else {
							r.audioData = append(r.audioData, audioBytes...)
						}
					} else {
						// Silent buffer - write zeros
						silentBytes := dataSize
						if wfx.WBitsPerSample == 32 {
							silentBytes = dataSize / 2 // 32-bit to 16-bit
						}
						r.audioData = append(r.audioData, make([]byte, silentBytes)...)
					}
					r.mu.Unlock()
				}

				acc.ReleaseBuffer(numFramesAvailable)

				// Check for more packets
				if err := acc.GetNextPacketSize(&packetLength); err != nil {
					break
				}
			}
		}
	}
}

// StopRecording stops recording and saves to file
func (r *WASAPILoopbackRecorder) StopRecording() (string, error) {
	r.mu.Lock()
	if !r.isRecording {
		r.mu.Unlock()
		log.Printf("[WASAPI] Error: StopRecording called but not currently recording")
		return "", fmt.Errorf("未在錄音中")
	}
	r.isRecording = false
	r.mu.Unlock()

	// Signal stop
	close(r.stopChan)

	// Wait for the loop to finish
	time.Sleep(200 * time.Millisecond)

	// Save to temporary WAV file first
	tempPath := r.outputPath + ".temp.wav"
	originalPath := r.outputPath
	r.outputPath = tempPath
	
	err := r.saveWAV()
	if err != nil {
		return "", fmt.Errorf("failed to save WAV: %w", err)
	}

	// Convert to 16kHz using ffmpeg
	finalPath, err := r.convertTo16kHz(tempPath, originalPath)
	if err != nil {
		// If conversion fails, use original file
		fmt.Printf("Warning: ffmpeg conversion failed, using original file: %v\n", err)
		os.Rename(tempPath, originalPath)
		return originalPath, nil
	}

	// Clean up temp file
	os.Remove(tempPath)

	return finalPath, nil
}

// convertTo16kHz converts audio file to 16kHz WAV format using ffmpeg
func (r *WASAPILoopbackRecorder) convertTo16kHz(inputPath, outputPath string) (string, error) {
	// Check if ffmpeg is available
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	// ffmpeg -i input.wav -ar 16000 -ac 1 -sample_fmt s16 output.wav
	cmd := exec.Command("ffmpeg", 
		"-i", inputPath,
		"-ar", "16000",      // 16kHz sample rate
		"-ac", "1",          // mono
		"-sample_fmt", "s16", // 16-bit PCM
		"-y",                // overwrite output file
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %w, output: %s", err, string(output))
	}

	fmt.Printf("Converted audio to 16kHz: %s\n", outputPath)
	return outputPath, nil
}

// saveWAV saves recorded audio data to WAV file
func (r *WASAPILoopbackRecorder) saveWAV() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create output directory
	dir := filepath.Dir(r.outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(r.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Get format info (always output 16-bit PCM)
	channels := uint16(2)
	sampleRate := uint32(48000)
	bitsPerSample := uint16(16) // Always 16-bit output
	
	if r.format != nil {
		channels = r.format.NChannels
		sampleRate = r.format.NSamplesPerSec
	}

	// Write WAV header
	dataSize := len(r.audioData)
	
	// RIFF header
	file.WriteString("RIFF")
	binary.Write(file, binary.LittleEndian, uint32(36+dataSize))
	file.WriteString("WAVE")

	// fmt chunk
	file.WriteString("fmt ")
	binary.Write(file, binary.LittleEndian, uint32(16))
	binary.Write(file, binary.LittleEndian, uint16(1)) // PCM
	binary.Write(file, binary.LittleEndian, channels)
	binary.Write(file, binary.LittleEndian, sampleRate)
	binary.Write(file, binary.LittleEndian, sampleRate*uint32(channels)*uint32(bitsPerSample)/8) // Byte rate
	binary.Write(file, binary.LittleEndian, channels*bitsPerSample/8) // Block align
	binary.Write(file, binary.LittleEndian, bitsPerSample)

	// data chunk
	file.WriteString("data")
	binary.Write(file, binary.LittleEndian, uint32(dataSize))
	file.Write(r.audioData)

	fmt.Printf("Saved %d bytes of audio data\n", dataSize)
	return nil
}

// IsRecording returns whether currently recording
func (r *WASAPILoopbackRecorder) IsRecording() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.isRecording
}

// GetDuration returns recording duration in seconds
func (r *WASAPILoopbackRecorder) GetDuration() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.isRecording {
		return 0
	}
	return time.Since(r.startTime).Seconds()
}
