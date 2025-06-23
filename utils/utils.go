package utils

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

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

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

func CopyToTmuxBufferStdin(text string) error {
	cmd := exec.Command("tmux", "load-buffer", "-")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
