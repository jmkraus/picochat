package types_test

import (
	"picochat/paths"
	"picochat/types"
	"testing"
)

func TestNewHistory(t *testing.T) {
	h := types.NewHistory("hello world")

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
	h := types.NewHistory("init")

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
	h := types.NewHistory("persist me")
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
