package clipb

import (
	"fmt"
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

func TestReadClipboard(t *testing.T) {
	origReadClipboardAll := readClipboardAll
	origReadTmuxBuffer := readTmuxBuffer
	t.Cleanup(func() {
		readClipboardAll = origReadClipboardAll
		readTmuxBuffer = origReadTmuxBuffer
	})

	t.Run("clipboard success trimmed", func(t *testing.T) {
		t.Setenv("TMUX", "")
		readClipboardAll = func() (string, error) {
			return "  hello  ", nil
		}
		readTmuxBuffer = func() (string, error) {
			t.Fatal("tmux read should not be called")
			return "", nil
		}

		got, err := ReadClipboard()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "hello" {
			t.Fatalf("got %q, want %q", got, "hello")
		}
	})

	t.Run("clipboard error without tmux", func(t *testing.T) {
		t.Setenv("TMUX", "")
		readClipboardAll = func() (string, error) {
			return "", fmt.Errorf("read failed")
		}

		_, err := ReadClipboard()
		if err == nil || !strings.Contains(err.Error(), "clipboard read failed") {
			t.Fatalf("expected clipboard read failed error, got %v", err)
		}
	})

	t.Run("tmux fallback on clipboard error", func(t *testing.T) {
		t.Setenv("TMUX", "1")
		readClipboardAll = func() (string, error) {
			return "", fmt.Errorf("read failed")
		}
		readTmuxBuffer = func() (string, error) {
			return "  from tmux  ", nil
		}

		got, err := ReadClipboard()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "from tmux" {
			t.Fatalf("got %q, want %q", got, "from tmux")
		}
	})

	t.Run("tmux fallback error", func(t *testing.T) {
		t.Setenv("TMUX", "1")
		readClipboardAll = func() (string, error) {
			return "", fmt.Errorf("read failed")
		}
		readTmuxBuffer = func() (string, error) {
			return "", fmt.Errorf("tmux failed")
		}

		_, err := ReadClipboard()
		if err == nil || !strings.Contains(err.Error(), "tmux clipboard read failed") {
			t.Fatalf("expected tmux clipboard read failed error, got %v", err)
		}
	})

	t.Run("empty clipboard returns empty error", func(t *testing.T) {
		t.Setenv("TMUX", "")
		readClipboardAll = func() (string, error) {
			return "   ", nil
		}

		_, err := ReadClipboard()
		if err == nil || !strings.Contains(err.Error(), "clipboard is empty") {
			t.Fatalf("expected clipboard empty error, got %v", err)
		}
	})
}

func TestWriteClipboard(t *testing.T) {
	origWriteClipboardAll := writeClipboardAll
	origWriteTmuxBuffer := writeTmuxBuffer
	t.Cleanup(func() {
		writeClipboardAll = origWriteClipboardAll
		writeTmuxBuffer = origWriteTmuxBuffer
	})

	t.Run("clipboard write error", func(t *testing.T) {
		t.Setenv("TMUX", "")
		writeClipboardAll = func(string) error {
			return fmt.Errorf("write failed")
		}

		err := WriteClipboard("hello")
		if err == nil || !strings.Contains(err.Error(), "clipboard write failed") {
			t.Fatalf("expected clipboard write failed error, got %v", err)
		}
	})

	t.Run("clipboard write success without tmux", func(t *testing.T) {
		t.Setenv("TMUX", "")
		writeClipboardAll = func(string) error {
			return nil
		}
		writeTmuxBuffer = func(string) error {
			t.Fatal("tmux write should not be called")
			return nil
		}

		if err := WriteClipboard("hello"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("tmux write error", func(t *testing.T) {
		t.Setenv("TMUX", "1")
		writeClipboardAll = func(string) error {
			return nil
		}
		writeTmuxBuffer = func(string) error {
			return fmt.Errorf("tmux write failed")
		}

		err := WriteClipboard("hello")
		if err == nil || !strings.Contains(err.Error(), "tmux clipboard write failed") {
			t.Fatalf("expected tmux clipboard write failed error, got %v", err)
		}
	})

	t.Run("tmux write success", func(t *testing.T) {
		t.Setenv("TMUX", "1")
		writeClipboardAll = func(string) error {
			return nil
		}
		writeTmuxBuffer = func(string) error {
			return nil
		}

		if err := WriteClipboard("hello"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
