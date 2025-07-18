package command

import (
	"bufio"
	"fmt"
	"io"
	"picochat/config"
	"picochat/messages"
	"picochat/requests"
	"picochat/utils"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

type CommandResult struct {
	Output string
	Error  error
	Quit   bool
	Prompt string
	Repeat bool
}

func HandleCommand(commandLine string, history *messages.ChatHistory, input io.Reader) CommandResult {
	cfg := config.Get()

	cmd, args := parseCommandArgs(commandLine)
	switch cmd {
	case "done":
		return CommandResult{Output: "Use this command for terminating a multi-line input."}
	case "cancel":
		return CommandResult{Output: "Use this command for cancelling a multi-line input."}
	case "bye":
		return CommandResult{Output: "Chat has ended.", Quit: true}
	case "save":
		name, err := history.SaveHistoryToFile(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("history save failed: %w", err)}
		}
		return CommandResult{Output: fmt.Sprintf("History saved as '%s'", name)}
	case "load":
		filename := args
		if filename == "" {
			fmt.Print("Enter filename to load: ")
			reader := bufio.NewReader(input)
			inputLine, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(inputLine)
		}

		if len(filename) > 0 {
			loaded, err := messages.LoadHistoryFromFile(filename)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("history load failed: %w", err)}
			}
			history.Replace(loaded.Get())
			return CommandResult{Output: "History loaded successfully."}
		} else {
			return CommandResult{Output: "Load cancelled."}
		}
	case "info":
		serverVersion, err := requests.GetServerVersion(cfg.URL)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("fetching server version failed: %w", err)}
		}

		model := fmt.Sprintf("Current model is '%s'", cfg.Model)
		messages := fmt.Sprintf("Context has %d messages (max. %d)", history.Len(), history.MaxCtx())
		tokens := fmt.Sprintf("Context token estimation: %d", history.EstimateTokens())
		server := fmt.Sprintf("Server version is %s", serverVersion)

		return CommandResult{Output: model + "\n" + messages + "\n" + tokens + "\n" + server}
	case "message":
		lastAnswer := utils.StripReasoning(history.GetLast().Content)
		return CommandResult{Output: lastAnswer}
	case "list":
		files, err := utils.ListHistoryFiles()
		if err != nil {
			return CommandResult{Error: fmt.Errorf("listing history files failed: %w", err)}
		}
		return CommandResult{Output: files}
	case "copy":
		lastAnswer := utils.StripReasoning(history.GetLast().Content)
		if args == "code" {
			lastAnswer = utils.ExtractCodeBlock(lastAnswer)
		}
		err := clipboard.WriteAll(lastAnswer)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("clipboard failed: %w", err)}
		}
		if utils.IsTmuxSession() {
			err := utils.CopyToTmuxBufferStdin(lastAnswer)
			if err != nil {
				return CommandResult{Error: fmt.Errorf("tmux clipboard failed: %w", err)}
			}
		}
		return CommandResult{Output: "Last answer written to clipboard."}
	case "paste":
		text, err := clipboard.ReadAll()
		if err != nil {
			return CommandResult{Error: fmt.Errorf("clipboard read failed: %w", err)}
		}
		text = strings.TrimSpace(text)
		if text == "" {
			return CommandResult{Error: fmt.Errorf("clipboard is empty")}
		}

		return CommandResult{
			Output: fmt.Sprintf("Pasted %d characters from clipboard.", len(text)),
			Prompt: text,
		}
	case "retry":
		history.Discard()
		return CommandResult{Output: "Repeating last chat history user prompt.", Repeat: true}
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
		err = applyToConfig("model", model)
		return CommandResult{Output: fmt.Sprintf("Switched model to '%s'.", model)}
	case "set":
		if args == "" {
			cfg := config.Get()

			list := []string{
				fmt.Sprintf("context = %d", cfg.Context),
				fmt.Sprintf("temperature = %.2f", cfg.Temperature),
				fmt.Sprintf("top_p = %.2f", cfg.TopP),
			}

			return CommandResult{Output: utils.FormatList(list, "Config settings", false)}
		}

		key, value, err := ParseArgs(args)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("parse args failed: %w", err)}
		}
		err = applyToConfig(key, value)
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
				return CommandResult{Error: fmt.Errorf("update context size failed: %w", err)}
			}
		}
		return CommandResult{Output: fmt.Sprintf("Config updated: %s = %v", key, value)}
	case "clear":
		history.ClearExceptSystemPrompt()
		return CommandResult{Output: "History cleared (system prompt retained)."}
	case "help", "?":
		return CommandResult{Output: HelpText(args)}
	default:
		return CommandResult{Error: fmt.Errorf("unknown command")}
	}
}

func parseCommandArgs(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}
	cmd := strings.TrimSpace(parts[0])

	// replace special abbreviation for done
	if cmd == "///" {
		cmd = "/done"
	}

	// normalize
	cmd = strings.ToLower(strings.TrimPrefix(cmd, "/"))

	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(strings.Join(parts[1:], " "))
	}
	return cmd, arg
}
