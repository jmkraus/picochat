package backend

import (
	"picochat/config"
	"picochat/messages"
	"strings"
)

type ChatInput struct {
	Model       string
	Messages    []messages.Message
	Temperature *float64
	TopP        *float64
	Reasoning   bool
	Effort      string
	Format      map[string]any
}

type ChatChunk struct {
	Thinking string
	Content  string
	Done     bool
}

type ChatFinal struct {
	Reasoning string
	Content   string
}

type Client interface {
	ChatStream(input ChatInput, onChunk func(ChatChunk) error) (ChatFinal, error)
	GetAvailableModels() ([]string, error)
	GetServerVersion() (string, error)
}

func New(cfg *config.Config) Client {
	baseURL := strings.TrimRight(cfg.URL, "/")
	switch strings.ToLower(strings.TrimSpace(cfg.Backend)) {
	case "responses":
		return &openAIResponsesClient{
			baseURL: baseURL,
			apiKey:  cfg.APIKey,
		}
	case "openai":
		return &openAIClient{
			baseURL: baseURL,
			apiKey:  cfg.APIKey,
		}
	default:
		return &ollamaClient{
			baseURL: baseURL,
		}
	}
}
