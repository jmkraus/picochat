package command

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"picochat/clipb"
	"picochat/config"
	"picochat/messages"
	"picochat/paths"
	"picochat/requests"
	"picochat/utils"
	"strconv"
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
		return CommandResult{Error: fmt.Errorf("config read failed: %w", err)}
	}

	cmd, args := parseCommandArgs(commandLine)
	switch cmd {
	case "test":
		err := utils.CreateTestFile(cfg.URL)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("test file save failed: %w", err)}
		}
		return CommandResult{Info: "Test file was created.", Quit: true}
	case "bye":
		return CommandResult{Info: "Chat has ended.", Quit: true}
	case "save":
		name, err := history.SaveHistoryToFile(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("history save failed: %w", err)}
		}
		return CommandResult{Info: fmt.Sprintf("History saved as '%s'", name)}
	case "load":
		filename, err := getHistoryFilename(args, input)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("history load failed: %w", err)}
		}

		if len(filename) > 0 {
			loaded, err := messages.LoadHistoryFromFile(filename)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("history load failed: %w", err)}
			}
			history.Replace(loaded.Get())
			return CommandResult{Info: "History loaded successfully."}
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
			return CommandResult{Error: fmt.Errorf("fetching server version failed: %w", err)}
		}

		list := []string{
			fmt.Sprintf("Configuration file used: %s", cfg.ConfigPath),
			fmt.Sprintf("Response output format: %s", cfg.OutputFmt),
			fmt.Sprintf("Current model is '%s'", cfg.Model),
			fmt.Sprintf("Context has %d messages (max. %d)", history.Len(), history.MaxCtx()),
			fmt.Sprintf("Context token estimation: %.0f", math.Ceil(history.EstimateTokens())),
			fmt.Sprintf("Server version is %s", serverVersion),
		}

		return CommandResult{Output: utils.FormatList(list, "Server info", false)}
	case "message":
		switch args {
		case "system", "user", "assistant":
			msg, found := history.GetLastRole(args)
			if found {
				return CommandResult{Output: msg.Content}
			} else {
				return CommandResult{Info: fmt.Sprintf("No element for role type '%s' found.", args)}
			}
		default:
			msg := history.GetLast().Content
			return CommandResult{Output: msg}

		}
	case "list":
		files, err := utils.ListHistoryFiles()
		if err != nil {
			return CommandResult{Error: fmt.Errorf("listing history files failed: %w", err)}
		}
		return CommandResult{Output: files}
	case "copy":
		lastAnswer := history.GetLast().Content
		if args == "think" {
			lastAnswer = history.GetLast().Raw
		}
		if args == "code" {
			codeBlock, found := messages.ExtractCodeBlock(lastAnswer)
			if found {
				lastAnswer = codeBlock
			} else {
				return CommandResult{Info: "Nothing to copy."}
			}
		}
		err := clipb.WriteClipboard(lastAnswer)
		if err != nil {
			return CommandResult{Error: err}
		}
		return CommandResult{Info: "Last answer written to clipboard."}
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
				return CommandResult{Error: fmt.Errorf("models list failed: %w", err)}
			}
			return CommandResult{Output: models}
		}

		index, err := strconv.Atoi(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("value not an integer")}
		}
		model, ok := utils.GetModelsByIndex(index)
		if !ok {
			return CommandResult{Error: fmt.Errorf("no value for given index found")}
		}
		cfg.Model = model
		return CommandResult{Info: fmt.Sprintf("Switched model to '%s'.", model)}
	case "set":
		if args == "" {
			list := []string{
				fmt.Sprintf("model = %s", cfg.Model),
				fmt.Sprintf("context = %d", cfg.Context),
				fmt.Sprintf("temperature = %.2f", cfg.Temperature),
				fmt.Sprintf("top_p = %.2f", cfg.TopP),
				fmt.Sprintf("reasoning = %t", cfg.Reasoning),
			}

			return CommandResult{Output: utils.FormatList(list, "Config settings", false)}
		}

		key, value, err := parseArgs(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("parse args failed: %w", err)}
		}
		err = config.Set(key, value)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("apply to config failed: %w", err)}
		}
		if key == "context" {
			intVal, ok := value.(int)
			if !ok {
				return CommandResult{Error: fmt.Errorf("value not an integer")}
			}
			err := history.SetContextSize(intVal)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("set context size failed: %w", err)}
			}
		}
		return CommandResult{Info: fmt.Sprintf("Config updated: %s = %v", key, value)}
	case "clear":
		history.ClearExceptSystemPrompt()
		return CommandResult{Info: "History cleared (system prompt retained)."}
	case "help":
		return CommandResult{Info: HelpText(args)}
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

// getHistoryFilename does the check of the filename for loading history sessions.
// It detects if an index (prefixed by #), a filename or none is given as arg.
// If the arg is empty, it requests for an input.
//
// Parameters:
//
//	f (string) - the argument of the /load command
//	input (io.Reader) - optional input stream for unit tests
//
// Returns:
//
//	string - selected filename of the history session
//	error  - Error msg if anything went wrong
func getHistoryFilename(f string, input io.Reader) (string, error) {
	filename, found := strings.CutPrefix(f, "#")
	if found {
		return getFilenameByIndex(filename)
	}

	if filename == "" {
		var err error
		filename, err = promptForFilename(input)
		if err != nil {
			return "", err
		}
	}

	return filename, nil
}

// getFilenameByIndex retrieves the filename corresponding to a given index string.
//
// Parameters:
//
//	indexStr - string representation of the index.
//
// Returns:
//
//	string - the filename associated with the index
//	error  - error if index is invalid or not found
func getFilenameByIndex(indexStr string) (string, error) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", fmt.Errorf("value not an integer")
	}
	fname, ok := utils.GetHistoryByIndex(index)
	if !ok {
		return "", fmt.Errorf("no value for given index found")
	}
	return fname, nil
}

// promptForFilename prompts the user to enter a filename to load.
//
// Parameters:
//
//	input - io.Reader used for reading user input.
//
// Returns:
//
//	string - the filename entered by the user
//	error  - error if reading input fails
func promptForFilename(input io.Reader) (string, error) {
	fmt.Print("Enter filename to load: ")
	reader := bufio.NewReader(input)
	inputLine, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("input read failed: %w", err)
	}
	filename := strings.TrimSpace(inputLine)
	return filename, nil
}
