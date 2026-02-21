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
	"picochat/utils"
	"picochat/version"
)

type Session struct {
	Config  *config.Config
	History *messages.ChatHistory
	Quiet   bool
}

func sendPrompt(session *Session, prompt string) {
	if err := session.History.AddUser(prompt, session.Config.ImagePath); err != nil {
		console.Error(fmt.Sprintf("%v", err))
		return
	}

	session.Config.ImagePath = "" // store once in history and forget
	runChat(session)
}

func repeatPrompt(session *Session) {
	if session.History.Len() < 2 {
		console.Warn("chat history is empty")
		return
	}

	lastEntry := session.History.GetLast()
	if lastEntry.Role != messages.RoleUser {
		console.Warn("last entry in history is not a user prompt")
		return
	}

	runChat(session)
}

func runChat(session *Session) {
	stop := make(chan struct{})
	go console.StartSpinner(session.Quiet, stop)
	defer console.StopSpinner(session.Quiet, stop)

	result, err := chat.HandleChat(session.Config, session.History, stop)
	if err != nil {
		console.Error(fmt.Sprintf("%v", err))
		return
	}

	if err := format.RenderResult(
		os.Stdout,
		result,
		session.Config.OutputFmt,
		session.Quiet,
	); err != nil {
		console.Error(fmt.Sprintf("output failed: %v", err))
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

	if *args.Output == "" {
		cfg.OutputFmt = "plain" // default
	} else {
		f, ok := format.AllowedKeys(*args.Output)
		cfg.OutputFmt = f
		if !ok {
			cfg.OutputFmt = "plain"
			console.Warn("unknown output format - fallback to plain")
		}
	}

	if *args.Format != "" {
		cfg.OutputFmt = "plain" // There can be only one
		schema, err := utils.LoadSchemaFromFile(*args.Format)
		if err != nil {
			console.Error(fmt.Sprintf("load json schema file failed: %v", err))
			os.Exit(1)
		}
		cfg.SchemaFmt = schema
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
			console.Info(fmt.Sprintf("Configuration overridden by model='%s'", cfg.Model))
		}
		console.Info("PicoChat started.")
	}

	for {
		if !session.Quiet {
			fmt.Printf("\n%s", console.Prompt+console.Shadow)
			fmt.Printf(console.CursorToColumn, console.PromptWidth()+1)
		}

		input := console.ReadMultilineInput()
		if input.Error != nil {
			console.Error(input.Error.Error())
			continue
		}

		if input.Aborted {
			fmt.Println()
			console.Info("Input canceled.")
			continue
		}

		if input.Text == "" && !input.IsCommand {
			if input.EOF {
				// we come from stdin pipe
				console.NewLine(session.Quiet)
				break
			} else {
				continue
			}
		}

		if input.IsCommand {
			console.NewLine(session.Quiet)
			result := command.HandleCommand(input.Text, history, os.Stdin)
			if result.Error != nil {
				console.Error(fmt.Sprintf("command handler failed: %v", result.Error))
			}
			console.AddCommand(input.Text)
			if result.Info != "" && !session.Quiet {
				console.Info(result.Info)
			}
			if result.Output != "" {
				fmt.Println(result.Output)
			}

			if result.Quit {
				break
			}
			if result.Repeat {
				repeatPrompt(session)
			} else if result.Pasted != "" {
				// start the request with pasted content from clipboard
				sendPrompt(session, result.Pasted)
			}

			if input.EOF {
				// we come from stdin pipe
				break
			} else {
				continue
			}
		}

		sendPrompt(session, input.Text)

		if input.EOF {
			break
		}

	}
}
