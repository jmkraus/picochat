package utils

import (
	"fmt"
	"os"
	"picochat/paths"
	"picochat/requests"
	"strings"
)

func ListHistoryFiles() ([]string, error) {
	basePath, err := paths.GetHistoryDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		if !entry.IsDir() && HasSuffix(entry.Name()) {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

// showAvailableModels lists all models via /tags API call
func ShowAvailableModels(baseUrl string) (string, error) {
	models, err := requests.GetAvailableModels(baseUrl)
	if err != nil {
		return "", fmt.Errorf("Failed to get models: %v", err)
	}

	if len(models) == 0 {
		return "", fmt.Errorf("No models available.")
	}

	output := "Available models:\n"
	for _, name := range models {
		output += fmt.Sprintf(" - %s\n", name)
	}

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
