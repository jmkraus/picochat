package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"picochat/config"
	"picochat/messages"
	"picochat/requests"
	"strings"
	"time"
)

func HandleChat(cfg *config.Config, history *messages.ChatHistory) error {
	reqBody := messages.ChatRequest{
		Model:    cfg.Model,
		Messages: history.Messages,
		Stream:   true,
		Options: &messages.ChatOptions{
			Temperature: cfg.Temperature,
			TopP:        cfg.TopP,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	chatURL, err := requests.CleanUrl(cfg.URL, "chat")
	if err != nil {
		return err
	}

	start := time.Now()
	response, err := http.Post(chatURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("http request error for model %s: %w", cfg.Model, err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var fullReply strings.Builder

	for {
		var res messages.StreamResponse
		if err := decoder.Decode(&res); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("stream decoding error for model %s: %w", cfg.Model, err)
		}

		if res.Message.Content != "" {
			fmt.Print(res.Message.Content)
			fullReply.WriteString(res.Message.Content)
		}

		if res.Done {
			t := elapsedTime(start)
			fmt.Printf("\nElapsed (mm:ss): %s; Tokens: prompt_eval_count: %d, eval_count: %d", t, res.PromptEvalCount, res.EvalCount)
			break
		}
	}

	if fullReply.Len() == 0 {
		return fmt.Errorf("no content received from model %s — possible config issue or invalid model?", cfg.Model)
	}

	err = history.Add(messages.RoleAssistant, fullReply.String())
	if err != nil {
		return fmt.Errorf("could not add message to history: %w", err)
	}
	return nil
}

func elapsedTime(t time.Time) string {
	elapsed := time.Since(t)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
