package backend

import "strings"

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
