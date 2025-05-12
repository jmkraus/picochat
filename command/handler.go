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

		if len(filename) > 0 {
			if err := history.LoadFromFile(filename); err != nil {
				return types.CommandResult{Output: "Load failed: " + err.Error()}
			}
			return types.CommandResult{Output: "History loaded successfully."}
		} else {
			return types.CommandResult{Output: "Load cancelled."}
		}
	case "/show":
		return types.CommandResult{Output: fmt.Sprintf("History has %d messages.", history.Len())}
	case "/list":
		files, err := utils.ListHistoryFiles()
		if err != nil {
			return types.CommandResult{Output: "Listing failed: " + err.Error()}
		}
		if len(files) == 0 {
			return types.CommandResult{Output: "No history files found."}
		}
		return types.CommandResult{Output: "Available history files:\n- " + strings.Join(files, "\n- ")}
	case "/clear":
		history.ClearExceptSystemPrompt()
		return types.CommandResult{Output: "History cleared (system prompt retained)."}
	case "/help":
		commands := []string{
			" /bye   - Exit the chat",
			" /done  - Terminate the input",
			" /save  - Save current chat history to a file",
			" /load  - Load chat history from a file",
			" /list  - List available saved history files",
			" /show  - Show number of messages in history",
			" /clear - Clear history and reinitialize with system prompt",
			" /help  - Show available commands",
		}
		return types.CommandResult{Output: strings.Join(commands, "\n")}
	default:
		return types.CommandResult{Output: "Unknown command."}
	}
}
