package command

import (
	"bufio"
	"fmt"
	"io"
	"picochat/config"
	"picochat/requests"
	"picochat/types"
	"picochat/utils"
	"strings"

	"github.com/atotto/clipboard"
)

func Handle(commandLine string, history *types.ChatHistory, input io.Reader) types.CommandResult {
	cfg := config.Get()

	cmd, args := parseCommandArgs(commandLine)
	switch cmd {
	case "/done", "///":
		return types.CommandResult{Output: "Use this command for terminating a multi-line input."}
	case "/bye":
		return types.CommandResult{Output: "Chat has ended.", Quit: true}
	case "/save":
		name, err := history.SaveHistoryToFile(args)
		if err != nil {
			return types.CommandResult{Output: "Save failed: " + err.Error()}
		}
		return types.CommandResult{Output: fmt.Sprintf("History saved as %s", name)}
	case "/load":
		filename := args
		if filename == "" {
			fmt.Print("Enter filename to load: ")
			reader := bufio.NewReader(input)
			inputLine, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(inputLine)
		}

		if len(filename) > 0 {
			loaded, err := types.LoadHistoryFromFile(filename)
			if err != nil {
				return types.CommandResult{Output: "Load failed: " + err.Error()}
			}
			history.Replace(loaded.Get())
			return types.CommandResult{Output: "History loaded successfully."}
		} else {
			return types.CommandResult{Output: "Load cancelled."}
		}
	case "/show":
		serverVersion, err := requests.GetServerVersion(cfg.URL)
		if err != nil {
			return types.CommandResult{Output: "Fetching server version failed: " + err.Error()}
		}

		server := fmt.Sprintf("Server version is %s", serverVersion)
		messages := fmt.Sprintf("History has %d messages (max. %d)", history.Len(), history.Max())
		model := fmt.Sprintf("Current model is '%s'", cfg.Model)

		return types.CommandResult{Output: model + "\n" + messages + "\n" + server}
	case "/list":
		files, err := utils.ListHistoryFiles()
		if err != nil {
			return types.CommandResult{Output: "Listing failed: " + err.Error()}
		}
		return types.CommandResult{Output: files}
	case "/copy":
		lastAnswer := utils.StripReasoning(history.GetLast().Content)
		if args == "code" {
			lastAnswer = utils.ExtractCodeBlock(lastAnswer)
		}
		err := clipboard.WriteAll(lastAnswer)
		if err != nil {
			return types.CommandResult{Output: "Clipboard failed: " + err.Error()}
		}
		return types.CommandResult{Output: "Last answer written to clipboard."}
	case "/models":
		models, err := utils.ShowAvailableModels(cfg.URL)
		if err != nil {
			return types.CommandResult{Output: "Models failed: " + err.Error()}
		}
		return types.CommandResult{Output: models}
	case "/clear":
		history.ClearExceptSystemPrompt()
		return types.CommandResult{Output: "History cleared (system prompt retained)."}
	case "/help", "/?":
		return types.CommandResult{Output: HelpText(args)}
	default:
		return types.CommandResult{Output: "Unknown command."}
	}
}

func parseCommandArgs(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}
	cmd := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = strings.Join(parts[1:], " ")
	}
	return cmd, arg
}
