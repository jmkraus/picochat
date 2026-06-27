package backend

import "strings"

// Ollama chat API
func newOllamaClient(baseURL string) Client {
	return &ollamaClient{
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// OpenAI Chat Completions API
func newOpenAIClient(baseURL, apiKey string) Client {
	return &openAIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}

// OpenAI Responses API
func newOpenAIResponsesClient(baseURL, apiKey string) Client {
	return &openAIResponsesClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}
