package clipb

import (
	"os"
	"strings"
	"testing"
)

func TestIsTmuxSession(t *testing.T) {
	// Save and restore previous value
	prev := os.Getenv("TMUX")
	defer os.Setenv("TMUX", prev)

	os.Unsetenv("TMUX")
	if isTmuxSession() {
		t.Error("Expected false when TMUX not set")
	}

	os.Setenv("TMUX", "/tmp/tmux-1234/default,1234,0")
	if !isTmuxSession() {
		t.Error("Expected true when TMUX is set")
	}
}

func TestCopyToTmuxBufferStdin(t *testing.T) {
	if !isTmuxSession() {
		t.Skip("Skipping test because not in tmux session")
	}
	err := copyToTmuxBufferStdin("Hello tmux")
	if err != nil {
		t.Errorf("Failed to copy to tmux buffer: %v", err)
	}
}

func TestCopyFromTmuxBufferStdout(t *testing.T) {
	if !isTmuxSession() {
		t.Skip("Skipping test because not in tmux session")
	}

	expected := "Hello tmux stdout"
	if err := copyToTmuxBufferStdin(expected); err != nil {
		t.Fatalf("failed to prepare tmux buffer: %v", err)
	}

	got, err := copyFromTmuxBufferStdout()
	if err != nil {
		t.Fatalf("failed to read from tmux buffer: %v", err)
	}
	if strings.TrimSpace(got) != expected {
		t.Fatalf("unexpected tmux buffer content: got %q, want %q", strings.TrimSpace(got), expected)
	}
}
