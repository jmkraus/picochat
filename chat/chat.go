package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"picochat/config"
	"picochat/requests"
	"picochat/types"
	"strings"
	"time"
)

func HandleChat(cfg *config.Config, history *types.ChatHistory) error {
	reqBody := types.ChatRequest{
		Model:    cfg.Model,
		Messages: history.Messages,
		Stream:   true,
		Options: &types.ChatOptions{
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
	resp, err := http.Post(chatURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var fullReply strings.Builder

	for {
		var res types.StreamResponse
		if err := decoder.Decode(&res); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("stream error: %w", err)
		}

		fmt.Print(res.Message.Content)
		fullReply.WriteString(res.Message.Content)

		if res.Done {
			elapsed := time.Since(start)
			minutes := int(elapsed.Minutes())
			seconds := int(elapsed.Seconds()) % 60
			t := fmt.Sprintf("%02d:%02d", minutes, seconds)
			fmt.Println()
			fmt.Printf("Elapsed (mm:ss): %s; Tokens: prompt_eval_count=%d, eval_count=%d", t, res.PromptEvalCount, res.EvalCount)
			break
		}
	}

	history.Add("assistant", fullReply.String())
	return nil
}
