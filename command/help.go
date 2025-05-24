// help.go
package command

import "strings"

func HelpText() string {
	commands := []string{
		"Available Commands:",
		"  /done       Terminate the input",
		"  /show       Show number of messages in history",
		"  /load       Load a session",
		"  /save       Save your current session",
		"  /list       List saved sessions",
		"  /models     List available (downloaded) models",
		"  /clear      Clear session context",
		"  /bye        Exit",
		"  /?, /help   Show available commands",
	}
	return strings.Join(commands, "\n")
}
