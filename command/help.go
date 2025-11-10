package command

import (
	"fmt"
	"strings"
)

var helpTopics = map[string][]string{
	"": {
		"  [Ctrl] + D      Submit multiline input (EOF)",
		"  [Esc]           Cancel multiline input and return to prompt",
		"  [Up] / [Down]   Browse prompt command history",
		"  /copy, /c       Copy the last answer to clipboard",
		"  /paste, /v      Get clipboard content as user input and send",
		"  /info           Show system information",
		"  /message        Show last message again",
		"  /load           Load a session",
		"  /save           Save current session",
		"  /list           List saved sessions",
		"  /models         List available (downloaded) models",
		"  /clear          Clear session context",
		"  /set            Set session variables (key=value)",
		"  /retry          Sends chat history again, but without last answer",
		"  /bye            Exit",
		"  /help, /?       Show available commands",
	},
	"copy": {
		"  /copy         Copy the last answer to clipboard",
		"  /copy think   Copy the last answer to clipboard & retain reasoning",
		"  /copy code    Copy only code snippets between ``` to clipboard",
	},
	"load": {
		"  /load <filename>   Load the history file with name <filename>",
		"  /load #<number>    Load the history file with index <number>",
		"  If <filename> is omitted, the filename is requested by input line.",
		"  If no filename is entered, the load is canceled.",
		"  To use the index load, enter '/list' command first & check index.",
	},
	"message": {
		"  /message           Shows the last entry in the chat history",
		"  /message <role>    Shows the last entry of the given role",
		"  Valid roles: system, user, assistant",
	},
	"models": {
		"  /models            Lists the available models of the LLM server",
		"  /models <number>   Loads the model by index <number>",
		"  To use the load option, list the available models first & check index.",
	},
	"save": {
		"  /save <filename>   Save the history file with name <filename>",
		"  If <filename> is omitted, the filename is set as current timestamp.",
	},
	"set": {
		"  /set               Show available parameters and current settings",
		"  /set <key=value>   Set the parameter <key> to new setting <value>",
		"  Example: /set temperature=0.7",
	},
}

func HelpText(topic string) string {
	topic = strings.ToLower(strings.TrimSpace(topic))

	if lines, ok := helpTopics[topic]; ok {
		help := append([]string{"Available commands:"}, lines...)
		return strings.Join(help, "\n")
	}
	return fmt.Sprintf("No help available for: %s. Use /help for a list of commands.", topic)
}
