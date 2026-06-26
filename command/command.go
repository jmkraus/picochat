package command

import (
	"fmt"
	"io"
	"math"
	"path/filepath"
	"picochat/backend"
	"picochat/clipb"
	"picochat/config"
	"picochat/console"
	"picochat/envs"
	"picochat/messages"
	"picochat/output"
	"picochat/paths"
	"picochat/utils"
	"strings"
	"unicode/utf8"
)

type CommandResult struct {
	Output string
	Info   string
	Warn   string
	Error  error
	Quit   bool
	Pasted string
	Retry  bool
}

var readClipboard = clipb.ReadClipboard

// HandleCommand processes a command line input, performs the requested action,
// and returns a CommandResult containing the outcome of the command.
//
// Parameters:
//
//	commandLine - the raw command line string entered by the user.
//	history     - the chat history to operate on.
//	input       - io.Reader (default: os.Stdin) used for unit tests
//
// Returns:
//
//	CommandResult - a struct containing output, error, quit flag, prompt,
//	and repeat flag for the command.
func HandleCommand(commandLine string, history *messages.ChatHistory, input io.Reader) CommandResult {
	cfg, _, err := config.Get()
	if err != nil {
		return CommandResult{Error: fmt.Errorf("read config failed: %w", err)}
	}

	cmd, args := parseCommandArgs(commandLine)
	switch cmd {
	case "hello":
		hello := "Hello, are you there?"
		return CommandResult{Output: hello, Pasted: hello}
	case "test":
		models, err := backend.New(cfg).GetAvailableModels()
		if err != nil {
			return CommandResult{Error: fmt.Errorf("get models failed: %w", err)}
		}

		err = utils.CreateTestFile(models)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("save test file failed: %w", err)}
		}
		return CommandResult{Info: "Test file created.", Quit: true}
	case "bye":
		return CommandResult{Info: "Chat has ended.", Quit: true}
	case "save":
		overwrite := false
		if args != "" {
			historyPath, err := paths.GetHistoryPath()
			if err != nil {
				return CommandResult{Error: fmt.Errorf("history path not found: %w", err)}
			}
			targetName := paths.EnsureSuffix(filepath.Base(args), paths.HistorySuffix)
			targetPath := filepath.Join(historyPath, targetName)
			if paths.FileExists(targetPath) {
				overwrite, err = askConfirmation(fmt.Sprintf("File %q already exists. Overwrite?", targetName), input)
				if err != nil {
					return CommandResult{Error: fmt.Errorf("overwrite confirmation failed: %w", err)}
				}
				if !overwrite {
					return CommandResult{Warn: "Save canceled."}
				}
			}
		}

		filename, err := messages.SaveHistoryToFile(args, history.Get(), overwrite)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("save history failed: %w", err)}
		}
		return CommandResult{Info: fmt.Sprintf("History saved as file %q.", filename)}
	case "load":
		if args == "" {
			files, err := utils.ListHistoryFiles()
			if err != nil {
				return CommandResult{Error: fmt.Errorf("list history files failed: %w", err)}
			}
			fmt.Println(files)
		}

		filename, err := getHistoryFilename(args, input)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("get history filename failed: %w", err)}
		}

		if filename != "" {
			filename = paths.EnsureSuffix(filename, paths.HistorySuffix)
			loaded, err := messages.LoadHistoryFromFile(filename)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("load history failed: %w", err)}
			}
			history.Replace(loaded.Get())
			return CommandResult{Info: fmt.Sprintf("History file %q loaded.", filename)}
		} else {
			return CommandResult{Warn: "Load canceled."}
		}
	case "image":
		if args != "" {
			if !paths.FileExists(args) {
				return CommandResult{Error: fmt.Errorf("image file not found")}
			}
			cfg.ImagePath = args
			return CommandResult{Info: fmt.Sprintf("Image file path set to: %s", cfg.ImagePath)}
		}
		return CommandResult{Error: fmt.Errorf("no image file path provided")}
	case "info":
		serverVersion, err := backend.New(cfg).GetServerVersion()
		if err != nil {
			serverVersion = fmt.Sprintf("%s connection error", console.ErrPrefix)
		}

		list := []string{
			fmt.Sprintf("Configuration file: %s", cfg.ConfigPath),
			fmt.Sprintf("Backend API used: %s", cfg.Backend),
			fmt.Sprintf("Response output format: %s", cfg.OutputFmt),
			fmt.Sprintf("Current model is %q", cfg.Model),
			fmt.Sprintf("Context has %d messages (max. %d)", history.Len(), history.MaxCtx()),
			fmt.Sprintf("Context token estimation: %.0f", math.Ceil(history.EstimateTokens())),
			fmt.Sprintf("Server version: %s", serverVersion),
		}

		return CommandResult{Output: utils.FormatList(list, "Server info", false)}
	case "trim":
		args := strings.TrimPrefix(args, "#") // accept and ignore # prefix
		if args == "" {
			return CommandResult{Error: fmt.Errorf("missing index argument")}
		}
		index, err := parseIndex(args)
		if err != nil {
			return CommandResult{Error: err}
		}
		ok := history.Trim(index)
		if ok {
			return CommandResult{Info: "Chat history has been truncated."}
		}
		return CommandResult{Warn: "Chat history not changed."}
	case "message":
		ok := false
		if args, ok = strings.CutPrefix(args, "#"); ok {
			msg, err := getMessageByIndex(args, history)
			if err != nil {
				return CommandResult{Error: err}
			}
			return CommandResult{Output: msg}
		}

		switch args {
		case "all":
			conversation := output.FormatConversation(history.Get())
			return CommandResult{Output: conversation}
		case messages.RoleAssistant, messages.RoleUser, messages.RoleSystem:
			msg, found := history.GetLastRole(args)
			if found {
				return CommandResult{Output: msg.Content}
			}
			return CommandResult{Warn: fmt.Sprintf("No element for role %q found.", args)}
		default:
			msg := history.GetLast().Content
			return CommandResult{Output: msg}

		}
	case "copy":
		payload, err := resolveCopyPayload(args, history)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("copy message data failed: %w", err)}
		}
		if payload.Text == "" {
			return CommandResult{Warn: payload.Info}
		}
		if err := clipb.WriteClipboard(payload.Text); err != nil {
			return CommandResult{Error: err}
		}
		return CommandResult{Info: payload.Info}
	case "paste":
		tpl, err := config.GetTemplate(args)
		if err != nil {
			return CommandResult{Error: err}
		}
		clip, err := readClipboard()
		if err != nil {
			return CommandResult{Error: err}
		}
		text := clip
		if tpl != "" {
			text = fmt.Sprintf("%s\n\n%s", tpl, clip)
		}
		count := utf8.RuneCountInString(clip)
		return CommandResult{
			Info:   fmt.Sprintf("Pasted %d characters from clipboard.", count),
			Pasted: text,
		}
	case "retry":
		history.Discard()
		if history.IsEmpty() {
			return CommandResult{Error: fmt.Errorf("chat history is empty")}
		}
		if !history.CheckIfLastEntryIsRole(messages.RoleUser) {
			return CommandResult{Error: fmt.Errorf("last entry in chat history is not a user prompt")}
		}
		return CommandResult{Info: "Repeating last chat history user prompt.", Retry: true}
	case "models":
		if args == "" {
			models, err := backend.New(cfg).GetAvailableModels()
			if err != nil {
				return CommandResult{Error: fmt.Errorf("get models failed: %w", err)}
			}

			list, err := utils.ListAvailableModels(models)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("list models failed: %w", err)}
			}
			return CommandResult{Output: list}
		}

		args := strings.TrimPrefix(args, "#") // accept and ignore # prefix
		index, err := parseIndex(args)
		if err != nil {
			return CommandResult{Error: err}
		}
		model, ok := utils.GetModelsByIndex(index)
		if !ok {
			return CommandResult{Error: fmt.Errorf("no value for given index found")}
		}
		cfg.Model = model
		return CommandResult{Info: fmt.Sprintf("Switched model to %q.", model)}
	case "set":
		if args == "" {
			list := []string{
				fmt.Sprintf("context = %d", cfg.Context),
				fmt.Sprintf("temperature = %.2f", cfg.Temperature),
				fmt.Sprintf("top_p = %.2f", cfg.Top_p),
				fmt.Sprintf("reasoning = %t", cfg.Reasoning),
				fmt.Sprintf("effort = %s", cfg.Effort),
			}
			return CommandResult{Output: utils.FormatList(list, "Config settings", false)}
		}

		key, value, err := parseKeyVal(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("parse args failed: %w", err)}
		}

		var warnings []string
		warnings, err = config.Set(key, value)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("apply to config failed: %w", err)}
		}

		var warn string
		if len(warnings) > 0 {
			if len(warnings) == 1 {
				warn = warnings[0]
			} else {
				warn = fmt.Sprintf("%s (+%d more)", warnings[0], len(warnings)-1)
			}
		}

		if key == "context" {
			err := history.SetContextSize(cfg.Context)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("set context size failed: %w", err)}
			}
		}
		return CommandResult{Info: fmt.Sprintf("Config updated for %s.", key), Warn: warn}
	case "clear":
		history.ClearExceptSystemPrompt()
		return CommandResult{Info: "History cleared (system prompt retained)."}
	case "help":
		switch args {
		case "env", "envs":
			return CommandResult{Output: envs.ListEnvVars()}
		case "tpl", "templates":
			return CommandResult{Output: config.ListTemplates()}
		default:
			return CommandResult{Output: HelpText(args)}
		}
	default:
		return CommandResult{Error: fmt.Errorf("unknown command")}
	}
}

// parseCommandArgs splits the input string into a command and its arguments,
// normalizes the command, and handles special abbreviations.
//
// Parameters:
//
//	input - the raw input string entered by the user.
//
// Returns:
//
//	string - the normalized command string without leading slash.
//	string - the remaining arguments as a single string.
func parseCommandArgs(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}
	cmd := strings.TrimSpace(parts[0])

	// normalize
	cmd = strings.ToLower(strings.TrimPrefix(cmd, "/"))

	// command replacements
	switch cmd {
	case "c":
		cmd = "copy"
	case "v":
		cmd = "paste"
	case "?":
		cmd = "help"
	case "exit", "quit":
		cmd = "bye"
	case "hallo":
		cmd = "hello" // just in case...
	}

	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(strings.Join(parts[1:], " "))
	}
	return cmd, arg
}
