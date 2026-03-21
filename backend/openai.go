package backend

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"picochat/messages"
	"strings"
	"time"
)

type openAIClient struct {
	baseURL string
	apiKey  string
}

type openAIResponsesRequest struct {
	Model       string            `json:"model"`
	Input       []openAIInputItem `json:"input"`
	Stream      bool              `json:"stream"`
	Temperature float64           `json:"temperature,omitempty"`
	TopP        float64           `json:"top_p,omitempty"`
}

type openAIInputItem struct {
	Role    string            `json:"role"`
	Content []openAIInputPart `json:"content"`
}

type openAIInputPart struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type openAIModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (c *openAIClient) ChatStream(input ChatInput, onChunk func(ChatChunk) error) (ChatFinal, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI API key")
	}
	if strings.TrimSpace(c.baseURL) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI base URL")
	}

	payload := openAIResponsesRequest{
		Model:       input.Model,
		Input:       mapMessagesToOpenAIInput(input.Messages),
		Stream:      true,
		Temperature: input.Temperature,
		TopP:        input.TopP,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return ChatFinal{}, fmt.Errorf("marshal json failed: %w", err)
	}

	endpoint, err := joinURL(c.baseURL, "/v1/responses")
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

	var fullThinking strings.Builder
	var fullContent strings.Builder

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return ChatFinal{}, fmt.Errorf("read stream failed: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			if onChunk != nil {
				if err := onChunk(ChatChunk{Done: true}); err != nil {
					return ChatFinal{}, err
				}
			}
			break
		}

		thinking, content, done, err := parseOpenAIEvent(data)
		if err != nil {
			return ChatFinal{}, err
		}

		if thinking != "" {
			fullThinking.WriteString(thinking)
		}
		if content != "" {
			fullContent.WriteString(content)
		}

		if onChunk != nil {
			if err := onChunk(ChatChunk{
				Thinking: thinking,
				Content:  content,
				Done:     done,
			}); err != nil {
				return ChatFinal{}, err
			}
		}
	}

	return ChatFinal{
		Reasoning: fullThinking.String(),
		Content:   fullContent.String(),
	}, nil
}

func (c *openAIClient) GetAvailableModels() ([]string, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return nil, fmt.Errorf("missing OpenAI API key")
	}

	endpoint, err := joinURL(c.baseURL, "/v1/models")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch models failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(msg))
	}

	var result openAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	models := make([]string, 0, len(result.Data))
	for _, v := range result.Data {
		if v.ID != "" {
			models = append(models, v.ID)
		}
	}
	return models, nil
}

func (c *openAIClient) GetServerVersion() (string, error) {
	return "OpenAI Responses API", nil
}

func parseOpenAIEvent(data string) (thinking, content string, done bool, err error) {
	var evt map[string]any
	if err := json.Unmarshal([]byte(data), &evt); err != nil {
		return "", "", false, fmt.Errorf("decode response failed: %w", err)
	}

	eventType, _ := evt["type"].(string)
	switch eventType {
	case "response.output_text.delta":
		content, _ = evt["delta"].(string)
	case "response.reasoning_summary_text.delta", "response.reasoning.delta":
		thinking, _ = evt["delta"].(string)
	case "response.completed":
		done = true
	}

	return thinking, content, done, nil
}

func mapMessagesToOpenAIInput(in []messages.Message) []openAIInputItem {
	out := make([]openAIInputItem, 0, len(in))
	for _, msg := range in {
		parts := make([]openAIInputPart, 0, 1+len(msg.Images))
		if strings.TrimSpace(msg.Content) != "" {
			parts = append(parts, openAIInputPart{
				Type: "input_text",
				Text: msg.Content,
			})
		}
		for _, img := range msg.Images {
			if strings.TrimSpace(img) == "" {
				continue
			}
			parts = append(parts, openAIInputPart{
				Type:     "input_image",
				ImageURL: "data:image/png;base64," + img,
			})
		}
		out = append(out, openAIInputItem{
			Role:    msg.Role,
			Content: parts,
		})
	}
	return out
}

func joinURL(baseURL, path string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", fmt.Errorf("invalid http url string %w", err)
	}
	u.Path = strings.TrimRight(u.Path, "/") + path
	return u.String(), nil
}
