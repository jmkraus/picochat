package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"picochat/messages"
	"strings"
)

type openAIClient struct {
	baseURL string
	apiKey  string
}

type openAIChatCompletionsRequest struct {
	Model       string              `json:"model"`
	Messages    []openAIChatMessage `json:"messages"`
	Stream      bool                `json:"stream"`
	Temperature *float64            `json:"temperature,omitempty"`
	TopP        *float64            `json:"top_p,omitempty"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type openAIStreamEvent struct {
	Choices []struct {
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			Reasoning        string `json:"reasoning"`
		} `json:"delta"`
		FinishReason any `json:"finish_reason"`
	} `json:"choices"`
}

// ChatStream sends a streaming chat completion request to an OpenAI-compatible endpoint.
//
// Parameters:
//
//	input (ChatInput)               - normalized chat payload
//	onChunk (func(ChatChunk) error) - callback for streamed chunks
//
// Returns:
//
//	ChatFinal - accumulated reasoning and content
//	error     - error if request/stream handling fails
func (c *openAIClient) ChatStream(input ChatInput, onChunk func(ChatChunk) error) (ChatFinal, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI API key")
	}
	if strings.TrimSpace(c.baseURL) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI base URL")
	}

	payload := openAIChatCompletionsRequest{
		Model:       input.Model,
		Messages:    mapMessagesToOpenAIChatMessages(input.Messages),
		Stream:      true,
		Temperature: input.Temperature,
		TopP:        input.TopP,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return ChatFinal{}, fmt.Errorf("marshal json failed: %w", err)
	}

	endpoint, err := buildOpenAIURL(c.baseURL, "chat/completions")
	if err != nil {
		return ChatFinal{}, err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return ChatFinal{}, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := (&http.Client{Timeout: 0}).Do(req)
	if err != nil {
		return ChatFinal{}, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return ChatFinal{}, fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(msg))
	}

	return consumeSSEStream(resp.Body, parseOpenAIChatCompletionEvent, onChunk)
}

// GetAvailableModels fetches models from the OpenAI-compatible models endpoint.
//
// Parameters:
//
//	none
//
// Returns:
//
//	[]string - list of model IDs
//	error    - error if request or decoding fails
func (c *openAIClient) GetAvailableModels() ([]string, error) {
	return fetchOpenAIModels(c.baseURL, c.apiKey)
}

// GetServerVersion returns a static descriptor for this backend protocol.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - protocol descriptor
//	error  - always nil
func (c *openAIClient) GetServerVersion() (string, error) {
	return "unknown (using OpenAI Chat Completions API)", nil
}

// parseOpenAIChatCompletionEvent parses one SSE event payload and extracts
// incremental reasoning/content plus completion state.
//
// Parameters:
//
//	data (string) - SSE event payload (JSON)
//
// Returns:
//
//	string - reasoning delta
//	string - content delta
//	bool   - done flag
//	error  - error if decoding fails
func parseOpenAIChatCompletionEvent(data string) (thinking, content string, done bool, err error) {
	var evt openAIStreamEvent
	if err := json.Unmarshal([]byte(data), &evt); err != nil {
		return "", "", false, fmt.Errorf("decode response failed: %w", err)
	}
	if len(evt.Choices) == 0 {
		return "", "", false, nil
	}

	delta := evt.Choices[0].Delta
	content = delta.Content
	if delta.ReasoningContent != "" {
		thinking = delta.ReasoningContent
	} else {
		thinking = delta.Reasoning
	}

	if evt.Choices[0].FinishReason != nil {
		done = true
	}
	return thinking, content, done, nil
}

// mapMessagesToOpenAIChatMessages maps internal messages to chat-completions format.
//
// Parameters:
//
//	in ([]messages.Message) - internal chat history messages
//
// Returns:
//
//	[]openAIChatMessage - mapped request messages
func mapMessagesToOpenAIChatMessages(in []messages.Message) []openAIChatMessage {
	out := make([]openAIChatMessage, 0, len(in))
	for _, msg := range in {
		if len(msg.Images) == 0 {
			out = append(out, openAIChatMessage{Role: msg.Role, Content: msg.Content})
			continue
		}

		parts := make([]map[string]any, 0, 1+len(msg.Images))
		if strings.TrimSpace(msg.Content) != "" {
			parts = append(parts, map[string]any{
				"type": "text",
				"text": msg.Content,
			})
		}
		for _, img := range msg.Images {
			if strings.TrimSpace(img) == "" {
				continue
			}
			parts = append(parts, map[string]any{
				"type": "image_url",
				"image_url": map[string]any{
					"url": img,
				},
			})
		}

		out = append(out, openAIChatMessage{Role: msg.Role, Content: parts})
	}
	return out
}
