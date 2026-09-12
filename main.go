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

// Central data for running instance of PicoChat
type Instance struct {
	Config   *config.Config
	History  *messages.ChatHistory
	Sessions *messages.SessionManager
	Quiet    bool
}

const (
	defaultOutputFormat = "plain"
)

// sendPrompt appends a user message to history and starts a chat run.
//
// Parameters:
//
//	instance (*Instance) - active runtime instance
//	prompt  (string)     - user input prompt
//
// Returns:
//
//	none
func sendPrompt(instance *Instance, prompt string) {
	if err := instance.History.AddUser(prompt, instance.Config.ImagePath); err != nil {
		console.Error(err)
		return
	}

	instance.Config.ImagePath = "" // store once in history and forget
	runChat(instance)
}

// retryPrompt triggers a new chat run based on existing history.
//
// Parameters:
//
//	instance (*Instance) - active runtime instance
//
// Returns:
//
//	none
func retryPrompt(instance *Instance) {
	runChat(instance)
}

// getActiveHistory retrieves the pointer to the active history from the
// session manager and stores it in the instance.
//
// Parameters:
//
//	instance (*Instance) - active runtime instance
//
// Returns:
//
//	error - error if the session manager is unavailable
func getActiveHistory(instance *Instance) error {
	if instance.Sessions == nil {
		return fmt.Errorf("session manager is unavailable")
	}

	history, err := instance.Sessions.Active()
	if err != nil {
		return err
	}
	instance.History = history
	return nil
}

// runChat sends the prepared chat request and renders the final result.
//
// Parameters:
//
//	instance (*Instance) - active runtime instance
//
// Returns:
//
//	none
func runChat(instance *Instance) {
	stop := make(chan struct{})
	go console.StartSpinner(instance.Quiet, stop)
	defer console.StopSpinner(instance.Quiet, stop)

	result, err := chat.HandleChat(instance.Config, instance.History, stop)
	if err != nil {
		console.Error(err)
		return
	}

	if err := output.RenderResult(
		os.Stdout,
		result,
		instance.Config.OutputFmt,
		instance.Quiet,
	); err != nil {
		console.Error(fmt.Errorf("output failed: %w", err))
	}
}

// initInstanceFromArgs parses CLI args, loads config, applies overrides,
// initializes history, and returns a prepared instance.
//
// Parameters:
//
//	none
//
// Returns:
//
//	bool      - true if caller should print version and exit
//	*Instance - prepared runtime instance
//	[]string  - startup warnings if any
//	error     - error if startup initialization fails
func initInstanceFromArgs() (bool, *Instance, []string, error) {
	args.Parse()

	if *args.ShowVersion {
		return true, nil, nil, nil
	}

	config.Init(*args.ConfigPath)
	cfg, warn, err := config.Get()
	if err != nil {
		return false, nil, nil, fmt.Errorf("load configuration failed: %w", err)
	}

	if *args.Quiet {
		// only override config if arg actively set
		cfg.Quiet = true
	}

	if *args.Output == "" {
		cfg.OutputFmt = defaultOutputFormat
	} else {
		f, ok := output.AllowedKeys(*args.Output)
		cfg.OutputFmt = f
		if !ok {
			cfg.OutputFmt = defaultOutputFormat
			warn = append(warn, fmt.Sprintf("unknown output format - fallback to %s", defaultOutputFormat))
		}
	}

	if *args.Schema != "" {
		cfg.OutputFmt = defaultOutputFormat // Schema overrules output format
		schema, err := utils.LoadSchemaFromFile(*args.Schema)
		if err != nil {
			return false, nil, nil, fmt.Errorf("load json schema file failed: %w", err)
		}
		cfg.SchemaFmt = schema
	}

	if *args.Model != "" {
		cfg.Model = *args.Model
	}

	if *args.Image != "" {
		cfg.ImagePath = *args.Image
		if !paths.FileExists(cfg.ImagePath) {
			return false, nil, nil, fmt.Errorf("image file not found")
		}
	}

	var history *messages.ChatHistory
	if *args.HistoryFile != "" {
		history, err = messages.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			return false, nil, nil, fmt.Errorf("load history failed: %w", err)
		}
		ctxSize := max(cfg.Context, history.Len()+1)
		ctxSize = min(ctxSize, config.MaxContext)
		if err := history.SetContextSize(ctxSize); err != nil {
			return false, nil, nil, fmt.Errorf("set context size failed: %w", err)
		}
	} else {
		history = messages.NewHistory(cfg.Prompt, cfg.Context)
	}

	manager := messages.NewSessionManager()
	_, err = manager.Add(history)
	if err != nil {
		return false, nil, nil, fmt.Errorf("init session manager failed: %w", err)
	}

	instance := &Instance{
		Config:   cfg,
		History:  history,
		Sessions: manager,
		Quiet:    cfg.Quiet,
	}

	return false, instance, warn, nil
}

func main() {
	showVersion, instance, warn, err := initInstanceFromArgs()
	if showVersion {
		fmt.Printf("picochat version is %s\n", version.Version)
		os.Exit(0)
	}
	if err != nil {
		console.Error(err)
		os.Exit(1)
	}

	printNewLine := func() {
		if !instance.Quiet {
			fmt.Println()
		}
	}

	if !instance.Quiet {
		console.Warns(warn)
		if *args.Model != "" {
			console.Info(fmt.Sprintf("Using model from CLI override: %q.", instance.Config.Model))
		}
		console.Info("PicoChat started.")
	}

	for {
		printNewLine()
		prompt := console.NumberedPrompt(instance.Sessions.ActiveIndex())
		if !instance.Quiet {
			fmt.Print(prompt + console.ShadowText)
			console.SetCursorPos(console.PromptWidth() + 1)
		}

		input := console.ReadMultilineInputWithPrompt(prompt)
		if input.Error != nil {
			console.Error(input.Error)
			continue
		}

		if input.Aborted {
			printNewLine()
			if !instance.Quiet {
				console.Warn("Input canceled.")
			}
			continue
		}

		if input.Text == "" && !input.IsCommand {
			if input.EOF {
				// we come from stdin pipe
				printNewLine()
				break
			} else {
				continue
			}
		}

		if input.IsCommand {
			fmt.Println() // newline even in quiet mode
			result := command.HandleCommand(input.Text, instance.History, instance.Sessions, os.Stdin)
			console.AddCommand(input.Text)
			if result.Error != nil {
				console.Error(fmt.Errorf("command handler error: %w", result.Error))
				continue
			}
			if result.SessionChanged {
				if err := getActiveHistory(instance); err != nil {
					console.Error(fmt.Errorf("get active session failed: %w", err))
					continue
				}
			}
			if !instance.Quiet {
				console.Warn(result.Warn)
				console.Info(result.Info)
			}
			if result.Output != "" {
				fmt.Println(result.Output)
			}
			if result.Quit {
				break
			}
			if result.Retry {
				retryPrompt(instance)
			} else if result.Pasted != "" {
				// start the request with pasted content from clipboard
				sendPrompt(instance, result.Pasted)
			}

			if input.EOF {
				// we come from stdin pipe
				break
			} else {
				continue
			}
		}

		sendPrompt(instance, input.Text)

		if input.EOF {
			break
		}
	}
}
