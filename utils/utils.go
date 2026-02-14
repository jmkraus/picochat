package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"picochat/paths"
	"picochat/requests"
	"runtime"
	"strings"
	"unicode"
)

// Global cache for list selections (/list, /models)
var (
	ModelsList  []string
	HistoryList []string
)

// Indexed list
var (
	ModelsMap  = make(map[int]string)
	HistoryMap = make(map[int]string)
)

// ListHistoryFiles shows all saved sessions in 'history' dir
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - formatted list of history files
//	error
func ListHistoryFiles() (string, error) {
	historyPath, err := paths.GetHistoryPath()
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(historyPath)
	if err != nil {
		return "", err
	}

	var history []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".chat") {
			history = append(history, entry.Name())
		}
	}

	if len(history) == 0 {
		return "", fmt.Errorf("no history files found.")
	}

	setHistoryList(history)
	return FormatList(history, "history files", true), nil
}

// ShowAvailableModels lists all models via /tags API call
//
// Parameters:
//
//	baseUrl (string) - Base URL for the API
//
// Returns:
//
//	string - formatted list of available models
//	error
func ShowAvailableModels(baseUrl string) (string, error) {
	models, err := requests.GetAvailableModels(baseUrl)
	if err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no models available.")
	}

	setModelsList(models)
	return FormatList(models, "Language models", true), nil
}

// Filled by /list command
//
// Parameters:
//
//	list ([]string) - List of stored history sessions
//
// Returns:
//
//	none
func setHistoryList(list []string) {
	HistoryList = list
	HistoryMap = make(map[int]string)
	for i, name := range list {
		HistoryMap[i+1] = name
	}
}

// GetHistoryByIndex retrieves a history session by its index
//
// Parameters:
//
//	i (int) - Index of the history session
//
// Returns:
//
//	string - session name
//	bool   - boolean indicating success
func GetHistoryByIndex(i int) (string, bool) {
	val, ok := HistoryMap[i]
	return val, ok
}

// Filled by /models command
//
// Parameters:
//
//	list ([]string) - List of available Models
//
// Returns:
//
//	none
func setModelsList(list []string) {
	ModelsList = list
	ModelsMap = make(map[int]string)
	for i, name := range list {
		ModelsMap[i+1] = name
	}
}

// GetModelsByIndex retrieves a model by its index
//
// Parameters:
//
//	i (int) - Index of the model
//
// Returns:
//
//	string - model name
//	bool   - boolean indicating success
func GetModelsByIndex(i int) (string, bool) {
	val, ok := ModelsMap[i]
	return val, ok
}

// FormatList formats a list of items with an optional heading and numbering
//
// Parameters:
//
//	content ([]string) - Items to format
//	heading (string)   - Heading for the list
//	numbered (bool)    - Whether to number the items
//
// Returns:
//
//	string - formatted list
func FormatList(content []string, heading string, numbered bool) string {
	if len(content) == 0 {
		return fmt.Sprintf("no %s found.", strings.ToLower(heading))
	}

	var lines []string
	for i, item := range content {
		if numbered {
			lines = append(lines, fmt.Sprintf("(%02d) %s", i+1, item))
		} else {
			lines = append(lines, " - "+item)
		}
	}

	return fmt.Sprintf("%s:\n%s", capitalize(heading), strings.Join(lines, "\n"))
}

// capitalize capitalizes the first letter of a string
//
// Parameters:
//
//	s (string) - Input string
//
// Returns:
//
//	string - String with first letter capitalized
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CreateTestFile writes a batch script where each downloaded model is called with a sample text.
// The text contains Western and non-Western characters as well as an Emoji.
// This is for internal test purposes only and runs with the '/test' command.
//
// Parameters:
//
//	baseUrl (string) - Base URL for the API
//
// Returns:
//
//	none
func CreateTestFile(baseUrl string) error {
	var (
		def   string
		cmd   string
		ext   string
		first string
	)

	const (
		text = "中国是一个拥有悠久历史的文明古国。 Can you translate this into English? 😊"
	)

	switch runtime.GOOS {
	case "windows":
		def = fmt.Sprintf("set \"TEXT=%s\"", text)
		cmd = "echo %%TEXT%% | picochat -model %s"
		ext = "cmd"
		first = "@echo off\nchcp 65001"
	default:
		def = fmt.Sprintf("TEXT='%s'", text)
		cmd = "echo \"$TEXT\" | picochat -model %s"
		ext = "sh"
		first = "#!/bin/sh"
	}

	models, err := requests.GetAvailableModels(baseUrl)
	if err != nil {
		return err
	}

	if len(models) == 0 {
		return fmt.Errorf("no models available")
	}

	var rows []string
	rows = append(rows, first)
	rows = append(rows, def)
	for _, m := range models {
		rows = append(rows, fmt.Sprintf(cmd, m))
	}

	filename := fmt.Sprintf("test.%s", ext)
	if err := os.WriteFile(filename, []byte(strings.Join(rows, "\n")), 0755); err != nil {
		return fmt.Errorf("write test file %s failed: %w", filename, err)
	}

	return nil
}

// ImageToBase64 reads an image file and converts it to a base64 string
//
// Parameter:
//
//	path (string) - the path to the image
//
// Returns:
//
//	string - the base64 encoded image
//	error  - error if any
func ImageToBase64(path string) (string, error) {
	// Expand homedir if applicable
	fullPath, err := paths.ExpandHomeDir(path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

// LoadSchemaFromFile loads a json schema string file and
// transforms it into a json object representation.
// ollama also allows for a simple "json" string for
// arbitrary output in json format.
//
// Parameters:
//
//	path (string) - the path to the json schema text file
//
// Returns:
//
//	any   - the json object or a "json" string
//	error - error if any
func LoadSchemaFromFile(path string) (any, error) {
	// Expand homedir if applicable
	fullPath, err := paths.ExpandHomeDir(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	// Special case: file only contains "json"
	// This is the arbitrary and random structured json output
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "json" || trimmed == `"json"` {
		return "json", nil
	}

	// Default: parse file content as json object
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("invalid json schema: %w", err)
	}

	return schema, nil
}
