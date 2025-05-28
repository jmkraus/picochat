package command

import "strings"

var helpTopics = map[string][]string{
	"": {
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
		"  /help, /?   Show available commands",
	},
	"test": {
		"TEST:",
		"  /test: Example output",
	},
}

func HelpText(topic string) string {
	topic = strings.ToLower(strings.TrimSpace(topic))

	if lines, ok := helpTopics[topic]; ok {
		return strings.Join(lines, "\n")
	}
	return "No help available for: " + topic + "\nUse /help for a list of commands."
}
