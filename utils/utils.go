package utils

import (
	"fmt"
	"picochat/requests"
	"strings"
	"unicode"
)

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
