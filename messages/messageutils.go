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
//	filename string
//
// Returns:
//
//	*ChatHistory
//	error
func LoadHistoryFromFile(filename string) (*ChatHistory, error) {
	historyPath, err := paths.GetHistoryPath()
	if err != nil {
		return nil, fmt.Errorf("history path not found: %w", err)
	}
	if historyPath == "" {
		return nil, fmt.Errorf("empty history path returned")
	}

	filename = paths.EnsureSuffix(filepath.Base(filename), paths.HistorySuffix)
	fullPath := filepath.Join(historyPath, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read file %s failed: %w", fullPath, err)
	}

	var messages []Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("parse json in file %s failed: %w", filename, err)
	}

	return &ChatHistory{Messages: messages}, nil
}

// SaveHistoryToFile writes the chat history to a file in the history
// directory and returns the filename.
//
// Parameters:
//
//	filename (string)      - optional filename (or timestamp if omitted)
//	history (*ChatHistory) - the complete current chat history
//
// Returns:
//
//	string - the actual filename
//	error  - error if any
func SaveHistoryToFile(filename string, history *ChatHistory) (string, error) {
	if strings.HasPrefix(filename, "#") {
		return "", fmt.Errorf("filename must not start with '#'")
	}

	historyPath, err := paths.GetHistoryPath()
	if err != nil {
		return "", fmt.Errorf("history path not found: %w", err)
	}
	if historyPath == "" {
		return "", fmt.Errorf("empty history path returned")
	}

	if filename == "" {
		filename = paths.EnsureSuffix(time.Now().Format("2006-01-02_15-04-05"), paths.HistorySuffix)
	} else {
		filename = paths.EnsureSuffix(filepath.Base(filename), paths.HistorySuffix)
	}
	fullPath := filepath.Join(historyPath, filename)

	if !paths.FileExists(fullPath) {
		data, err := json.MarshalIndent(history.Messages, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal messages failed: %w", err)
		}

		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			return "", fmt.Errorf("write file %s failed: %w", filename, err)
		}

		return filename, nil
	} else {
		return "", fmt.Errorf("filename already exists")
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
