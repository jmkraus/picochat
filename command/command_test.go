package command

import (
	"os"
	"path/filepath"
	// "picochat/command"
	"fmt"
	"picochat/messages"
	"picochat/paths"
	"strings"
	"testing"
)

func TestHandleClear(t *testing.T) {
	h := messages.NewHistory("initial system prompt", 50)
	h.AddUser("hello", "")

	result := HandleCommand("/clear", h, strings.NewReader(""))
	if !h.IsEmpty() {
		t.Error("expected history to be cleared except system prompt")
	}
	if !strings.Contains(result.Info, "cleared") {
		t.Errorf("unexpected output: %s", result.Info)
	}
}

func TestHandleHelp(t *testing.T) {
	h := messages.NewHistory("prompt", 50)
	result := HandleCommand("/help", h, strings.NewReader(""))
	if !strings.Contains(result.Output, "/save") {
		t.Errorf("expected help to contain /save, got: %s", result.Info)
	}
}

func TestHandleLoad_WithFilename(t *testing.T) {
	const dummyChat = `[
  {
    "role": "system",
    "content": "You are a LLM"
  },
  {
    "role": "user",
    "content": "Hello, are you there?"
  },
  {
    "role": "assistant",
    "content": "How can I assist you today?"
  }
]`

	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	err := os.WriteFile(filepath.Join(tmpDir, "dummy.chat"), []byte(dummyChat), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy.chat: %v", err)
	}

	h := messages.NewHistory("system prompt", 50)
	input := strings.NewReader("dummy.chat\n")

	result := HandleCommand("/load", h, input)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if result.Info != "History file \"dummy.chat\" loaded." {
		t.Fatalf("unexpected info: %q", result.Info)
	}

	loaded := h.Get()
	if len(loaded) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(loaded))
	}

	if loaded[0].Role != messages.RoleSystem || loaded[0].Content != "You are a LLM" {
		t.Fatalf("unexpected message[0]: %#v", loaded[0])
	}
	if loaded[1].Role != messages.RoleUser || loaded[1].Content != "Hello, are you there?" {
		t.Fatalf("unexpected message[1]: %#v", loaded[1])
	}
	if loaded[2].Role != messages.RoleAssistant || loaded[2].Content != "How can I assist you today?" {
		t.Fatalf("unexpected message[2]: %#v", loaded[2])
	}
}

func TestHandleCommand_Paste_UsesRuneCount(t *testing.T) {
	prevReadClipboard := readClipboard
	t.Cleanup(func() {
		readClipboard = prevReadClipboard
	})

	readClipboard = func() (string, error) {
		return "Hello, 世界", nil
	}

	history := messages.NewHistory("sys", 10)
	result := HandleCommand("/paste", history, strings.NewReader(""))

	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}
	if result.Pasted != "Hello, 世界" {
		t.Fatalf("unexpected pasted payload: %q", result.Pasted)
	}
	wantInfo := "Pasted 9 characters from clipboard."
	if result.Info != wantInfo {
		t.Fatalf("info = %q, want %q", result.Info, wantInfo)
	}
}

func TestHandleCommand_Paste_ReadClipboardError(t *testing.T) {
	prevReadClipboard := readClipboard
	t.Cleanup(func() {
		readClipboard = prevReadClipboard
	})

	readClipboard = func() (string, error) {
		return "", fmt.Errorf("clipboard read failed")
	}

	history := messages.NewHistory("sys", 10)
	result := HandleCommand("/paste", history, strings.NewReader(""))

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(result.Error.Error(), "clipboard read failed") {
		t.Fatalf("unexpected error: %v", result.Error)
	}
}
