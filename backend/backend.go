package backend

import (
	"picochat/config"
	"picochat/messages"
)

type ChatInput struct {
	Model       string
	Messages    []messages.Message
	Temperature float64
	TopP        float64
	Reasoning   bool
	Format      any
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
	// Keep existing behavior for now: ollama API semantics.
	return newOllamaClient(cfg.URL)
}
