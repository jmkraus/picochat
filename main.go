package main

import (
	"fmt"
	"os"
	"picochat/args"
	"picochat/chat"
	"picochat/command"
	"picochat/config"
	"picochat/console"
	"picochat/format"
	"picochat/messages"
	"picochat/paths"
	"picochat/version"
)

type Session struct {
	Config  *config.Config
	History *messages.ChatHistory
	Quiet   bool
}

func sendPrompt(session *Session, prompt string, image string) {
	stop := make(chan struct{})
	go console.StartSpinner(session.Quiet, stop)

	err := session.History.Add(messages.RoleUser, prompt, image)
	if err != nil {
		console.StopSpinner(session.Quiet, stop)
		console.Error(fmt.Sprintf("%v", err))
		return
	}

	result, err := chat.HandleChat(session.Config, session.History, stop)
	if err != nil {
		console.StopSpinner(session.Quiet, stop)
		console.Error(fmt.Sprintf("%v", err))
		return
	}

	outputResult(result, session)
}

func repeatPrompt(session *Session) {
	if session.History.Len() < 2 {
		console.Warn("chat history is empty")
		return
	}

	stop := make(chan struct{})
	go console.StartSpinner(session.Quiet, stop)

	lastEntry := session.History.GetLast()
	if lastEntry.Role != messages.RoleUser {
		console.Warn("last entry in history is not a user prompt")
		return
	}

	result, err := chat.HandleChat(session.Config, session.History, stop)
	if err != nil {
		console.StopSpinner(session.Quiet, stop)
		console.Error(fmt.Sprintf("%v", err))
		return
	}

	outputResult(result, session)
}

func outputResult(result *chat.ChatResult, session *Session) {
	if session.Quiet {
		return
	}

	switch session.Config.OutputFmt {
	case "json":
		// if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		// 	console.Error(fmt.Sprintf("json encode failed: %v", err))
		// }
	case "yaml":
		// if err := yaml.NewEncoder(os.Stdout).Encode(result); err != nil {
		// 	console.Error(fmt.Sprintf("yaml encode failed: %v", err))
		// }
	default:
		msg := fmt.Sprintf(
			"\nElapsed (mm:ss): %s | Tokens/sec: %.1f",
			result.Elapsed,
			result.TokensPS,
		)
		console.Info(msg)
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

	if *args.Format == "" {
		// default
		cfg.OutputFmt = "plain"
	} else {
		format, ok := format.AllowedKeys(*args.Format)
		if ok {
			cfg.OutputFmt = format
		} else {
			cfg.OutputFmt = "plain"
			console.Warn("unknown output format - fallback to plain")
		}
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

	session := &Session{
		Config:  cfg,
		History: history,
		Quiet:   cfg.Quiet,
	}

	if !session.Quiet {
		if *args.Model != "" {
			console.Info(fmt.Sprintf("Configuration overridden by model='%s'", *args.Model))
		}
		console.Info("PicoChat started. Help with '/?'")
	}

	for {
		if !session.Quiet {
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
				repeatPrompt(session)
			} else if result.Pasted != "" {
				// start the request with pasted content from clipboard
				sendPrompt(session, result.Pasted, cfg.ImagePath)
			}
			if input.EOF {
				// we come from stdin pipe
				break
			} else {
				continue
			}
		}

		sendPrompt(session, input.Text, cfg.ImagePath)

		if input.EOF {
			break
		}

	}
}
