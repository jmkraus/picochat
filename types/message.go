package types

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"picochat/paths"
	"picochat/utils"
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
	basePath := paths.GetHistoryDir()
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
	basePath := paths.GetHistoryDir()
	filename = utils.EnsureSuffix(filename)
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

func (h *ChatHistory) ClearExceptSystemPrompt() {
	if len(h.Messages) > 1 {
		h.Messages = h.Messages[:1]
	}
}

func (h *ChatHistory) Len() int {
	return len(h.Messages)
}

func (h *ChatHistory) IsEmpty() bool {
	return len(h.Messages) == 1
}
