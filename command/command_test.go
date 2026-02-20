package command_test

import (
	"os"
	"path/filepath"
	"picochat/command"
	"picochat/messages"
	"picochat/paths"
	"strings"
	"testing"
)

func TestHandleClear(t *testing.T) {
	h := messages.NewHistory("initial system prompt", 50)
	h.AddUser("hello", "")

	result := command.HandleCommand("/clear", h, strings.NewReader(""))
	if !h.IsEmpty() {
		t.Error("expected history to be cleared except system prompt")
	}
	if !strings.Contains(result.Info, "cleared") {
		t.Errorf("unexpected output: %s", result.Info)
	}
}

func TestHandleHelp(t *testing.T) {
	h := messages.NewHistory("prompt", 50)
	result := command.HandleCommand("/help", h, strings.NewReader(""))
	if !strings.Contains(result.Info, "/save") {
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

	tmp := t.TempDir()

	paths.OverrideHistoryPath(tmp)
	t.Cleanup(func() {
		paths.OverrideHistoryPath("")
	})

	err := os.WriteFile(filepath.Join(tmp, "dummy.chat"), []byte(dummyChat), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy.chat: %v", err)
	}

	h := messages.NewHistory("system prompt", 50)
	input := strings.NewReader("dummy.chat\n")

	result := command.HandleCommand("/load", h, input)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if result.Info != "History loaded successfully." {
		t.Fatalf("unexpected info: %q", result.Info)
	}

	loaded := h.Get()
	if len(loaded) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(loaded))
	}

	if loaded[0].Role != "system" || loaded[0].Content != "You are a LLM" {
		t.Fatalf("unexpected message[0]: %#v", loaded[0])
	}
	if loaded[1].Role != "user" || loaded[1].Content != "Hello, are you there?" {
		t.Fatalf("unexpected message[1]: %#v", loaded[1])
	}
	if loaded[2].Role != "assistant" || loaded[2].Content != "How can I assist you today?" {
		t.Fatalf("unexpected message[2]: %#v", loaded[2])
	}
}
