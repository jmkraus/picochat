package command

import (
	"bufio"
	"fmt"
	"io"
	"picochat/envs"
	"picochat/messages"
	"picochat/output"
	"picochat/utils"
	"picochat/vartypes"
	"strings"
)

// parseKeyVal parses a string of the form "key=value" and returns
// the canonical lowercase JSON key, converted value, and error. User input
// is accepted case-insensitively.
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

	fieldCfg, ok := envs.EnvSpecByField(key)
	if !ok || !fieldCfg.Runtime {
		return "", nil, fmt.Errorf("unsupported config key %q", key)
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
//	error - error if any
func parseIndex(indexStr string) (int, error) {
	indexAny, err := vartypes.Convert(vartypes.VarInt, indexStr)
	if err != nil {
		return 0, fmt.Errorf("value not an integer")
	}
	return indexAny.(int), nil
}

// getMessageByIndex retrieves a history message by index and formats it with
// a bold header line in the form "(index:role)".
//
// Parameters:
//
//	args (string) - the message index as string
//	history (*messages.ChatHistory) - chat history used for index lookup
//
// Returns:
//
//	string - formatted message including header and content
//	error  - error if index parsing or lookup fails
func getMessageByIndex(args string, history *messages.ChatHistory) (string, error) {
	index, err := parseIndex(args)
	if err != nil {
		return "", fmt.Errorf("get message failed: %w", err)
	}
	msg, err := history.GetByIndex(index)
	if err != nil {
		return "", fmt.Errorf("get message failed: %w", err)
	}
	return output.FormatMessage(msg, index, true, false), nil
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
//	error  - error if any
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

// askConfirmation asks a yes/no question and parses the input as boolean.
//
// Parameters:
//
//	question (string) - prompt text shown before the fixed "(y/n)" suffix
//	input (io.Reader) - input stream used for reading user input
//
// Returns:
//
//	bool  - parsed confirmation result
//	error - error if input cannot be read or parsed as boolean
func askConfirmation(question string, input io.Reader) (bool, error) {
	fmt.Printf("%s (y/N): ", question)

	reader := bufio.NewReader(input)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("input read failed: %w", err)
	}

	value := strings.TrimSpace(line)
	if value == "" {
		return false, nil
	}
	converted, err := vartypes.Convert(vartypes.VarBool, value)
	if err != nil {
		return false, fmt.Errorf("invalid confirmation value: %w", err)
	}

	return converted.(bool), nil
}
