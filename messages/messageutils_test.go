package messages

import (
	"math"
	"picochat/paths"
	"strings"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	h := NewHistory("persist me", 5)
	if err := h.add("user", "", "data", ""); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	filename, err := SaveHistoryToFile("", h)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if !strings.HasSuffix(filename, paths.HistorySuffix) {
		t.Fatalf("expected filename suffix %q, got %q", paths.HistorySuffix, filename)
	}

	loaded, err := LoadHistoryFromFile(filename)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Len() != h.Len() {
		t.Errorf("loaded history length mismatch: got %d, want %d", loaded.Len(), h.Len())
	}
}

func TestSaveAndLoad_DropsReasoning(t *testing.T) {
	h := NewHistory("persist me", 5)
	if err := h.add("assistant", "internal reasoning that must not be persisted", "visible content", ""); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	filename, err := SaveHistoryToFile("", h)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadHistoryFromFile(filename)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Len() != 2 {
		t.Fatalf("expected loaded history length 2, got %d", loaded.Len())
	}
	if loaded.Messages[1].Content != "visible content" {
		t.Errorf("expected loaded content to match; got %q", loaded.Messages[1].Content)
	}
	if loaded.Messages[1].Reasoning != "" {
		t.Errorf("expected reasoning to be empty after load, got %q", loaded.Messages[1].Reasoning)
	}
}

func TestSaveHistoryToFile_RejectsInvalidFilename(t *testing.T) {
	h := NewHistory("persist me", 5)

	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	if _, err := SaveHistoryToFile("#12", h); err == nil {
		t.Fatal("expected error for filename starting with '#', got nil")
	}
}

func TestLoadHistoryFromFile_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	if _, err := LoadHistoryFromFile("missing"); err == nil {
		t.Fatal("expected error for missing history file, got nil")
	}
}

func TestCalculateTokens(t *testing.T) {
	if got := CalculateTokens(""); got != 0 {
		t.Fatalf("expected 0 for empty string, got %.1f", got)
	}

	got := CalculateTokens("one two three")
	want := 3.9
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("expected %.1f, got %.1f", want, got)
	}
}
