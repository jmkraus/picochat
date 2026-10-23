package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// postStreamingJSON posts a JSON payload to an OpenAI-compatible endpoint
// and consumes the response as an SSE stream.
func postStreamingJSON(
	baseURL string,
	apiKey string,
	endpoint string,
	payload any,
	parse parseEventFn,
	onChunk func(ChatChunk) error,
) (ChatFinal, error) {
	if strings.TrimSpace(apiKey) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI API key")
	}
	if strings.TrimSpace(baseURL) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI base URL")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return ChatFinal{}, fmt.Errorf("marshal json failed: %w", err)
	}

	url, err := buildOpenAIURL(baseURL, endpoint)
	if err != nil {
		return ChatFinal{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return ChatFinal{}, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := (&http.Client{Timeout: 0}).Do(req)
	if err != nil {
		return ChatFinal{}, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return ChatFinal{}, fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(msg))
	}

	return consumeSSEStream(resp.Body, parse, onChunk)
}
