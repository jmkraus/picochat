package backend

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"picochat/messages"
)

type openAIResponsesClient struct {
	baseURL string
	apiKey  string
}

type responsesRequest struct {
	Model       string               `json:"model"`
	Input       []responsesInputItem `json:"input"`
	Stream      bool                 `json:"stream"`
	Temperature float64              `json:"temperature,omitempty"`
	TopP        float64              `json:"top_p,omitempty"`
	Text        *responsesText       `json:"text,omitempty"`
}

type responsesText struct {
	Format map[string]any `json:"format,omitempty"`
}

type responsesInputItem struct {
	Role    string               `json:"role"`
	Content []responsesInputPart `json:"content"`
}

type responsesInputPart struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type responsesModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (c *openAIResponsesClient) ChatStream(input ChatInput, onChunk func(ChatChunk) error) (ChatFinal, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI API key")
	}
	if strings.TrimSpace(c.baseURL) == "" {
		return ChatFinal{}, fmt.Errorf("missing OpenAI base URL")
	}

	reqPayload := responsesRequest{
		Model:       input.Model,
		Input:       mapMessagesToResponsesInput(input.Messages),
		Stream:      true,
		Temperature: input.Temperature,
		TopP:        input.TopP,
	}
	if len(input.Format) != 0 {
		reqPayload.Text = &responsesText{Format: input.Format}
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return ChatFinal{}, fmt.Errorf("marshal json failed: %w", err)
	}

	endpoint, err := buildOpenAIURL(c.baseURL, "responses")
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

		thinking, content, done, err := parseResponsesEvent(data)
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

	return ChatFinal{Reasoning: fullThinking.String(), Content: fullContent.String()}, nil
}

func (c *openAIResponsesClient) GetAvailableModels() ([]string, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return nil, fmt.Errorf("missing OpenAI API key")
	}

	endpoint, err := buildOpenAIURL(c.baseURL, "models")
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

	var result responsesModelsResponse
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

func (c *openAIResponsesClient) GetServerVersion() (string, error) {
	return "OpenAI Responses API", nil
}

func mapMessagesToResponsesInput(in []messages.Message) []responsesInputItem {
	out := make([]responsesInputItem, 0, len(in))

	for _, msg := range in {
		parts := make([]responsesInputPart, 0, 1+len(msg.Images))

		if strings.TrimSpace(msg.Content) != "" {
			parts = append(parts, responsesInputPart{
				Type: "input_text",
				Text: msg.Content,
			})
		}

		for _, img := range msg.Images {
			if strings.TrimSpace(img) == "" {
				continue
			}
			parts = append(parts, responsesInputPart{
				Type:     "input_image",
				ImageURL: "data:image/png;base64," + img,
			})
		}

		out = append(out, responsesInputItem{
			Role:    msg.Role,
			Content: parts,
		})
	}

	return out
}

func parseResponsesEvent(data string) (thinking, content string, done bool, err error) {
	var evt map[string]any
	if err := json.Unmarshal([]byte(data), &evt); err != nil {
		return "", "", false, fmt.Errorf("decode response failed: %w", err)
	}

	eventType, _ := evt["type"].(string)
	switch eventType {
	case "response.output_text.delta":
		content, _ = evt["delta"].(string)
	case "response.reasoning.delta", "response.reasoning_summary_text.delta":
		thinking, _ = evt["delta"].(string)
	case "response.completed":
		done = true
	}

	return thinking, content, done, nil
}
