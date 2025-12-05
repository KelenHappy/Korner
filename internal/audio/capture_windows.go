//go:build windows
// +build windows

package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CaptureSystemAudio uses ffmpeg to capture system audio (loopback)
func CaptureSystemAudio(duration time.Duration, outputPath string) error {
	// Check if ffmpeg is available
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found. Please install ffmpeg to capture system audio")
	}

	// Create output directory
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use ffmpeg to capture audio from default audio device (loopback)
	// -f dshow: DirectShow (Windows)
	// -i audio="Stereo Mix": Capture from stereo mix (system audio)
	// -t: Duration
	cmd := exec.Command("ffmpeg",
		"-f", "dshow",
		"-i", "audio=Stereo Mix",
		"-t", fmt.Sprintf("%.0f", duration.Seconds()),
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		outputPath,
	)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ListAudioDevices lists available audio devices using ffmpeg
func ListAudioDevices() ([]string, error) {
	cmd := exec.Command("ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, _ := cmd.CombinedOutput()
	
	// ffmpeg outputs device list to stderr
	return []string{string(output)}, nil
}
