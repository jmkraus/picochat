package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"picochat/config"
	"picochat/console"
	"picochat/messages"
	"picochat/requests"
	"strings"
	"time"
)

// HandleChat sends a chat request to the configured model, streams the response,
// updates the chat history, and returns a summary message with elapsed time
// and token speed.
//
// Parameters:
//
//	cfg      - (optional) config data for unit tests, can be nil by default
//	history  - chat history to send and update
//	stop     - channel used to stop the spinner when the first token arrives
//
// Returns:
//
//	string  - summary message with elapsed time and token speed
//	error   - error encountered during the request or processing
func HandleChat(cfg *config.Config, history *messages.ChatHistory, stop chan struct{}) (string, error) {
	cfg, err := getConfig(cfg)
	if err != nil {
		return "", fmt.Errorf("config load failed: %w", err)
	}

	// evaluate reasoning
	var reasoning *messages.Reasoning
	if cfg.Reasoning {
		reasoning = &messages.Reasoning{Effort: "medium"}
	} else {
		reasoning = nil
	}

	reqBody := messages.ChatRequest{
		Model:     cfg.Model,
		Messages:  history.Messages,
		Stream:    true,
		Reasoning: reasoning,
		Think:     cfg.Reasoning,
		Options: &messages.ChatOptions{
			Temperature: cfg.Temperature,
			TopP:        cfg.TopP,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}

	chatURL, err := requests.CleanUrl(cfg.URL, "chat")
	if err != nil {
		return "", err
	}

	start := time.Now()
	response, err := http.Post(chatURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("http request error for model %s: %w", cfg.Model, err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var fullReply strings.Builder

	seconds := 0
	elapsed := "--:--"
	firstToken := true
	for {
		var res messages.StreamResponse
		if err := decoder.Decode(&res); err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("stream decoding error for model %s: %w", cfg.Model, err)
		}

		if res.Message.Content != "" {
			if firstToken {
				console.StopSpinner(cfg.Quiet, stop)
				firstToken = false
			}
			fmt.Print(res.Message.Content)
			fullReply.WriteString(res.Message.Content)
		}

		if res.Done {
			seconds, elapsed = elapsedTime(start)
			fmt.Println()
			break
		}
	}

	if fullReply.Len() == 0 {
		console.StopSpinner(cfg.Quiet, stop)
		return "", fmt.Errorf("no content received from model %s — possible config issue or invalid model?", cfg.Model)
	}

	cleanMsg := messages.TrimEmptyLines(fullReply.String())
	cfg.ImagePath = "" ////IMAGES
	err = history.Add(messages.RoleAssistant, cleanMsg, cfg.ImagePath)
	if err != nil {
		return "", fmt.Errorf("could not add message to history: %w", err)
	}
	speed := tokenSpeed(seconds, fullReply.String())
	msg := fmt.Sprintf("\nElapsed (mm:ss): %s | Tokens/sec: %.1f", elapsed, speed)
	return msg, nil
}

// elapsedTime returns the elapsed time in seconds and a formatted
// "MM:SS" string.  All calculations are performed in whole seconds
// to avoid floating‑point rounding differences.
//
// Parameters:
//
//	t - start time
//
// Returns:
//
//	int    - total elapsed seconds
//	string - formatted elapsed time "MM:SS"
func elapsedTime(t time.Time) (int, string) {
	elapsed := time.Since(t)

	totalSeconds := int(elapsed.Seconds())

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return totalSeconds, fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// tokenSpeed calculates the average number of tokens processed per
// unit of time (t).  It returns 0 when t is zero to avoid division by
// zero.  The result is rounded to one decimal place.
//
// Parameters:
//
//	t (int)    - elapsed time in seconds
//	s (string) - string containing the full reply
//
// Returns:
//
//	float64 - tokens per second
func tokenSpeed(t int, s string) float64 {
	if (t == 0) || (s == "") {
		return 0
	}

	tokens := messages.CalculateTokens(s)
	speed := float64(tokens) / float64(t)
	roundFactor := 10.0

	return math.Round(speed*roundFactor) / roundFactor
}

// getConfig is a factory function to ensure a properly loaded config.
//
// Parameters:
//
//	cfg - a config struct or nil
//
// Returns:
//
//	config.Config - a struct filled with config data
//	error         - an error if anything went wrong
func getConfig(cfg *config.Config) (*config.Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	return config.Get()
}
