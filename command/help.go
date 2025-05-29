package command

import "strings"

var helpTopics = map[string][]string{
	"": {
		"Available Commands:",
		"  /done, ///   Terminate the input",
		"  /copy        Copy the last answer to clipboard",
		"  /show        Show system information",
		"  /load        Load a session",
		"  /save        Save current session",
		"  /list        List saved sessions",
		"  /models      List available (downloaded) models",
		"  /clear       Clear session context",
		"  /bye         Exit",
		"  /help, /?    Show available commands",
	},
	"copy": {
		"Available Commands:",
		"  /copy        Copy the last answer to clipboard",
		"  /copy code   Copy only code between ``` to clipboard",
	},
	"load": {
		"Standard command:",
		"  /load <filename>   Load the history file with name <filename>",
		"  If <filename> is omitted, the filename is requested by input line.",
		"  If no filename is entered, the load is cancelled.",
	},
	"save": {
		"Standard command:",
		"  /save <filename>   Save the history file with name <filename>",
		"  If <filename> is omitted, the filename is set as current timestamp.",
	},
}

func HelpText(topic string) string {
	topic = strings.ToLower(strings.TrimSpace(topic))

	if lines, ok := helpTopics[topic]; ok {
		return strings.Join(lines, "\n")
	}
	return "No help available for: " + topic + "\nUse /help for a list of commands."
}
