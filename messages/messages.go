package messages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"picochat/console"
	"picochat/paths"
	"picochat/utils"
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
	Role      string   `json:"role"`
	Thinking  string   `json:"thinking,omitempty"`
	Content   string   `json:"content"`
	Images    []string `json:"images,omitempty"` ////IMAGES
	Reasoning string   `json:"-"`
}

type ChatHistory struct {
	Messages          []Message
	MaxContext        int
	MaxContextReached bool
}

// NewHistory creates a new ChatHistory with a system prompt and maximum context size.
//
// Parameters:
//
//	systemPrompt (string) - The system prompt as specified in the config file
//	maxContext (int)      - Maximum number of messages in the chat session
//
// Returns:
//
//	*ChatHistory
func NewHistory(systemPrompt string, maxContext int) *ChatHistory {
	return &ChatHistory{
		Messages:   []Message{{Role: RoleSystem, Content: systemPrompt}},
		MaxContext: maxContext,
	}
}

// Add appends a message with the specified role and content to the history.
// It handles role validation and compression.
//
// Parameters:
//
//	role (string)    - Role of the stored message (System, User, Assistant)
//	content (string) - Message body
//
// Returns:
//
//	error
func (h *ChatHistory) add(role, reasoning, content, image string) error {
	switch role {
	case RoleSystem, RoleUser, RoleAssistant:
		if role != RoleAssistant {
			// safety net - just in case
			reasoning = ""
		}

		////IMAGES
		var img []string
		if image != "" {
			b64, err := utils.ImageToBase64(image)
			if err != nil {
				return fmt.Errorf("convert image to base64 failed: %w", err)
			}
			img = append(img, b64)
		}

		h.Messages = append(h.Messages, Message{Role: role, Reasoning: reasoning, Content: content, Images: img})

		if h.MaxContext > 0 {
			h.Compress(h.MaxContext)
		}

		return nil
	default:
		return fmt.Errorf("invalid role '%s'", role)
	}
}

func (h *ChatHistory) AddUser(content, image string) error {
	return h.add(RoleUser, "", content, image)
}

func (h *ChatHistory) AddAssistant(reasoning, content string) error {
	return h.add(RoleAssistant, reasoning, content, "")
}

// Discard removes the last assistant message from the history if present.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func (h *ChatHistory) Discard() {
	if h.Len() <= 1 {
		return // only system prompt in history
	}

	lastIndex := h.Len() - 1
	if h.Messages[lastIndex].Role == RoleAssistant {
		h.Messages = h.Messages[:lastIndex]
	}
}

// Get returns the slice of messages in the history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	[]Message - Array of struct containing the full message history
func (h *ChatHistory) Get() []Message {
	return h.Messages
}

// GetLast returns the most recent message.
//
// Parameters:
//
//	none
//
// Returns:
//
//	Message - Struct containing the last message in the history session
func (h *ChatHistory) GetLast() Message {
	return h.Messages[h.Len()-1]
}

// GetLastRole searches backward for the last message with the specified
// role and returns it along with a boolean indicating success.
//
// Parameters:
//
//	role (string) - Role of the message (System, User, Assistant)
//
// Returns:
//
//	Message - Struct containing the message
//	bool    - Status if the search was successful (true or false)
func (h *ChatHistory) GetLastRole(role string) (Message, bool) {
	for i := h.Len() - 1; i >= 0; i-- {
		if h.Messages[i].Role == role {
			return h.Messages[i], true
		}
	}
	return Message{}, false
}

// Replace replaces the entire message slice with newMessages.
//
// Parameters:
//
//	newMessages []Message
//
// Returns:
//
//	none
func (h *ChatHistory) Replace(newMessages []Message) {
	h.Messages = newMessages
}

// SaveHistoryToFile writes the chat history to a file in the history
// directory and returns the filename.
//
// Parameters:
//
//	filename string
//
// Returns:
//
//	string
//	error
func (h *ChatHistory) SaveHistoryToFile(filename string) (string, error) {
	if strings.HasPrefix(filename, "#") {
		return "", fmt.Errorf("filename must not start with '#'")
	}

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

	if !paths.FileExists(fullPath) {
		data, err := json.MarshalIndent(h.Messages, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal messages failed: %w", err)
		}

		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			return "", fmt.Errorf("could not write file %s: %w", filename, err)
		}

		return filename, nil
	} else {
		return "", fmt.Errorf("filename already exists")
	}

}

// LoadHistoryFromFile reads a chat history from a file and returns
// a ChatHistory instance.
//
// Parameters:
//
//	filename string
//
// Returns:
//
//	*ChatHistory
//	error
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

// ClearExceptSystemPrompt removes all messages except the system prompt
// and resets context flag.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func (h *ChatHistory) ClearExceptSystemPrompt() {
	if h.Len() > 1 {
		h.Messages = h.Messages[:1]
	}
	h.MaxContextReached = false
}

// SetContextSize sets the maximum context size and trims history if necessary.
//
// Parameters:
//
//	max (int) - new maximum context size for the current history session
//
// Returns:
//
//	error
func (h *ChatHistory) SetContextSize(max int) error {
	if max < 3 || max > 100 {
		return fmt.Errorf("context size must be between 3 and 100")
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

// Compress reduces the history to the specified maximum number of messages,
// keeping the system prompt.
//
// Parameters:
//
//	max int
//
// Returns:
//
//	none
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

// Len returns the number of messages in the history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	int
func (h *ChatHistory) Len() int {
	return len(h.Messages)
}

// MaxCtx returns the maximum context size.
//
// Parameters:
//
//	none
//
// Returns:
//
//	int
func (h *ChatHistory) MaxCtx() int {
	return h.MaxContext
}

// IsEmpty checks if the history contains only the system prompt.
//
// Parameters:
//
//	none
//
// Returns:
//
//	bool
func (h *ChatHistory) IsEmpty() bool {
	return h.Len() == 1
}

// EstimateTokens estimates the total token count of the history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	float64
func (h *ChatHistory) EstimateTokens() float64 {
	total := 0.0
	for _, msg := range h.Messages {
		// use full data (incl. reasoning) if available
		text := msg.Reasoning + msg.Content
		total += CalculateTokens(text)
	}
	return total
}

// CalculateTokens estimates the number of tokens in a string based on word count.
//
// Parameters:
//
//	s (string) - the input string.
//
// Returns:
//
//	float64 - the estimated token count.
func CalculateTokens(s string) float64 {
	words := strings.Fields(s)
	return float64(len(words)) * 1.3
}
