package backend

import (
	"encoding/json"
	"fmt"
	"picochat/messages"
	"strings"
)

type openAIBase struct {
	baseURL string
	apiKey  string
}

func (b openAIBase) GetAvailableModels() ([]string, error) {
	return fetchOpenAIModels(b.baseURL, b.apiKey)
}

type openAIClient struct {
	openAIBase
}

type openAIChatCompletionsRequest struct {
	Model          string                `json:"model"`
	Messages       []openAIChatMessage   `json:"messages"`
	Stream         bool                  `json:"stream"`
	Temperature    *float64              `json:"temperature,omitempty"`
	TopP           *float64              `json:"top_p,omitempty"`
	ResponseFormat *openAIResponseFormat `json:"response_format,omitempty"`
}

type openAIResponseFormat struct {
	Type       string                  `json:"type"`
	JSONSchema *openAIJSONSchemaFormat `json:"json_schema,omitempty"`
}

type openAIJSONSchemaFormat struct {
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"`
	Strict bool           `json:"strict"`
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
	payload := openAIChatCompletionsRequest{
		Model:          input.Model,
		Messages:       mapMessagesToOpenAIChatMessages(input.Messages),
		Stream:         true,
		Temperature:    input.Temperature,
		TopP:           input.TopP,
		ResponseFormat: buildOpenAIResponseFormat(input.Format),
	}

	return postStreamingJSON(
		c.baseURL,
		c.apiKey,
		"chat/completions",
		payload,
		parseOpenAIChatCompletionEvent,
		onChunk,
	)
}

// buildOpenAIResponseFormat builds the Chat Completions response_format
// payload, or nil when no schema is provided.
func buildOpenAIResponseFormat(schema map[string]any) *openAIResponseFormat {
	if len(schema) == 0 {
		return nil
	}

	return &openAIResponseFormat{
		Type: "json_schema",
		JSONSchema: &openAIJSONSchemaFormat{
			Name:   "user",
			Schema: schema,
			Strict: false,
		},
	}
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

// parseOpenAIChatCompletionEvent extracts reasoning, content, and completion
// state from a Chat Completions SSE event.
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

// mapMessagesToOpenAIChatMessages maps internal messages to Chat Completions format.
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
