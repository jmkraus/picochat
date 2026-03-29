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

// OpenAI Chat Completions API
func newOpenAIClient(baseURL, apiKey string) Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.openai.com"
	}
	return &openAIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}

// OpenAI Responses API
func newOpenAIResponsesClient(baseURL, apiKey string) Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.openai.com"
	}
	return &openAIResponsesClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}
