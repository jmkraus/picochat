package messages

import (
	"regexp"
	"strings"
)

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

func ExtractCodeBlock(s string) (string, bool) {
	re := regexp.MustCompile("(?s)```\\w*\\n(.*?)```")
	match := re.FindStringSubmatch(s)
	if len(match) >= 2 {
		return match[1], true
	}
	return "", false
}

func CalculateTokens(s string) int {
	words := strings.Fields(s)
	return int(float64(len(words)) * 1.3)
}
