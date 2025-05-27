// help.go
package command

import "strings"

func HelpText() string {
	commands := []string{
		"Available Commands:",
		"  /done, ///  Terminate the input",
		"  /copy       Copy the last answer to clipboard",
		"  /show       Show number of messages in history",
		"  /load       Load a session",
		"  /save       Save current session",
		"  /list       List saved sessions",
		"  /models     List available (downloaded) models",
		"  /clear      Clear session context",
		"  /bye        Exit",
		"  /?, /help   Show available commands",
	}
	return strings.Join(commands, "\n")
}
