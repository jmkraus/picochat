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
	maxHistory   int
	contextLimit bool
}

func NewHistory(systemPrompt string, maxHistory int) *ChatHistory {
	return &ChatHistory{
		Messages:   []Message{{Role: "system", Content: systemPrompt}},
		maxHistory: maxHistory,
	}
}

func (h *ChatHistory) Add(role, content string) {
	h.Messages = append(h.Messages, Message{Role: role, Content: content})

	if h.maxHistory > 0 {
		h.Compress(h.maxHistory)
	}
}

func (h *ChatHistory) Get() []Message {
	return h.Messages
}

func (h *ChatHistory) Replace(newMessages []Message) {
	h.Messages = newMessages
}

func (h *ChatHistory) SaveHistoryToFile() (string, error) {
	basePath, err := paths.GetHistoryDir()
	if err != nil {
		return "", err
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
	h.contextLimit = false
}

func (h *ChatHistory) Compress(max int) {
	if len(h.Messages) < max {
		return
	}

	if !h.contextLimit {
		log.Println("Message history limit of", h.maxHistory, "reached.")
		h.contextLimit = true
	}

	keep := h.Messages[len(h.Messages)-(max-1):]
	h.Messages = append(h.Messages[:1], keep...)
}

func (h *ChatHistory) Len() int {
	return len(h.Messages)
}

func (h *ChatHistory) Max() int {
	return h.maxHistory
}

func (h *ChatHistory) IsEmpty() bool {
	return len(h.Messages) == 1
}
