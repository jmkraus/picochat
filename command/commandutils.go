package command

import (
	"bufio"
	"fmt"
	"io"
	"picochat/utils"
	"regexp"
	"strconv"
	"strings"
)

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
	filename, found := strings.CutPrefix(f, "#")
	if found {
		return getFilenameByIndex(filename)
	}

	if filename == "" {
		var err error
		filename, err = promptForFilename(input)
		if err != nil {
			return "", err
		}
	}

	return filename, nil
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
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", fmt.Errorf("value not an integer")
	}
	fname, ok := utils.GetHistoryByIndex(index)
	if !ok {
		return "", fmt.Errorf("no value for given index found")
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
	fmt.Print("Enter filename to load: ")
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
