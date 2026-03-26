package backend

import "strings"

// Ollama chat API
func newOllamaClient(baseURL string) Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://localhost:11434"
	}
	return &ollamaClient{
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}
