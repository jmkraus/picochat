package backend

import "strings"

// NewOpenAIClient creates a backend client for OpenAI Responses API.
// It is intentionally not wired into New(cfg) yet.
func NewOpenAIClient(baseURL, apiKey string) Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.openai.com"
	}
	return &openAIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}
