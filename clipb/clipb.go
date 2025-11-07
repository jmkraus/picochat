package clipb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

// isTmuxSession returns whether the current process is running inside a tmux session.
//
// Parameters:
//
//	none
//
// Returns:
//
//	bool - app is running in a tmux session (true / false)
func isTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

// copyToTmuxBufferStdin loads the given text into the tmux buffer via stdin.
//
// Parameters:
//
//	text string - the text to load into the tmux buffer
//
// Returns:
//
//	error - any error encountered while writing the tmux buffer
func copyToTmuxBufferStdin(text string) error {
	cmd := exec.Command("tmux", "load-buffer", "-")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// ReadClipboard reads the current contents of the system clipboard.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the clipboard contents
//	error  - any error encountered while reading the clipboard
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

// WriteClipboard writes the given text to the system clipboard and, if running inside tmux,
// also copies it to the tmux buffer.
//
// Parameters:
//
//	text string - the text to write to the clipboard
//
// Returns:
//
//	error - any error encountered while writing to the clipboard or tmux buffer
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
