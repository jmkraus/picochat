package messages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"picochat/paths"
	"strings"
	"time"
)

// LoadHistoryFromFile reads a chat history from a file and returns
// a ChatHistory instance.
//
// Parameters:
//
//	fileName string
//
// Returns:
//
//	*ChatHistory
//	error
func LoadHistoryFromFile(fileName string) (*ChatHistory, error) {
	historyPath, err := paths.GetHistoryPath()
	if err != nil {
		return nil, fmt.Errorf("history path not found: %w", err)
	}
	if historyPath == "" {
		return nil, fmt.Errorf("empty history path returned")
	}

	fileName = paths.EnsureSuffix(filepath.Base(fileName), paths.HistorySuffix)
	fullPath := filepath.Join(historyPath, fileName)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read file %s failed: %w", fullPath, err)
	}

	var messages []Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("parse json in file %s failed: %w", fileName, err)
	}

	return &ChatHistory{Messages: messages}, nil
}

// SaveHistoryToFile writes the chat history to a file in the history
// directory and returns the fileName.
//
// Parameters:
//
//	fileName (string)      - optional fileName (or timestamp if omitted)
//	history (*ChatHistory) - the complete current chat history
//
// Returns:
//
//	string - the actual fileName
//	error  - error if any
func SaveHistoryToFile(fileName string, history *ChatHistory) (string, error) {
	if strings.HasPrefix(fileName, "#") {
		return "", fmt.Errorf("fileName must not start with '#'")
	}

	historyPath, err := paths.GetHistoryPath()
	if err != nil {
		return "", fmt.Errorf("history path not found: %w", err)
	}
	if historyPath == "" {
		return "", fmt.Errorf("empty history path returned")
	}

	if fileName == "" {
		fileName = paths.EnsureSuffix(time.Now().Format("2006-01-02_15-04-05"), paths.HistorySuffix)
	} else {
		fileName = paths.EnsureSuffix(filepath.Base(fileName), paths.HistorySuffix)
	}
	fullPath := filepath.Join(historyPath, fileName)

	if !paths.FileExists(fullPath) {
		data, err := json.MarshalIndent(history.Messages, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal messages failed: %w", err)
		}

		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			return "", fmt.Errorf("write file %s failed: %w", fileName, err)
		}

		return fileName, nil
	} else {
		return "", fmt.Errorf("fileName already exists")
	}
}

// CalculateTokens estimates the number of tokens in a string based on word count.
//
// Parameters:
//
//	s (string) - the input string.
//
// Returns:
//
//	float64 - the estimated token count.
func CalculateTokens(s string) float64 {
	words := strings.Fields(s)
	return float64(len(words)) * 1.3
}
