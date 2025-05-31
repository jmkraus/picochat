package utils

import (
	"fmt"
	"os"
	"picochat/paths"
	"picochat/requests"
	"regexp"
	"strings"
)

// ListHistoryFiles shows all saved sessions in 'history' dir
func ListHistoryFiles() (string, error) {
	basePath, err := paths.GetHistoryPath()
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

	return FormatList(result, "history files", true), nil
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

	return FormatList(models, "models", true), nil
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

func FormatList(content []string, heading string, numbered bool) string {
	if len(content) == 0 {
		return fmt.Sprintf("No %s found.", strings.ToLower(heading))
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

func StripReasoning(answer string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleaned := re.ReplaceAllString(answer, "")
	return cleaned
}

func ExtractCodeBlock(s string) string {
	re := regexp.MustCompile("(?s)```\\w*\\n(.*?)```")
	match := re.FindStringSubmatch(s)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}
