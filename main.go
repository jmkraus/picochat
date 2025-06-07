package main

import (
	"bufio"
	"fmt"
	"os"
	"picochat/args"
	"picochat/chat"
	"picochat/command"
	"picochat/config"
	"picochat/types"
	"picochat/version"
	"strings"
)

func sendUserPrompt(prompt string, cfg *config.Config, history *types.ChatHistory) {
	history.Add("user", prompt)
	if err := chat.HandleChat(cfg, history); err != nil {
		fmt.Fprintf(os.Stderr, "Chat error: %v", err)
	}
}

func readMultilineInput() (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if firstLine && strings.HasPrefix(trimmed, "/") {
			return trimmed, true // input, isCommand
		}

		if trimmed == "/done" || trimmed == "///" {
			break
		}

		lines = append(lines, line)
		firstLine = false
	}

	return strings.Join(lines, "\n"), false
}

func main() {
	args.Parse()

	if *args.ShowVersion {
		fmt.Println("picochat version is", version.Version)
		os.Exit(0)
	}

	err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while loading configuration: %v", err)
		os.Exit(1)
	}
	cfg := config.Get()

	var history *types.ChatHistory
	if *args.HistoryFile != "" {
		history, err = types.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not load history: %v", err)
			os.Exit(1)
		}
	} else {
		history = types.NewHistory(cfg.Prompt, cfg.Context)
	}

	fmt.Println("Chat with Pico AI started. Help with '/?'.")

	for {
		fmt.Print("\n>>> ")

		input, isCommand := readMultilineInput()

		if isCommand {
			result := command.Handle(input, history, os.Stdin)
			if result.Output != "" {
				fmt.Println(result.Output)
			}
			if result.Quit {
				break
			}
			if result.Prompt != "" {
				sendUserPrompt(result.Prompt, cfg, history)
			}
			continue
		}

		sendUserPrompt(input, cfg, history)
	}
}
