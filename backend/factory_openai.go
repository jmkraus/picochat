package backend

import "strings"

// NewOpenAIClient creates a backend client for OpenAI-compatible
// Chat Completions API.
func newOpenAIClient(baseURL, apiKey string) Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.openai.com"
	}
	return &openAIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}
