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
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found. Please install ffmpeg to capture system audio")
	}

	// Create output directory
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Try to find stereo mix or loopback device
	devices, err := ListAudioDevices()
	var loopbackDevice string
	if err == nil && len(devices) > 0 {
		// Look for common loopback device names
		for _, device := range devices {
			if device == "Stereo Mix" || device == "立體聲混音" || 
			   device == "What U Hear" || device == "Wave Out Mix" {
				loopbackDevice = device
				break
			}
		}
	}

	// If no loopback device found, try default "Stereo Mix"
	if loopbackDevice == "" {
		loopbackDevice = "Stereo Mix"
	}

	// Use ffmpeg to capture audio from loopback device
	// -f dshow: DirectShow (Windows)
	// -i audio="device": Capture from specified device
	// -t: Duration
	// -rtbufsize: Increase buffer size to prevent drops
	cmd := exec.Command(ffmpegPath,
		"-f", "dshow",
		"-rtbufsize", "100M",
		"-i", fmt.Sprintf("audio=%s", loopbackDevice),
		"-t", fmt.Sprintf("%.0f", duration.Seconds()),
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		"-y", // Overwrite output file
		outputPath,
	)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// AudioDevice represents an audio device with its name and identifier
type AudioDevice struct {
	Name       string // Display name
	Identifier string // Device identifier for ffmpeg
}

// ListAudioDevices lists available audio devices using ffmpeg
func ListAudioDevices() ([]string, error) {
	devices, err := ListAudioDevicesDetailed()
	if err != nil {
		return nil, err
	}
	
	// Return identifiers for backward compatibility
	result := make([]string, len(devices))
	for i, dev := range devices {
		result[i] = dev.Identifier
	}
	return result, nil
}

// ListAudioDevicesDetailed lists available audio devices with full details
func ListAudioDevicesDetailed() ([]AudioDevice, error) {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found")
	}

	cmd := exec.Command(ffmpegPath, "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, _ := cmd.CombinedOutput()
	
	// Parse the output to extract device names and alternative names
	// ffmpeg outputs device list to stderr in format:
	// [dshow @ ...] "Device Name" (audio)
	// [dshow @ ...] Alternative name "@device_cm_..."
	devices := []AudioDevice{}
	lines := splitLines(string(output))
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		// Look for audio device lines
		if contains(line, "(audio)") && contains(line, "\"") {
			// Extract device name
			start := indexOf(line, "\"")
			if start >= 0 {
				end := indexOf(line[start+1:], "\"")
				if end >= 0 {
					deviceName := line[start+1 : start+1+end]
					
					// Look for alternative name in next line
					identifier := deviceName // Default to name
					if i+1 < len(lines) && contains(lines[i+1], "Alternative name") {
						altLine := lines[i+1]
						altStart := indexOf(altLine, "\"")
						if altStart >= 0 {
							altEnd := indexOf(altLine[altStart+1:], "\"")
							if altEnd >= 0 {
								identifier = altLine[altStart+1 : altStart+1+altEnd]
							}
						}
					}
					
					devices = append(devices, AudioDevice{
						Name:       deviceName,
						Identifier: identifier,
					})
				}
			}
		}
	}
	
	return devices, nil
}

// Helper functions
func splitLines(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == '\n' || c == '\r' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func contains(s, substr string) bool {
	return indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
