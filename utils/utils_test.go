package utils

import (
	"os"
	"testing"
)

func TestFormatList_WithBullets(t *testing.T) {
	items := []string{"first.chat", "second.chat"}
	expected := "Available history files:\n - first.chat\n - second.chat"

	result := FormatList(items, "history files", false)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestFormatList_WithNumbers(t *testing.T) {
	items := []string{"model-a", "model-b", "model-c"}
	expected := "Available models:\n(01) model-a\n(02) model-b\n(03) model-c"

	result := FormatList(items, "models", true)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestFormatList_Empty(t *testing.T) {
	items := []string{}
	expected := "No items found."

	result := FormatList(items, "items", false)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestHasSuffix(t *testing.T) {
	if !HasSuffix("test.chat") {
		t.Error("Expected HasSuffix to return true for 'test.chat'")
	}
	if HasSuffix("test.txt") {
		t.Error("Expected HasSuffix to return false for 'test.txt'")
	}
}

func TestEnsureSuffix(t *testing.T) {
	if EnsureSuffix("foo") != "foo.chat" {
		t.Errorf("Expected suffix to be appended")
	}
	if EnsureSuffix("bar.chat") != "bar.chat" {
		t.Errorf("Suffix should not be duplicated")
	}
}

func TestStripReasoning(t *testing.T) {
	input := "Answer here <think>reasoning inside</think> end."
	expected := "Answer here  end."

	result := StripReasoning(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestExtractCodeBlock(t *testing.T) {
	text := "Some explanation.\n```go\nfmt.Println(\"hi\")\n```"
	code := ExtractCodeBlock(text)
	if code != "fmt.Println(\"hi\")\n" {
		t.Errorf("Unexpected extracted code block: %q", code)
	}
}

func TestExtractCodeBlock_Empty(t *testing.T) {
	text := "No code block here"
	code := ExtractCodeBlock(text)
	if code != "" {
		t.Errorf("Expected empty string, got %q", code)
	}
}

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
