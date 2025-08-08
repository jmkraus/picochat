package messages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"picochat/console"
	"picochat/paths"
	"strings"
	"time"
)

const suffix string = ".chat"

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Raw     string `json:"-"`
}

type ChatHistory struct {
	Messages          []Message
	MaxContext        int
	MaxContextReached bool
}

func NewHistory(systemPrompt string, maxContext int) *ChatHistory {
	return &ChatHistory{
		Messages:   []Message{{Role: RoleSystem, Content: systemPrompt}},
		MaxContext: maxContext,
	}
}

func (h *ChatHistory) Add(role, content string) error {
	switch role {
	case RoleSystem, RoleUser, RoleAssistant:
		clearContent := StripReasoning(content)
		h.Messages = append(h.Messages, Message{Role: role, Content: clearContent, Raw: content})

		if h.MaxContext > 0 {
			h.Compress(h.MaxContext)
		}

		return nil
	default:
		return fmt.Errorf("invalid role '%s'", role)
	}
}

func (h *ChatHistory) Discard() {
	if h.Len() <= 1 {
		return // only system prompt in history
	}

	lastIndex := h.Len() - 1
	if h.Messages[lastIndex].Role == RoleAssistant {
		h.Messages = h.Messages[:lastIndex]
	}
}

func (h *ChatHistory) Get() []Message {
	return h.Messages
}

func (h *ChatHistory) GetLast() Message {
	return h.Get()[h.Len()-1]
}

func (h *ChatHistory) Replace(newMessages []Message) {
	h.Messages = newMessages
}

func (h *ChatHistory) SaveHistoryToFile(filename string) (string, error) {
	historyPath, err := paths.GetHistoryPath()
	if err != nil {
		return "", fmt.Errorf("history path not found: %w", err)
	}
	if historyPath == "" {
		return "", fmt.Errorf("empty history path returned")
	}

	if filename == "" {
		filename = paths.EnsureSuffix(time.Now().Format("2006-01-02_15-04-05"), suffix)
	} else {
		filename = paths.EnsureSuffix(filepath.Base(filename), suffix)
	}
	fullPath := filepath.Join(historyPath, filename)

	data, err := json.MarshalIndent(h.Messages, "", "  ")
	if err != nil {
		return "", fmt.Errorf("could not marshal messages: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("could not write file %s: %w", filename, err)
	}

	return filename, nil
}

func LoadHistoryFromFile(filename string) (*ChatHistory, error) {
	historyPath, err := paths.GetHistoryPath()
	if err != nil {
		return nil, fmt.Errorf("history path not found: %w", err)
	}
	if historyPath == "" {
		return nil, fmt.Errorf("empty history path returned")
	}

	filename = paths.EnsureSuffix(filepath.Base(filename), suffix)
	fullPath := filepath.Join(historyPath, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", fullPath, err)
	}

	var messages []Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("could not parse json in file %s: %w", filename, err)
	}

	return &ChatHistory{Messages: messages}, nil
}

func (h *ChatHistory) ClearExceptSystemPrompt() {
	if h.Len() > 1 {
		h.Messages = h.Messages[:1]
	}
	h.MaxContextReached = false
}

func (h *ChatHistory) SetContextSize(max int) error {
	if max < 5 || max > 100 {
		return fmt.Errorf("context size must be between 5 and 100")
	}
	if h.MaxContext == max {
		return nil
	}

	h.MaxContext = max

	if h.Len() >= max {
		start := h.Len() - (max - 1)
		if start < 1 {
			start = 1
		}
		h.Messages = append([]Message{h.Messages[0]}, h.Messages[start:]...)
	}

	h.MaxContextReached = h.Len() >= h.MaxContext
	return nil
}

func (h *ChatHistory) Compress(max int) {
	if h.Len() < max {
		return
	}

	if !h.MaxContextReached {
		console.Warn(fmt.Sprintf("Context size limit of %d reached.", h.MaxContext))
		h.MaxContextReached = true
	}

	keep := h.Messages[h.Len()-(max-1):]
	h.Messages = append(h.Messages[:1], keep...)
}

func (h *ChatHistory) Len() int {
	return len(h.Messages)
}

func (h *ChatHistory) MaxCtx() int {
	return h.MaxContext
}

func (h *ChatHistory) IsEmpty() bool {
	return h.Len() == 1
}

func (h *ChatHistory) EstimateTokens() int {
	total := 0
	for _, msg := range h.Messages {
		total += calculateTokens(msg.Content)
	}
	return total
}

func calculateTokens(text string) int {
	words := strings.Fields(text)
	return int(float64(len(words)) * 1.3)
}
