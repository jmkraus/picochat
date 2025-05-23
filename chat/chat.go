package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"picochat/config"
	"picochat/requests"
	"picochat/types"
	"strings"
)

func HandleChat(cfg *config.Config, history *types.ChatHistory) error {
	reqBody := types.ChatRequest{
		Model:    cfg.Model,
		Messages: history.Messages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	chatURL, err := requests.CleanUrl(cfg.URL, "chat")
	if err != nil {
		return err
	}

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
			if res.PromptEvalCount != 0 && res.EvalCount != 0 {
				log.Printf("Token stats: prompt_eval_count=%d, eval_count=%d", res.PromptEvalCount, res.EvalCount)
			}
			break
		}
	}

	history.Add("assistant", fullReply.String())
	return nil
}
