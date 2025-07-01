package main

import (
	"bufio"
	"fmt"
	"os"
	"picochat/args"
	"picochat/chat"
	"picochat/command"
	"picochat/config"
	"picochat/console"
	"picochat/messages"
	"picochat/version"
	"strings"
)

func sendPrompt(prompt string, cfg *config.Config, history *messages.ChatHistory) {
	history.Add(messages.RoleUser, prompt)
	if err := chat.HandleChat(cfg, history); err != nil {
		console.Error(err)
	}
}

func repeatPrompt(cfg *config.Config, history *messages.ChatHistory) {
	if history.Len() < 2 {
		console.Warn("chat history is empty.")
		return
	}

	lastEntry := history.GetLast()
	if lastEntry.Role != messages.RoleUser {
		console.Warn("last entry in history is not a user prompt.")
		return
	}

	if err := chat.HandleChat(cfg, history); err != nil {
		console.Error(err)
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
		console.Info(fmt.Sprintf("picochat version is %s", version.Version))
		os.Exit(0)
	}

	err := config.Load()
	if err != nil {
		console.Errorf("load configuration failed: %v", err)
		os.Exit(1)
	}
	cfg := config.Get()

	var history *messages.ChatHistory
	if *args.HistoryFile != "" {
		history, err = messages.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			console.Errorf("load history failed: %v", err)
			os.Exit(1)
		}
	} else {
		history = messages.NewHistory(cfg.Prompt, cfg.Context)
	}

	console.Info("Chat with Pico AI started. Help with '/?'")

	for {
		fmt.Print("\n>>> ")

		input, isCommand := readMultilineInput()

		if isCommand {
			result := command.HandleCommand(input, history, os.Stdin)
			if result.Output != "" {
				console.Info(result.Output)
			}
			if result.Error != nil {
				console.Errorf("command handler failed: %v", result.Error)
			}
			if result.Quit {
				break
			}
			if result.Repeat {
				repeatPrompt(cfg, history)
			} else if result.Prompt != "" {
				sendPrompt(result.Prompt, cfg, history)
			}
			continue
		}

		sendPrompt(input, cfg, history)
	}
}
