package types

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"picochat/paths"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory struct {
	Messages []Message
}

func NewHistory(systemPrompt string) *ChatHistory {
	return &ChatHistory{
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
		},
	}
}

func (h *ChatHistory) Add(role, content string) {
	h.Messages = append(h.Messages, Message{Role: role, Content: content})
}

func (h *ChatHistory) Get() []Message {
	return h.Messages
}

func (h *ChatHistory) Replace(newMessages []Message) {
	h.Messages = newMessages
}

func (h *ChatHistory) SaveToFile() (string, error) {
	basePath := filepath.Join(filepath.Dir(paths.GetConfigPath()), "history")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", fmt.Errorf("could not create history dir: %w", err)
	}

	filename := time.Now().Format("2006-01-02_15-04-05") + ".chat"
	fullPath := filepath.Join(basePath, filename)

	data, err := json.MarshalIndent(h.Messages, "", "  ")
	if err != nil {
		return "", fmt.Errorf("could not marshal messages: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("could not write file: %w", err)
	}

	return filename, nil
}

func (h *ChatHistory) LoadFromFile(filename string) error {
	basePath := filepath.Join(filepath.Dir(paths.GetConfigPath()), "history")
	fullPath := filepath.Join(basePath, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	var msgs []Message
	if err := json.Unmarshal(data, &msgs); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	h.Replace(msgs)
	return nil
}
