package utils

import (
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
		ext   string
		first string
	)

	const (
		teststr = "中国是一个拥有悠久历史的文明古国。 Can you translate this for me? 😊"
		cmdline = "echo \"%s\" | picochat -model %s"
	)

	switch runtime.GOOS {
	case "windows":
		ext = "cmd"
		first = "@echo off\nchcp 65001"
	default:
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
	for _, m := range models {
		rows = append(rows, fmt.Sprintf(cmdline, teststr, m))
	}

	filename := fmt.Sprintf("test.%s", ext)
	if err := os.WriteFile(filename, []byte(strings.Join(rows, "\n")), 0755); err != nil {
		return fmt.Errorf("could not write test file: %w", err)
	}

	return nil
}
