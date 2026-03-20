package chat

import (
	"fmt"
	"math"
	"picochat/messages"
	"regexp"
	"strings"
	"time"
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

// elapsedTime returns the elapsed time in seconds and a formatted
// "MM:SS" string.  All calculations are performed in whole seconds
// to avoid floating-point rounding differences.
//
// Parameters:
//
//	t - start time
//
// Returns:
//
//	int    - total elapsed seconds
//	string - formatted elapsed time "0m 0s"
func elapsedTime(t time.Time) (int, string) {
	elapsed := time.Since(t)

	totalSeconds := int(elapsed.Seconds())

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	if minutes == 0 {
		return totalSeconds, fmt.Sprintf("%ds", seconds)
	}

	return totalSeconds, fmt.Sprintf("%dm %ds", minutes, seconds)
}

// tokenSpeed calculates the average number of tokens processed per
// unit of time (t).  It returns 0 when t is zero to avoid division by
// zero.  The result is rounded to one decimal place.
//
// Parameters:
//
//	t (int)    - elapsed time in seconds
//	s (string) - string containing the full reply
//
// Returns:
//
//	float64 - tokens per second
func tokenSpeed(t int, s string) float64 {
	if (t == 0) || (s == "") {
		return 0
	}

	tokens := messages.CalculateTokens(s)
	speed := float64(tokens) / float64(t)
	roundFactor := 10.0

	return math.Round(speed*roundFactor) / roundFactor
}
