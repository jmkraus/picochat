package command

import (
	"bufio"
	"fmt"
	"io"
	"picochat/envs"
	"picochat/messages"
	"picochat/utils"
	"picochat/vartypes"
	"regexp"
	"strings"
)

type copyPayload struct {
	Text string
	Info string
}

// parseKeyVal parses a string of the form "key=value" and returns
// the key, converted value, and error.
//
// Parameters:
//
//	args (string) - the input string to parse
//
// Returns:
//
//	string - the parsed key in lower case
//	any    - the parsed and converted value
//	error  - error if any
func parseKeyVal(args string) (string, any, error) {
	parts := strings.SplitN(args, "=", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid format, expected key=value")
	}

	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", nil, fmt.Errorf("invalid format, missing key")
	}
	if value == "" {
		return "", nil, fmt.Errorf("invalid format, missing value")
	}

	fieldCfg, ok := envs.ConfigByField(key)
	if !ok || !fieldCfg.Runtime {
		return "", nil, fmt.Errorf("unsupported config key '%s'", key)
	}

	convertedValue, err := vartypes.Convert(fieldCfg.Type, value)
	if err != nil {
		return "", nil, fmt.Errorf("convert type for key %s failed: %w", key, err)
	}

	return key, convertedValue, nil
}

// parseIndex parses an integer index from string input.
//
// Parameters:
//
//	indexStr - string representation of the index.
//
// Returns:
//
//	int   - parsed index
//	error - error if input is not an integer
func parseIndex(indexStr string) (int, error) {
	indexAny, err := vartypes.Convert(vartypes.VarInt, indexStr)
	if err != nil {
		return 0, fmt.Errorf("value not an integer")
	}
	return indexAny.(int), nil
}

// getHistoryFilename does the check of the filename for loading history sessions.
// It detects if an index (prefixed by #), a filename or none is given as arg.
// If the arg is empty, it requests for an input.
//
// Parameters:
//
//	f (string) - the argument of the /load command
//	input (io.Reader) - optional input stream for unit tests
//
// Returns:
//
//	string - selected filename of the history session
//	error  - Error msg if anything went wrong
func getHistoryFilename(f string, input io.Reader) (string, error) {
	if f == "" {
		var err error
		f, err = promptForFilename(input)
		if err != nil {
			return "", err
		}
	}

	if n, ok := strings.CutPrefix(f, "#"); ok {
		return getFilenameByIndex(n)
	}
	return f, nil
}

// getFilenameByIndex retrieves the filename corresponding to a given index string.
//
// Parameters:
//
//	indexStr - string representation of the index.
//
// Returns:
//
//	string - the filename associated with the index
//	error  - error if index is invalid or not found
func getFilenameByIndex(indexStr string) (string, error) {
	index, err := parseIndex(indexStr)
	if err != nil {
		return "", err
	}
	fname, ok := utils.GetHistoryByIndex(index)
	if !ok {
		return "", fmt.Errorf("index %d out of bounds", index)
	}
	return fname, nil
}

// promptForFilename prompts the user to enter a filename to load.
//
// Parameters:
//
//	input - io.Reader used for reading user input.
//
// Returns:
//
//	string - the filename entered by the user
//	error  - error if reading input fails
func promptForFilename(input io.Reader) (string, error) {
	fmt.Print("\nEnter filename or #<index> to load: ")
	reader := bufio.NewReader(input)
	inputLine, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("input read failed: %w", err)
	}
	filename := strings.TrimSpace(inputLine)
	return filename, nil
}

// extractCodeBlock extracts the first code block from a string
// formatted with triple backticks.
//
// Parameters:
//
//	s (string) - the input string containing code blocks.
//
// Returns:
//
//	string - the extracted code block content.
//	bool - true if a code block was found, false otherwise.
func extractCodeBlock(s string) (string, bool) {
	re := regexp.MustCompile("(?s)```\\w*\\n(.*?)```")
	match := re.FindStringSubmatch(s)
	if len(match) >= 2 {
		return match[1], true
	}
	return "", false
}

// encloseThinkingTags adds tags around a given string to
// identify it as the reasoning part of the text
//
// Parameters:
//
//	s (string) - the string to be tagged
//
// Returns:
//
//	string - the tagged text
func encloseThinkingTags(s string) string {
	return fmt.Sprintf("<think>\n%s\n</think>\n\n", s)
}

// resolveCopyPayload determines which text should be copied based on the given
// /copy arguments and returns the text plus the corresponding status info.
//
// Parameters:
//
//	args (string) - argument passed to the /copy command (e.g. "#3", "assistant", "code")
//	history (*messages.ChatHistory) - chat history used as source for message lookup
//
// Returns:
//
//	copyPayload - resolved text and info message for clipboard handling
//	error       - error if argument is invalid or index lookup fails
func resolveCopyPayload(args string, history *messages.ChatHistory) (copyPayload, error) {
	nothing := "Nothing to copy."
	if indexStr, ok := strings.CutPrefix(args, "#"); ok {
		index, err := parseIndex(indexStr)
		if err != nil {
			return copyPayload{}, err
		}
		msg, err := history.GetByIndex(index)
		if err != nil {
			return copyPayload{}, err
		}
		return copyPayload{
			Text: msg.Content,
			Info: fmt.Sprintf("Message #%d written to clipboard", index),
		}, nil
	}

	if args == "" {
		args = messages.RoleAssistant
	}

	switch args {
	case messages.RoleAssistant, messages.RoleUser, messages.RoleSystem:
		lastMessage, found := history.GetLastRole(args)
		if !found || lastMessage.Content == "" {
			return copyPayload{Info: nothing}, nil
		}
		return copyPayload{
			Text: lastMessage.Content,
			Info: fmt.Sprintf("Last %s prompt written to clipboard.", args),
		}, nil

	case "think":
		lastMessage, found := history.GetLastRole(messages.RoleAssistant)
		if !found || (lastMessage.Content == "" && lastMessage.Reasoning == "") {
			return copyPayload{Info: nothing}, nil
		}
		return copyPayload{
			Text: encloseThinkingTags(lastMessage.Reasoning) + lastMessage.Content,
			Info: "Last assistant prompt (with thinking) written to clipboard.",
		}, nil

	case "code":
		lastMessage, found := history.GetLastRole(messages.RoleAssistant)
		if !found || lastMessage.Content == "" {
			return copyPayload{Info: nothing}, nil
		}
		codeBlock, found := extractCodeBlock(lastMessage.Content)
		if !found {
			return copyPayload{Info: nothing}, nil
		}
		return copyPayload{
			Text: codeBlock,
			Info: "First code block written to clipboard.",
		}, nil

	default:
		return copyPayload{}, fmt.Errorf("unknown copy argument")
	}
}
