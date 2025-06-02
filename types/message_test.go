package types_test

import (
	"fmt"
	"picochat/paths"
	"picochat/types"
	"testing"
)

func TestNewHistory(t *testing.T) {
	h := types.NewHistory("hello world", 5)

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
	h := types.NewHistory("init", 5)

	h.Add("user", "Hello!")
	h.Add("assistant", "Hi there!")

	if h.Len() != 3 {
		t.Errorf("expected 3 messages, got %d", h.Len())
	}

	h.ClearExceptSystemPrompt()

	if h.Len() != 1 || h.Get()[0].Role != "system" {
		t.Errorf("clear should leave only system prompt")
	}
}

func TestSaveAndLoad(t *testing.T) {
	h := types.NewHistory("persist me", 5)
	h.Add("user", "data")

	// Force temporary directory
	tmpDir := t.TempDir()
	paths.OverrideHistoryPath(tmpDir)

	filename, err := h.SaveHistoryToFile("")
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := types.LoadHistoryFromFile(filename)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Len() != h.Len() {
		t.Errorf("loaded history length mismatch: got %d, want %d", loaded.Len(), h.Len())
	}
}

func TestCompressKeepsLimitAndPrompt(t *testing.T) {
	h := types.NewHistory("system prompt", 5)

	for i := 1; i <= 10; i++ {
		h.Add("user", fmt.Sprintf("msg %d", i))
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
	h := types.NewHistory("system prompt", 10)

	// Add 9 messages to fill context
	for i := 1; i < 10; i++ {
		h.Add("user", fmt.Sprintf("msg %d", i))
	}

	// Invalid values
	if err := h.SetContextSize(3); err == nil {
		t.Errorf("Expected error for context size < 5, got nil")
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

	// Reduction to 5
	if err := h.SetContextSize(5); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(h.Messages) != 5 {
		t.Errorf("Expected 5 messages after trimming, got %d", len(h.Messages))
	}
	if h.Messages[0].Role != "system" {
		t.Errorf("First message should be system prompt")
	}
	if h.Messages[1].Content != "msg 6" {
		t.Errorf("Expected oldest kept message to be 'msg 6', got '%s'", h.Messages[1].Content)
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
	h := types.NewHistory("system prompt", 5)

	// Add 4 User/Assistant messages → +1 System = 5
	for i := 1; i <= 4; i++ {
		h.Add("user", "Hello")
		h.Add("assistant", "Hi there")
	}

	if h.Len() != 5 {
		t.Errorf("Expected history length to be 5, got %d", h.Len())
	}

	if !h.MaxContextReached {
		t.Errorf("Expected context limit to be reached")
	}

	// Now add another message - this should trigger Compress()
	h.Add("user", "Overflow message")

	if h.Len() != 5 {
		t.Errorf("After compress, expected length 5, got %d", h.Len())
	}

	// Ensure that System prompt is kept
	first := h.Get()[0]
	if first.Role != "system" {
		t.Errorf("Expected first message to be system prompt, got %s", first.Role)
	}
}
