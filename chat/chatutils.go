package chat

import (
	"regexp"
	"strings"
)

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
	found := false
	if reasoning, content, found = strings.Cut(s, "</think>"); found {
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

	start, end := 0, len(lines)
	for start < end && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	for start < end && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	return strings.Join(lines[start:end], "\n")
}
