package utils

import (
	"fmt"
	"os"
	"picochat/paths"
	"picochat/requests"
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
// Parameters: None
// Returns: (string, error)
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
// Parameters: baseUrl string
// Returns: (string, error)
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
// Parameters: list []string
// Returns: None
func setHistoryList(list []string) {
	HistoryList = list
	HistoryMap = make(map[int]string)
	for i, name := range list {
		HistoryMap[i+1] = name
	}
}

func GetHistoryByIndex(i int) (string, bool) {
	val, ok := HistoryMap[i]
	return val, ok
}

// Filled by /models command
// Parameters: list []string
// Returns: None
func setModelsList(list []string) {
	ModelsList = list
	ModelsMap = make(map[int]string)
	for i, name := range list {
		ModelsMap[i+1] = name
	}
}

func GetModelsByIndex(i int) (string, bool) {
	val, ok := ModelsMap[i]
	return val, ok
}

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

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
