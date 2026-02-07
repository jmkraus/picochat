package chat

import (
	"regexp"
	"strings"
)

// stripReasoning removes <think>...</think> tags from the input string
// and trims empty lines.
//
// Parameters:
//
//	s (string) - the input string.
//
// Returns:
//
//	string - the cleaned string.
func stripReasoning(s string) string {
	// 1st case: correct pair <think>...</think>
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	if re.MatchString(s) {
		cleaned := re.ReplaceAllString(s, "")
		return trimEmptyLines(cleaned)
	}

	// 2nd case: only </think> exists
	if _, after, found := strings.Cut(s, "</think>"); found {
		return trimEmptyLines(after)
	}

	// 3rd case: no changes
	return s
}

// splitReasoning separates reasoning and content if merged in the content part.
//
// Parameters:
//
//	s (string) - the input string.
//
// Returns:
//
//	reasoning (string) - the reasoning part
//	content (string)   - the content part (answer)
func splitReasoning(s string) (reasoning string, content string) {
	// 1st case: correct pair <think>...</think>
	re := regexp.MustCompile(`(?s)<think>(.*?)</think>`)
	if matches := re.FindStringSubmatch(s); len(matches) > 0 {
		reasoning = matches[1]
		content = re.ReplaceAllString(s, "")
		return trimEmptyLines(reasoning), trimEmptyLines(content)
	}

	// 2nd case: only </think> exists
	if before, after, found := strings.Cut(s, "</think>"); found {
		reasoning = before
		content = after
		return trimEmptyLines(reasoning), trimEmptyLines(content)
	}

	// 3rd case: no reasoning
	return "", trimEmptyLines(s)
}

// trimEmptyLines removes leading and trailing empty lines from the input string.
//
// Parameters:
//
//	s (string) - the input string.
//
// Returns:
//
//	string - the string without leading or trailing empty lines.
func trimEmptyLines(s string) string {
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
