package clipb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

func isTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

func copyToTmuxBufferStdin(text string) error {
	cmd := exec.Command("tmux", "load-buffer", "-")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func ReadClipboard() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("clipboard read failed: %w", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("clipboard is empty")
	}
	return text, nil
}

func WriteClipboard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return fmt.Errorf("clipboard write failed: %w", err)
	}
	if isTmuxSession() {
		err := copyToTmuxBufferStdin(text)
		if err != nil {
			return fmt.Errorf("tmux clipboard write failed: %w", err)
		}
	}
	return nil
}
