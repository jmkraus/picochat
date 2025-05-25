package types

import (
	"encoding/json"
	"fmt"
	"log"
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
	Messages     []Message
	Limit        int
	limitReached bool
}

func NewHistory(systemPrompt string, maxLimit int) *ChatHistory {
	return &ChatHistory{
		Messages: []Message{{Role: "system", Content: systemPrompt}},
		Limit:    maxLimit,
	}
}

func (h *ChatHistory) Add(role, content string) {
	h.Messages = append(h.Messages, Message{Role: role, Content: content})

	if h.Limit > 0 {
		h.Compress(h.Limit)
	}
}

func (h *ChatHistory) Get() []Message {
	return h.Messages
}

func (h *ChatHistory) Replace(newMessages []Message) {
	h.Messages = newMessages
}

func (h *ChatHistory) SaveHistoryToFile(filename string) (string, error) {
	basePath, err := paths.GetHistoryDir()
	if err != nil {
		return "", err
	}
	if filename == "" {
		filename = utils.EnsureSuffix(time.Now().Format("2006-01-02_15-04-05"))
	} else {
		filename = utils.EnsureSuffix(filename)
	}
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

func LoadHistoryFromFile(filename string) (*ChatHistory, error) {
	filename = utils.EnsureSuffix(filename)
	historyDir, err := paths.GetHistoryDir()
	if err != nil {
		return nil, fmt.Errorf("history dir not found.")
	}
	fullPath := filepath.Join(historyDir, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", fullPath, err)
	}

	var messages []Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}

	return &ChatHistory{Messages: messages}, nil
}

func (h *ChatHistory) ClearExceptSystemPrompt() {
	if len(h.Messages) > 1 {
		h.Messages = h.Messages[:1]
	}
	h.limitReached = false
}

func (h *ChatHistory) Compress(max int) {
	if len(h.Messages) < max {
		return
	}

	if !h.limitReached {
		log.Println("Message history limit of", h.Limit, "reached.")
		h.limitReached = true
	}

	keep := h.Messages[len(h.Messages)-(max-1):]
	h.Messages = append(h.Messages[:1], keep...)
}

func (h *ChatHistory) Len() int {
	return len(h.Messages)
}

func (h *ChatHistory) Max() int {
	return h.Limit
}

func (h *ChatHistory) IsEmpty() bool {
	return len(h.Messages) == 1
}
