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

	// Temporäres Verzeichnis erzwingen
	tmpDir := t.TempDir()
	paths.OverrideHistoryDir(tmpDir)

	filename, err := h.SaveHistoryToFile()
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
