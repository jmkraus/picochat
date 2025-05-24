package command

import (
	"bufio"
	"fmt"
	"io"
	"picochat/config"
	"picochat/types"
	"picochat/utils"
	"strings"
)

func Handle(cmd string, history *types.ChatHistory, input io.Reader) types.CommandResult {
	cfg := config.Get()

	switch cmd {
	case "/bye":
		return types.CommandResult{Quit: true}
	case "/save":
		name, err := history.SaveHistoryToFile()
		if err != nil {
			return types.CommandResult{Output: "Save failed: " + err.Error()}
		}
		return types.CommandResult{Output: fmt.Sprintf("History saved as %s", name)}
	case "/load":
		fmt.Print("Enter filename to load: ")
		reader := bufio.NewReader(input)
		inputLine, _ := reader.ReadString('\n')
		filename := strings.TrimSpace(inputLine)

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
		return types.CommandResult{Output: fmt.Sprintf("History has %d messages (max. %d).", history.Len(), history.Max())}
	case "/list":
		files, err := utils.ListHistoryFiles()
		if err != nil {
			return types.CommandResult{Output: "Listing failed: " + err.Error()}
		}
		if len(files) == 0 {
			return types.CommandResult{Output: "No history files found."}
		}
		return types.CommandResult{Output: "Available history files:\n- " + strings.Join(files, "\n- ")}
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
		return types.CommandResult{Output: strings.Join(commands, "\n")}
	default:
		return types.CommandResult{Output: "Unknown command."}
	}
}
