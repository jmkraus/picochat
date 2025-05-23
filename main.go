package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"picochat/args"
	"picochat/command"
	"picochat/config"
	"picochat/requests"
	"picochat/types"
	"picochat/version"
	"strings"
)

func readMultilineInput() (string, bool, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if firstLine && strings.HasPrefix(trimmed, "/") {
			return trimmed, true, false // input, isCommand, quit=false
		}

		if trimmed == "/done" {
			break
		}

		lines = append(lines, line)
		firstLine = false
	}

	return strings.Join(lines, "\n"), false, false
}

// main.go

func handleChat(cfg *config.Config, history *types.ChatHistory) error {
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

func main() {
	args.Parse()

	if *args.ShowVersion {
		fmt.Println("picochat version is", version.Version)
		os.Exit(0)
	}

	err := config.Load()
	if err != nil {
		log.Fatalf("Error while loading configuration: %v", err)
	}
	cfg := config.Get()

	var history *types.ChatHistory
	if *args.HistoryFile != "" {
		history, err = types.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			log.Fatalf("Could not load history: %v", err)
		}
	} else {
		history = types.NewHistory(cfg.Prompt, cfg.Context)
	}

	log.Println("Chat with PicoAI started. Help with '/?'.")

	for {
		fmt.Print("\n>>> ")

		input, isCommand, quit := readMultilineInput()
		if quit {
			break
		}

		if isCommand {
			result := command.Handle(input, history, os.Stdin)
			if result.Output != "" {
				fmt.Println(result.Output)
			}
			if result.Quit {
				log.Println("Chat has ended.")
				break
			}
			continue
		}

		history.Add("user", input)
		if err := handleChat(cfg, history); err != nil {
			fmt.Println("Chat error:", err)
		}
	}
}
