package clipb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

func CopyToTmuxBufferStdin(text string) error {
	cmd := exec.Command("tmux", "load-buffer", "-")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func GetFromClipboard() (string, error) {
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

func PutToClipboard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return fmt.Errorf("clipboard failed: %w", err)
	}
	if IsTmuxSession() {
		err := CopyToTmuxBufferStdin(text)
		if err != nil {
			return fmt.Errorf("tmux clipboard failed: %w", err)
		}
	}
	return nil
}
