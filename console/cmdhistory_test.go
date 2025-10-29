package console

import "testing"

func TestCommandHistory(t *testing.T) {
	h := &CommandHistory{}

	// --- Test Add ---
	h.add("/help")
	h.add("/models")
	h.add("/info")

	if len(h.entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(h.entries))
	}

	// Duplicate add check
	h.add("/info")
	if len(h.entries) != 3 {
		t.Fatalf("duplicate add should not increase length, got %d", len(h.entries))
	}

	// --- Test Prev ---
	prev := h.prev()
	if prev != "/info" {
		t.Errorf("expected /info, got %s", prev)
	}

	prev = h.prev()
	if prev != "/models" {
		t.Errorf("expected /models, got %s", prev)
	}

	prev = h.prev()
	if prev != "/help" {
		t.Errorf("expected /help, got %s", prev)
	}

	// Prev at beginning should stay on first
	prev = h.prev()
	if prev != "/help" {
		t.Errorf("expected /help when at beginning, got %s", prev)
	}

	// --- Test Next ---
	next := h.next()
	if next != "/models" {
		t.Errorf("expected /models, got %s", next)
	}

	next = h.next()
	if next != "/info" {
		t.Errorf("expected /info, got %s", next)
	}

	// Next beyond end → should reset index and return empty
	next = h.next()
	if next != "" {
		t.Errorf("expected empty string after end, got %s", next)
	}
}
