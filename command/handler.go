package command

import (
	"bufio"
	"fmt"
	"os"
	"picochat/types"
	"picochat/utils"
	"strings"
)

func Handle(cmd string, history *types.ChatHistory) types.CommandResult {
	switch cmd {
	case "/bye":
		return types.CommandResult{Quit: true}
	case "/save":
		name, err := history.SaveToFile()
		if err != nil {
			return types.CommandResult{Output: "Save failed: " + err.Error()}
		}
		return types.CommandResult{Output: fmt.Sprintf("History saved as %s", name)}
	case "/load":
		fmt.Print("Enter filename to load: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		filename := strings.TrimSpace(input)

		if err := history.LoadFromFile(filename); err != nil {
			return types.CommandResult{Output: "Load failed: " + err.Error()}
		}
		return types.CommandResult{Output: "History loaded successfully."}
	case "/show":
		return types.CommandResult{Output: fmt.Sprintf("History has %d messages.", len(history.Get()))}
	case "/list":
		files, err := utils.ListHistoryFiles()
		if err != nil {
			return types.CommandResult{Output: "Listing failed: " + err.Error()}
		}
		if len(files) == 0 {
			return types.CommandResult{Output: "No history files found."}
		}
		return types.CommandResult{Output: "Available history files:\n- " + strings.Join(files, "\n- ")}
	default:
		return types.CommandResult{Output: "Unknown command."}
	}
}
