package messages

import (
	"regexp"
	"strings"
)

// StripReasoning removes <think>...</think> tags from the input string and trims empty lines.
// Parameters:
//
//	s (string) - the input string.
//
// Returns:
//
//	string - the cleaned string.
func StripReasoning(s string) string {
	// 1st case: correct pair <think>...</think>
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	if re.MatchString(s) {
		cleaned := re.ReplaceAllString(s, "")
		return TrimEmptyLines(cleaned)
	}

	// 2nd case: only </think> exists
	if idx := strings.Index(s, "</think>"); idx != -1 {
		cleaned := s[idx+len("</think>"):]
		return TrimEmptyLines(cleaned)
	}

	// 3rd case: no changes
	return s
}

// TrimEmptyLines removes leading and trailing empty lines from the input string.
// Parameters:
//
//	s (string) - the input string.
//
// Returns:
//
//	string - the string without leading or trailing empty lines.
func TrimEmptyLines(s string) string {
	lines := strings.Split(s, "\n")

	// before
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	// after
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}

// ExtractCodeBlock extracts the first code block from a string formatted with triple backticks.
// Parameters:
//
//	s (string) - the input string containing code blocks.
//
// Returns:
//
//	string - the extracted code block content.
//	bool - true if a code block was found, false otherwise.
func ExtractCodeBlock(s string) (string, bool) {
	re := regexp.MustCompile("(?s)```\\w*\\n(.*?)```")
	match := re.FindStringSubmatch(s)
	if len(match) >= 2 {
		return match[1], true
	}
	return "", false
}

// CalculateTokens estimates the number of tokens in a string based on word count.
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
