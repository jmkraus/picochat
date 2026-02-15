package command

import (
	"fmt"
	"strings"
)

var helpTopics = map[string][]string{
	"": {
		"  [Ctrl]+D           Submit multiline input (EOF)",
		"  [Esc], [Ctrl]+C    Cancel multiline input and return to prompt",
		"  [Up] / [Down]      Browse prompt history",
		"  /copy, /c          Copy the last answer to clipboard",
		"  /paste, /v         Paste clipboard contents as user input and send",
		"  /info              Show system information",
		"  /message           Show the last message again (e.g., after load)",
		"  /load              Load chat history from a file",
		"  /save              Save current chat history to a file",
		"  /list              List available saved history files",
		"  /models            List downloaded models (and switch models)",
		"  /clear             Clear session context",
		"  /set               Set session variables (key=value)",
		"  /image             Set image file path",
		"  /retry             Resend the chat history, excluding the last answer",
		"  /bye               Quit PicoChat",
		"  /help, /?          Show available commands",
	},
	"copy": {
		"  /copy              Copy the last answer to clipboard",
		"  /copy code         Copy first code snippet enclosed in ``` to clipboard",
		"  /copy think        Copy the last answer to clipboard & retain reasoning",
		"  /copy <role>       Copy the last entry of the given role to clipboard",
		"  Valid roles: system, user, assistant",
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
