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

type ollamaClient struct {
	baseURL string
}

type ollamaChatRequest struct {
	Model     string             `json:"model"`
	Messages  []messages.Message `json:"messages"`
	Reasoning *ollamaReasoning   `json:"reasoning,omitempty"`
	Options   *ollamaOptions     `json:"options,omitempty"`
	Stream    bool               `json:"stream"`
	Think     bool               `json:"think,omitempty"`
	Format    map[string]any     `json:"format,omitempty"`
}

type ollamaReasoning struct {
	Effort string `json:"effort,omitempty"`
}

type ollamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

type ollamaStreamResponse struct {
	Message struct {
		Thinking string `json:"thinking,omitempty"`
		Content  string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

type ollamaModelTag struct {
	Name string `json:"name"`
}

type ollamaModelTagsResponse struct {
	Models []ollamaModelTag `json:"models"`
}

type ollamaServerVersion struct {
	Version string `json:"version"`
}

func (c *ollamaClient) ChatStream(input ChatInput, onChunk func(ChatChunk) error) (ChatFinal, error) {
	var reasoning *ollamaReasoning
	if input.Reasoning {
		reasoning = &ollamaReasoning{Effort: "medium"}
	}

	ollamaMessages := normalizeOllamaImages(input.Messages)

	reqBody := ollamaChatRequest{
		Model:     input.Model,
		Messages:  ollamaMessages,
		Reasoning: reasoning,
		Options: &ollamaOptions{
			Temperature: input.Temperature,
			TopP:        input.TopP,
		},
		Stream: true,
		Think:  input.Reasoning,
		Format: input.Format,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return ChatFinal{}, fmt.Errorf("marshal json failed: %w", err)
	}

	chatURL, err := buildOllamaURL(c.baseURL, "chat")
	if err != nil {
		return ChatFinal{}, err
	}

	response, err := http.Post(chatURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ChatFinal{}, fmt.Errorf("http request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return ChatFinal{}, fmt.Errorf("non-200 response: %d - %s", response.StatusCode, string(body))
	}

	decoder := json.NewDecoder(response.Body)
	var fullThinking strings.Builder
	var fullContent strings.Builder

	for {
		var res ollamaStreamResponse
		if err := decoder.Decode(&res); err != nil {
			if err == io.EOF {
				break
			}
			return ChatFinal{}, fmt.Errorf("decode response failed: %w", err)
		}

		if res.Message.Thinking != "" {
			fullThinking.WriteString(res.Message.Thinking)
		}
		if res.Message.Content != "" {
			fullContent.WriteString(res.Message.Content)
		}

		if onChunk != nil {
			if err := onChunk(ChatChunk{
				Thinking: res.Message.Thinking,
				Content:  res.Message.Content,
				Done:     res.Done,
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

func (c *ollamaClient) GetAvailableModels() ([]string, error) {
	tagsURL, err := buildOllamaURL(c.baseURL, "tags")
	if err != nil {
		return nil, fmt.Errorf("fetch models failed: %w", err)
	}

	resp, err := http.Get(tagsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch models failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(body))
	}

	var result ollamaModelTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	var names []string
	for _, model := range result.Models {
		names = append(names, model.Name)
	}
	return names, nil
}

func (c *ollamaClient) GetServerVersion() (string, error) {
	versionURL, err := buildOllamaURL(c.baseURL, "version")
	if err != nil {
		return "", fmt.Errorf("fetch version failed: %w", err)
	}

	resp, err := http.Get(versionURL)
	if err != nil {
		return "", fmt.Errorf("fetch version failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(body))
	}

	var v ollamaServerVersion
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", fmt.Errorf("invalid server version response: %w", err)
	}

	if v.Version == "" {
		return "", fmt.Errorf("server response did not include a version")
	}

	return v.Version, nil
}

func normalizeOllamaImages(in []messages.Message) []messages.Message {
	out := make([]messages.Message, len(in))
	copy(out, in)

	for i := range out {
		if len(out[i].Images) == 0 {
			continue
		}

		imgs := make([]string, len(out[i].Images))
		for j, img := range out[i].Images {
			imgs[j] = stripDataURLPrefix(img)
		}
		out[i].Images = imgs
	}

	return out
}

func stripDataURLPrefix(s string) string {
	if !strings.HasPrefix(s, "data:") {
		return s
	}

	const marker = ";base64,"
	if idx := strings.Index(s, marker); idx >= 0 {
		return s[idx+len(marker):]
	}
	return s
}
