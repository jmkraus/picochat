package command

import (
	"fmt"
	"picochat/types"
)

func Handle(cmd string, history *types.ChatHistory) types.CommandResult {
	switch cmd {
	case "/bye":
		return types.CommandResult{Quit: true}
	case "/save":
		// Placeholder
		return types.CommandResult{Output: "Saving history is not implemented yet."}
	case "/load":
		// Placeholder
		return types.CommandResult{Output: "Loading history is not implemented yet."}
	case "/show":
		return types.CommandResult{Output: fmt.Sprintf("History has %d messages.", len(history.Get()))}
	default:
		return types.CommandResult{Output: "Unknown command."}
	}
}
