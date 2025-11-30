package history

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Conversation represents a single conversation entry
type Conversation struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	Question       string    `json:"question"`
	Answer         string    `json:"answer"`
	ScreenshotPath string    `json:"screenshot_path,omitempty"`
	Provider       string    `json:"provider"`
	Model          string    `json:"model,omitempty"`
}

// Manager handles conversation history
type Manager struct {
	historyDir string
}

// NewManager creates a new history manager
func NewManager() (*Manager, error) {
	// Get executable directory
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Create history directory
	historyDir := filepath.Join(exeDir, "history")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("create history directory: %w", err)
	}

	return &Manager{
		historyDir: historyDir,
	}, nil
}

// Save saves a conversation to history
func (m *Manager) Save(conv Conversation) error {
	if conv.ID == "" {
		conv.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if conv.Timestamp.IsZero() {
		conv.Timestamp = time.Now()
	}

	// Create filename based on date
	filename := fmt.Sprintf("%s.json", conv.Timestamp.Format("2006-01-02"))
	filepath := filepath.Join(m.historyDir, filename)

	// Read existing conversations for this day
	var conversations []Conversation
	if data, err := ioutil.ReadFile(filepath); err == nil {
		if err := json.Unmarshal(data, &conversations); err != nil {
			log.Printf("Warning: failed to parse existing history: %v", err)
			conversations = []Conversation{}
		}
	}

	// Append new conversation
	conversations = append(conversations, conv)

	// Save back to file
	data, err := json.MarshalIndent(conversations, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal conversations: %w", err)
	}

	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("write history file: %w", err)
	}

	log.Printf("Saved conversation to history: %s", filepath)
	return nil
}

// GetToday returns all conversations from today
func (m *Manager) GetToday() ([]Conversation, error) {
	filename := fmt.Sprintf("%s.json", time.Now().Format("2006-01-02"))
	return m.loadFile(filename)
}

// GetByDate returns all conversations from a specific date
func (m *Manager) GetByDate(date time.Time) ([]Conversation, error) {
	filename := fmt.Sprintf("%s.json", date.Format("2006-01-02"))
	return m.loadFile(filename)
}

// GetRecent returns the most recent N conversations
func (m *Manager) GetRecent(limit int) ([]Conversation, error) {
	files, err := ioutil.ReadDir(m.historyDir)
	if err != nil {
		return nil, fmt.Errorf("read history directory: %w", err)
	}

	// Sort files by name (date) in descending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})

	var allConversations []Conversation
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		conversations, err := m.loadFile(file.Name())
		if err != nil {
			log.Printf("Warning: failed to load %s: %v", file.Name(), err)
			continue
		}

		allConversations = append(allConversations, conversations...)

		if len(allConversations) >= limit {
			break
		}
	}

	// Sort by timestamp descending
	sort.Slice(allConversations, func(i, j int) bool {
		return allConversations[i].Timestamp.After(allConversations[j].Timestamp)
	})

	if len(allConversations) > limit {
		allConversations = allConversations[:limit]
	}

	return allConversations, nil
}

// GetAll returns all conversations
func (m *Manager) GetAll() ([]Conversation, error) {
	files, err := ioutil.ReadDir(m.historyDir)
	if err != nil {
		return nil, fmt.Errorf("read history directory: %w", err)
	}

	var allConversations []Conversation
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		conversations, err := m.loadFile(file.Name())
		if err != nil {
			log.Printf("Warning: failed to load %s: %v", file.Name(), err)
			continue
		}

		allConversations = append(allConversations, conversations...)
	}

	// Sort by timestamp descending
	sort.Slice(allConversations, func(i, j int) bool {
		return allConversations[i].Timestamp.After(allConversations[j].Timestamp)
	})

	return allConversations, nil
}

// Delete deletes a conversation by ID
func (m *Manager) Delete(id string) error {
	files, err := ioutil.ReadDir(m.historyDir)
	if err != nil {
		return fmt.Errorf("read history directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filepath := filepath.Join(m.historyDir, file.Name())
		conversations, err := m.loadFile(file.Name())
		if err != nil {
			continue
		}

		// Filter out the conversation with matching ID
		var filtered []Conversation
		found := false
		for _, conv := range conversations {
			if conv.ID == id {
				found = true
				continue
			}
			filtered = append(filtered, conv)
		}

		if found {
			// Save filtered conversations back
			if len(filtered) == 0 {
				// Delete file if empty
				return os.Remove(filepath)
			}

			data, err := json.MarshalIndent(filtered, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal conversations: %w", err)
			}

			return ioutil.WriteFile(filepath, data, 0644)
		}
	}

	return fmt.Errorf("conversation not found: %s", id)
}

// Clear deletes all history
func (m *Manager) Clear() error {
	files, err := ioutil.ReadDir(m.historyDir)
	if err != nil {
		return fmt.Errorf("read history directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := filepath.Join(m.historyDir, file.Name())
		if err := os.Remove(filepath); err != nil {
			log.Printf("Warning: failed to delete %s: %v", filepath, err)
		}
	}

	return nil
}

// loadFile loads conversations from a specific file
func (m *Manager) loadFile(filename string) ([]Conversation, error) {
	filepath := filepath.Join(m.historyDir, filename)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Conversation{}, nil
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	var conversations []Conversation
	if err := json.Unmarshal(data, &conversations); err != nil {
		return nil, fmt.Errorf("unmarshal conversations: %w", err)
	}

	return conversations, nil
}

// ExportToText exports all conversations to a text file
func (m *Manager) ExportToText(outputPath string) error {
	conversations, err := m.GetAll()
	if err != nil {
		return err
	}

	var content string
	content += "=================================================\n"
	content += "           Korner Chat History\n"
	content += "=================================================\n\n"

	for i, conv := range conversations {
		content += fmt.Sprintf("Chat #%d\n", i+1)
		content += fmt.Sprintf("Time: %s\n", conv.Timestamp.Format("2006-01-02 15:04:05"))
		content += fmt.Sprintf("Support: %s\n", conv.Provider)
		if conv.Model != "" {
			content += fmt.Sprintf("Model: %s\n", conv.Model)
		}
		if conv.ScreenshotPath != "" {
			content += fmt.Sprintf("Screenshot: %s\n", conv.ScreenshotPath)
		}
		content += fmt.Sprintf("\nQuestion:\n%s\n\n", conv.Question)
		content += fmt.Sprintf("Answer:\n%s\n\n", conv.Answer)
		content += "-------------------------------------------------\n\n"
	}

	return ioutil.WriteFile(outputPath, []byte(content), 0644)
}
