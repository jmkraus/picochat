package command

import (
	"fmt"
	"strings"
)

var helpTopics = map[string][]string{
	"": {
		"  [Ctrl]+D           Submit multiline input (EOF)",
		"  [Esc], [Ctrl]+C    Cancel multiline input and return to prompt",
		"  [Up] / [Down]      Browse prompt history (commands only)",
		"  /copy, /c          Copy selected answer to clipboard",
		"  /paste, /v         Paste clipboard content as user input and send",
		"  /info              Show system information",
		"  /trim              Remove all elements after given index",
		"  /message           Show message(s) from chat history",
		"  /load              Load chat history from file",
		"  /save              Save current chat history to file",
		"  /models            List downloaded models (and switch models)",
		"  /clear             Clear chat history (retaining system prompt)",
		"  /set               Set session variables (key=value)",
		"  /image             Set image file path",
		"  /retry             Resend the chat history excluding last answer",
		"  /bye               Quit PicoChat",
		"  /help, /?          Show available commands",
		"",
		"  /? envs            Show environment variable status table",
		"  /? templates       Show template key and description table",
	},
	"copy": {
		"  /copy              Copy the last answer to clipboard",
		"  /copy code         Copy first code snippet enclosed in ``` to clipboard",
		"  /copy think        Copy the last answer to clipboard & retain reasoning",
		"  /copy #<number>    Copy the message with index <number> to clipboard",
		"  /copy <role>       Copy the last entry of the given role to clipboard",
		"  Valid roles: system, user, assistant",
	},
	"paste": {
		"  /paste             Paste clipboard content as user prompt and send",
		"  /paste <key>       Prepend template text to pasted clipboard content and send",
	},
	"load": {
		"  /load              Show list of history files and request filename",
		"  /load <filename>   Load the history file with name <filename>",
		"  /load #<number>    Load the history file with index <number>",
		"  If no filename is entered, the load is canceled.",
	},
	"message": {
		"  /message           Show the last entry in the chat history",
		"  /message all       Show full conversation with color coded roles",
		"  /message #<number> Show the message with index <number>",
		"  /message <role>    Show the last entry of the given role",
		"  Valid roles: system, user, assistant",
	},
	"models": {
		"  /models            List the available models of the LLM server",
		"  /models <number>   Load the model by index <number>",
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
		"",
		"  Value ranges:",
		"  context            3..100",
		"  temperature        0..2",
		"  top_p              0..1",
		"  effort             none, low, medium, high",
	},
}

func HelpText(topic string) string {
	topic = strings.ToLower(strings.TrimSpace(topic))
	header := "Available commands:"

	if lines, ok := helpTopics[topic]; ok {
		if topic != "" {
			header = fmt.Sprintf("Details for %s", topic)
		}
		help := append([]string{header}, lines...)
		return strings.Join(help, "\n")
	}
	return fmt.Sprintf("No help available for: %s. Use /? for a list of commands.", topic)
}
