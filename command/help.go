package command

import "strings"

var helpTopics = map[string][]string{
	"": {
		"Available Commands:",
		"  /done, ///   Terminate the input message and send",
		"  /cancel      Cancel multi-line input and return to prompt",
		"  /copy        Copy the last answer to clipboard",
		"  /paste       Get clipboard content as user input and send",
		"  /info        Show system information",
		"  /load        Load a session",
		"  /save        Save current session",
		"  /list        List saved sessions",
		"  /models      List available (downloaded) models",
		"  /clear       Clear session context",
		"  /set         Set session variables",
		"  /retry       Sends chat history again, but without last answer",
		"  /bye         Exit",
		"  /help, /?    Show available commands",
	},
	"copy": {
		"Available Commands:",
		"  /copy        Copy the last answer to clipboard",
		"  /copy code   Copy only code between ``` to clipboard",
	},
	"models": {
		"Available Commands:",
		"  /models            Lists the available models of the LLM server",
		"  /models <index>    Loads the model by index number",
		"  To use the load option, list the available models first & check index",
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
