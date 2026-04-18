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
	"picochat/output"
	"picochat/paths"
	"picochat/utils"
	"picochat/version"
)

type Session struct {
	Config  *config.Config
	History *messages.ChatHistory
	Quiet   bool
}

// sendPrompt appends a user message to history and starts a chat run.
//
// Parameters:
//
//	session (*Session) - active runtime session
//	prompt  (string)   - user input prompt
//
// Returns:
//
//	none
func sendPrompt(session *Session, prompt string) {
	if err := session.History.AddUser(prompt, session.Config.ImagePath); err != nil {
		console.Error(err.Error())
		return
	}

	session.Config.ImagePath = "" // store once in history and forget
	runChat(session)
}

// repeatPrompt triggers a new chat run based on existing history.
//
// Parameters:
//
//	session (*Session) - active runtime session
//
// Returns:
//
//	none
func repeatPrompt(session *Session) {
	runChat(session)
}

// runChat sends the prepared chat request and renders the final result.
//
// Parameters:
//
//	session (*Session) - active runtime session
//
// Returns:
//
//	none
func runChat(session *Session) {
	stop := make(chan struct{})
	go console.StartSpinner(session.Quiet, stop)
	defer console.StopSpinner(session.Quiet, stop)

	result, err := chat.HandleChat(session.Config, session.History, stop)
	if err != nil {
		console.Error(err.Error())
		return
	}

	if err := output.RenderResult(
		os.Stdout,
		result,
		session.Config.OutputFmt,
		session.Quiet,
	); err != nil {
		console.Error(fmt.Sprintf("output failed: %v", err))
	}
}

// initSessionFromArgs parses CLI args, loads config, applies overrides,
// initializes history, and returns a prepared session.
//
// Parameters:
//
//	none
//
// Returns:
//
//	*Session - prepared runtime session
//	[]string - startup warnings (if any)
//	bool     - true if caller should print version and exit
//	error    - error if startup initialization fails
func initSessionFromArgs() (*Session, []string, bool, error) {
	args.Parse()

	if *args.ShowVersion {
		return nil, nil, true, nil
	}

	config.Init(*args.ConfigPath)
	cfg, warn, err := config.Get()
	if err != nil {
		return nil, nil, false, fmt.Errorf("load configuration failed: %w", err)
	}

	if *args.Quiet {
		// only override config if arg actively set
		cfg.Quiet = true
	}

	if *args.Output == "" {
		cfg.OutputFmt = "plain" // default
	} else {
		f, ok := output.AllowedKeys(*args.Output)
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
			return nil, nil, false, fmt.Errorf("load json schema file failed: %w", err)
		}
		cfg.SchemaFmt = schema
	}

	if *args.Model != "" {
		cfg.Model = *args.Model
	}

	if *args.Image != "" {
		cfg.ImagePath = *args.Image
		if !paths.FileExists(cfg.ImagePath) {
			return nil, nil, false, fmt.Errorf("image file not found")
		}
	}

	var history *messages.ChatHistory
	if *args.HistoryFile != "" {
		history, err = messages.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			return nil, nil, false, fmt.Errorf("load history failed: %w", err)
		}
	} else {
		history = messages.NewHistory(cfg.Prompt, cfg.Context)
	}

	session := &Session{
		Config:  cfg,
		History: history,
		Quiet:   cfg.Quiet,
	}

	return session, warn, false, nil
}

func main() {
	session, warn, showVersion, err := initSessionFromArgs()
	if showVersion {
		fmt.Printf("picochat version is %s", version.Version)
		fmt.Println()
		os.Exit(0)
	}
	if err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}

	if len(warn) > 0 {
		console.Warns(warn)
	}

	if !session.Quiet {
		if *args.Model != "" {
			console.Info(fmt.Sprintf("Using model from CLI override: %s.", session.Config.Model))
		}
		console.Info("PicoChat started.")
	}

	for {
		if !session.Quiet {
			fmt.Println()
			fmt.Print(console.Prompt + console.Shadow)
			console.SetCursorPos(console.PromptWidth() + 1)
		}

		input := console.ReadMultilineInput()
		if input.Error != nil {
			console.Error(input.Error.Error())
			continue
		}

		if input.Aborted {
			if !session.Quiet {
				console.NewLine(false)
				console.Warn("Input canceled.")
			}
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
			result := command.HandleCommand(input.Text, session.History, os.Stdin)
			if result.Error != nil {
				console.Error(fmt.Sprintf("command handler failed: %v", result.Error))
				continue
			}
			if result.Warn != "" {
				console.Warn(result.Warn)
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
