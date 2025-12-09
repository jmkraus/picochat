package messages_test

import (
	"fmt"
	"picochat/messages"
	"picochat/paths"
	"testing"
)

func TestNewHistory(t *testing.T) {
	h := messages.NewHistory("hello world", 5)

	if len(h.Get()) != 1 {
		t.Fatalf("expected 1 message, got %d", len(h.Get()))
	}

	if h.Get()[0].Role != "system" {
		t.Errorf("expected first message role to be 'system', got %s", h.Get()[0].Role)
	}

	if !h.IsEmpty() {
		t.Errorf("expected IsEmpty() to return true when only system prompt is present")
	}
}

func TestAddAndClear(t *testing.T) {
	h := messages.NewHistory("init", 5)

	h.Add("user", "Hello!", "")
	h.Add("assistant", "Hi there!", "")

	if h.Len() != 3 {
		t.Errorf("expected 3 messages, got %d", h.Len())
	}

	h.ClearExceptSystemPrompt()

	if h.Len() != 1 || h.Get()[0].Role != "system" {
		t.Errorf("clear should leave only system prompt")
	}
}

func TestChatHistory_Add_InvalidRole(t *testing.T) {
	h := &messages.ChatHistory{}

	initialLen := h.Len()

	err := h.Add("alien", "we come in peace", "")

	if err == nil {
		t.Fatalf("expected error for invalid role, got nil")
	}

	if h.Len() != initialLen {
		t.Errorf("expected message length to remain %d, got %d", initialLen, h.Len())
	}

	expectedErr := "invalid role 'alien'"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message:\nwant: %q\ngot:  %q", expectedErr, err.Error())
	}
}

func TestSaveAndLoad(t *testing.T) {
	h := messages.NewHistory("persist me", 5)
	h.Add("user", "data", "")

	// Force temporary directory
	tmpDir := t.TempDir()
	paths.OverrideHistoryPath(tmpDir)

	filename, err := h.SaveHistoryToFile("")
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := messages.LoadHistoryFromFile(filename)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Len() != h.Len() {
		t.Errorf("loaded history length mismatch: got %d, want %d", loaded.Len(), h.Len())
	}
}

func TestCompressKeepsLimitAndPrompt(t *testing.T) {
	h := messages.NewHistory("system prompt", 5)

	for i := 1; i <= 10; i++ {
		h.Add("user", fmt.Sprintf("msg %d", i), "")
	}

	if h.Len() != 5 {
		t.Errorf("expected 5 messages, got %d", h.Len())
	}

	if h.Get()[0].Role != "system" {
		t.Error("expected first message to be system prompt")
	}

	first := h.Get()[1]
	if first.Content != "msg 7" {
		t.Errorf("expected first message to be msg 7, got %s", first.Content)
	}

	last := h.Get()[h.Len()-1]
	if last.Content != "msg 10" {
		t.Errorf("expected last message to be msg 10, got %s", last.Content)
	}
}

func TestSetContextSize(t *testing.T) {
	h := messages.NewHistory("system prompt", 10)

	// Add 9 messages to fill context
	for i := 1; i < 10; i++ {
		h.Add("user", fmt.Sprintf("msg %d", i), "")
	}

	// Invalid values
	if err := h.SetContextSize(2); err == nil {
		t.Errorf("Expected error for context size < 3, got nil")
	}
	if err := h.SetContextSize(200); err == nil {
		t.Errorf("Expected error for context size > 100, got nil")
	}

	// No effect at same size
	if err := h.SetContextSize(10); err != nil {
		t.Errorf("Expected no error for unchanged context, got %v", err)
	}
	if len(h.Messages) != 10 {
		t.Errorf("Expected 10 messages, got %d", len(h.Messages))
	}

	// Reduction to 3
	if err := h.SetContextSize(3); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(h.Messages) != 3 {
		t.Errorf("Expected 3 messages after trimming, got %d", len(h.Messages))
	}
	if h.Messages[0].Role != "system" {
		t.Errorf("First message should be system prompt")
	}
	if h.Messages[1].Content != "msg 8" {
		t.Errorf("Expected oldest kept message to be 'msg 8', got '%s'", h.Messages[1].Content)
	}

	// Increasing size (no changes for existing messages)
	prevLen := len(h.Messages)
	if err := h.SetContextSize(8); err != nil {
		t.Errorf("Unexpected error on enlargement: %v", err)
	}
	if len(h.Messages) != prevLen {
		t.Errorf("Message list should remain same on context enlargement")
	}
}

func TestAddAndCompress(t *testing.T) {
	h := messages.NewHistory("system prompt", 5)

	// Add 4 User/Assistant messages → +1 System = 5
	for i := 1; i <= 4; i++ {
		h.Add("user", "Hello", "")
		h.Add("assistant", "Hi there", "")
	}

	if h.Len() != 5 {
		t.Errorf("Expected history length to be 5, got %d", h.Len())
	}

	if !h.MaxContextReached {
		t.Errorf("Expected context limit to be reached")
	}

	// Now add another message - this should trigger Compress()
	h.Add("user", "Overflow message", "")

	if h.Len() != 5 {
		t.Errorf("After compress, expected length 5, got %d", h.Len())
	}

	// Ensure that System prompt is kept
	first := h.Get()[0]
	if first.Role != "system" {
		t.Errorf("Expected first message to be system prompt, got %s", first.Role)
	}
}

func TestReplaceMessages(t *testing.T) {
	h := messages.NewHistory("init", 5)
	h.Add("user", "first", "")

	newMessages := []messages.Message{
		{Role: "system", Content: "replaced system"},
		{Role: "user", Content: "replaced message"},
	}
	h.Replace(newMessages)

	if h.Len() != 2 {
		t.Errorf("expected 2 messages after replace, got %d", h.Len())
	}

	if h.Get()[0].Content != "replaced system" {
		t.Errorf("expected first message to be replaced system, got %s", h.Get()[0].Content)
	}
}

func TestEstimateTokens(t *testing.T) {
	h := messages.NewHistory("short prompt", 5)
	h.Add("user", "This is a short message", "")
	h.Add("assistant", "This is a slightly longer reply that includes a few more words.", "")

	tokens := h.EstimateTokens()
	if tokens <= 0 {
		t.Errorf("expected positive token count, got %.1f", tokens)
	}
}

func TestGetLastMessage(t *testing.T) {
	h := messages.NewHistory("sys", 5)
	last := h.GetLast()
	if last.Role != "system" {
		t.Errorf("expected last message to be system at init, got %s", last.Role)
	}

	h.Add("user", "msg", "")
	last = h.GetLast()
	if last.Role != "user" || last.Content != "msg" {
		t.Errorf("unexpected last message: %+v", last)
	}
}

func TestGetLastRole(t *testing.T) {
	h := messages.NewHistory("init", 5)
	h.Add("user", "hello", "")
	h.Add("assistant", "hi", "")
	h.Add("user", "bye", "")

	t.Run("find last user", func(t *testing.T) {
		msg, ok := h.GetLastRole("user")
		if !ok {
			t.Fatalf("expected to find user message, got none")
		}
		if msg.Content != "bye" {
			t.Errorf("expected last user message 'bye', got %q", msg.Content)
		}
	})

	t.Run("find last assistant", func(t *testing.T) {
		msg, ok := h.GetLastRole("assistant")
		if !ok {
			t.Fatalf("expected to find assistant message, got none")
		}
		if msg.Content != "hi" {
			t.Errorf("expected assistant message 'hi', got %q", msg.Content)
		}
	})

	t.Run("role not found", func(t *testing.T) {
		_, ok := h.GetLastRole("foobar")
		if ok {
			t.Errorf("expected no match for 'foobar', but got one")
		}
	})

	t.Run("empty history", func(t *testing.T) {
		empty := messages.NewHistory("empty", 5)
		_, ok := empty.GetLastRole("user")
		if ok {
			t.Errorf("expected no match in empty history, but got one")
		}
	})
}

func TestIsEmpty(t *testing.T) {
	h := messages.NewHistory("only system", 5)
	if !h.IsEmpty() {
		t.Errorf("expected history to be empty (only system prompt)")
	}
	h.Add("user", "hi", "")
	if h.IsEmpty() {
		t.Errorf("expected history to be non-empty after user message")
	}
}
