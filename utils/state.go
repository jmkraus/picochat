package utils

import (
	"fmt"
	"os"
	"picochat/paths"
	"picochat/requests"
	"strings"
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

	SetHistoryList(history)
	return FormatList(history, "history files", true), nil
}

// ShowAvailableModels lists all models via /tags API call
func ShowAvailableModels(baseUrl string) (string, error) {
	models, err := requests.GetAvailableModels(baseUrl)
	if err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no models available.")
	}

	SetModelsList(models)
	return FormatList(models, "models", true), nil
}

// Filled by /list command
func SetHistoryList(list []string) {
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
func SetModelsList(list []string) {
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

	return fmt.Sprintf("Available %s:\n%s", strings.ToLower(heading), strings.Join(lines, "\n"))
}
