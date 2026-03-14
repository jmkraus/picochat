package command

import (
	"fmt"
	"io"
	"math"
	"picochat/clipb"
	"picochat/config"
	"picochat/console"
	"picochat/envs"
	"picochat/messages"
	"picochat/paths"
	"picochat/requests"
	"picochat/utils"
	"strings"
)

type CommandResult struct {
	Output string
	Info   string
	Error  error
	Quit   bool
	Pasted string
	Repeat bool
}

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
	cfg, err := config.Get()
	if err != nil {
		return CommandResult{Error: fmt.Errorf("read config failed: %w", err)}
	}

	cmd, args := parseCommandArgs(commandLine)
	switch cmd {
	case "hello":
		hello := "Hello, are you there?"
		return CommandResult{Output: hello, Pasted: hello}
	case "test":
		err := utils.CreateTestFile(cfg.URL)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("save test file failed: %w", err)}
		}
		return CommandResult{Info: "Test file created.", Quit: true}
	case "bye":
		return CommandResult{Info: "Chat has ended.", Quit: true}
	case "save":
		name, err := history.SaveHistoryToFile(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("save history failed: %w", err)}
		}
		return CommandResult{Info: fmt.Sprintf("History saved as %s.", name)}
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
			return CommandResult{Info: fmt.Sprintf("History file %s loaded.", filename)}
		} else {
			return CommandResult{Info: "Load canceled."}
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
		serverVersion, err := requests.GetServerVersion(cfg.URL)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("fetch server version failed: %w", err)}
		}

		list := []string{
			fmt.Sprintf("Configuration file: %s", cfg.ConfigPath),
			fmt.Sprintf("Response output format: %s", cfg.OutputFmt),
			fmt.Sprintf("Current model is '%s'", cfg.Model),
			fmt.Sprintf("Context has %d messages (max. %d)", history.Len(), history.MaxCtx()),
			fmt.Sprintf("Context token estimation: %.0f", math.Ceil(history.EstimateTokens())),
			fmt.Sprintf("Server version is %s", serverVersion),
		}

		return CommandResult{Output: utils.FormatList(list, "Server info", false)}
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
			conversation := history.FullConversation()
			return CommandResult{Output: conversation}
		case messages.RoleAssistant, messages.RoleUser, messages.RoleSystem:
			msg, found := history.GetLastRole(args)
			if found {
				return CommandResult{Output: msg.Content}
			}
			return CommandResult{Info: fmt.Sprintf("No element for role '%s' found.", args)}
		default:
			msg := history.GetLast().Content
			return CommandResult{Output: msg}

		}
	case "copy":
		payload, err := resolveCopyPayload(args, history)
		if err != nil {
			return CommandResult{Error: err}
		}
		if payload.Text == "" {
			return CommandResult{Info: payload.Info}
		}
		if err := clipb.WriteClipboard(payload.Text); err != nil {
			return CommandResult{Error: err}
		}
		return CommandResult{Info: payload.Info}
	case "paste":
		text, err := clipb.ReadClipboard()
		if err != nil {
			return CommandResult{Error: err}
		}
		return CommandResult{
			Info:   fmt.Sprintf("Pasted %d characters from clipboard.", len(text)),
			Pasted: text,
		}
	case "retry":
		history.Discard()
		return CommandResult{Info: "Repeating last chat history user prompt.", Repeat: true}
	case "models":
		if args == "" {
			models, err := utils.ShowAvailableModels(cfg.URL)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("list models failed: %w", err)}
			}
			return CommandResult{Output: models}
		}

		args, _ := strings.CutPrefix(args, "#") // accept and ignore # prefix
		index, err := parseIndex(args)
		if err != nil {
			return CommandResult{Error: err}
		}
		model, ok := utils.GetModelsByIndex(index)
		if !ok {
			return CommandResult{Error: fmt.Errorf("no value for given index found")}
		}
		cfg.Model = model
		return CommandResult{Info: fmt.Sprintf("Switched model to '%s'.", model)}
	case "envs":
		envSetup := envs.ConfigEnvVarsMarkdownTable()
		return CommandResult{Output: envSetup}
	case "set":
		if args == "" {
			list := []string{
				fmt.Sprintf("model = %s", cfg.Model),
				fmt.Sprintf("context = %d", cfg.Context),
				fmt.Sprintf("temperature = %.2f", cfg.Temperature),
				fmt.Sprintf("top_p = %.2f", cfg.Top_p),
				fmt.Sprintf("reasoning = %t", cfg.Reasoning),
			}
			return CommandResult{Output: utils.FormatList(list, "Config settings", false)}
		}

		key, value, err := parseKeyVal(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("parse args failed: %w", err)}
		}
		err = config.Set(key, value)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("apply to config failed: %w", err)}
		}
		if key == "context" {
			intVal, _ := value.(int)
			err := history.SetContextSize(intVal)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("set context size failed: %w", err)}
			}
		}
		return CommandResult{Info: fmt.Sprintf("Config updated: %s = %v.", key, value)}
	case "clear":
		history.ClearExceptSystemPrompt()
		return CommandResult{Info: "History cleared (system prompt retained)."}
	case "help":
		return CommandResult{Output: HelpText(args)}
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
//	cmd  - the normalized command string without leading slash.
//	arg  - the remaining arguments as a single string.
func parseCommandArgs(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}
	cmd := strings.TrimSpace(parts[0])

	// replace special abbreviations
	switch cmd {
	case "/c":
		cmd = "/copy"
	case "/v":
		cmd = "/paste"
	case "/?":
		cmd = "/help"
	}

	// normalize
	cmd = strings.ToLower(strings.TrimPrefix(cmd, "/"))

	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(strings.Join(parts[1:], " "))
	}
	return cmd, arg
}

func getMessageByIndex(args string, history *messages.ChatHistory) (string, error) {
	index, err := parseIndex(args)
	if err != nil {
		return "", fmt.Errorf("get message failed: %w", err)
	}
	msg, err := history.GetByIndex(index)
	if err != nil {
		return "", fmt.Errorf("get message failed: %w", err)
	}
	headerText := fmt.Sprintf("(%d:%s)", index, msg.Role)
	header := console.Style(console.Bold, headerText)
	return fmt.Sprintf("%s\n%s", header, msg.Content), nil
}
