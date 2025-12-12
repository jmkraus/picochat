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
	"picochat/paths"
	"picochat/version"
)

func sendPrompt(prompt string, image string, quiet bool, history *messages.ChatHistory) {
	stop := make(chan struct{})
	go console.StartSpinner(quiet, stop)

	err := history.Add(messages.RoleUser, prompt, image)
	if err != nil {
		console.StopSpinner(quiet, stop)
		console.Error(fmt.Sprintf("%v", err))
		return
	}

	msg, err := chat.HandleChat(nil, history, stop)
	if err != nil {
		console.StopSpinner(quiet, stop)
		console.Error(fmt.Sprintf("%v", err))
	} else {
		if !quiet {
			console.Info(msg)
		}
	}
}

func repeatPrompt(quiet bool, history *messages.ChatHistory) {
	if history.Len() < 2 {
		console.Warn("chat history is empty")
		return
	}

	stop := make(chan struct{})
	go console.StartSpinner(quiet, stop)

	lastEntry := history.GetLast()
	if lastEntry.Role != messages.RoleUser {
		console.Warn("last entry in history is not a user prompt")
		return
	}

	msg, err := chat.HandleChat(nil, history, stop)
	if err != nil {
		console.StopSpinner(quiet, stop)
		console.Error(fmt.Sprintf("%v", err))
	} else {
		if !quiet {
			console.Info(msg)
		}
	}
}

func main() {
	args.Parse()

	if *args.ShowVersion {
		console.Info(fmt.Sprintf("picochat version is %s", version.Version))
		os.Exit(0)
	}

	cfg, err := config.Get()
	if err != nil {
		console.Error(fmt.Sprintf("load configuration failed: %v", err))
		os.Exit(1)
	}

	if *args.Quiet {
		// only override config if arg actively set
		cfg.Quiet = true
	}

	if *args.Model != "" {
		cfg.Model = *args.Model
	}

	if *args.Image != "" {
		cfg.ImagePath = *args.Image
		if !paths.FileExists(cfg.ImagePath) {
			console.Warn("image file not found")
		}
	}

	var history *messages.ChatHistory
	if *args.HistoryFile != "" {
		history, err = messages.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			console.Error(fmt.Sprintf("load history failed: %v", err))
			os.Exit(1)
		}
	} else {
		history = messages.NewHistory(cfg.Prompt, cfg.Context)
	}

	if !cfg.Quiet {
		console.Info(fmt.Sprintf("Configuration file used: %s", cfg.FilePath))
		if *args.Model != "" {
			console.Info(fmt.Sprintf("Configuration overridden by model='%s'", *args.Model))
		}
		console.Info("PicoChat started. Help with '/?'")
	}

	for {
		if !cfg.Quiet {
			fmt.Printf("\n%s", console.Prompt)
		}

		input := console.ReadMultilineInput()
		if input.Error != nil {
			console.Error(fmt.Sprintf("%v", input.Error))
		}

		if input.Aborted {
			console.Info("\nInput canceled.")
			continue
		}

		if input.Text == "" && !input.IsCommand {
			continue
		}

		if input.IsCommand {
			fmt.Println()
			result := command.HandleCommand(input.Text, history, os.Stdin)
			if result.Error != nil {
				console.Error(fmt.Sprintf("command handler failed: %v", result.Error))
			}
			console.AddCommand(input.Text)
			if result.Output != "" {
				console.Info(result.Output)
			}
			if result.Quit {
				break
			}
			if result.Repeat {
				repeatPrompt(cfg.Quiet, history)
			} else if result.Pasted != "" {
				// start the request with pasted content from clipboard
				sendPrompt(result.Pasted, cfg.ImagePath, cfg.Quiet, history)
			}
			if input.EOF {
				// we come from stdin pipe
				break
			} else {
				continue
			}
		}

		sendPrompt(input.Text, cfg.ImagePath, cfg.Quiet, history)

		if input.EOF {
			break
		}

	}
}
