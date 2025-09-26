package main

import (
	"fmt"
	"os"
	"picochat/args"
	"picochat/chat"
	"picochat/command"
	"picochat/config"
	"picochat/console"
	"picochat/messages"
	"picochat/version"
)

func sendPrompt(prompt string, cfg *config.Config, history *messages.ChatHistory) {
	stop := make(chan struct{})
	go console.StartSpinner(stop)

	history.Add(messages.RoleUser, prompt)

	msg, err := chat.HandleChat(cfg, history, stop)
	if err != nil {
		console.Error(err)
	} else {
		console.Info(msg)
	}
}

func repeatPrompt(cfg *config.Config, history *messages.ChatHistory) {
	if history.Len() < 2 {
		console.Warn("chat history is empty")
		return
	}

	stop := make(chan struct{})
	go console.StartSpinner(stop)

	lastEntry := history.GetLast()
	if lastEntry.Role != messages.RoleUser {
		console.Warn("last entry in history is not a user prompt")
		return
	}

	msg, err := chat.HandleChat(cfg, history, stop)
	if err != nil {
		console.Error(err)
	} else {
		console.Info(msg)
	}
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

	console.Info("PicoChat started. Help with '/?'")

	for {
		fmt.Print("\n>>> ")

		input, isCommand := console.ReadMultilineInput()

		if input == "" && !isCommand {
			console.Info("\nInput canceled.")
			continue
		}

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
