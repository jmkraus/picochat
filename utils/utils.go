package utils

import (
	"regexp"
	"strings"
)

func StripReasoning(input string) string {
	// 1st case: correct pair <think>...</think>
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	if re.MatchString(input) {
		cleaned := re.ReplaceAllString(input, "")
		return TrimEmptyLines(cleaned)
	}

	// 2nd case: only </think> exists
	if idx := strings.Index(input, "</think>"); idx != -1 {
		cleaned := input[idx+len("</think>"):]
		return TrimEmptyLines(cleaned)
	}

	// 3rd case: no changes
	return input
}

func TrimEmptyLines(s string) string {
	lines := strings.Split(s, "\n")

	// Anfang: leere Zeilen entfernen
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	// Ende: leere Zeilen entfernen
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}

func ExtractCodeBlock(s string) string {
	re := regexp.MustCompile("(?s)```\\w*\\n(.*?)```")
	match := re.FindStringSubmatch(s)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}
