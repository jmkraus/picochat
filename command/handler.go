package command

import (
	"bufio"
	"fmt"
	"io"
	"picochat/config"
	"picochat/messages"
	"picochat/requests"
	"picochat/utils"
	"strings"

	"github.com/atotto/clipboard"
)

type CommandResult struct {
	Output string
	Error  string
	Quit   bool
	Prompt string
	Repeat bool
}

func HandleCommand(commandLine string, history *messages.ChatHistory, input io.Reader) CommandResult {
	cfg := config.Get()

	cmd, args := parseCommandArgs(commandLine)
	switch cmd {
	case "/done", "///":
		return CommandResult{Output: "Use this command for terminating a multi-line input."}
	case "/bye":
		return CommandResult{Output: "Chat has ended.", Quit: true}
	case "/save":
		name, err := history.SaveHistoryToFile(args)
		if err != nil {
			return CommandResult{Error: "history save failed: " + err.Error()}
		}
		return CommandResult{Output: fmt.Sprintf("History saved as '%s'", name)}
	case "/load":
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
				return CommandResult{Error: "history load failed: " + err.Error()}
			}
			history.Replace(loaded.Get())
			return CommandResult{Output: "History loaded successfully."}
		} else {
			return CommandResult{Output: "Load cancelled."}
		}
	case "/show":
		serverVersion, err := requests.GetServerVersion(cfg.URL)
		if err != nil {
			return CommandResult{Error: "fetching server version failed: " + err.Error()}
		}

		model := fmt.Sprintf("Current model is '%s'", cfg.Model)
		messages := fmt.Sprintf("Context has %d messages (max. %d)", history.Len(), history.MaxCtx())
		server := fmt.Sprintf("Server version is %s", serverVersion)

		return CommandResult{Output: model + "\n" + messages + "\n" + server}
	case "/list":
		files, err := utils.ListHistoryFiles()
		if err != nil {
			return CommandResult{Error: "listing history files failed: " + err.Error()}
		}
		return CommandResult{Output: files}
	case "/copy":
		lastAnswer := utils.StripReasoning(history.GetLast().Content)
		if args == "code" {
			lastAnswer = utils.ExtractCodeBlock(lastAnswer)
		}
		err := clipboard.WriteAll(lastAnswer)
		if err != nil {
			return CommandResult{Error: "clipboard failed: " + err.Error()}
		}
		if utils.IsTmuxSession() {
			err := utils.CopyToTmuxBufferStdin(lastAnswer)
			if err != nil {
				return CommandResult{Error: "tmux clipboard failed: " + err.Error()}
			}
		}
		return CommandResult{Output: "Last answer written to clipboard."}
	case "/paste":
		text, err := clipboard.ReadAll()
		if err != nil {
			return CommandResult{Error: "clipboard read failed: " + err.Error()}
		}
		text = strings.TrimSpace(text)
		if text == "" {
			return CommandResult{Error: "clipboard is empty."}
		}

		return CommandResult{
			Output: fmt.Sprintf("Pasted %d characters from clipboard.", len(text)),
			Prompt: text,
		}
	case "/discard":
		history.Discard()
		return CommandResult{Output: "Last answer removed from chat history."}
	case "/retry":
		return CommandResult{Output: "Repeating last chat history content.", Repeat: true}
	case "/models":
		if args == "" {
			models, err := utils.ShowAvailableModels(cfg.URL)
			if err != nil {
				return CommandResult{Error: "models list failed: " + err.Error()}
			}
			return CommandResult{Output: models}
		}

		return CommandResult{Output: "Not implemented."}
	case "/set":
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
			return CommandResult{Error: "parse args failed: " + err.Error()}
		}
		err = applyToConfig(key, value)
		if err != nil {
			return CommandResult{Error: "apply to config failed: " + err.Error()}
		}
		if key == "context" {
			intVal, ok := value.(int)
			if !ok {
				return CommandResult{Error: "value not an integer"}
			}
			err := history.SetContextSize(intVal)
			if err != nil {
				return CommandResult{Error: "update context size failed: " + err.Error()}
			}
		}
		return CommandResult{Output: fmt.Sprintf("Config updated: %s = %v", key, value)}
	case "/clear":
		history.ClearExceptSystemPrompt()
		return CommandResult{Output: "History cleared (system prompt retained)."}
	case "/help", "/?":
		return CommandResult{Output: HelpText(args)}
	default:
		return CommandResult{Error: "unknown command"}
	}
}

func parseCommandArgs(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}
	cmd := strings.TrimSpace(parts[0])
	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(strings.Join(parts[1:], " "))
	}
	return cmd, arg
}
