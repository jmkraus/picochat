package clipb

import (
	"os"
	"testing"
)

func TestIsTmuxSession(t *testing.T) {
	// Save and restore previous value
	prev := os.Getenv("TMUX")
	defer os.Setenv("TMUX", prev)

	os.Unsetenv("TMUX")
	if IsTmuxSession() {
		t.Error("Expected false when TMUX not set")
	}

	os.Setenv("TMUX", "/tmp/tmux-1234/default,1234,0")
	if !IsTmuxSession() {
		t.Error("Expected true when TMUX is set")
	}
}

func TestCopyToTmuxBufferStdin(t *testing.T) {
	if !IsTmuxSession() {
		t.Skip("Skipping test because not in tmux session")
	}
	err := CopyToTmuxBufferStdin("Hello tmux")
	if err != nil {
		t.Errorf("Failed to copy to tmux buffer: %v", err)
	}
}
