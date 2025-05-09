package command

import (
	"picochat/types"
	"strings"
)

// Example structure for a chat message (depending on your project, you may need to import or adjust it)

// Handle processes slash commands like /bye, /show, /load, /save
// Return: (Output to user, should-program-be-ended)
func Handle(cmd string, history *[]types.Message) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(cmd)) {
	case "/bye":
		return "", true

	case "/show":
		// later implement
		return "(not yet implemented: /show)", false

	case "/save":
		// later implement
		return "(not yet implemented: /save)", false

	case "/load":
		// later implement
		return "(not yet implemented: /load)", false

	default:
		return "Unknown command: " + cmd, false
	}
}
