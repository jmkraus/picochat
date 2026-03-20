package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
//	string     - summary message with elapsed time and token speed
//	ChatResult - a struct containing output, elapsed time, and estimated tokens/s
//	error      - error if any
func HandleChat(cfg *config.Config, history *messages.ChatHistory, stop chan struct{}) (*ChatResult, error) {
	cfg, err := getConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	// evaluate reasoning
	var reasoning *Reasoning
	if cfg.Reasoning {
		reasoning = &Reasoning{Effort: "medium"}
	} else {
		reasoning = nil
	}

	reqBody := ChatRequest{
		Model:     cfg.Model,
		Messages:  history.Messages,
		Stream:    true,
		Reasoning: reasoning,
		Think:     cfg.Reasoning,
		Options: &ChatOptions{
			Temperature: cfg.Temperature,
			Top_p:       cfg.Top_p,
		},
		Format: cfg.SchemaFmt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal json failed: %w", err)
	}

	chatURL, err := requests.BuildCleanUrl(cfg.URL, "chat")
	if err != nil {
		return nil, err
	}

	start := time.Now()
	response, err := http.Post(chatURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("non-200 response: %d - %s", response.StatusCode, string(body))
	}

	decoder := json.NewDecoder(response.Body)
	var fullThinking strings.Builder
	var fullContent strings.Builder

	seconds := 0
	elapsed := "0s"
	firstToken := true
	firstContent := true
	streamPlain := cfg.OutputFmt == "plain"

	for {
		var res StreamResponse
		if err := decoder.Decode(&res); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decode response failed: %w", err)
		}

		if (res.Message.Content != "" ||
			(res.Message.Thinking != "" && cfg.Reasoning)) &&
			firstToken &&
			streamPlain {
			console.StopSpinner(cfg.Quiet, stop)
			firstToken = false
		}

		// Reasoning
		if res.Message.Thinking != "" {
			fullThinking.WriteString(res.Message.Thinking)
			if streamPlain && cfg.Reasoning {
				console.ColorPrint(console.Gray256, res.Message.Thinking)
			}
		}

		// Content
		if res.Message.Content != "" {
			fullContent.WriteString(res.Message.Content)
			if streamPlain {
				if firstContent {
					if fullThinking.Len() != 0 && cfg.Reasoning {
						fmt.Println()
					}
					firstContent = false
				}
				fmt.Print(res.Message.Content)
			}
		}

		if res.Done {
			if firstToken {
				console.StopSpinner(cfg.Quiet, stop)
			}
			seconds, elapsed = elapsedTime(start)
			break
		}
	}

	if fullContent.Len() == 0 {
		console.StopSpinner(cfg.Quiet, stop)
		return nil, fmt.Errorf("no content received from model %s", cfg.Model)
	}

	cleanReasoning, cleanContent := postProcessingChat(fullThinking.String(), fullContent.String())
	err = history.AddAssistant(cleanReasoning, cleanContent)
	if err != nil {
		return nil, fmt.Errorf("add message to history failed: %w", err)
	}
	speed := tokenSpeed(seconds, cleanReasoning+cleanContent)

	return &ChatResult{Output: cleanContent, Elapsed: elapsed, TokensPS: speed}, nil
}

// postProcessingChat separates reasoning part from content and cleans the text
// by stripping empty lines or extracting embedded thinking from content.
//
// Parameters:
//
//	thinking (string) - the thinking part of the response
//	content (string)  - the content part of the response
//
// Returns:
//
//	string - the cleaned reasoning part
//	string - the cleaned content part
func postProcessingChat(thinking, content string) (string, string) {
	// Case 1: thinking contains data
	// This should be the default for ollama reasoning models
	if thinking != "" {
		cleanReasoning := trimEmptyLines(thinking)
		cleanContent := trimEmptyLines(content)
		return cleanReasoning, cleanContent
	}

	// Case 2: Check if content contains <think> tags
	// This is the case for AI servers which embed thinking in the content part
	// If no think tag is found, the function returns an empty Reasoning part
	return splitReasoning(content)
}

// getConfig is a factory function to ensure a properly loaded config.
//
// Parameters:
//
//	cfg - a config struct or nil
//
// Returns:
//
//	config.Config - a struct with config data
//	error         - error if any
func getConfig(cfg *config.Config) (*config.Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	return config.Get()
}
