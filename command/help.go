package command

import "strings"

func HelpText(helpArgs string) string {
	commands := []string{}

	helpArgs = strings.ToLower(strings.TrimSpace(helpArgs))
	switch helpArgs {
	case "", "help", "?":
		commands = []string{
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
	case "test":
		commands = []string{
			"TEST:",
			"  /test: Example output",
		}
	default:
		commands = []string{
			"No help available for: " + helpArgs,
			"Use /help for a list of commands.",
		}
	}
	return strings.Join(commands, "\n")
}
