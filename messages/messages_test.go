package messages

import (
	"fmt"
	"testing"
)

func TestNewHistory(t *testing.T) {
	h := NewHistory("hello world", 5)

	if len(h.Get()) != 1 {
		t.Fatalf("expected 1 message, got %d", len(h.Get()))
	}

	if h.Get()[0].Role != RoleSystem {
		t.Errorf("expected first message role to be 'system', got %s", h.Get()[0].Role)
	}

	if !h.IsEmpty() {
		t.Errorf("expected IsEmpty() to return true when only system prompt is present")
	}
}

func TestAddAndClear(t *testing.T) {
	h := NewHistory("init", 5)

	h.add(RoleUser, "", "Hello!", "")
	h.add(RoleAssistant, "", "Hi there!", "")

	if h.Len() != 3 {
		t.Errorf("expected 3 messages, got %d", h.Len())
	}

	h.ClearExceptSystemPrompt()

	if h.Len() != 1 || h.Get()[0].Role != RoleSystem {
		t.Errorf("clear should leave only system prompt")
	}
}

func TestDiscard(t *testing.T) {
	t.Run("only system prompt - no change", func(t *testing.T) {
		h := NewHistory("sys", 5)

		h.Discard()

		if h.Len() != 1 {
			t.Fatalf("expected length 1, got %d", h.Len())
		}
		if h.Get()[0].Role != RoleSystem {
			t.Fatalf("expected first role %q, got %q", RoleSystem, h.Get()[0].Role)
		}
	})

	t.Run("last message is assistant - remove last", func(t *testing.T) {
		h := NewHistory("sys", 5)
		_ = h.add(RoleUser, "", "hello", "")
		_ = h.add(RoleAssistant, "", "hi", "")

		h.Discard()

		if h.Len() != 2 {
			t.Fatalf("expected length 2 after discard, got %d", h.Len())
		}
		last := h.GetLast()
		if last.Role != RoleUser {
			t.Fatalf("expected last role %q, got %q", RoleUser, last.Role)
		}
	})

	t.Run("last message is not assistant - no change", func(t *testing.T) {
		h := NewHistory("sys", 5)
		_ = h.add(RoleUser, "", "hello", "")

		h.Discard()

		if h.Len() != 2 {
			t.Fatalf("expected length 2, got %d", h.Len())
		}
		last := h.GetLast()
		if last.Role != RoleUser {
			t.Fatalf("expected last role %q, got %q", RoleUser, last.Role)
		}
	})
}

func TestChatHistory_add_InvalidRole(t *testing.T) {
	h := &ChatHistory{}

	initialLen := h.Len()

	err := h.add("alien", "", "we come in peace", "")

	if err == nil {
		t.Fatalf("expected error for invalid role, got nil")
	}

	if h.Len() != initialLen {
		t.Errorf("expected message length to remain %d, got %d", initialLen, h.Len())
	}

	expectedErr := "invalid role \"alien\""
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message:\nwant: %q\ngot:  %q", expectedErr, err.Error())
	}
}

func TestCompressKeepsLimitAndPrompt(t *testing.T) {
	h := NewHistory("system prompt", 5)

	for i := 1; i <= 10; i++ {
		h.add(RoleUser, "", fmt.Sprintf("msg %d", i), "")
	}

	if h.Len() != 5 {
		t.Errorf("expected 5 messages, got %d", h.Len())
	}

	if h.Get()[0].Role != RoleSystem {
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
	h := NewHistory("system prompt", 10)

	// add 9 messages to fill context
	for i := 1; i < 10; i++ {
		h.add(RoleUser, "", fmt.Sprintf("msg %d", i), "")
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
	if h.Messages[0].Role != RoleSystem {
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
	h := NewHistory("system prompt", 5)

	// add 4 User/Assistant messages → +1 System = 5
	for i := 1; i <= 4; i++ {
		h.add(RoleUser, "", "Hello", "")
		h.add(RoleAssistant, "", "Hi there", "")
	}

	if h.Len() != 5 {
		t.Errorf("Expected history length to be 5, got %d", h.Len())
	}

	if !h.MaxContextReached {
		t.Errorf("Expected context limit to be reached")
	}

	// Now add another message - this should trigger Compress()
	h.add(RoleUser, "", "Overflow message", "")

	if h.Len() != 5 {
		t.Errorf("After compress, expected length 5, got %d", h.Len())
	}

	// Ensure that System prompt is kept
	first := h.Get()[0]
	if first.Role != RoleSystem {
		t.Errorf("Expected first message to be system prompt, got %s", first.Role)
	}
}

func TestReplaceMessages(t *testing.T) {
	h := NewHistory("init", 5)
	h.add(RoleUser, "", "first", "")

	newMessages := []Message{
		{Role: RoleSystem, Content: "replaced system"},
		{Role: RoleUser, Content: "replaced message"},
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
	h := NewHistory("short prompt", 5)
	h.add(RoleUser, "", "This is a short message", "")
	h.add(RoleAssistant, "", "This is a slightly longer reply that includes a few more words.", "")

	tokens := h.EstimateTokens()
	if tokens <= 0 {
		t.Errorf("expected positive token count, got %.1f", tokens)
	}
}

func TestAddUser(t *testing.T) {
	h := NewHistory("sys", 5)

	if err := h.AddUser("hello", ""); err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	last := h.GetLast()
	if last.Role != RoleUser {
		t.Fatalf("last role = %q, want %q", last.Role, RoleUser)
	}
	if last.Content != "hello" {
		t.Fatalf("last content = %q, want %q", last.Content, "hello")
	}
}

func TestAddAssistant(t *testing.T) {
	h := NewHistory("sys", 5)

	if err := h.AddAssistant("internal reasoning", "visible answer"); err != nil {
		t.Fatalf("AddAssistant failed: %v", err)
	}

	last := h.GetLast()
	if last.Role != RoleAssistant {
		t.Fatalf("last role = %q, want %q", last.Role, RoleAssistant)
	}
	if last.Reasoning != "internal reasoning" {
		t.Fatalf("last reasoning = %q, want %q", last.Reasoning, "internal reasoning")
	}
	if last.Content != "visible answer" {
		t.Fatalf("last content = %q, want %q", last.Content, "visible answer")
	}
}

func TestGetByIndex(t *testing.T) {
	h := NewHistory("sys", 5)
	_ = h.add(RoleUser, "", "hello", "")

	t.Run("valid index", func(t *testing.T) {
		msg, err := h.GetByIndex(1)
		if err != nil {
			t.Fatalf("GetByIndex returned error: %v", err)
		}
		if msg.Role != RoleUser || msg.Content != "hello" {
			t.Fatalf("unexpected message: %+v", msg)
		}
	})

	t.Run("negative index", func(t *testing.T) {
		_, err := h.GetByIndex(-1)
		if err == nil {
			t.Fatal("expected error for negative index, got nil")
		}
	})

	t.Run("out of bounds index", func(t *testing.T) {
		_, err := h.GetByIndex(99)
		if err == nil {
			t.Fatal("expected error for out-of-bounds index, got nil")
		}
	})
}

func TestMaxCtx(t *testing.T) {
	h := NewHistory("sys", 42)
	if h.MaxCtx() != 42 {
		t.Fatalf("MaxCtx = %d, want %d", h.MaxCtx(), 42)
	}
}

func TestGetLastMessage(t *testing.T) {
	h := NewHistory("sys", 5)
	last := h.GetLast()
	if last.Role != RoleSystem {
		t.Errorf("expected last message to be system at init, got %s", last.Role)
	}

	h.add(RoleUser, "", "msg", "")
	last = h.GetLast()
	if last.Role != RoleUser || last.Content != "msg" {
		t.Errorf("unexpected last message: %+v", last)
	}
}

func TestGetLastRole(t *testing.T) {
	h := NewHistory("init", 5)
	h.add(RoleUser, "", "hello", "")
	h.add(RoleAssistant, "", "hi", "")
	h.add(RoleUser, "", "bye", "")

	t.Run("find last user", func(t *testing.T) {
		msg, ok := h.GetLastRole(RoleUser)
		if !ok {
			t.Fatalf("expected to find user message, got none")
		}
		if msg.Content != "bye" {
			t.Errorf("expected last user message 'bye', got %q", msg.Content)
		}
	})

	t.Run("find last assistant", func(t *testing.T) {
		msg, ok := h.GetLastRole(RoleAssistant)
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
		empty := NewHistory("empty", 5)
		_, ok := empty.GetLastRole(RoleUser)
		if ok {
			t.Errorf("expected no match in empty history, but got one")
		}
	})
}

func TestIsEmpty(t *testing.T) {
	h := NewHistory("only system", 5)
	if !h.IsEmpty() {
		t.Errorf("expected history to be empty (only system prompt)")
	}
	h.add(RoleUser, "", "hi", "")
	if h.IsEmpty() {
		t.Errorf("expected history to be non-empty after user message")
	}
}

func TestEstimateTokens_IncludesReasoning(t *testing.T) {
	h := NewHistory("sys", 5)

	// Baseline: only system prompt
	base := h.EstimateTokens()

	// Reasoning-only message (empty content) should raise number of tokens
	// because EstimateTokens counts Reasoning + Content.
	if err := h.add(RoleAssistant, "This is hidden reasoning text", "", ""); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	withReasoning := h.EstimateTokens()
	if withReasoning <= base {
		t.Errorf("expected tokens to increase when reasoning is present; base=%.1f withReasoning=%.1f", base, withReasoning)
	}
}

func TestKeep(t *testing.T) {
	newHistory := func() *ChatHistory {
		h := NewHistory("sys", 10)
		_ = h.add(RoleUser, "", "u1", "")
		_ = h.add(RoleAssistant, "", "a1", "")
		_ = h.add(RoleUser, "", "u2", "")
		return h
	}

	t.Run("negative index", func(t *testing.T) {
		h := newHistory()
		before := h.Len()
		if ok := h.Keep(-1); ok {
			t.Fatal("expected false for negative index, got true")
		}
		if h.Len() != before {
			t.Fatalf("history length changed: got %d, want %d", h.Len(), before)
		}
	})

	t.Run("out of bounds index", func(t *testing.T) {
		h := newHistory()
		before := h.Len()
		if ok := h.Keep(h.Len()); ok {
			t.Fatal("expected false for out-of-bounds index, got true")
		}
		if h.Len() != before {
			t.Fatalf("history length changed: got %d, want %d", h.Len(), before)
		}
	})

	t.Run("keep system only", func(t *testing.T) {
		h := newHistory()
		if ok := h.Keep(0); !ok {
			t.Fatal("expected true for index 0, got false")
		}
		if h.Len() != 1 {
			t.Fatalf("expected length 1, got %d", h.Len())
		}
		if h.GetLast().Role != RoleSystem {
			t.Fatalf("expected last role %q, got %q", RoleSystem, h.GetLast().Role)
		}
	})

	t.Run("keep up to middle index", func(t *testing.T) {
		h := newHistory()
		if ok := h.Keep(2); !ok {
			t.Fatal("expected true for valid index, got false")
		}
		if h.Len() != 3 {
			t.Fatalf("expected length 3, got %d", h.Len())
		}
		last := h.GetLast()
		if last.Role != RoleAssistant || last.Content != "a1" {
			t.Fatalf("unexpected last message after keep: %+v", last)
		}
	})
}
