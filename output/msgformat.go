package output

import (
	"fmt"
	"picochat/console"
	"picochat/messages"
	"strings"
)

// FormatMessage formats a single chat message with optional index header and
// role-based color output.
//
// Parameters:
//
//	msg (messages.Message) - the message to format
//	index (int)            - the message index shown in the header
//	header (bool)          - include role/index header if true
//	color (bool)           - apply role-based colors if true
//
// Returns:
//
//	string - the formatted message text
func FormatMessage(msg messages.Message, index int, header, color bool) string {
	headerText := ""
	if header {
		headerText = fmt.Sprintf("(%d:%s)\n", index, msg.Role)
		headerText = console.Style(console.Bold, headerText)
	}

	output := fmt.Sprintf("%s%s", headerText, msg.Content)

	if color {
		switch msg.Role {
		case "system":
			output = console.Colorize(console.Magenta, output)
		case "user":
			output = console.Colorize(console.Cyan, output)
		case "assistant":
			// nothing to do here
		}
	}
	return output
}

// FormatConversation returns a color coded full conversation text.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the full conversation text (without reasoning)
func FormatConversation(msgs []messages.Message) string {
	var builder strings.Builder

	for index, msg := range msgs {
		formatted := FormatMessage(msg, index, true, true)
		builder.WriteString(formatted)
		builder.WriteString("\n\n") // Maintains the double newline between messages
	}

	return builder.String()
}
