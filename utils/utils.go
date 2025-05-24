package utils

import (
	"fmt"
	"os"
	"picochat/paths"
	"picochat/requests"
	"strings"
)

// ListHistoryFiles shows all saved sessions in 'history' dir
func ListHistoryFiles() (string, error) {
	basePath, err := paths.GetHistoryDir()
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return "", err
	}

	var result []string
	for _, entry := range entries {
		if !entry.IsDir() && HasSuffix(entry.Name()) {
			result = append(result, entry.Name())
		}
	}

	if len(result) == 0 {
		return "", fmt.Errorf("No history files found.")
	}

	output := "Available history files:\n- " + strings.Join(result, "\n- ")

	return output, nil
}

// showAvailableModels lists all models via /tags API call
func ShowAvailableModels(baseUrl string) (string, error) {
	models, err := requests.GetAvailableModels(baseUrl)
	if err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("No models available.")
	}

	modelLines := make([]string, len(models))
	for i, name := range models {
		modelLines[i] = fmt.Sprintf(" - %s", name)
	}

	output := "Available models:\n" + strings.Join(modelLines, "\n")

	return output, nil
}

const fileSuffix = ".chat"

func HasSuffix(filename string) bool {
	return strings.HasSuffix(filename, fileSuffix)
}

func EnsureSuffix(filename string) string {
	if !HasSuffix(filename) {
		return filename + fileSuffix
	}
	return filename
}
